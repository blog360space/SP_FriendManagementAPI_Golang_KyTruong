package common

import (
	"github.com/mcnijman/go-emailaddress"
)

func getIndex(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1
}

func removeItemInStringSlice(slice []string, index int) []string {
	slice[index] = slice[len(slice)-1]
	slice[len(slice)-1] = ""
	slice = slice[:len(slice)-1]

	return slice
}

func IsValidEmail(email string) bool {
	_, err := emailaddress.Parse(email)

	return err == nil
}

func GetIndex(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1
}

func RemoveItemInStringSlice(slice []string, index int) []string {
	slice[index] = slice[len(slice)-1]
	slice[len(slice)-1] = ""
	slice = slice[:len(slice)-1]

	return slice
}
