package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
)

const (
	// CLOUDCORE cloud core config name
	CLOUDCORE = ".cloudcore"
	// CLOUD cloud container config name
	CLOUD = ".cloud"
	// CLOUDIGNORE cloud ignore config name
	CLOUDIGNORE = ".cloudignore"
)

type (
	// Provider for the cloud like Amazon, Rackspace etc.
	Provider struct {
		Name   string `json:"name"`
		Key    string `json:"key"`
		Secret string `json:"secret"`
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
	// Containers is config for .cloud
	Containers struct {
		Containers []Container `json:"containers"`
	}
)

// IsExists checks if config exists
func IsExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// Create config
func Create(filename string, v interface{}) error {
	template, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(filename, template, 0644); err != nil {
		return err
	}
	return nil
}

func initConfigs() {
	// core
	u, err := user.Current()
	if err != nil {
		panic(err)
	}
	cloudcorepath := path.Join(u.HomeDir, CLOUDCORE)
	if ok := IsExists(cloudcorepath); !ok {
		core := Core{
			Providers: []Provider{
				{
					Name:   "CloudFiles",
					Key:    "mykeyhere",
					Secret: "mysecrethere",
				},
			},
		}
		if err := Create(cloudcorepath, core); err != nil {
			panic(err)
		}
		fmt.Println("Initializing file:\t", cloudcorepath)
	}
	// cloud
	if ok := IsExists(CLOUD); !ok {
		cloud := Containers{
			Containers: []Container{
				{
					Provider: "CloudFiles",
					Name:     "containername",
				},
			},
		}
		if err := Create(CLOUD, cloud); err != nil {
			panic(err)
		}
		fmt.Println("Initializing file:\t", CLOUD)
	}
	// TODO cloudignore
}
