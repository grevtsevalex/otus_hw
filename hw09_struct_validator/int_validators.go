package hw09structvalidator

import (
	"strconv"
	"strings"
)

func IntValidator(field string, value int, tag string) (ValidationErrors, error) {
	res := make(ValidationErrors, 0)
	validatorTags := strings.Split(tag, vTagjValidatorsSep)
	for _, vt := range validatorTags {
		validateTagData := strings.Split(vt, vTagSep)
		validatorName := validateTagData[0]
		args := make([]string, 0)

		if len(validateTagData) > 1 {
			args = strings.Split(validateTagData[1], vTagValueSep)
		}

		if validatorName == "in" {
			allowed := make([]int, 0, len(args))
			for _, arg := range args {
				v, _ := strconv.Atoi(arg)
				allowed = append(allowed, v)
			}
			err := IntInValidator(value, allowed)
			if err != nil {
				res = append(res, ValidationError{Field: field, Err: err})
			}
		}

		if validatorName == "max" {
			max, _ := strconv.Atoi(args[0])
			err := MaxValidator(value, max)
			if err != nil {
				res = append(res, ValidationError{Field: field, Err: err})
			}
		}

		if validatorName == "min" {
			min, _ := strconv.Atoi(args[0])
			err := MinValidator(value, min)
			if err != nil {
				res = append(res, ValidationError{Field: field, Err: err})
			}
		}
	}

	return res, nil
}

func MinValidator(value int, min int) error {
	if value < min {
		return ErrIntMinValidation
	}
	return nil
}

func MaxValidator(value int, max int) error {
	if value > max {
		return ErrIntMaxValidation
	}
	return nil
}

func IntInValidator(value int, allowed []int) error {
	for _, field := range allowed {
		if field == value {
			return nil
		}
	}

	return ErrIntInValidation
}
