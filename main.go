package main

import (
	"fmt"
	"os"

	model "github.com/chenyunda218/openapi-go-generator/openapi"
	"github.com/dave/jennifer/jen"
	"github.com/getkin/kin-openapi/openapi3"
	"gopkg.in/yaml.v3"
)

const defaultRoot = "./openapigogenerator"
const defaultPackageName = "openapigogenerator"
const defaultInputFile = "./openapi.yaml"

func main() {
	root := defaultRoot
	packageName := defaultPackageName
	inputFile := defaultInputFile
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
	f := jen.NewFile("main")
	f.Func().Id("main").Params().Block(
		jen.Qual("fmt", "Println").Call(jen.Lit("Hello, world")),
	)
	fmt.Printf("%#v", f)
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile(inputFile)
	for k, v := range doc.Components.Schemas {
		fmt.Println(k)
		fmt.Println(v.Value.Type)
	}
}
