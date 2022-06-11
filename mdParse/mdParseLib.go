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
	emEl bool
	comEl *comEl
	parEl *parEl
	tblEl *tblEl
	imgEl *imgEl
	ulEl *uList
	olEl *oList
}

type parEl struct {
	typ int
	fin bool
	subEl []parSubEl
	txtbuf []byte
}

type comEl struct {
	txt string
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
	ftn
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

func dispHtmlEl(num int)(str string) {

	var htmlDisp [15]string

	htmlDisp[0] = "br"
	htmlDisp[1] = "par"
	htmlDisp[2] = "hr"
	htmlDisp[3] = "span"
	htmlDisp[4] = "ul"
	htmlDisp[5] = "ol"
	htmlDisp[6] = "li"
	htmlDisp[7] = "img"
	htmlDisp[8] = "ftn"
	htmlDisp[9] = "h1"
	htmlDisp[10] = "h2"
	htmlDisp[11] = "h3"
	htmlDisp[12] = "h4"
	htmlDisp[13] = "h5"
	htmlDisp[14] = "h6"

	if num > len(htmlDisp)-1 {return ""}
	return htmlDisp[num]
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

	var fch, sch, tch byte

	mdP.istate = NE
//	for lin:=0; lin<len(mdP.linList); lin++ {
	maxLin := 30
	if len(mdP.linList) < maxLin {maxLin = len(mdP.linList)}

	fmt.Println("*************** lines **************************")
	for lin:=0; lin<maxLin; lin++ {

		linSt := mdP.linList[lin].linSt
		linEnd := mdP.linList[lin].linEnd
		linLen := linEnd - linSt

		fch = (*mdP.inBuf)[mdP.linList[lin].linSt]
		sch = 0
		tch = 0
		if linLen > 1 {sch = (*mdP.inBuf)[mdP.linList[lin].linSt + 1]}
		if linLen > 2 {tch = (*mdP.inBuf)[mdP.linList[lin].linSt + 2]}

		fmt.Printf("*** line %d len %d: state: %d fch: %c sch: %c tch: %c\n", lin, linLen, mdP.istate, fch, sch, tch)
		switch fch {
			case '\n':
				// end of par?
				// end of header?
				// is cr only char?
			fmt.Printf("cr istate: %d linSt: %d linEnd %d\n", mdP.istate, linSt, linEnd)
				if linEnd > linSt + 1 {return fmt.Errorf("line %d: text after cr", lin)}

				switch mdP.istate {
				case PAR:
					err = mdP.checkParEnd(lin)
					if err != nil {return fmt.Errorf("line %d ParEnd %v", lin, err)}
					err = mdP.checkBR()
					if err != nil {return fmt.Errorf("line %d checkBR %v", lin, err)}
					mdP.istate = BR

				case NE:
						// insert BR element
					err = mdP.checkBR()
					if err != nil {return fmt.Errorf("line %d checkBR %v", lin, err)}
					mdP.istate = BR

				case BR:
					err = mdP.checkBR()
					if err != nil {return fmt.Errorf("line %d checkBR %v", lin, err)}
					mdP.istate = BR

				default:
					return fmt.Errorf("line %d istate %d cr", lin, mdP.istate)
				}

			case '#':
				// heading
				err = mdP.checkHeading(lin)
				if err != nil {return fmt.Errorf("line %d: %v", lin, err)}
				mdP.istate = NE

			case '-':
				// horizontal ruler ---
				mdP.checkHr()
				// unordered list -
				mdP.checkUnList(lin)

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
					mdP.checkUnList(lin)
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
				mdP.checkUnList(lin)

			case '~':
				// strike-through
				mdP.checkStrike()

			case ' ':
				// ws
				// nested list
				mdP.checkUnList(lin)
/*
				FNCh, err := mdP.checkWs(lin)
				if err != nil { return fmt.Errorf("lin %d checkWs: %v",lin, err)}
				fmt.Printf("FNCh: %c\n",FNCh)
				if FNCh == '*' {
					mdP.checkUnList(lin)
					break
				}
*/
			case '!':
				// image
				mdP.checkImage()

			case '[':
				// check comment
				res, err := mdP.checkComment(lin)
				if err != nil { return fmt.Errorf("lin %d comment: %v",lin, err)}
				if res {break}
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
					fmt.Println(string((*mdP.inBuf)[linSt:linEnd]))
					mdP.checkPar(lin)
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

fmt.Printf("EOL: '%c' '%c'\n", buf[linLen-2], buf[linLen-1])
	if (buf[linLen-2] == ' ') && (buf[linLen-1] == ' ') { return true}
	return false
}

func (mdP *mdParseObj) checkWs(lin int)(fch byte, err error) {

	linSt := mdP.linList[lin].linSt
	linEnd := mdP.linList[lin].linEnd
//	linLen := linEnd - linSt

	buf := (*mdP.inBuf)
	fnchpos := 0;
	numWs :=0
	for i:=linSt; i < linEnd; i++ {
		if buf[i] == ' ' {
			numWs++
		} else {
			fnchpos = i
			break
		}
	}
	if fnchpos == 0 { return 0, fmt.Errorf("line %d all ws", lin) }

	fch = buf[fnchpos]
fmt.Printf("fch: %c numWs: %d\n", fch, numWs)

	return fch, nil
}

func (mdP *mdParseObj) checkComment(lin int)(res bool, err error) {
// method that checks whether line is a comment

	var el structEl
	var comEl comEl

	linSt := mdP.linList[lin].linSt
	linEnd := mdP.linList[lin].linEnd
	linLen := linEnd - linSt
	if linLen < 10 {return false, nil}
	buf := (*mdP.inBuf)
	if buf[linSt + 1] != '/' {return false, nil}
	if buf[linSt + 2] != '/' {return false, nil}
	if buf[linSt + 3] != ']' {return false, nil}
	if buf[linSt + 4] != ':' {return false, nil}
	if buf[linSt + 5] != ' ' {return false, nil}
	if buf[linSt + 6] != '*' {return false, nil}
	if buf[linSt + 7] != ' ' {return false, nil}
	if buf[linSt + 8] != '(' {return false, fmt.Errorf("lin %d comment: no '(' found!", lin)}

	closPar :=0
	for i:= linEnd; i> linSt+8; i-- {
		if buf[linSt + 8] != ')' {
			closPar = i
			break
		}
	}
	if closPar == 0 {return false, fmt.Errorf("lin %d no ')' found!", lin)}

	comEl.txt = string(buf[linSt+9:closPar-1])
	el.comEl = &comEl
	mdP.elList = append(mdP.elList, el)


	return true, nil
}

func (mdP *mdParseObj) checkPar(lin int)(res bool, err error) {

	var el structEl
	var parEl parEl
	var subEl parSubEl

fmt.Println("*** checkpar start")
	linSt := mdP.linList[lin].linSt
	linEnd := mdP.linList[lin].linEnd
	buf := (*mdP.inBuf)

	crEOL := checkEOL(buf[linSt:linEnd-1])
	if crEOL {
		parEl.fin = true
		mdP.istate = NE
		subEl.txt = string(buf[linSt:linEnd+1])
	} else {
		parEl.fin = false
		mdP.istate = PAR
		subEl.txt = string(buf[linSt:linEnd])
	}

	parEl.subEl = append(parEl.subEl, subEl)
	el.parEl = &parEl
	mdP.elList = append(mdP.elList, el)
	mdP.istate = PAR

	return true, nil
}

func (mdP *mdParseObj) checkParEnd(lin int) (err error){

fmt.Println("*** par end")
	elLast := len(mdP.elList) -1
	lastEl := mdP.elList[elLast]
	// check for empty line
	if lastEl.emEl {return nil}
	// see whether last element was a par el
	if lastEl.parEl == nil {return nil}
	// parEl pointer points to a parEl
	ParEl := *lastEl.parEl
	lastPar := len(ParEl.subEl) -1
	subEl := ParEl.subEl[lastPar]
	subEl.txt +="\n"
	ParEl.fin = true
	mdP.istate = NE
	return nil
}

/*
func (mdP *mdParseObj) checkHeadEnd(lin int) (err error){

fmt.Println("*** heading end")
	elLast := len(mdP.elList) -1
	lastEl := mdP.elList[elLast]
	HdEl := *lastEl.hdEl
	lastPar := len(HdEl.subEl) -1
	subEl := HdEl.subEl[lastPar]
	subEl.txt +="\n"
	HdEl.fin = true
	mdP.istate = NE
	return nil
}
*/

func (mdP *mdParseObj) checkHeading(lin int) (err error){
// function that parses a line for headings
	var el structEl
	var parEl parEl
//	listEl := mdP.elList[el]
	linSt := mdP.linList[lin].linSt
	linEnd := mdP.linList[lin].linEnd
	buf := (*mdP.inBuf)
fmt.Printf("heading buffer: %s\n", buf[linSt:linEnd])
	hd := 0
	hdEnd := linSt
	for i:= 0; i < 7; i++ {
		if buf[i] != '#' {
			if buf[i] != ' ' {
				return fmt.Errorf("line %d: no ws after #!", lin)
			}
			hdEnd = i
			break
		}
		hd++
	}
	// check the end of the line

	crEOL := checkEOL(buf[hdEnd+1:linEnd])
	// last char is cr. ergo paragraph not finished

	txtstr := string(buf[hdEnd+1:linEnd])
	mdP.istate = PAR
	hdStr := fmt.Sprintf("h%d", hd)
	hdtyp := 0
	switch hd {
		case 1:
			hdtyp = h1
		case 2:
			hdtyp = h2
		case 3:
			hdtyp = h3
		case 4:
			hdtyp = h4
		case 5:
			hdtyp = h5
		case 6:
			hdtyp = h6
		default:
	}

	fmt.Printf("header: %d %s text: \"%s\" \n", hdtyp, hdStr, txtstr)
	parEl.typ = hdtyp
	parEl.txtbuf = buf[hdEnd+1:len(buf)-1]
	parEl.fin = false
	if crEOL {parEl.fin = true}
	el.parEl = &parEl
	mdP.elList = append(mdP.elList, el)
	mdP.istate = PAR
	return nil
}

func (mdP *mdParseObj) checkHr() {

}

func (mdP *mdParseObj) checkBold(){

}

func (mdP *mdParseObj) checkItalics() {

}

// lll
func (mdP *mdParseObj) checkUnList(lin int) {

	var el structEl
	var ulEl uList
	var parEl parEl
	var subEl parSubEl

//	listEl := mdP.elList[el]
	linSt := mdP.linList[lin].linSt
	linEnd := mdP.linList[lin].linEnd
	buf := (*mdP.inBuf)
	parSt:=0
	istate := 0
	wsNum := 0
	for i:= linSt; i< linEnd; i++ {
		switch istate {
			case 0:
				switch buf[i] {
				case ' ':
					wsNum++
				case '*':
					istate = 1
				default:
				}

			case 1:
				switch buf[i] {
				case ' ':

				default:
					parSt = i
				}

			default:
		}
		if parSt > 0 {break}
	}
	// nest lev
	nestLev := wsNum/4

	crEOL := checkEOL(buf[parSt:linEnd-1])
	if crEOL {
		parEl.fin = true
		mdP.istate = NE
		subEl.txt = string(buf[parSt:linEnd+1])
	} else {
		parEl.fin = false
		mdP.istate = PAR
		subEl.txt = string(buf[parSt:linEnd])
	}
fmt.Printf("uel: %s\n", subEl.txt)
	parEl.subEl = append(parEl.subEl, subEl)
	ulEl.nest = nestLev
	ulEl.parEl = &parEl
	el.ulEl = &ulEl
	mdP.elList = append(mdP.elList, el)
	mdP.istate = PAR

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

func (mdP *mdParseObj) checkBR()(err error) {

	var el structEl

	el.emEl = true
	mdP.elList = append(mdP.elList, el)

	return nil
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

	fmt.Println("*********** El List ***********")
	fmt.Printf("  el nam typ  subels fin txt\n")
	for i:=0; i < len(mdP.elList); i++ {
		el := mdP.elList[i]
		fmt.Printf("el %3d: ", i)
		if el.emEl {
			fmt.Printf("eL: %t", el.emEl)
		}
		if el.comEl != nil {
			fmt.Printf( "com: %s", el.comEl.txt)
		}
		if el.parEl != nil {
			ParEl := *el.parEl
			fmt.Printf( "par %-5s: subels: %d status: %t\n", dispHtmlEl(ParEl.typ), len(ParEl.subEl), ParEl.fin)
			subLen := len(ParEl.subEl)
			if subLen == 1 {
				fmt.Printf("         subel 0: %s\n", ParEl.subEl[0].txt)
			} else {
				for i:=0; i< subLen; i++ {
					fmt.Printf("         subel %d: %s\n", i, ParEl.subEl[i].txt)
				}
			}
		}
		if el.tblEl != nil {
			fmt.Printf( "tbl: rows: %d  cols: %d ", el.tblEl.rows, el.tblEl.cols)
		}
		if el.imgEl != nil {
			fmt.Printf( "img h: %d w: %d src: %s", el.imgEl.height, el.imgEl.width, el.imgEl.src)
		}
		if el.ulEl != nil {
			fmt.Printf( "ul nest %d ", el.ulEl.nest)
			if el.ulEl.parEl != nil {
				ParEl := el.ulEl.parEl
				fmt.Printf( "par typ: %-5s subels: %d stat: %t", dispHtmlEl(ParEl.typ), len(ParEl.subEl), ParEl.fin)
				subLen := len(ParEl.subEl)
				if subLen == 1 {
					fmt.Printf("         subel 0: %s", ParEl.subEl[0].txt)
				} else {
					for i:=0; i< subLen; i++ {
						fmt.Printf("         subel %d: %s\n", i, ParEl.subEl[i].txt)
					}
				}
			} else {
				fmt.Printf( "ulEl par nil!")
			}
		}
		if el.olEl != nil {
			fmt.Printf( "ol nest %d ", el.olEl.nest)
			if el.olEl.parEl != nil {
				ParEl := el.olEl.parEl
				fmt.Printf( "par typ: %-5s subels: %d stat %t", dispHtmlEl(ParEl.typ), len(ParEl.subEl), ParEl.fin)
				subLen := len(ParEl.subEl)
				if subLen == 1 {
					fmt.Printf("         subel 0: %s", ParEl.subEl[0].txt)
				} else {
					for i:=0; i< subLen; i++ {
						fmt.Printf("         subel %d: %s\n", i, ParEl.subEl[i].txt)
					}
				}
			} else {
				fmt.Printf( "olEl par nil!")
			}
		}

		if !(el.parEl != nil) {fmt.Printf("\n")}
	}
}
