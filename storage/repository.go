package storage

import "net/url"

type MemRepository interface {
	Save(encodedURL string, url *url.URL) (id string, err error)
	GetByID(id string) (url *url.URL, err error)
}
