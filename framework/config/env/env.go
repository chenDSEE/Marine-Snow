package env

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

type envInfo struct {
	path string // path of environment variable file
	file string // name of environment variable file
	envs map[string]string
}

const (
	DEFAULT_MODE = "default"
	TEST_MODE    = "test"
	DEVELOP_MODE = "develop"
)

var environ envInfo

func init() {
	environ.path = "./config" // as a default dir path
	environ.file = ".env"

	environ.envs = make(map[string]string)
	environ.envs["ENV_MODE"] = "default"
}

// Priority: runtime-env > .env file > hard code env
func EnvInit(pathString string) bool {
	if len(pathString) > 0 {
		environ.path = pathString
	}

	/* update from .env file */
	// default fullPath is "./config/.env"
	pathToFile := path.Join(environ.path, environ.file)

	if _, err := os.Stat(pathToFile); os.IsNotExist(err) == false {
		file, err := os.Open(pathToFile)
		if err == nil {
			defer file.Close()
			fmt.Printf("Load .env from:[%s]\n", pathToFile)

			// load each line
			reader := bufio.NewReader(file)
			for {
				line, _, err := reader.ReadLine()
				if err == io.EOF {
					break
				}

				// .env data format:
				// KEY_STRING=VALUE_STRING
				pair := bytes.SplitN(line, []byte{'='}, 2)
				if len(pair) < 2 {
					// ignore error key-value pair
					continue
				}

				key := string(pair[0])
				val := string(pair[1])
				environ.envs[key] = val
			}
		}
	}

	/* update from .env OS environment variable */
	for _, entry := range os.Environ() {
		pair := strings.SplitN(entry, "=", 2)
		if len(pair) < 2 {
			continue
		}

		environ.envs[pair[0]] = pair[1]
	}

	return true
}

func Get(name string) string {
	if val, ok := environ.envs[name]; ok {
		return val
	}

	return ""
}

func IsExist(name string) bool {
	_, ok := environ.envs[name]
	return ok
}

func GetAll() map[string]string {
	return environ.envs
}

func GetBasicInfo() string {
	return fmt.Sprintf("From[%s], Mode[%s]", environ.path+environ.file, Get("ENV_MODE"))
}

func DumpAll() {
	fmt.Println("========== dump environ ==========")
	for key, val := range GetAll() {
		fmt.Printf("%-15s  =  %s\n", key, val)
	}
	fmt.Println("==================================")
}
