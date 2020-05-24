package utils

import "fmt"

func Logg(x interface{}) {
	fmt.Printf("%+v\n", x)
}