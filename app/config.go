package main

import (
	"github.com/caarlos0/env/v11"
)

// Global config object, safe to use once .ParseFromEnv() has been called (usually in func main())
var Config config

type config struct {
	// URL that is exposed to end-user (can be https://shareit.example.com for instance)
	//
	// If left empty, it will be guessed from headers, but can be an issue when a reverse proxy is
	// being used.
	PublicURL string `env:"PUBLIC_URL"`

	// File storage URL. This is where all the content is persisted.
	// We support file://, mem://, s3:// URLs.
	//
	// see https://gocloud.dev/howto/blob for more details on URLs schemes supported
	BlobStorageURL string `env:"BLOB_STORAGE_URL" envDefault:"file://./storage?create_dir=true"`
}

func (c *config) ParseFromEnv() error {
	// Load the app config from env
	cfg, err := env.ParseAs[config]()
	if err != nil {
		return err
	}

	*c = cfg
	return nil
}
