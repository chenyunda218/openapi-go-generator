package openapi

import (
	"fmt"
	"os"

	"github.com/chenyunda218/gwg"
)

const BODY_NAME = "openapi_go_gingenertor_body"

func (o Openapi) Generate(root string, packageName string) {
	os.MkdirAll(root, os.ModePerm)
	writer := gwg.Package{
		Name: packageName,
	}
	writer.AddImport(gwg.Import{Packages: []string{"github.com/gin-gonic/gin"}})
	for _, m := range o.FindModel() {
		writer.AddCode(m)
	}
	for _, e := range o.FindEnums() {
		writer.AddCode(e)
	}
	for _, i := range o.CreateInterfaces() {
		writer.AddCode(i)
		for _, binder := range CreateInterfaceBinders(i) {
			writer.AddCode(binder)
		}
	}
	writer.Wirte(packageName)
}

func (o Openapi) CreateInterfaces() (interfaces []gwg.Interface) {
	tags := o.Tags
	for _, tag := range tags {
		interfaces = append(interfaces, o.CreateInterface(tag.Name))
	}
	return interfaces
}

func CreateInterfaceBinders(i gwg.Interface) (binders []gwg.Func) {
	for _, m := range i.Methods {
		binders = append(binders, CreateBinder(m))
	}
	return binders
}

func CreateBinder(i gwg.Method) (binder gwg.Func) {
	binder.Name = i.Name + "Binder"
	binder.Parameters.Add(gwg.Pair{
		Left:  "c",
		Right: "*gin.Context",
	})
	for _, p := range i.Parameters.Pairs {
		if p.Left == BODY_NAME {
			fmt.Println(BODY_NAME)
			binder.AddLine(
				gwg.Line{Content: fmt.Sprintf("var %s %s", BODY_NAME, p.Right)},
				gwg.Line{Content: fmt.Sprintf("c.ShouldBindJSON(&%s)", BODY_NAME)},
			)
		}
	}
	fmt.Println(len(binder.Lines))
	return binder
}

func (o Openapi) CreateInterface(group string) (i gwg.Interface) {
	i.Name = group + "ApiInterface"
	for _, path := range o.Paths {
		for _, api := range path {
			if len(api.Tags) > 0 && api.Tags[0] == group {
				method := gwg.Method{
					Name: FirstToUpper(api.OperationId),
				}
				method.Parameters.Add(gwg.Pair{
					Left:  "c",
					Right: "*gin.Context",
				})
				for _, parameter := range api.Parameters {
					if parameter.Ref == nil {
						method.Parameters.Add(gwg.Pair{
							Left:  parameter.Name,
							Right: "string",
						})
					} else {
						p := o.GetParameter(RefObject(*parameter.Ref))
						method.Parameters.Add(gwg.Pair{
							Left:  p.Name,
							Right: "string",
						})
					}

				}
				if api.RequestBody != nil {
					if api.RequestBody.Content.Json != nil {
						method.Parameters.Add(gwg.Pair{
							Left:  BODY_NAME,
							Right: RefObject(*api.RequestBody.Content.Json.Schema.Ref),
						})
					}
				}
				i.AddMethod(method)
			}
		}
	}
	return i
}

func (o Openapi) GetParameter(name string) Parameter {
	for k, p := range o.Components.Parameters {
		if k == name {
			return p
		}
	}
	return Parameter{}
}

func (o Openapi) FindEnums() []gwg.Enums {
	var enums []gwg.Enums
	for name, schema := range o.Components.Schemas {
		if schema.Type == "string" && len(schema.Enum) > 0 {
			enums = append(enums, ConvertEnum(name, schema))
		}
	}
	return enums
}

func (o Openapi) FindModel() []gwg.Struct {
	if o.Components == nil {
		return []gwg.Struct{}
	}
	var models []gwg.Struct
	for name, schema := range o.Components.Schemas {
		if schema.AllOf != nil || schema.Type == "object" {
			models = append(models, ConvertSchema(name, schema))
		}
	}
	return models
}

func ConvertEnum(title string, s Schema) gwg.Enums {
	return gwg.Enums{
		Title:  title,
		Values: s.Enum,
	}
}

func ConvertSchema(name string, s Schema) gwg.Struct {
	if len(s.AllOf) != 0 {
		o := gwg.Struct{
			Name: name,
		}
		for _, c := range s.AllOf {
			if c.Ref != nil {
				o.AddCombination(RefObject(*c.Ref))
			}
		}
		for _, c := range s.AllOf {
			if c.Ref == nil {
				o.Properties = ConvertProperties(c.Properties, c.Required)
				return o
			}
		}
	}
	return gwg.Struct{
		Name:       name,
		Properties: ConvertProperties(s.Properties, s.Required),
	}
}

func ConvertProperties(p map[string]Schema, requiredList []string) []gwg.Property {
	var properties []gwg.Property
	for n, s := range p {
		required := false
		for _, k := range requiredList {
			if k == n {
				required = true
			}
		}
		properties = append(properties, ConvertProperty(n, s, required))
	}
	return properties
}

func ConvertProperty(label string, s Schema, required bool) gwg.Property {
	tags := []gwg.Tag{
		{Label: "json", Content: FirstToLower(label)},
	}
	if required {
		tags = append(tags, gwg.Tag{
			Label:   "binding",
			Content: "required",
		})
	}
	var t string = "string"

	switch s.Type {
	case "string":
		t = "string"
	case "integer":
		t = ConvertInteger(s.Format)
	case "number":
		t = ConvertNumber(s.Format)
	}
	if s.Ref != nil {
		t = RefObject(*s.Ref)
	}

	return gwg.Property{
		Label: FirstToUpper(label),
		Type:  t,
		Tags:  tags,
	}
}

func ConvertInteger(format string) string {
	switch format {
	case "int64":
		return "int64"
	case "int32":
		return "int32"
	default:
		return "int"
	}
}

func ConvertNumber(format string) string {
	switch format {
	case "double":
		return "float64"
	case "float":
		return "float32"
	default:
		return "float64"
	}
}
