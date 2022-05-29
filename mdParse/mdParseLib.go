//
// mdParseLib.go
// parse markdown file
// usage: parse file.go
//
// author: prr azul software
// date: 28 May 2022
// copyright prr azul software
//
package mdParseLib

import (
	"os"
	"fmt"
	 "google/gdoc/util"
)

type mdParseObj struct {
	filnam string
	inBuf *[]byte
	elList []mdEl
}

type mdEl struct {
	elSt int
	elEnd int
	typ byte
	fchar byte
}

func InitMdParse() (mdp *mdParseObj) {

	mdp = new(mdParseObj)
	return mdp
}

func (mdP *mdParseObj) ParseMdFile(inpfilnam string) (err error) {
// function that opens md file

	var outfilnam string

	// check whether input file is valid

	// check for period
	iper := -1
	for i:=0; i< len(inpfilnam); i++ {
		if inpfilnam[i] == '.' {
			iper = i
			break
		}
	}

	if iper < 0 {return fmt.Errorf(" error parsing input file name, no period for extension found!")}

	// check for md extension
	if inpfilnam[iper+1] != 'm' {return fmt.Errorf(" error extension not md!")}
	if inpfilnam[iper+2] != 'd' {return fmt.Errorf(" error extension not md!")}

	inpfil, err := os.Open(os.Args[1])
	defer inpfil.Close()
	if err != nil {return fmt.Errorf("os.Open: %v\n", err)}

	inpfilInfo,_ := inpfil.Stat()
	inpSize := inpfilInfo.Size()

	fmt.Printf("*** input file:  %s Size: %d\n", inpfilnam, inpSize)

	// create output file name
//	found := false
	for i:=0; i< len(inpfilnam); i++ {
		if inpfilnam[i] == '.' {
//			found = true
			outfilnam = string(inpfilnam[:i])
			break
		}
	}

	mdP.filnam = outfilnam

    bufp := make([]byte, inpSize)
    nb, _ := inpfil.Read(bufp)
	if nb != int(inpSize) {return fmt.Errorf("error could not read file!")}
	mdP.inBuf = &bufp

	fmt.Println(" **** parsing md file!")

	mdP.parseMdOne()
	return nil
}

func (mdP *mdParseObj) parseMdOne()(err error) {
	var el mdEl

	buf := *(mdP.inBuf)
	ilin := 0
	ist := 0
	for i:=0; i< len(buf); i++ {
		if buf[i] == '\n' {
			el.elSt = ist
			el.elEnd = i
			mdP.elList = append(mdP.elList, el)
			ist = i+1
			ilin++
		}
	}

	fmt.Printf("lines: %d elList: %d\n", ilin, len(mdP.elList))

	mdP.printElList()
	return nil
}


func (mdP *mdParseObj) parseMdTwo()(err error) {
	var fch byte

	for el:=0; el<len(mdP.elList); el++ {
		fch = (*mdP.inBuf)[mdP.elList[el].elSt]
		switch fch {
			case '\r':
				// end of par?
				mdP.checkParEnd()

			case '#':
				// heading
				mdP.checkHeading()

			case '-':
				// horizontal ruler ---
				mdP.checkHr()
				// unordered list -
				mdP.checkUnList()

			case '_':
				// horizontal ruler ___wsp/ret
				mdP.checkHr()
				// bold text wsp__text__wsp
				mdP.checkBold()
				// italics wsp_text_wsp
				mdP.checkItalics()

			case '*':
				// horizontal ruler ***wsp|text
				mdP.checkHr()
				// bold text wsp**text**wsp
				mdP.checkBold()
				// unordered list *wsp
				mdP.checkUnList()
				// italics *text*
				mdP.checkItalics()

			case '>':
				// block quotes
				mdP.checkBlock()

			case '+':
				// unordered list + wsp|par
				mdP.checkUnList()

			case '~':
				// strike-through
				mdP.checkStrike()

			case ' ':
				//bold italics wsp*/_

			case '!':
				// image
				mdP.checkImage()

			case '[':
				// link [text](
				mdP.checkLink()

			case '|':
				// table |text|
				mdP.checkTable()

			default:

				if utilLib.IsNumeric(fch) {
				// ordered list 1.
					mdP.checkOrList()
				}
				if utilLib.IsAlpha(fch) {
				// paragraph
					mdP.checkPar()
				}
		}

	}
	return nil
}

func (mdP *mdParseObj) checkPar() {

}

func (mdP *mdParseObj) checkParEnd() {

}

func (mdP *mdParseObj) checkHeading() {

}

func (mdP *mdParseObj) checkHr() {

}

func (mdP *mdParseObj) checkBold(){

}

func (mdP *mdParseObj) checkItalics() {

}

func (mdP *mdParseObj) checkUnList() {

}

func (mdP *mdParseObj) checkBlock() {

}

func (mdP *mdParseObj) checkStrike() {

}

func (mdP *mdParseObj) checkImage() {

}

func (mdP *mdParseObj) checkLink() {

}

func (mdP *mdParseObj) checkTable() {

}

func (mdP *mdParseObj) checkOrList() {

}


func (mdP *mdParseObj) printElList()() {

	fmt.Printf("el Start End Fch text\n")
	for el:=0; el<len(mdP.elList); el++ {
		fmt.Printf("el %2d %3d %3d ", el, mdP.elList[el].elSt, mdP.elList[el].elEnd)
		str:=string((*mdP.inBuf)[mdP.elList[el].elSt:mdP.elList[el].elEnd])
		fmt.Printf("%q:%s\n", (*mdP.inBuf)[mdP.elList[el].elSt], str)
	}

}

func (mdP *mdParseObj) cvtMdToHtml()(err error) {

	fmt.Printf("*** output file: %s\n", mdP.filnam + ".md")

	outfil, err := os.Create(mdP.filnam + ".md")
	if err != nil { return fmt.Errorf("os.Create: %v\n", err)}
	defer outfil.Close()

	return nil
}

func IsTyp(buf []byte)(typNam string, res bool) {
// function that checks whether an input line is a type definition.
// It returns the type name if the input line is a typpe definition.

	if len(buf) < 5 { return "", false }

	if buf[0] != 't' {return "", false}
	if buf[1] != 'y' {return "", false}
	if buf[2] != 'p' {return "", false}
	if buf[3] != 'e' {return "", false}
	if buf[4] != ' ' {return "", false}

	fnamSt := 0
	for i:= 5; i<len(buf); i++ {
		if buf[i] == ' ' {continue}
		if utilLib.IsAlpha(buf[i]) {
			fnamSt = i
			break
		}
	}

	if fnamSt == 0 {return "", false}

	fnamEnd := 0
	for i:= fnamSt; i<len(buf); i++ {
		if (buf[i] == ' ') {
			fnamEnd = i
			break
		}
	}

	if fnamEnd == 0 {return "", false}

	typNam = string(buf[fnamSt:fnamEnd])
	return typNam, true
}

func IsFunc(buf []byte)(funcNam string, res bool) {
// function that checks whether an input line is a function.
// It returns the function name if the input line is a  function.

	if len(buf) < 5 { return "", false }

	if buf[0] != 'f' {return "", false}
	if buf[1] != 'u' {return "", false}
	if buf[2] != 'n' {return "", false}
	if buf[3] != 'c' {return "", false}
	if buf[4] != ' ' {return "", false}

	fnamSt := 0
	for i:= 5; i<len(buf); i++ {
		if buf[i] == ' ' {continue}
		if buf[i] == '(' {return "", false}
		if utilLib.IsAlpha(buf[i]) {
			fnamSt = i
			break
		}
	}

	if fnamSt == 0 {return "", false}

	fnamEnd := 0
	for i:= fnamSt; i<len(buf); i++ {
		if (buf[i] == ' ') || (buf[i] == '(') {
			fnamEnd = i
			break
		}
	}

	if fnamEnd == 0 {return "", false}

	funcNam = string(buf[fnamSt:fnamEnd])
	return funcNam, true
}

func IsMethod(buf []byte)(methNam string, typNam string, res bool) {
// function that detemines whether a input line is a  method.
// if so, the function returns the method name and the name of the structure the method is associated with

	if len(buf) < 5 { return "","", false }

	if buf[0] != 'f' {return "", "", false}
	if buf[1] != 'u' {return "", "", false}
	if buf[2] != 'n' {return "", "", false}
	if buf[3] != 'c' {return "", "", false}
	if buf[4] != ' ' {return "", "", false}

	typSt := 0
	for i:= 5; i<len(buf); i++ {
		if buf[i] == ' ' {continue}
		if utilLib.IsAlpha(buf[i]) {return "", "", false}
		if buf[i] == '(' {
			typSt = i
			break
		}
	}

	if typSt == 0 {return "", "", false}

	typEnd := 0
	for i:= typSt; i<len(buf); i++ {
		if buf[i] == ')' {
			typEnd = i
			break
		}
	}

	if typEnd == 0 {return "", "", false}

	typNam = string(buf[typSt+1:typEnd-1])
	methSt := 0
	for i:= typEnd + 1; i<len(buf); i++ {
		if utilLib.IsAlpha(buf[i]) {
			methSt = i
			break
		}
	}

	if methSt == 0 {return "", "", false}

	methEnd := 0
	for i:= methSt; i<len(buf); i++ {
		if (buf[i] == ' ') || (buf[i] == '(') {
			methEnd = i
			break
		}
	}

	if methEnd == 0 {return "", "", false}

	methNam = string(buf[methSt:methEnd])

	return methNam, typNam, true
}

func IsComment (buf []byte)(desc string, res bool) {
// function that determines whether the input line is a comment line.
// If so, it returns the comment in the desc.

	if len(buf) < 2 {return "", false}
	if buf[0] != '/' {return "", false}
	if buf[1] != '/' {return "",  false}

	desc = string(buf[2:])

	return desc, true
}
