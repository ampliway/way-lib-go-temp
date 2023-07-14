package v1

type Config struct {
	StorageEndpoint        string `json:"storage_endpoint"`
	StorageAccessKeyID     string `json:"storage_access_key_id"`
	StorageSecretAccessKey string `json:"storage_secret_access_key"`
	StorageSecure          bool   `json:"storage_secure"`
	ExpirationDays         int    `json:"expiration_days"`
}
