// Package common provides functionality shared by all Remarkable devices
package common

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/fako1024/go-remarkable/device"
	"github.com/fako1024/gotools/shell"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
)

// Supported file extensions
const (
	ExtPDF  = ".pdf"
	ExtEPUB = ".epub"
)

// Device denotes a generic device (independent of exact model)
type Device struct {
	dataPath string
}

// NewDevice instantiates a new generic device
func NewDevice(dataPath string) *Device {
	return &Device{
		dataPath: dataPath,
	}
}

// Upload uploads a file (PDF / ePUB) to the device tree
func (d *Device) Upload(docs ...device.Document) error {

	for _, doc := range docs {

		// Grab the file extension and initialize a new random UUID
		ext := filepath.Ext(doc.Name)
		if !isValidExt(ext) {
			return fmt.Errorf("invalid extension: %s", ext)
		}
		id := uuid.New()

		// Write the document metadata
		if err := writeJSON(filepath.Join(d.dataPath, id.String()+".metadata"), &device.FileMetaData{
			LastModified: fmt.Sprintf("%d", int64(time.Nanosecond)*time.Now().UnixNano()/int64(time.Millisecond)),
			Type:         "DocumentType",
			Version:      1,
			VisibleName:  strings.TrimSuffix(filepath.Base(doc.Name), ext),
		}); err != nil {
			return err
		}

		// Write the device content metadata
		if err := writeJSON(filepath.Join(d.dataPath, id.String()+".content"), &device.FileContentData{
			ExtraMetadata: device.ExtraMetadata{},
			FileType:      ext[1:],
			Transform:     device.Transform{},
		}); err != nil {
			return err
		}

		// Write the file
		if err := ioutil.WriteFile(filepath.Join(d.dataPath, id.String()+ext), doc.Content, 0600); err != nil {
			return err
		}
	}

	// Restart xochitl to force a scan of new files
	if _, err := shell.Run("systemctl restart xochitl"); err != nil {
		return err
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func writeJSON(path string, v interface{}) error {
	data, err := jsoniter.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, data, 0600)
}

func isValidExt(ext string) bool {
	if ext == ExtPDF || ext == ExtEPUB {
		return true
	}

	return false
}
