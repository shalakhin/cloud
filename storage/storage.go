// Package storage provides abstract layer for cloud storages
// that should provide basically CRUD and Authentication methods.
package storage

type (
	// Storage interface for abstraction
	Storage interface {
		Create(container, filename string, data []byte) error
		Read(container, filename string) ([]byte, error)
		Update(container, filename string, data []byte) error
		// Upsert(filename string, data []byte) error
		Delete(container, filename string) error
		Authenticate() error
		GetContainer() (*Container, error)
	}

	// Provider for the cloud like Amazon, Rackspace etc.
	Provider struct {
		Provider string `json:"provider"`
		Name     string `json:"name"`
		Key      string `json:"key"`
		Secret   string `json:"secret,omitempty"`
		AuthURL  string `json:"auth_url,omitempty"`
	}

	// Core is ~/.cloudcore struct
	Core struct {
		Providers []Provider `json:"providers"`
	}

	// Container data to work with later
	Container struct {
		Provider string `json:"provider"`
		Name     string `json:"name"`
	}
	// Cloud is config for .cloud
	Cloud struct {
		Containers []Container `json:"containers"`
	}
)

const (
	// CLOUDFILES name for Rackspace CloudFiles storage
	CLOUDFILES = "CloudFiles"
	// S3 name for Amazon S3
	S3 = "S3"
)
