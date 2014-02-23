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

// GetIgnoreList returns pattern to ignore paths based on .cloudignore
func GetIgnoreList() []string {
	f, err := os.Open(cloudIgnoreFile)
	if err != nil {
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

// GetCore returns parsed .cloudcore struct
func GetCore() (storage.Core, error) {
	core := storage.Core{}
	// open file
	u, err := user.Current()
	if err != nil {
		return core, err
	}
	corepath := path.Join(u.HomeDir, cloudCoreFile)
	data, err := ioutil.ReadFile(corepath)
	if err != nil {
		return core, err
	}
	// parse it
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
	cloud := storage.Cloud{}
	// open file
	data, err := ioutil.ReadFile(cloudFile)
	if err != nil {
		return cloud, err
	}
	// parse
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
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// CreateConfig with filename and defined structure
func CreateConfig(filename string, v interface{}) error {
	template, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		panic(err)
	}
	if err = ioutil.WriteFile(filename, template, 0644); err != nil {
		panic(err)
	}
	return nil
}

func initConfigs() {
	// core
	u, err := user.Current()
	if err != nil {
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
	core, err := GetCore()
	if err != nil {
		panic(err)
	}
	fmt.Println("Found\t.cloudcore")
	// get .cloud
	cloud, err := GetCloud()
	if err != nil {
		panic(err)
	}
	fmt.Println("Found\t.cloud")
	// auth container
	// get container by name
	c, err := GetContainer(name, &cloud)
	if err != nil {
		panic(err)
	}
	// get container provider
	var s storage.Storage
	switch {
	case c.Provider == storage.CLOUDFILES:
		p, err := GetProvider(&c, &core)
		if err != nil {
			panic(err)
		}
		fmt.Println("Container found:\t", c.Name)
		if err != nil {
			panic(err)
		}
		s = &cloudfiles.Storage{
			Provider:   p,
			Container:  c,
			Info:       &cloudfiles.Info{},
			Connection: rs.RsConnection{},
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
	dir, err := os.Getwd()
	if err != nil {
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
	// TODO I am sure it can be done in more elegant and efficient way
	for _, v := range ignorelist {
		re, err := regexp.Compile("^" + v)
		if err != nil {
			panic(err)
		}
		// if match found file must be ignored
		if re.MatchString(filename) {
			return true
		}
	}
	return false
}
