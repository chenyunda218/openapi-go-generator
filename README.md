# openapi-go-gin-genertor

openapi gin api genertor.

## CONFIG

| Flat | Description     | Default              |
| ---- | --------------- | -------------------- |
| -i   | Input yaml file | ./openapi.yaml       |
| -o   | Output path     | ./openapigingenertor |
| -p   | Package name    | openapigingenertor   |

## Sample

```bash
git clone https://github.com/chenyunda218/openapi-go-gin-genertor
cd openapi-go-gin-genertor
go run main.go -o api -p api -i ./openapi.ayml
```

## Feature

- Generate go interface
- Generate gin router

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
| oneOf    | ✓         |
| required | ✓         |
| anyOf    | ✗         |
