package gmob

import (
	"errors"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
)

// Build Mongo Driver BSON document (M) from Map and Struct.
// For Struct input, all fields must have "bson" tag, and returned result will not include fields with zero values like false, "", and 0,
// if you want keep zero values use Map instead.
func Build(in interface{}) (result bson.M, err error) {
	val := getReflectValue(in)

	switch val.Kind() {
	case reflect.Struct, reflect.Map:
		result = bson.M{}
		buildResult(0, result, "", in)
	default:
		err = errors.New("Invalid input type")
	}

	return
}
