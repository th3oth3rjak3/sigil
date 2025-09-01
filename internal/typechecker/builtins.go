package typechecker

var builtinTypes = map[string]*BuiltinTypeInfo{
	"len": {
		Arity:      1,
		ParamTypes: []Type{&StringType{}},
		ReturnType: &NumberType{},
	},
	"print": {
		Arity:      -1,
		ParamTypes: []Type{&StringType{}},
		ReturnType: &VoidType{},
	},
	"println": {
		Arity:      -1,
		ParamTypes: []Type{&StringType{}},
		ReturnType: &VoidType{},
	},
	"string": {
		Arity:      1,
		ParamTypes: []Type{&UnknownType{}}, // Unknown means we allow any
		ReturnType: &StringType{},
	},
}
