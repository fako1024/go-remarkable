// Package common provides functionality shared by all Remarkable devices
package common

import (
	"fmt"
	"io/fs"
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

// FileTree returns a tree structure of all files on the device
func (d *Device) FileTree() (*device.Tree, error) {
	nodes := make(map[string]device.FileMetaData)
	err := filepath.WalkDir(d.dataPath, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ".metadata" {
			var fileMetaData device.FileMetaData
			if err := readJSON(s, &fileMetaData); err != nil {
				return err
			}

			nodes[d.Name()[:len(d.Name())-len(filepath.Ext(d.Name()))]] = fileMetaData
		}
		return nil
	})

	return buildTree(nodes), err
}

// Download retrieves a file (PDF / ePUB) from the device tree
func (d *Device) Download(id string) (*device.Document, error) {

	var fileMetaData device.FileMetaData
	if err := readJSON(filepath.Join(d.dataPath, id+".metadata"), &fileMetaData); err != nil {
		return nil, err
	}
	var fileContent device.FileContentData
	if err := readJSON(filepath.Join(d.dataPath, id+".content"), &fileContent); err != nil {
		return nil, err
	}

	//TODO: Only read pdfs * epubs, what to do with notebooks?

	data, err := ioutil.ReadFile(filepath.Clean(filepath.Join(d.dataPath, id+"."+fileContent.FileType)))
	if err != nil {
		return nil, err
	}

	return &device.Document{
		Name:    fileMetaData.VisibleName + "." + fileContent.FileType,
		Content: data,
	}, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func readJSON(path string, v interface{}) error {
	data, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		return err
	}

	return jsoniter.Unmarshal(data, v)
}

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

func buildTree(nodes map[string]device.FileMetaData) *device.Tree {

	tree := &device.Tree{
		Folders: make(map[string]*device.Folder),
	}
	buildLayer("", tree, nodes)

	return tree
}

func buildLayer(parentID string, folder *device.Folder, nodes map[string]device.FileMetaData) {
	for k, v := range nodes {
		if v.Parent == parentID {
			if v.Type == "CollectionType" {
				folder.Folders[v.VisibleName] = &device.Folder{
					ID:      k,
					Name:    v.VisibleName,
					Folders: make(map[string]*device.Folder),
				}
			} else {
				folder.Files = append(folder.Files, device.File{
					ID:   k,
					Name: v.VisibleName,
				})
			}
			delete(nodes, k)
		}
	}

	for _, subfolder := range folder.Folders {
		buildLayer(subfolder.ID, subfolder, nodes)
	}
}
