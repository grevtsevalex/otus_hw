package hw09structvalidator

import (
	"regexp"
	"strconv"
	"strings"
)

func StringValidator(field string, value string, tag string) (ValidationErrors, error) {
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
			err := InValidator(value, args)
			if err != nil {
				res = append(res, ValidationError{Field: field, Err: err})
			}
		}

		if validatorName == "len" {
			length, _ := strconv.Atoi(args[0])
			err := LengthValidator(value, length)
			if err != nil {
				res = append(res, ValidationError{Field: field, Err: err})
			}
		}

		if validatorName == "regexp" {
			err := RegExpValidator(value, args[0])
			if err != nil {
				res = append(res, ValidationError{Field: field, Err: err})
			}
		}
	}

	return res, nil
}

func InValidator(value string, allowed []string) error {
	for _, field := range allowed {
		if field == value {
			return nil
		}
	}
	return ErrInValidation
}

func RegExpValidator(value string, pattern string) error {
	rgx, err := regexp.Compile(pattern)
	res := rgx.MatchString(value)
	if err != nil {
		return ErrRegexpPattern
	}

	if res {
		return nil
	}

	return ErrRegexpValidation
}

func LengthValidator(value string, length int) error {
	if len(value) == length {
		return nil
	}

	return ErrLengthValidation
}
