package main

import (
	"fmt"
	"os"

	model "github.com/chenyunda218/openapi-go-generator/openapi"
	"gopkg.in/yaml.v3"
)

func main() {
	root := "./openapigogenerator"
	packageName := "openapigogenerator"
	inputFile := "./openapi.yaml"
	argsWithoutProg := os.Args[1:]
	for i, v := range argsWithoutProg {
		if v == "-o" && len(argsWithoutProg) > i+1 {
			root = argsWithoutProg[i+1]
		}
		if v == "-p" && len(argsWithoutProg) > i+1 {
			packageName = argsWithoutProg[i+1]
		}
		if v == "-i" && len(argsWithoutProg) > i+1 {
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
