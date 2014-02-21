package storage

// Storage interface for abstraction
type (
	Storage interface {
		Create(filename string, data []byte) error
		Read(filename string) ([]byte, error)
		Update(filename string, data []byte) error
		Upsert(filename string, data []byte) error
		Delete(filename string) error
		Authenticate() error
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
