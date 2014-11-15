// inspired by https://github.com/stathat/jconfig/blob/master/config.go

package configuration

import (
	"bytes"
	//"encoding/json"
	"log"
	"os"
)

// LoadConfig will load the appropriate configuration data from a JSON file
func LoadConfig(filename string) *Json {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("error loading config file %s: %s", filename, err)
	}
	defer f.Close()
	b := new(bytes.Buffer)
	_, err = b.ReadFrom(f)
	if err != nil {
		log.Fatalf("error loading config file %s: %s", filename, err)
	}
	result, err := NewJson(b.Bytes())
	if err != nil {
		log.Fatalf("error loading config file %s: %s", filename, err)
	}
	return result
}
