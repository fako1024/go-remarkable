package common

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/fako1024/go-remarkable/device"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
)

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
func (d *Device) Upload(name string, data []byte) error {

	// Grab the file extension and initialize a new random UUID
	ext := filepath.Ext(name)
	if !isValidExt(ext) {
		return fmt.Errorf("invalid extension: %s", ext)
	}
	id := uuid.New()

	// Write the document metadata
	if err := writeJSON(filepath.Join(d.dataPath, id.String()+".metadata"), &device.FileMetaData{
		LastModified: fmt.Sprintf("%d", int64(time.Nanosecond)*time.Now().UnixNano()/int64(time.Millisecond)),
		Type:         "DocumentType",
		Version:      1,
		VisibleName:  strings.TrimSuffix(filepath.Base(name), ext),
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
	return ioutil.WriteFile(filepath.Join(d.dataPath, id.String()+ext), data, 0644)
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func writeJSON(path string, v interface{}) error {
	data, err := jsoniter.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, data, 0644)
}

func isValidExt(ext string) bool {
	if ext == ExtPDF || ext == ExtEPUB {
		return true
	}

	return false
}
