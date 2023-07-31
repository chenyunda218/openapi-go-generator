package openapi

import (
	"fmt"
	"strings"

	"github.com/chenyunda218/gwg"
)

const BODY_NAME = "gin_body"
const GIN_CONTEXT_LABEL = "gin_context"
const GIN_ROUTER_LABEL = "gin_router"
const GWG_API_LABEL = "gwg_api_label"
const BUILDER_SUFFIX = "Builder"

func (o Openapi) Generate(root string, packageName string) {
	writer := gwg.Package{
		Name: packageName,
	}
	writer.AddImport(gwg.Import{Packages: []string{"github.com/gin-gonic/gin", "strconv"}})
	for _, m := range o.FindModel() {
		writer.AddCode(m)
	}
	for _, e := range o.FindEnums() {
		writer.AddCode(e)
	}
	for _, i := range o.CreateInterfaces() {
		writer.AddCode(i)
		for _, binder := range o.CreateInterfaceBinders(i) {
			writer.AddCode(binder)
		}
		writer.AddCode(o.CreateApiMounter(i))
	}
	// Type converter
	writer.AddCode(stringToInt32)
	writer.AddCode(stringToInt64)
	writer.AddCode(stringToFloat32)
	writer.AddCode(stringToFloat64)
	writer.Wirte(root)
}

func (o Openapi) CreateApiMounter(i gwg.Interface) (fs gwg.Func) {
	fs.Name = i.Name + "Mounter"
	fs.Parameters.Add(gwg.Pair{
		Left:  GIN_ROUTER_LABEL,
		Right: "*gin.Engine",
	}, gwg.Pair{Left: GWG_API_LABEL, Right: i.Name})
	for _, m := range i.Methods {
		path, method := o.GetMethodAndPathByOperationId(m.Name)
		fs.AddLine(gwg.Line{
			Content: fmt.Sprintf("%s.%s(\"%s\", %s(%s))",
				GIN_ROUTER_LABEL, strings.ToUpper(method), PathConverter(path), builderName(m.Name), GWG_API_LABEL),
		})
	}
	return fs
}

func (o Openapi) CreateInterfaces() (interfaces []gwg.Interface) {
	tags := o.Tags
	for _, tag := range tags {
		interfaces = append(interfaces, o.CreateInterface(tag.Name))
	}
	return interfaces
}

func (o Openapi) CreateInterfaceBinders(i gwg.Interface) (binders []gwg.Func) {
	for _, m := range i.Methods {
		binders = append(binders, o.CreateBinder(m, i.Name))
	}
	return binders
}

func builderName(name string) string {
	return name + BUILDER_SUFFIX
}

func (o Openapi) GetMethodAndPathByOperationId(operationId string) (string, string) {
	for path, methods := range o.Paths {
		for method, api := range methods {
			if api.OperationId == operationId {
				return path, method
			}
		}
	}
	return "", ""
}

func (o Openapi) CreateBinder(i gwg.Method, apiName string) (binder gwg.Func) {
	binder.Name = builderName(i.Name)
	binder.Parameters.Add(
		gwg.Pair{
			Left:  "api",
			Right: apiName,
		},
	)
	binder.Outputs.Pairs = append(binder.Outputs.Pairs, gwg.Pair{
		Right: "func(c *gin.Context)",
	})
	binder.AddLine(gwg.Line{Content: fmt.Sprintf("return func(%s *gin.Context) {", GIN_CONTEXT_LABEL)})
	var ps []string = []string{GIN_CONTEXT_LABEL}
	api := o.GetApiByOperationId(i.Name)
	var parameters []Parameter
	for _, p := range api.Parameters {
		if p.Ref != nil {
			parameters = append(parameters, o.GetParameter(RefObject(*p.Ref)))
		} else {
			parameters = append(parameters, p)
		}
	}
	for _, p := range parameters {
		if p.In == "path" {
			binder.AddLine(
				gwg.Line{Content: fmt.Sprintf("%s := %s.Param(\"%s\")", p.Name, GIN_CONTEXT_LABEL, p.Name)},
			)
		} else {
			binder.AddLine(
				gwg.Line{Content: fmt.Sprintf("%s := %s.Query(\"%s\")", p.Name, GIN_CONTEXT_LABEL, p.Name)},
			)
		}
		if p.Ref != nil {
			finded := o.GetParameter(RefObject(*p.Ref))
			ps = append(ps, BinderSchema(finded.Schema, finded.Name))
		} else {
			ps = append(ps, BinderSchema(p.Schema, p.Name))
		}
	}
	// Body binding
	if api.RequestBody != nil && api.RequestBody.Content.Json != nil {
		if api.RequestBody.Content.Json.Schema.Ref != nil {
			valueName := FirstToLower(RefObject(*api.RequestBody.Content.Json.Schema.Ref))
			binder.AddLine(
				gwg.Line{Content: fmt.Sprintf("var %s %s",
					valueName,
					FirstToUpper(valueName),
				)},
				gwg.Line{Content: fmt.Sprintf("%s.ShouldBindJSON(&%s)",
					GIN_CONTEXT_LABEL,
					valueName,
				)},
			)
			ps = append(ps, valueName)
		}
	}
	binder.AddLine(gwg.Line{
		Content: fmt.Sprintf("api.%s(%s)", i.Name, strings.Join(ps, ", ")),
	})

	binder.AddLine(gwg.Line{Content: "}"})
	return binder
}

func BinderSchema(s Schema, name string) string {
	if s.Ref != nil {
		return fmt.Sprintf("%s(%s)", RefObject(*s.Ref), name)
	}
	var convert string
	switch ConvertType(s) {
	case "string":
		return name
	case "int64":
		convert = stringToInt64.Name
	case "int":
		convert = stringToInt64.Name
	case "int32":
		convert = stringToInt32.Name
	case "float32":
		convert = stringToFloat32.Name
	case "float64":
		convert = stringToFloat64.Name
	}
	return fmt.Sprintf("%s(%s)", convert, name)
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
					Left:  GIN_CONTEXT_LABEL,
					Right: "*gin.Context",
				})
				for _, parameter := range api.Parameters {
					p := parameter
					if p.Ref != nil {
						p = o.GetParameter(RefObject(*parameter.Ref))
					}
					method.Parameters.Add(gwg.Pair{
						Left:  p.Name,
						Right: ConvertType(p.Schema),
					})
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

func (o Openapi) GetApiByOperationId(id string) *Api {
	for _, path := range o.Paths {
		for _, api := range path {
			if api.OperationId == id {
				return &api
			}
		}
	}
	return nil
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
		} else if len(schema.OneOf) != 0 && schema.Discriminator.PropertyName != "" {
			// TODO: oneOf
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
		tags = []gwg.Tag{
			{Label: "json", Content: FirstToLower(label) + ",omitempty"},
		}
	}
	if required {
		tags = append(tags, gwg.Tag{
			Label:   "binding",
			Content: "required",
		})
	}
	var t string
	t = ConvertType(s)
	if !required {
		t = "*" + t
	}
	return gwg.Property{
		Label: FirstToUpper(label),
		Type:  t,
		Tags:  tags,
	}
}

func ConvertType(s Schema) string {
	switch s.Type {
	case "string":
		return "string"
	case "integer":
		return ConvertInteger(s.Format)
	case "number":
		return ConvertNumber(s.Format)
	case "array":
		if s.Items == nil {
			panic("Array without times")
		}
		return fmt.Sprintf("[]%s", ConvertType(*s.Items))
	}

	if s.Ref != nil {
		return RefObject(*s.Ref)
	}
	return "string"
}

func ConvertInteger(format string) string {
	switch format {
	case "int64":
		return "int64"
	case "int32":
		return "int32"
	default:
		return "int64"
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
