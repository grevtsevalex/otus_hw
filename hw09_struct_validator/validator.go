package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

var (
	ErrIntMinValidation = errors.New("число меньше минимального")
	ErrIntMaxValidation = errors.New("число больше максимального")
	ErrIntInValidation  = errors.New("число не входит в множество допустимых")
	ErrRegexpValidation = errors.New("значение поля не соответствует шаблону")
	ErrRegexpPattern    = errors.New("ошибка при компиляции регулярного выражения")
	ErrInValidation     = errors.New("значение поля не входит в перечень допустимых")
	ErrLengthValidation = errors.New("длина поля не соответствует разрешенной")
	ErrNotStruct        = errors.New("валидируемое значение должно быть структурой")
)

func (v ValidationErrors) Error() string {
	res := strings.Builder{}

	for _, validationError := range v {
		res.WriteString(fmt.Sprintf("Ошибка в поле %s: ", validationError.Field))
		res.WriteString(validationError.Err.Error())
		res.WriteString("\n")
	}

	return res.String()
}

const (
	vTagSep            = ":"
	vTagValueSep       = ","
	vTagjValidatorsSep = "|"
	vTagName           = "validate"
)

func Validate(v interface{}) (ValidationErrors, error) {
	res := make(ValidationErrors, 0)
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Struct {
		return ValidationErrors{}, ErrNotStruct
	}

	for i := 0; i < rv.NumField(); i++ {
		field := rv.Field(i)
		typeField := rv.Type().Field(i)
		tagValue := typeField.Tag.Get(vTagName)
		fType := typeField.Type.String()
		fName := typeField.Name

		if !strings.Contains(tagValue, vTagSep) {
			continue
		}

		if field.Kind() == reflect.Slice {
			validationErrs, err := sliceValidator(fName, fType, field, tagValue)
			if err != nil {
				return ValidationErrors{}, err
			}
			res = append(res, validationErrs...)
			continue
		}

		switch fType {
		case "string":
			validationErrs, err := StringValidator(fName, field.String(), tagValue)
			if err != nil {
				return ValidationErrors{}, err
			}
			res = append(res, validationErrs...)
		case "int":
			validationErrs, err := IntValidator(fName, int(field.Int()), tagValue)
			if err != nil {
				return ValidationErrors{}, err
			}
			res = append(res, validationErrs...)
		default:
		}
	}

	return res, nil
}

func sliceValidator(fName string, fType string, field reflect.Value, tagV string) (ValidationErrors, error) {
	res := make(ValidationErrors, 0)
	switch fType {
	case "[]string":
		if field.Len() == 0 {
			defaultValue := ""
			validationErrs, err := StringValidator(fName, defaultValue, tagV)
			if err != nil {
				return ValidationErrors{}, err
			}
			res = append(res, validationErrs...)
			break
		}

		for i := 0; i < field.Len(); i++ {
			validationErrs, err := StringValidator(fName, field.Index(i).String(), tagV)
			if err != nil {
				return ValidationErrors{}, err
			}
			res = append(res, validationErrs...)
		}
	case "[]int":
		if field.Len() == 0 {
			defaultValue := 0
			validationErrs, err := IntValidator(fName, defaultValue, tagV)
			if err != nil {
				return ValidationErrors{}, err
			}
			res = append(res, validationErrs...)
			break
		}

		for i := 0; i < field.Len(); i++ {
			validationErrs, err := IntValidator(fName, int(field.Index(i).Int()), tagV)
			if err != nil {
				return ValidationErrors{}, err
			}
			res = append(res, validationErrs...)
		}
	default:
	}

	return res, nil
}
