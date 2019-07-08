package main

import (
	"fmt"

	"github.com/Comonut/vectorugo/server"

	"github.com/Comonut/vectorugo/store"
)

func main() {
	var onesVector = store.Ones("", 32)
	fmt.Print(onesVector.Sum())

	server.Init()

}
