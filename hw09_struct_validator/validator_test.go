package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code  int    `validate:"in:200,404,500"`
		Body  string `json:"omitempty"`
		Email string `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
	}

	MyStruct struct {
		Names []string `validate:"in:Sam,Bill,Jack"`
		Ages  []int    `validate:"min:18|max:50"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in:          MyStruct{Names: []string{"Sam", "Bill", "Jack"}, Ages: []int{18, 21, 22}},
			expectedErr: ValidationErrors{},
		},
		{
			in:          Response{Code: 200, Body: "Hi", Email: "sdfs@mail.ru"},
			expectedErr: ValidationErrors{},
		},
		{
			in:          Response{Code: 500, Body: "Hi", Email: "mail@m.ru"},
			expectedErr: ValidationErrors{},
		},
		{
			in:          App{Version: "12345"},
			expectedErr: ValidationErrors{},
		},
		{
			in: User{
				ID: "71b5263d-37e2-11ed-901f-00155d8ed20b", Age: 19, Email: "mail@mail.ru",
				Role: "stuff", Phones: []string{"79998361241", "79268361241"},
			},
			expectedErr: ValidationErrors{},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			errs, err := Validate(tt.in)
			require.NoError(t, err)
			require.Equal(t, tt.expectedErr, errs)
			_ = tt
		})
	}
}

func TestValidateErrors(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: Response{Code: 1, Body: "Hi", Email: "mailm.ru"},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Code", Err: ErrIntInValidation},
				ValidationError{Field: "Email", Err: ErrRegexpValidation},
			},
		},
		{
			in:          Response{Code: 1, Body: "Hi", Email: "mail@m.ru"},
			expectedErr: ValidationErrors{ValidationError{Field: "Code", Err: ErrIntInValidation}},
		},
		{
			in:          Response{Code: 200, Body: "Hi", Email: "mailm.ru"},
			expectedErr: ValidationErrors{ValidationError{Field: "Email", Err: ErrRegexpValidation}},
		},
		{
			in:          App{Version: "1234"},
			expectedErr: ValidationErrors{ValidationError{Field: "Version", Err: ErrLengthValidation}},
		},
		{
			in: User{
				ID: "71b5263d", Age: 17, Email: "mailmail.ru",
				Role: "stuff", Phones: []string{"7", "79268361241"},
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "ID", Err: ErrLengthValidation},
				ValidationError{Field: "Age", Err: ErrIntMinValidation},
				ValidationError{Field: "Email", Err: ErrRegexpValidation},
				ValidationError{Field: "Phones", Err: ErrLengthValidation},
			},
		},
		{
			in: MyStruct{Names: []string{}, Ages: []int{}},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Names", Err: ErrInValidation},
				ValidationError{Field: "Ages", Err: ErrIntMinValidation},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			errs, err := Validate(tt.in)
			require.NoError(t, err)
			require.Equal(t, tt.expectedErr, errs)
			_ = tt
		})
	}
}

func TestErrors(t *testing.T) {
	t.Run("wrong type", func(t *testing.T) {
		errs, err := Validate("string")

		require.Equal(t, ErrNotStruct, err)
		require.Equal(t, ValidationErrors{}, errs)
	})
}
