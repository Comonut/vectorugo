package main

import (
	"fmt"

	"github.com/comonut/vectorugo/store"
)

func main() {
	var onesVector = store.Ones("", 32)
	fmt.Print(onesVector.Sum())

}
