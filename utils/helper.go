package utils

import (
	"fmt"
)

func Typeof(a any) {
	fmt.Printf("the variable %v is of type of %T \n", a, a)
}

func Contains(x []string, v string) bool {
	for _, s := range x {
		if v == s {
			return true
		}
	}
	return false
}
