package main

import (
	"fmt"
	"os"
	gdocUtil "google/gdoc/gdocUtil"
)

func main () {
	var tstInpStr string

	if len(os.Args) > 1 {
		tstInpStr = os.Args[1]
	} else {
    	tstInpStr = "Step%0.%1:"
	}
	fmt.Println("****** testing glyph Format String  *****")
	fmt.Printf("  test string: \"%s\"\n", tstInpStr)

    glFmt, err := gdocUtil.ParseGlyphFormat(tstInpStr)
    if err != nil {
        fmt.Println("error Parsing tstInpStr %v", err)
        os.Exit(1)
    }

//	fmt.Println("*** printing glFmt output ***")
	gdocUtil.PrintGlFmt(glFmt)

	fmt.Println("********** success *****************")
}
