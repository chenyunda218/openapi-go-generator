# openapi-go-generator

openapi gin server api genertor. openapi-go-generator generate gin server and go interface. Auto binding Parameters. Included in path and in query Parameters.  
This project is just started. Please feel free to report any bugs and suggestion.

## CONFIG

| Flat | Description     | Default              |
| ---- | --------------- | -------------------- |
| -i   | Input yaml file | ./openapi.yaml       |
| -o   | Output path     | ./openapigingenertor |
| -p   | Package name    | openapigingenertor   |

## Schema Support Feature

| Type     | Supported |
| -------- | --------- |
| object   | ✓         |
| string   | ✓         |
| enum     | ✓         |
| boolean  | ✓         |
| int32    | ✓         |
| int64    | ✓         |
| float    | ✓         |
| double   | ✓         |
| array    | ✓         |
| required | ✓         |
| allOf    | ✓         |
| oneOf    | ✗         |

## Sample

```bash
git clone https://github.com/chenyunda218/openapi-go-generator
cd openapi-go-generator
go run main.go -o api -p api -i ./openapi.yaml
```

## Installation

Resource openapi yaml

```yaml
openapi: 3.0.3
info:
  title: Example
  description: |-
    Example
  contact:
    email: chenyunda218@gmail.com
  version: 0.0.1
servers:
  - url: http://localhost/api/v1
tags:
  - name: Pet
    description: Api of account
paths:
  /pets:
    post:
      tags:
        - Pet
      operationId: CreateCat
      requestBody:
        content:
          application/json:
            schema:
              $ref: "#/components/schemas/Cat"
      responses:
        "200":
          description: Updated
  /pets/{id}:
    get:
      tags:
        - Pet
      operationId: GetCat
      parameters:
        - name: id
          in: path
          schema:
            type: integer
            format: int64
          required: true
      responses:
        "200":
          description: Updated
components:
  schemas:
    Dog:
      type: object
      properties:
        bark:
          type: boolean
        breed:
          type: string
          enum: [Dingo, Husky, Retriever, Shepherd]
    Cat:
      type: object
      properties:
        hunts:
          type: boolean
        age:
          type: integer
```

Result

```go
package openapigogenerator

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Dog struct {
	Bark  string `json:"bark"`
	Breed string `json:"breed"`
}
type Cat struct {
	Hunts string `json:"hunts"`
	Age   int64  `json:"age"`
}
type PetApiInterface interface {
	CreateCat(gin_context *gin.Context, gin_body Cat)
	GetCat(gin_context *gin.Context, id int64)
}

func CreateCatBinder(api PetApiInterface) func(c *gin.Context) {
	return func(gin_context *gin.Context) {
		var cat Cat
		gin_context.ShouldBindJSON(&cat)
		api.CreateCat(gin_context, cat)
	}
}
func GetCatBinder(api PetApiInterface) func(c *gin.Context) {
	return func(gin_context *gin.Context) {
		id := gin_context.Param("id")
		api.GetCat(gin_context, stringToInt64(id))
	}
}
func PetApiInterfaceMounter(gin_router *gin.Engine, gwg_api_label PetApiInterface) {
	gin_router.POST("/pets", CreateCatBinder(gwg_api_label))
	gin_router.GET("/pets/:id", GetCatBinder(gwg_api_label))
}
func stringToInt32(s string) int32 {
	value, _ := strconv.ParseInt(s, 10, 32)
	return int32(value)
}
func stringToInt64(s string) int64 {
	value, _ := strconv.ParseInt(s, 10, 64)
	return value
}
func stringToFloat32(s string) float32 {
	value, _ := strconv.ParseFloat(s, 32)
	return float32(value)
}
func stringToFloat64(s string) float64 {
	value, _ := strconv.ParseFloat(s, 64)
	return value
}

```

You should implement PetApiInterface interface.

```go
package main

import (
	"fmt"
	"genen/api"

	"github.com/gin-gonic/gin"
)

type PetApi struct{}

func (PetApi) CreateCat(c *gin.Context, cat api.Cat) {
	fmt.Println(cat)
}

func (PetApi) GetCat(c *gin.Context, id int64) {
	fmt.Println(id)
}

func main() {
	router := gin.Default()
	api.PetApiInterfaceMounter(router, &PetApi{})
	router.Run(":8081")
}

```

## Feature

- Generate go interface
- Generate gin router
