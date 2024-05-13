package hw10programoptimization

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"strings"
)

type User struct {
	Email string
}

type DomainStat map[string]int

var ErrInvalidInputData = errors.New("неверный формат входных данных")

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		var user User
		if err := json.Unmarshal(sc.Bytes(), &user); err != nil {
			return result, ErrInvalidInputData
		}

		matched := strings.Contains(user.Email, domain)

		if !matched {
			continue
		}

		fullDomain := strings.ToLower(strings.Split(user.Email, "@")[1])
		result[strings.ToLower(fullDomain)]++
	}
	return result, nil
}
