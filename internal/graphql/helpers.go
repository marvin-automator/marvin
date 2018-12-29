package graphql

import (
	"fmt"
	"github.com/graphql-go/graphql"
	"reflect"
	"runtime"
	"strings"
)

var outputRegistry = map[string]graphql.Output{}

type typeTransformer struct {
	Transformer reflect.Value
	InputType   reflect.Type
	GraphType   graphql.Output
}

func (t typeTransformer) Transform(input interface{}) (interface{}, error) {
	iv := reflect.Indirect(reflect.ValueOf(input))
	if !iv.IsValid() {
		return nil, nil
	}

	if !iv.Type().AssignableTo(t.InputType) {
		f := runtime.FuncForPC(t.Transformer.Pointer())
		return nil, fmt.Errorf("transformer %v expected %v, got %v", f.Name(), t.InputType.Name(), iv.Type().Name())
	}

	result := t.Transformer.Call([]reflect.Value{iv})
	if len(result) == 2 {
		return result[0].Interface(), result[1].Interface().(error)
	}

	return result[0].Interface(), nil
}

var typeTransformers = map[string]typeTransformer{}

func RegisterTypeTransformer(transformer interface{}) {
	rt := reflect.TypeOf(transformer)
	if !(rt.Kind() == reflect.Func && rt.NumIn() == 1 && (rt.NumOut() == 1 || rt.NumOut() == 2 && rt.Out(2).Implements(reflect.TypeOf(error(nil))))) {
		panic(fmt.Errorf("RegisterTypeResolver expects a function with one input and one or two outputs. The second output, if present, must be an error."))
	}

	inType := rt.In(0)

	outType := rt.Out(0)

	tk := typeKey(inType)
	typeTransformers[tk] = typeTransformer{
		Transformer: reflect.ValueOf(transformer),
		GraphType:   getGraphType(outType),
		InputType:   inType,
	}
}

func CreateOutputTypeFromStruct(v interface{}) graphql.Output {
	rt := reflect.Indirect(reflect.ValueOf(v)).Type()
	if rt.Kind() != reflect.Struct {
		panic(fmt.Errorf("CreateOutputType expected to get struct value, but got %v", rt.Name()))
	}

	return outputTypeFromStructType(rt)
}

func outputTypeFromStructType(rt reflect.Type) graphql.Output {
	tk := typeKey(rt)
	if ot, ok := outputRegistry[tk]; ok {
		return ot
	}

	fields := make(graphql.Fields, 0)

	ot := graphql.NewObject(graphql.ObjectConfig{
		Name:   rt.Name(),
		Fields: fields,
	})

	outputRegistry[tk] = ot

	appendFields(fields, bindFields(rt))

	return ot
}

func typeKey(t reflect.Type) string {
	return t.PkgPath() + "@" + t.Name()
}

func BindFields(v interface{}) graphql.Fields {
	return bindFields(reflect.Indirect(reflect.ValueOf(v)).Type())
}

func bindFields(t reflect.Type) graphql.Fields {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	fields := make(map[string]*graphql.Field)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		tag := extractTag(field.Tag)
		if tag == "-" {
			continue
		}

		fieldType := field.Type

		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}

		var graphType graphql.Output
		var fieldTransformer typeTransformer

		tk := typeKey(fieldType)
		if transformer, ok := typeTransformers[tk]; ok {
			fieldTransformer = transformer
			graphType = fieldTransformer.GraphType
		} else if fieldType.Kind() == reflect.Struct {

			if tag == "" {
				structFields := bindFields(t.Field(i).Type)
				fields = appendFields(fields, structFields)
				continue
			} else {
				graphType = outputTypeFromStructType(fieldType)
			}
		}

		if tag == "" {
			continue
		}

		if graphType == nil {
			graphType = getGraphType(fieldType)
		}
		fields[tag] = &graphql.Field{
			Type: graphType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				val := extractValue(tag, p.Source)
				var err error
				if fieldTransformer.Transformer.IsValid() {
					val, err = fieldTransformer.Transform(val)
				}

				return val, err
			},
		}
	}
	return fields
}

func GetGraphType(v interface{}) graphql.Output {
	switch v.(type) {
	case reflect.Type:
		return getGraphType(v.(reflect.Type))
	case reflect.Value:
		return getGraphType(v.(reflect.Value).Type())
	default:
		return getGraphType(reflect.TypeOf(v))
	}
}

func getGraphType(tipe reflect.Type) graphql.Output {
	switch tipe.Kind() {
	case reflect.Ptr:
		return getGraphType(tipe.Elem())
	case reflect.String:
		return graphql.String
	case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return graphql.Int
	case reflect.Float32, reflect.Float64:
		return graphql.Float
	case reflect.Bool:
		return graphql.Boolean
	case reflect.Slice:
		return graphql.NewList(getGraphType(tipe.Elem()))
	case reflect.Struct:
		return outputTypeFromStructType(tipe)
	}

	return graphql.String
}

func appendFields(dest, origin graphql.Fields) graphql.Fields {
	for key, value := range origin {
		dest[key] = value
	}
	return dest
}

func extractValue(originTag string, obj interface{}) interface{} {
	val := reflect.Indirect(reflect.ValueOf(obj))

	for j := 0; j < val.NumField(); j++ {
		field := val.Type().Field(j)
		if field.Type.Kind() == reflect.Struct {
			res := extractValue(originTag, val.Field(j).Interface())
			if res != nil {
				return res
			}
		}

		if originTag == extractTag(field.Tag) {
			return reflect.Indirect(val.Field(j)).Interface()
		}
	}
	return nil
}

func extractTag(tag reflect.StructTag) string {
	t := tag.Get("json")
	if t != "" {
		t = strings.Split(t, ",")[0]
	}
	return t
}
