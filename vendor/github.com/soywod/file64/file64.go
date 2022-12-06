package file64

import (
	"io/ioutil"
	"encoding/base64"
)

func Encode(path string) (string, error) {
	buff, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buff), nil
}

func Decode(code string, dest string) error {
	buff, err := base64.StdEncoding.DecodeString(code)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(dest, buff, 07440); err != nil {
		return err
	}

	return nil
}