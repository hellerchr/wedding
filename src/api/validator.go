package api

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"gopkg.in/go-playground/validator.v9"
)

type Validator struct {
	validator *validator.Validate
}

type ValidationError struct {
	Field     string `json:"field"`
	Violation string `json:"constraintViolation"`
}

func NewValidator() *Validator {
	val := validator.New()

	// return form names as field names
	val.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("form"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	if err := val.RegisterValidation("optionalof", validateOptionalOf, true); err != nil {
		panic(err)
	}

	return &Validator{val}
}

func validateOptionalOf(fl validator.FieldLevel) bool {
	field := fl.Field()
	if field.String() == "" {
		return true
	}

	vals := strings.Fields(fl.Param())

	var v string

	switch field.Kind() {
	case reflect.String:
		v = field.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v = strconv.FormatInt(field.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v = strconv.FormatUint(field.Uint(), 10)
	default:
		panic(fmt.Sprintf("Bad field type %T", field.Interface()))
	}

	for i := 0; i < len(vals); i++ {
		if vals[i] == v {
			return true
		}
	}

	return false
}

func (cv *Validator) Validate(i interface{}) []ValidationError {
	if err := cv.validator.Struct(i); err != nil {
		errors := err.(validator.ValidationErrors)

		ev := make([]ValidationError, 0, len(errors))
		for _, e := range errors {
			ev = append(ev, ValidationError{
				Field:     e.Field(),
				Violation: e.ActualTag(),
			})
		}

		return ev
	}

	return nil
}
