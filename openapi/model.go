package openapi

type Components struct {
	Schemas    map[string]Schema    `yaml:"schemas"`
	Parameters map[string]Parameter `yaml:"parameters"`
}

type Schema struct {
	Description string            `yaml:"description"`
	Type        string            `yaml:"type"`
	Format      string            `yaml:"format"`
	Required    []string          `yaml:"required"`
	Properties  map[string]Schema `yaml:"properties"`
	Enum        []string          `yaml:"enum"`
	Ref         *string           `yaml:"$ref"`
	Items       *Schema           `yaml:"items"`
	AllOf       []Schema          `yaml:"allOf"`
	Default     string            `yaml:"default"`
	Maximum     string            `yaml:"maximum"`
	Minimum     string            `yaml:"minimum"`
}

type Openapi struct {
	Openapi    string                    `yaml:"openapi"`
	Info       Info                      `yaml:"info"`
	Tags       []Tag                     `yaml:"tags"`
	Paths      map[string]map[string]Api `yaml:"paths"`
	Components *Components               `yaml:"components"`
}

type Info struct {
	Title       string  `yaml:"title"`
	Description string  `yaml:"description"`
	Version     string  `yaml:"version"`
	Contact     Contact `yaml:"contact"`
}

type Contact struct {
	Email string `yaml:"email"`
}

type Tag struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

type Api struct {
	Description string              `yaml:"description"`
	Tags        []string            `yaml:"tags"`
	OperationId string              `yaml:"operationId"`
	Responses   map[string]Response `yaml:"responses"`
	RequestBody *Response           `yaml:"requestBody"`
	Parameters  []Parameter         `yaml:"parameters"`
}

type Response struct {
	Description string  `yaml:"description"`
	Content     Content `yaml:"content"`
}

type Content struct {
	Json *Json `yaml:"application/json"`
}

type Json struct {
	Schema Schema `yaml:"schema"`
}

type Parameter struct {
	Ref         *string `yaml:"$ref"`
	Name        string  `yaml:"name"`
	Description string  `yaml:"description"`
	In          string  `yaml:"in"`
	Schema      Schema  `yaml:"schema"`
	Required    bool    `yaml:"required"`
}

func (p Parameter) Args() (c string) {
	return c
}
