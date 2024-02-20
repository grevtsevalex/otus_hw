package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	res := strings.Builder{}

	if s == "" {
		return "", nil
	}

	const oneCode = 49 // код единицы в ascii
	prev := rune(oneCode)
	for _, c := range s {
		isDigit := unicode.IsDigit(c)
		prevIsDigit := unicode.IsDigit(prev)

		if isDigit && prevIsDigit {
			return "", ErrInvalidString
		}

		if isDigit {
			number, _ := strconv.Atoi(string(c))
			if number == 0 { // Если 0, удаляем предыдущий символ
				buf := res.String()
				prevCharLenght := len(string(prev))
				buf = buf[0 : len(buf)-prevCharLenght]
				res.Reset()
				res.WriteString(buf)
			} else {
				res.WriteString(strings.Repeat(string(prev), number-1))
			}
		} else {
			res.WriteRune(c)
		}
		prev = c
	}
	return res.String(), nil
}
