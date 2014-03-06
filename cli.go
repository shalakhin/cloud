package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/OShalakhin/cloud/storage"
	"github.com/OShalakhin/cloud/storage/cloudfiles"
	"github.com/ncw/swift/rs"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"regexp"
)

const (
	// cloudCoreFile cloud core config name
	cloudCoreFile = ".cloudcore"
	// cloudFile cloud container config name
	cloudFile = ".cloud"
	// cloudIgnoreFile ignore config name
	cloudIgnoreFile = ".cloudignore"
)

var (
	ignorelist = GetIgnoreList()
)

// GetConfigs returns configs parsed from .cloudcore, and .cloud
// func GetConfigs(name string) (*storage.Core, *storage.Cloud, *storage.Container, *storage.Container) {
// 	return
// }

// GetIgnoreList returns pattern to ignore paths based on .cloudignore
func GetIgnoreList() []string {
	var f *os.File
	var err error
	if f, err = os.Open(cloudIgnoreFile); err != nil {
		return []string{"a^"}
	}
	str := []string{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := scanner.Text()
		if text[0:2] != "//" {
			str = append(str, text)
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return str
}

// GetContainerURL gets URL for the container by name. If no container is assigned
// it uses the first container in the list
func GetContainerURL(name string) {
	// get .cloudcore
	var core storage.Core
	var err error
	if core, err = GetCore(); err != nil {
		panic(err)
	}
	// get .cloud
	var cloud storage.Cloud
	if cloud, err = GetCloud(); err != nil {
		panic(err)
	}
	// get container from .cloud
	var container storage.Container
	if container, err = GetContainer(name, &cloud); err != nil {
		panic(err)
	}
	var s storage.Storage
	var provider storage.Provider
	// define storage backend
	switch {
	case container.Provider == storage.CloudFiles:
		if provider, err = GetProvider(&container, &core); err != nil {
			panic(err)
		}
		fmt.Println("Container found:\t", container.Name)
		s = &cloudfiles.Storage{
			Provider:  provider,
			Container: container,
			Info:      &cloudfiles.Info{},
			Conn:      rs.RsConnection{},
		}
	}
	// Authenticate (after authentication it is possible to return URL
	if err = s.Authenticate(); err != nil {
		panic(err)
	}
	fmt.Println("Container url is:", s.GetURL().String())
}

// GetCore returns parsed .cloudcore struct
func GetCore() (storage.Core, error) {
	var core storage.Core
	var err error
	var u *user.User
	core = storage.Core{}
	if u, err = user.Current(); err != nil {
		return core, err
	}
	corepath := path.Join(u.HomeDir, cloudCoreFile)
	var data []byte
	if data, err = ioutil.ReadFile(corepath); err != nil {
		return core, err
	}
	if err = json.Unmarshal(data, &core); err != nil {
		return core, err
	}
	return core, nil
}

// GetProvider from container string from the core config
func GetProvider(container *storage.Container, core *storage.Core) (storage.Provider, error) {
	for _, v := range core.Providers {
		if v.Provider == container.Provider {
			fmt.Println("Provider:", v.Provider)
			return v, nil
		}
	}
	return storage.Provider{}, errors.New("no provider with such name found")
}

// GetCloud returns parsed .cloud struct
func GetCloud() (storage.Cloud, error) {
	var data []byte
	var err error
	cloud := storage.Cloud{}
	if data, err = ioutil.ReadFile(cloudFile); err != nil {
		return cloud, err
	}
	if err = json.Unmarshal(data, &cloud); err != nil {
		return cloud, err
	}
	return cloud, nil
}

// GetContainer by name
func GetContainer(name string, cloud *storage.Cloud) (storage.Container, error) {
	if len(cloud.Containers) == 0 {
		return storage.Container{}, errors.New(".cloud has no containers")
	}
	if name == "" {
		fmt.Println("First container to be used:", cloud.Containers[0].Name)
		return cloud.Containers[0], nil
	}
	for _, v := range cloud.Containers {
		if v.Name == name {
			return v, nil
		}
	}
	return storage.Container{}, errors.New("no container found")
}

// IsExists checks if config exists
func IsExists(filename string) bool {
	var err error
	if _, err = os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// CreateConfig with filename and defined structure
func CreateConfig(filename string, v interface{}) error {
	var err error
	var template []byte
	if template, err = json.MarshalIndent(v, "", "    "); err != nil {
		return err
	}
	if err = ioutil.WriteFile(filename, template, 0644); err != nil {
		return err
	}
	return nil
}

func initConfigs() {
	// core
	var u *user.User
	var err error
	if u, err = user.Current(); err != nil {
		panic(err)
	}
	cloudcorepath := path.Join(u.HomeDir, cloudCoreFile)
	if ok := IsExists(cloudcorepath); !ok {
		core := storage.Core{
			Providers: []storage.Provider{
				{
					Provider: "CloudFiles",
					Name:     "myaccountname",
					Key:      "mykeyhere",
					Secret:   "mysecrethere",
					AuthURL:  "LON",
				},
			},
		}
		if err := CreateConfig(cloudcorepath, core); err != nil {
			panic(err)
		}
		fmt.Println("Initializing file:\t", cloudcorepath)
	}
	// cloud
	if ok := IsExists(cloudFile); !ok {
		cloud := storage.Cloud{
			Containers: []storage.Container{
				{
					Provider: "CloudFiles",
					Name:     "containername",
				},
			},
		}
		if err := CreateConfig(cloudFile, cloud); err != nil {
			panic(err)
		}
		fmt.Println("Initializing file:\t", cloudFile)
	}
	// cloudignore
	if ok := IsExists(cloudIgnoreFile); !ok {
		if err = ioutil.WriteFile(cloudIgnoreFile, []byte("// Put here what to ignore. Syntax like .gitignore\n.cloud\n.cloudignore\n"), 0644); err != nil {
			panic(err)
		}
		fmt.Println("Initializing file:\t", cloudIgnoreFile)
	}
}

// Sync folder with defined container name string (default is the first container in the list in local .cloud)
func Sync(name string) {
	fmt.Printf("Beginning sync with container: \"%s\"\n", name)
	// get .cloudcore
	var core storage.Core
	var err error
	if core, err = GetCore(); err != nil {
		panic(err)
	}
	fmt.Println("Found\t.cloudcore")
	// get .cloud
	var cloud storage.Cloud
	if cloud, err = GetCloud(); err != nil {
		panic(err)
	}
	fmt.Println("Found\t.cloud")
	var c storage.Container
	if c, err = GetContainer(name, &cloud); err != nil {
		panic(err)
	}
	// get container provider
	var s storage.Storage
	var p storage.Provider
	switch {
	case c.Provider == storage.CloudFiles:
		if p, err = GetProvider(&c, &core); err != nil {
			panic(err)
		}
		fmt.Println("Container found:\t", c.Name)
		s = &cloudfiles.Storage{
			Provider:  p,
			Container: c,
			Info:      &cloudfiles.Info{},
			Conn:      rs.RsConnection{},
		}
	default:
		// TODO what would be better to write here
		fmt.Println("Something went wrong!")
		return
	}
	// Authenticate
	if err = s.Authenticate(); err != nil {
		panic(err)
	}
	fmt.Println("Authenticated")
	// walk files upload file to the cloud
	var dir string
	if dir, err = os.Getwd(); err != nil {
		panic(err)
	}
	if err = filepath.Walk(dir, func(filename string, info os.FileInfo, err error) error {
		// upload only files, not directories
		if !info.IsDir() {
			fp, err := filepath.Rel(dir, filename)
			if err != nil {
				panic(err)
			}
			if !IsIgnored(fp) {
				data, err := ioutil.ReadFile(fp)
				if err != nil {
					panic(err)
				}
				fmt.Println("Sync\t", fp)
				if err = s.Create(fp, data); err != nil {
					panic(err)
				}
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}
}

// IsIgnored path or not. Data is taken from .cloudignore
func IsIgnored(filename string) bool {
	var err error
	for _, v := range ignorelist {
		var re *regexp.Regexp
		if re, err = regexp.Compile("^" + v); err != nil {
			panic(err)
		}
		if re.MatchString(filename) {
			return true
		}
	}
	return false
}
