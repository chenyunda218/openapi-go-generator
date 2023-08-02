package openapi

import "github.com/chenyunda218/gwg"

var stringToInt64 gwg.Func = gwg.Func{
	Name: "stringToInt64",
	Parameters: gwg.Parameters{
		Pairs: []gwg.Pair{{Left: "s", Right: "string"}},
	},
	Outputs: gwg.Outputs{Pairs: []gwg.Pair{{Right: "int64"}}},
	Lines: []gwg.Line{
		{Content: "if value, err := strconv.ParseInt(s, 10, 64); err != nil {"},
		{Content: "return value"},
		{Content: "}"},
		{Content: "return 0"},
	},
}

var stringToInt32 gwg.Func = gwg.Func{
	Name: "stringToInt32",
	Parameters: gwg.Parameters{
		Pairs: []gwg.Pair{{Left: "s", Right: "string"}},
	},
	Outputs: gwg.Outputs{Pairs: []gwg.Pair{{Right: "int32"}}},
	Lines: []gwg.Line{
		{Content: "if value, err := strconv.ParseInt(s, 10, 32); err != nil {"},
		{Content: "return int32(value)"},
		{Content: "}"},
		{Content: "return 0"},
	},
}

var stringToFloat32 gwg.Func = gwg.Func{
	Name: "stringToFloat32",
	Parameters: gwg.Parameters{
		Pairs: []gwg.Pair{{Left: "s", Right: "string"}},
	},
	Outputs: gwg.Outputs{Pairs: []gwg.Pair{{Right: "float32"}}},
	Lines: []gwg.Line{
		{Content: "if value, err := strconv.ParseFloat(s, 32); err != nil {"},
		{Content: "return float32(value)"},
		{Content: "}"},
		{Content: "return 0"},
	},
}

var stringToFloat64 gwg.Func = gwg.Func{
	Name: "stringToFloat64",
	Parameters: gwg.Parameters{
		Pairs: []gwg.Pair{{Left: "s", Right: "string"}},
	},
	Outputs: gwg.Outputs{Pairs: []gwg.Pair{{Right: "float64"}}},
	Lines: []gwg.Line{
		{Content: "if value, err := strconv.ParseFloat(s, 64); err != nil {"},
		{Content: "return value"},
		{Content: "}"},
		{Content: "return 0"},
	},
}
