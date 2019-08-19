package internal

import (
	"crypto/rand"
	"fmt"
	"strings"

	apierrs "k8s.io/apimachinery/pkg/api/errors"
)

const (
	CharactersUpper   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	CharactersLower   = "abcdefghijklmnopqrstuvwxyz"
	CharactersNumeric = "0123456789"
	CharactersSpecial = "~=+%^*()[]{}!#$?|"
)

func ContainsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func RemoveString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}

func IsInList(item string, items ...string) bool {
	for _, element := range items {
		if element == item {
			return true
		}
	}
	return false
}

func IgnoreNotFound(err error) error {
	if apierrs.IsNotFound(err) {
		return nil
	}
	return err
}

func RandomString(length int, charSet ...string) (string, error) {
	letters := strings.Join(charSet, "")
	bytes, err := GenerateRandomBytes(length)
	if err != nil {
		return "", fmt.Errorf("unable to generate random string: %s", err)
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes), nil
}

func GenerateRandomBytes(length int) ([]byte, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("unable to generate random bytes: %s", err)
	}

	return b, nil
}

func CoalesceString(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
