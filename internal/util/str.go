package util

import (
	"fmt"
	"io/ioutil"

	"gitlab.com/l3montree/microservices/libs/orchardclient"
)

func Str(s string) *string {
	return &s
}

func MustReadFile(filepath string) []byte {
	bytes, err := ioutil.ReadFile(filepath)
	orchardclient.FailOnError(err, fmt.Sprintf("could not read file: %s", filepath))
	return bytes
}
