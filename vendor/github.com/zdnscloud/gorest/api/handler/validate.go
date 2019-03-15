package handler

import (
	"reflect"
	"strings"

	"github.com/zdnscloud/cement/reflector"
	"github.com/zdnscloud/gorest/types"
)

type Operation string

var (
	Create = Operation("create")
	Update = Operation("update")
)

func CheckObjectFields(ctx *types.APIContext, obj types.Object) *types.APIError {
	structVal, ok := reflector.GetStructFromPointer(obj)
	if ok == false {
		return types.NewAPIError(types.ServerError, "get object structure but return "+structVal.Kind().String())
	}

	_, err := getStructValue(ctx, ctx.Schema, structVal)
	return err
}

func getStructValue(ctx *types.APIContext, schema *types.Schema, structVal reflect.Value) (map[string]interface{}, *types.APIError) {
	fieldValues := map[string]interface{}{}
	structTyp := structVal.Type()
	if schema == nil {
		schema = ctx.Schemas.Schema(ctx.Version, strings.ToLower(structTyp.Name()))
		if schema == nil {
			return nil, types.NewAPIError(types.NotFound, "no found schema "+strings.ToLower(structTyp.Name()))
		}
	}

	for i := 0; i < structVal.NumField(); i++ {
		field := structTyp.Field(i)
		if field.PkgPath != "" {
			continue
		}

		jsonName := types.GetJsonName(field)
		if jsonName == "-" {
			continue
		}

		if field.Anonymous && jsonName == "" {
			t := field.Type
			if t.Kind() == reflect.Ptr {
				t = t.Elem()
			}
			if t.Kind() == reflect.Struct {
				fieldVal := structVal.FieldByName(field.Name)
				if _, err := getStructValue(ctx, ctx.Schema, fieldVal); err != nil {
					return nil, err
				}
			}
			continue
		}

		fieldName := jsonName
		if fieldName == "" {
			fieldName = strings.ToLower(field.Name)
			if strings.HasSuffix(fieldName, "ID") {
				fieldName = strings.TrimSuffix(fieldName, "ID") + "Id"
			}
		}

		if types.BlacklistNames[fieldName] {
			continue
		}

		value, err := getFieldValue(ctx, field.Type, structVal.FieldByName(field.Name))
		if err != nil {
			return nil, err
		}

		schemaField := schema.ResourceFields[fieldName]
		if valueIsNil(value) && schemaField.Required {
			return nil, types.NewAPIError(types.MissingRequired, "field "+fieldName+" must be set when create")
		}

		fieldValues[fieldName] = value
	}

	return fieldValues, nil
}

func valueIsNil(value interface{}) bool {
	if value == nil || value == "" || value == 0 {
		return true
	}

	val := reflect.ValueOf(value)
	typ := val.Type()
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	switch typ.Kind() {
	case reflect.Map, reflect.Slice:
		return val.IsNil()
	default:
		return false
	}
}

func getFieldValue(ctx *types.APIContext, fieldType reflect.Type, fieldVal reflect.Value) (interface{}, *types.APIError) {
	if fieldType.Kind() == reflect.Ptr {
		fieldType = fieldType.Elem()
	}

	switch fieldType.Kind() {
	case reflect.Struct:
		return getStructValue(ctx, nil, fieldVal)
	case reflect.Slice:
		return getSliceValue(ctx, fieldVal)
	case reflect.Map:
		return getMapValue(ctx, fieldVal)
	default:
		return fieldVal.Interface(), nil
	}
}

func getSliceValue(ctx *types.APIContext, fieldValSlice reflect.Value) (interface{}, *types.APIError) {
	var values []interface{}
	for i := 0; i < fieldValSlice.Len(); i++ {
		fieldVal := fieldValSlice.Index(i)
		if val, err := getFieldValue(ctx, fieldVal.Type(), fieldVal); err != nil {
			return nil, err
		} else {
			values = append(values, val)
		}
	}

	return values, nil
}

func getMapValue(ctx *types.APIContext, fieldValMap reflect.Value) (interface{}, *types.APIError) {
	values := map[string]interface{}{}
	for _, key := range fieldValMap.MapKeys() {
		val := fieldValMap.MapIndex(key)
		val = reflect.ValueOf(val.Interface())
		if fieldVal, err := getFieldValue(ctx, val.Type(), val); err != nil {
			return nil, err
		} else {
			values[key.String()] = fieldVal
		}
	}

	return values, nil
}
