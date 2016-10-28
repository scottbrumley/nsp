package main

import (
	"github.com/scottbrumley/nsp"
	"fmt"
)

func main() {
	myParms := nsp.GetParams()
	fmt.Print(nsp.SshCommand(myParms))
}
