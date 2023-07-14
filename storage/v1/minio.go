package v1

import (
	"encoding/xml"
	"log"
	"net/url"
	"time"

	"github.com/ampliway/way-lib-go/helper/id"
	"github.com/ampliway/way-lib-go/helper/reflection"
	"github.com/ampliway/way-lib-go/storage"
	"github.com/minio/minio-go"
	"github.com/minio/minio-go/v7/pkg/lifecycle"
)

var _ storage.V1 = (*Minio)(nil)

type Minio struct {
	client     *minio.Client
	bucketName string
	id         id.ID
}

func New(cfg *Config, id id.ID) (*Minio, error) {
	client, err := minio.New(
		cfg.StorageEndpoint,
		cfg.StorageAccessKeyID,
		cfg.StorageSecretAccessKey,
		cfg.StorageSecure,
	)
	if err != nil {
		log.Fatalln(err)
	}

	bucketName := reflection.AppNamePkg()
	exist, err := client.BucketExists(bucketName)
	if err != nil {
		log.Fatalln(err)
	}

	if !exist {
		err := client.MakeBucket(bucketName, "")
		if err != nil {
			log.Fatalln(err)
		}
	}

	if cfg.ExpirationDays > 0 {
		config := lifecycle.NewConfiguration()
		config.Rules = []lifecycle.Rule{
			{
				ID:     "expiration",
				Status: "Enabled",
				Expiration: lifecycle.Expiration{
					Days: lifecycle.ExpirationDays(cfg.ExpirationDays),
				},
			},
		}

		buf, err := xml.Marshal(config)
		if err != nil {
			log.Fatalln(err)
		}

		err = client.SetBucketLifecycle(bucketName, string(buf))
		if err != nil {
			log.Fatalln(err)
		}
	}

	return &Minio{
		client:     client,
		bucketName: bucketName,
		id:         id,
	}, nil
}

func (m *Minio) Save(config *storage.SaveConfig) (string, error) {
	if config == nil {
		return "", errConfigNull
	}

	if config.FilePath == "" {
		return "", errConfigFilePathEmpty
	}

	if config.Name == "" {
		config.Name = m.id.Random()
	}

	putOptions := minio.PutObjectOptions{
		UserMetadata:    config.UserMetadata,
		ContentType:     config.ContentType,
		ContentEncoding: config.ContentEncoding,
	}

	_, err := m.client.FPutObject(m.bucketName, config.Name, config.FilePath, putOptions)
	if err != nil {
		return config.Name, err
	}

	return config.Name, nil
}

func (m *Minio) Delete(objectName string) error {
	return m.client.RemoveObject(m.bucketName, objectName)
}

func (m *Minio) Link(objectName string, expiration time.Duration) (string, error) {
	url, err := m.client.Presign("GET", m.bucketName, objectName, expiration, url.Values{})
	if err != nil {
		return "", err
	}

	return url.String(), nil
}
