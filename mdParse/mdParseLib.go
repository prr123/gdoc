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
	linList []mdLin
	elList []structEl
	istate int
}

type mdLin struct {
	linSt int
	linEnd int
	typ byte
	fchar byte
}

type structEl struct {
	parEl *parEl
	tblEl *tblEl
	imgEl *imgEl
	ulEl *uList
	olEl *oList
}

type parEl struct {
	typ int
	subEl []parSubEl
}

type tblEl struct {
	rows int
	cols int
	caption string
}

type imgEl struct {
	width int
	height int
	src string
	alt string
	title string
}

type uList struct {
	nest int
	parEl *parEl
}

type oList struct {
	nest int
	parEl *parEl
}

type parSubEl struct {
	elSt int
	elEnd int
	txt string
	txtTyp []int
}

const (
	NE = iota
	BR
	PAR
	UL0
	OL0
	UL1
	OL1
	UL2
	OL2
)

//html elements
const (
	br = iota
	par
	hr
	span
	ul
	ol
	li
	img
	ftnote
	h1
	h2
	h3
	h4
	h5
	h6
)

// text attributes
const (
	bold = iota
	italic
	strike
	html
)

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
	var lin mdLin

	buf := *(mdP.inBuf)
	ilin := 0
	ist := 0
	for i:=0; i< len(buf); i++ {
		if buf[i] == '\n' {
			lin.linSt = ist
			lin.linEnd = i
			mdP.linList = append(mdP.linList, lin)
			ist = i+1
			ilin++
		}
	}

	fmt.Printf("lines: %d elList: %d\n", ilin, len(mdP.linList))

	mdP.printLinList()
	err = mdP.parseMdTwo()
	if err != nil {
		fmt.Printf("error %v\n", err)
	}
	return nil
}


func (mdP *mdParseObj) parseMdTwo()(err error) {
// function that parses the linList and create an el List
//	var el  structEl

	mdP.istate = NE
//	for lin:=0; lin<len(mdP.linList); lin++ {
	for lin:=0; lin<10; lin++ {
		linSt := mdP.linList[lin].linSt
		linEnd := mdP.linList[lin].linEnd

		fch := (*mdP.inBuf)[mdP.linList[lin].linSt]
		sch := (*mdP.inBuf)[mdP.linList[lin].linSt + 1]
		tch := (*mdP.inBuf)[mdP.linList[lin].linSt + 2]
		fmt.Printf("*** line %d: state: %d\n", lin, mdP.istate)
		switch fch {
			case '\r':
				// end of par?
				// is cr only char?
		fmt.Printf("cr istate: %d linSt: %d linEnd %d\n", mdP.istate, linSt, linEnd)
				if linSt+1 == linEnd {
					if mdP.istate == PAR {
						mdP.checkParEnd()
						mdP.istate = NE
					}
				} else {
					return fmt.Errorf("line %d: text after cr", lin)
				}
			case '#':
				// heading
				err = mdP.checkHeading(lin)
				if err != nil {
					return fmt.Errorf("line %d: %v", lin, err)
				}

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
				// check whether next char is whitespace
				if sch == ' ' {
					// unordered list *wsp
					fmt.Println("*** unordered list item")
					mdP.checkUnList()
					break
				}

				if sch == '*' {
					if utilLib.IsAlpha(tch) {
						// bold text wsp**text**wsp
						fmt.Println("*** start bold")
						mdP.checkBold()
					}
					if tch == '*' {
						fmt.Println("*** horizontal ruler")
						// horizontal ruler ***wsp|text
						mdP.checkHr()
					}
					break
				}
				if utilLib.IsAlpha(sch) {
					fmt.Println("*** start italics")
					// italics *text*
					mdP.checkItalics()
					break
				}
			case '>':
				fmt.Println("*** start blockquote")

				// block quotes
				mdP.checkBlock()

			case '<':
				// html start <>
				mdP.checkHtml()

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
					fmt.Println("*** par start ***")
					mdP.istate = PAR
					str := string((*mdP.inBuf)[linSt:linEnd-1])
					fmt.Println(str)
					mdP.checkPar()
				}
		}

	}
	mdP.printElList()
	return nil
}

func checkEOL(buf []byte)(res bool) {
// function to test whether a line is completed with a md cr

	linLen := len(buf)
	if linLen < 2 {return false}

	if (buf[linLen] == ' ') && (buf[linLen -1] == ' ') { return true}
	return false
}

func (mdP *mdParseObj) checkPar() {
	fmt.Println("*** par start")
}

func (mdP *mdParseObj) checkParEnd() {
	fmt.Println("*** par end")
}

func (mdP *mdParseObj) checkHeading(lin int) (err error){
// function that parses a line for headings
	var el structEl
	var parEl parEl
//	listEl := mdP.elList[el]
	linSt := mdP.linList[lin].linSt
	linEnd := mdP.linList[lin].linEnd
	buf := (*mdP.inBuf)[linSt:linEnd]
	fmt.Printf("buffer: %s\n", buf[:])
	hd := 0
	hdEnd := linSt
	for i:= 0; i < 7; i++ {
		if buf[i] != '#' {
			if buf[i] != ' ' {
				return fmt.Errorf("no ws after #!")
			}
			hdEnd = i
			break
		}
		hd++
	}
	// check the end of the line
//	crEOL := checkEOL(buf[hdEnd+1:len(buf)-1])
	// last char is cr. ergo paragraph not finished

	txtstr := string(buf[hdEnd+1:len(buf)-1])
	mdP.istate = PAR
	hdStr := fmt.Sprintf("h%d", hd)
	hdtyp := 0
	switch hd {
		case 1:
			hdtyp = h1
		case 2:
			hdtyp = h2
		default:
	}

	fmt.Printf("header: %d %s text: \"%s\" \n", hdtyp, hdStr, txtstr)
	parEl.typ = hdtyp
	el.parEl = &parEl
	mdP.elList = append(mdP.elList, el)

	return nil
}

func (mdP *mdParseObj) checkHr() {

}

func (mdP *mdParseObj) checkBold(){

}

func (mdP *mdParseObj) checkItalics() {

}

func (mdP *mdParseObj) checkUnList() {

}

func (mdP *mdParseObj) checkHtml() {

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


func (mdP *mdParseObj) printLinList()() {

	fmt.Printf("line Start End Fch text\n")
	for el:=0; el<len(mdP.linList); el++ {
		fmt.Printf("el %2d %3d %3d ", el, mdP.linList[el].linSt, mdP.linList[el].linEnd)
		str:=string((*mdP.inBuf)[mdP.linList[el].linSt:mdP.linList[el].linEnd])
		fmt.Printf("%q:%s\n", (*mdP.inBuf)[mdP.linList[el].linSt], str)
	}

}

func (mdP *mdParseObj) cvtMdToHtml()(err error) {
// function that converts the structural Element
	fmt.Printf("*** input file: %s\n", mdP.filnam + ".md")

	outfil, err := os.Create(mdP.filnam + ".html")
	if err != nil { return fmt.Errorf("os.Create: %v\n", err)}
	defer outfil.Close()

	return nil
}

func (mdP *mdParseObj) printElList () {
// function that prints out the structural element list

	fmt.Println("*********** El List *****")
	fmt.Printf(" el  typ subel\n")
	for i:=0; i < len(mdP.elList); i++ {
		el := mdP.elList[i]
		fmt.Printf("  %d : ", i)
		if el.parEl != nil {
			fmt.Printf( "par %d %d ", el.parEl.typ, len(el.parEl.subEl))
		}
		if el.tblEl != nil {
			fmt.Printf( "tbl %d %d ", el.tblEl.rows, el.tblEl.cols)
		}
		if el.imgEl != nil {
			fmt.Printf( "img %d %d %s", el.imgEl.height, el.imgEl.width, el.imgEl.src)
		}
		if el.ulEl != nil {
			fmt.Printf( "ul %d ", el.ulEl.nest)	
			if el.ulEl.parEl != nil {
				fmt.Printf( "par %d %d ", el.ulEl.parEl.typ, len(el.ulEl.parEl.subEl))
			} else {
				fmt.Printf( "ulEl par nil!")
			}
		}
		if el.olEl != nil {
			fmt.Printf( "ol %d ", el.olEl.nest)
			if el.olEl.parEl != nil {
				fmt.Printf( "par %d %d ", el.olEl.parEl.typ, len(el.olEl.parEl.subEl))
			} else {
				fmt.Printf( "olEl par nil!")
			}
		}

		fmt.Printf("\n")
	}
}
