package gmob

import (
	"reflect"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// struct tag
// to build document result key
// Ex:
// type Entity struct {
//     Name `bson:"name"`
// }
const keyTag = "bson"

// get reflect value
// with pointer check
func getReflectValue(in interface{}) reflect.Value {
	val := reflect.ValueOf(in)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	return val
}

// build BSON result recursively
func buildResult(depth int, result bson.M, key string, in interface{}) {
	v := getReflectValue(in)
	if !v.IsValid() {
		return
	}

	newDepth := depth + 1
	t := reflect.TypeOf(v.Interface())

	switch t.Kind() {
	case reflect.String, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		result[key] = v.Interface()
	case reflect.Slice, reflect.Array:
		// keep original primitive object id type
		if t, ok := v.Interface().(primitive.ObjectID); ok {
			result[key] = t
			return
		}

		vResult := bson.A{}
		for i := 0; i < v.Len(); i++ {
			tempResult := bson.M{}
			tempKey := strconv.Itoa(i)
			buildResult(newDepth, tempResult, tempKey, v.Index(i).Interface())
			vResult = append(vResult, tempResult[tempKey])
		}
		result[key] = vResult
	case reflect.Map:
		vResult := bson.M{}
		if depth == 0 {
			vResult = result
		}

		for _, k := range v.MapKeys() {
			buildResult(newDepth, vResult, k.String(), v.MapIndex(k).Interface())
		}
		if depth > 0 {
			result[key] = vResult
		}
	case reflect.Struct:
		// keep original time type
		if t, ok := v.Interface().(time.Time); ok {
			result[key] = t
			return
		}

		vResult := bson.M{}
		if depth == 0 {
			vResult = result
		}

		for i := 0; i < v.NumField(); i++ {
			typeField := v.Type().Field(i)
			bsonValue := typeField.Tag.Get(keyTag)
			bsonValueSplit := strings.SplitN(bsonValue, ",", 2)
			if bsonValue != "" {
				if v.Field(i).IsZero() {
					continue
				}

				typeFieldKey := bsonValueSplit[0]
				buildResult(newDepth, vResult, typeFieldKey, v.Field(i).Interface())
			}
		}

		if depth > 0 {
			result[key] = vResult
		}
	}
}
