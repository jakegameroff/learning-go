package cloud_database

import (
	"io"
)

type Service interface {
	Upload(key string, body io.Reader) error
	Download(key string) (io.ReadCloser, error)
	Delete(key string) error
}
