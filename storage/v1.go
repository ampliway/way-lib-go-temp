package storage

import "time"

const (
	MODULE_NAME = "storage"
)

type V1 interface {
	Save(config *SaveConfig) (string, error)
	Delete(objectName string) error
	Link(objectName string, expiration time.Duration) (string, error)
}

type SaveConfig struct {
	Name            string
	FilePath        string
	ContentType     string
	ContentEncoding string
	UserMetadata    map[string]string
}
