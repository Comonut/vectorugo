package main

import (
	"fmt"

	"github.com/comonut/vectorugo/store"
)

func main() {
	fmt.Print(*store.Random("", 32))

}
