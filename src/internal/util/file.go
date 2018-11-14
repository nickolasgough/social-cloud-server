package util

import (
	"fmt"
	"strings"
)


func ParseContentType(filename string) string {
	fileparts := strings.Split(filename, ".")
	return fmt.Sprintf("image/%s", fileparts[len(fileparts)-1])
}