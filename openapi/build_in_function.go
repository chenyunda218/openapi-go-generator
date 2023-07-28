package openapi

import "github.com/chenyunda218/gwg"

var stringToInt64 gwg.Func = gwg.Func{
	Name: "stringToInt64",
	Parameters: gwg.Parameters{
		Pairs: []gwg.Pair{{Left: "s", Right: "string"}},
	},
	Outputs: gwg.Outputs{Pairs: []gwg.Pair{{Right: "int64"}}},
	Lines: []gwg.Line{
		{Content: "value, _ := strconv.ParseInt(s, 10, 64)"},
		{Content: "return value"},
	},
}

var stringToInt32 gwg.Func = gwg.Func{
	Name: "stringToInt32",
	Parameters: gwg.Parameters{
		Pairs: []gwg.Pair{{Left: "s", Right: "string"}},
	},
	Outputs: gwg.Outputs{Pairs: []gwg.Pair{{Right: "int32"}}},
	Lines: []gwg.Line{
		{Content: "value, _ := strconv.ParseInt(s, 10, 32)"},
		{Content: "return int32(value)"},
	},
}

var stringToFloat32 gwg.Func = gwg.Func{
	Name: "stringToFloat32",
	Parameters: gwg.Parameters{
		Pairs: []gwg.Pair{{Left: "s", Right: "string"}},
	},
	Outputs: gwg.Outputs{Pairs: []gwg.Pair{{Right: "float32"}}},
	Lines: []gwg.Line{
		{Content: "value, _ := strconv.ParseFloat(s, 32)"},
		{Content: "return float32(value)"},
	},
}

var stringToFloat64 gwg.Func = gwg.Func{
	Name: "stringToFloat64",
	Parameters: gwg.Parameters{
		Pairs: []gwg.Pair{{Left: "s", Right: "string"}},
	},
	Outputs: gwg.Outputs{Pairs: []gwg.Pair{{Right: "float64"}}},
	Lines: []gwg.Line{
		{Content: "value, _ := strconv.ParseFloat(s, 64)"},
		{Content: "return value"},
	},
}
