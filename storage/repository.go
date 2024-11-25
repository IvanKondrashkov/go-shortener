package storage

import (
	"net/url"

	"github.com/google/uuid"
)

type MemRepository interface {
	Save(id uuid.UUID, url *url.URL) (res uuid.UUID, err error)
	GetByID(id uuid.UUID) (res *url.URL, err error)
}
