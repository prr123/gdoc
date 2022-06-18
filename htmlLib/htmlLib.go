// htmlLIb.go
//
// author prr, azul software
// date: 18 June 20022
// copyrigth 2022 prr azul software
// license see github
//
// library to create basic html syntax
// read from yaml file

package htmlLib

import (
	"fmt"
	"os"
)

func creHtmlBase(outfil *os.File)(err error) {

	outfil.WriteString("<html>")


	return nil
}
