package helpers

import (
	"encoding/json"
	"fmt"

	"github.com/gocarina/gocsv"
	"gopkg.in/yaml.v3"
)

func PrintJSON(obj interface{}) error {
	b, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func PrintYAML(obj interface{}) error {
	b, err := yaml.Marshal(obj)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func PrintCSV(obj interface{}) error {
	csv, err := gocsv.MarshalBytes(obj)
	if err != nil {
		return err
	}
	fmt.Println(string(csv))
	return nil
}
