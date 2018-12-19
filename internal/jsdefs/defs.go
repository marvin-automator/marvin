package jsdefs

import (
	"fmt"
	"github.com/marvin-automator/marvin/actions"
	"reflect"
	"strings"
)

var marvinDefs = map[string]map[string]interface{}{
	"input": {
		"!type": "fn(name: string, description: string) -> string",
		"!doc": "Define an input variable that users of your template can set when creating a chore. In the execution phaze, this function will return the value the user set.",
	},
	"isSetupPhase": {
		"!type": "bool",
		"!doc": "Whether the template is executing in the setup phase.",
	},
}

func GetDefs() map[string]map[string]interface{} {
	defs := map[string]map[string]interface{}{}
	typeMap := make(map[string]map[string]string)

	for _, p := range actions.Registry.Providers() {
		pdefs := baseDefFromInfo(p.Info())

		for _, g := range p.Groups() {
			gdefs := baseDefFromInfo(g.Info())

			for _, a := range g.Actions() {
				var adef map[string]string
				adef, typeMap = makeActionDef(a, typeMap)
				gdefs[a.Info().Name] = adef
			}

			pdefs[g.Info().Name] = gdefs
		}

		defs[p.Info().Name] = pdefs
	}

	defs["!define"] = make(map[string]interface{}, len(typeMap))
	for name, typ := range typeMap {
		defs["!define"][name] = typ
	}

	return defs
}

func baseDefFromInfo(i actions.BaseInfo) map[string]interface{} {
	return map[string]interface{}{
		"!description": i.Description,
	}
}

func makeActionDef(a actions.Action, typeMap map[string]map[string]string) (map[string]string, map[string]map[string]string) {
	ftype, otherTypes := getFunctionType(a, typeMap)
	return map[string]string{
		"!description": a.Info().Description,
		"!type": ftype,
	}, otherTypes
}

func getFunctionType(a actions.Action, typeMap map[string]map[string]string) (string, map[string]map[string]string) {
	inTypeName, typeMap := getType(a.Info().InputType, typeMap)
	outTypeName, typeMap := getType(a.Info().OutputType, typeMap)

	var ftype string
	if a.Info().IsTrigger {
		callbackName := a.Info().Name + "Callback"
		typeMap[callbackName] = map[string]string{
			"!type": fmt.Sprintf("fn(event: %v)", outTypeName),
		}

		ftype = fmt.Sprintf("fn(input: %v, callback: %v)", inTypeName, callbackName)
	} else {
		ftype = fmt.Sprintf("fn(%v) -> %v", inTypeName, outTypeName)
	}

	return ftype, typeMap
}

func getType(t reflect.Type, typeMap map[string]map[string]string) (string, map[string]map[string]string) {
	switch t.Kind() {
	case reflect.Ptr:
		return getType(t.Elem(), typeMap)
	case reflect.String:
		return "string", typeMap
	case reflect.Bool:
		return "bool", typeMap
	case reflect.Int, reflect.Uint, reflect.Int8, reflect.Uint8, reflect.Int16, reflect.Uint16,
	reflect.Int32, reflect.Uint32, reflect.Int64, reflect.Uint64, reflect.Float64, reflect.Float32:
		return "number", typeMap
	case reflect.Slice:
		tn, typeMap := getType(t.Elem(), typeMap)
		return "[" + tn + "]", typeMap
	case reflect.Map:
		return "Object", typeMap
	case reflect.Struct:
		return getStructType(t, typeMap)
	default:
		return "string", typeMap
	}
}

func getStructType(t reflect.Type, typeMap map[string]map[string]string) (string, map[string]map[string]string) {
	if t.Name() == "Time" && t.PkgPath() == "time" {
		return "+Date", typeMap
	}

	name := typeName(t)

	if _, ok := typeMap[name]; ok {
		return name, typeMap
	}

	typeDef, typeMap := getStructTypeFields(t, typeMap)
	typeMap[name] = typeDef

	return name, typeMap
}

func getStructTypeFields(t reflect.Type, typeMap map[string]map[string]string) (map[string]string, map[string]map[string]string) {
	typeDef := make(map[string]string, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		t := f.Tag.Get("json")

		// If the field is unexported or explicitly omitted, skip it.
		if f.PkgPath != "" || t == "-" {
			continue
		}

		if t == "" {
			ftype := f.Type
			if ftype.Kind() == reflect.Ptr {
				ftype = ftype.Elem()
			}

			if ftype.Kind() != reflect.Struct {
				continue
			}


			var anonymousFields map[string]string

			ftname := typeName(ftype)
			if fdef, ok := typeMap[ftname]; ok {
				anonymousFields = fdef
			} else {
				anonymousFields, typeMap = getStructTypeFields(ftype, typeMap)
			}

			for name, value := range anonymousFields {
				typeDef[name] = value
			}
		}

		var ft string
		ft, typeMap = getType(f.Type, typeMap)
		typeDef[t] = ft
	}

	return typeDef, typeMap
}

func typeName(t reflect.Type) string {
	pp := strings.Split(t.PkgPath(), "/")
	pkg := pp[len(pp) - 1]
	return pkg + "." + t.Name()
}