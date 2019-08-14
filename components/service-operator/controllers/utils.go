package controllers

import (
	apierrs "k8s.io/apimachinery/pkg/api/errors"
)

const (
	charactersUpper   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	charactersLower   = "abcdefghijklmnopqrstuvwxyz"
	charactersNumeric = "0123456789"
	charactersSpecial = "~=+%^*()[]{}!#$?|"
)

func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func removeString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}

func isInList(item string, items ...string) bool {
	for _, element := range items {
		if element == item {
			return true
		}
	}
	return false
}

func ignoreNotFound(err error) error {
	if apierrs.IsNotFound(err) {
		return nil
	}
	return err
}
