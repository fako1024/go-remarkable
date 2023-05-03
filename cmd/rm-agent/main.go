package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"

	"github.com/fako1024/go-remarkable/device"
	"github.com/fako1024/go-remarkable/device/rm2"
	"github.com/fako1024/go-remarkable/internal/images"
	"github.com/fako1024/incept"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/sirupsen/logrus"
)

type config struct {
	listenPort string
}

func main() {

	// Read / parse configuration from command-line
	cfg := parseConfig()
	log := logrus.StandardLogger()

	// Instantiate a new Remarkable device and ensure it is closed properly
	r, err := rm2.New()
	if err != nil {
		log.Fatalf("error instantiating Remarkable device: %s", err)
	}
	defer func() {
		if err := r.Close(); err != nil {
			log.Error(err)
		}
	}()

	i, err := incept.New()
	if err != nil {
		log.Fatalf("error instantiating self-reload functionality: %s", err)
	}

	// Instantiate new router
	router := fiber.New(fiber.Config{
		BodyLimit: 20 * 1024 * 1024, // 20 MB
	})
	router.Use(logger.New())
	router.Get("/screen", handleFrame(r))
	router.Get("/stream", handleStream(r))
	router.Get("/upload", handleUploadForm)
	router.Get("/tree", handleTree(r))
	router.Get("/download/:id", handleDownload(r))
	router.Post("/upload", handleUpload(r))
	router.Post("/update", func(c *fiber.Ctx) error {
		return i.Update(c.Body())
	})

	// Run the web server
	for {
		log.Error(router.Listen(":" + cfg.listenPort))
	}
}

func handleFrame(r device.Remarkable) func(c *fiber.Ctx) error {

	width, height := r.Dimensions()

	return func(c *fiber.Ctx) error {

		// Set quality
		opt := jpeg.Options{Quality: 80}
		if q := c.Context().QueryArgs().GetUintOrZero("quality"); q != 0 {
			if q < 0 || q > 100 {
				return fmt.Errorf("invalid quality: %d", q)
			}
			opt.Quality = q
		}

		// Get a single frame from the device
		data, err := r.Frame()
		if err != nil {
			return fmt.Errorf("error reading frame from framebuffer device: %w", err)
		}

		// Handle aspect ratio / orientation
		var img *image.Gray
		if c.Query("portrait") == "true" {
			bufRot := make([]byte, len(data))
			images.Transpose(bufRot, data, width, height)
			img = image.NewGray(image.Rect(0, 0, height, width))
			img.Pix = bufRot
		} else {
			img = image.NewGray(image.Rect(0, 0, width, height))
			img.Pix = data
		}

		// Encode and send the image
		c.Context().SetContentType("image/jpeg")
		imgBuf := new(bytes.Buffer)
		if err = jpeg.Encode(imgBuf, img, &opt); err != nil {
			return fmt.Errorf("error encoding frame: %w", err)
		}
		c.Response().SetBody(imgBuf.Bytes())

		return nil
	}
}

func handleStream(r device.Remarkable) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) (err error) {

		// Continuously stream the frames
		c.Response().SetBodyStreamWriter(func(w *bufio.Writer) {
			if err = r.NewStream(w); err != nil {
				return
			}
		})

		return err
	}
}

func handleTree(r device.Remarkable) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		tree, err := r.FileTree()
		if err != nil {
			return err
		}

		return c.JSON(tree.Flatten())
	}
}

func handleDownload(r device.Remarkable) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) (err error) {
		file, err := r.Download(c.Params("id"))
		if err != nil {
			return err
		}

		c.Attachment(file.Name)
		c.Response().SetBodyRaw(file.Content)
		return nil
	}
}

func handleUploadForm(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
	return c.SendString(staticHTML)
}

func handleUpload(r device.Remarkable) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) (err error) {

		form, err := c.MultipartForm()
		if err != nil {
			return err
		}

		var docs device.Documents
		for _, file := range form.File["file"] {
			f, err := file.Open()
			if err != nil {
				return err
			}
			defer func() {
				cerr := f.Close()
				if err == nil {
					err = cerr
				}
			}()

			fileData, err := ioutil.ReadAll(f)
			if err != nil {
				return err
			}
			docs = append(docs, device.Document{
				Name:    file.Filename,
				Content: fileData,
			})
		}

		if err := r.Upload(docs...); err != nil {
			return err
		}

		return c.RedirectBack("/upload")
	}
}

func parseConfig() (cfg config) {
	flag.StringVar(&cfg.listenPort, "l", "8090", "Port to listen on")
	flag.Parse()

	return
}
