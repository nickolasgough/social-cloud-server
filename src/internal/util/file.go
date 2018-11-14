package util

import (
	"fmt"
	"strings"
	"bytes"
	"image"
	"image/png"
	"image/jpeg"
	"cloud.google.com/go/storage"
)


func DecodeImageFile(filename string, imagebytes []byte) (string, image.Image, error) {
	contentType := ParseContentType(filename)
	var imagefile image.Image
	var err error
	if contentType == "image/png" {
		imagefile, err = png.Decode(bytes.NewReader(imagebytes))
	} else if contentType == "image/jpg" {
		imagefile, err = jpeg.Decode(bytes.NewReader(imagebytes))
	} else {
		err = fmt.Errorf("social-cloud: unsupported file type: %s\n", contentType)
	}
	return contentType, imagefile, err
}


func EncodeImageFile(writer *storage.Writer, imagefile image.Image) error {
	var err error
	if writer.ContentType == "image/png" {
		err = png.Encode(writer, imagefile)
	} else if writer.ContentType == "image/jpg" {
		err = jpeg.Encode(writer, imagefile, nil)
	} else {
		err = fmt.Errorf("social-cloud: unsupported file type: %s\n", writer.ContentType)
	}
	return err
}


func ParseContentType(filename string) string {
	fileparts := strings.Split(filename, ".")
	return fmt.Sprintf("image/%s", fileparts[len(fileparts)-1])
}