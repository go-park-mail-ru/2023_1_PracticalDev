package utils

import (
	"bytes"
	"mime/multipart"
)

type File struct {
	Name  string
	Bytes []byte
}

func CreateMultipartFormBody(values map[string]string,
	files map[string]File) (body *bytes.Buffer, contentType string, err error) {
	body = new(bytes.Buffer)
	mw := multipart.NewWriter(body)
	defer mw.Close()

	for key, value := range files {
		w, err := mw.CreateFormFile(key, value.Name)
		if err != nil {
			return nil, "", err
		}
		if _, err = w.Write(value.Bytes); err != nil {
			return nil, "", err
		}
	}

	for key, value := range values {
		err = mw.WriteField(key, value)
		if err != nil {
			return nil, "", err
		}
	}

	return body, mw.FormDataContentType(), nil
}
