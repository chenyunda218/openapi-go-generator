package main

import (
	"fmt"
	"os"

	model "github.com/chenyunda218/oopenapi-go-gin-genertor/openapi"
	"gopkg.in/yaml.v3"
)

func main() {
	root := "./openapigingenertor"
	packageName := "openapigingenertor"
	inputFile := "./openapi.yaml"
	argsWithoutProg := os.Args[1:]
	for i, v := range argsWithoutProg {
		if v == "-o" {
			root = argsWithoutProg[i+1]
		}
		if v == "-p" {
			packageName = argsWithoutProg[i+1]
		}
		if v == "-i" {
			inputFile = argsWithoutProg[i+1]
		}
	}

	var api map[string]interface{}
	dat, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Println(err)
	}
	yaml.Unmarshal(dat, &api)
	var openapi model.Openapi
	yaml.Unmarshal(dat, &openapi)
	openapi.Generate(root, packageName)

}
