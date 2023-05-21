package utils

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/constants"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
)

type File struct {
	Name  string
	Bytes []byte
}

func ReadBody(r *http.Request, log *zap.Logger) ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error(constants.FailedReadRequestBody, zap.Error(err))
		return nil, errors.Wrap(pkgErrors.ErrReadBody, err.Error())
	}

	err = r.Body.Close()
	if err != nil {
		log.Error(constants.FailedCloseRequestBody, zap.Error(err))
	}

	return body, nil
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
