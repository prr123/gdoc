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
	hrEl bool
	comEl *comEl
	parEl *parEl
	tblEl *tblEl
	imgEl *imgEl
	ulEl *uList
	olEl *oList
	bkEl *bkEl
	errEl *errEl
}

type parEl struct {
	typ int
	fin bool
	txtSt int
	txtEnd int
	txt string
	subEl []parSubEl
}

type parSubEl struct {
	elSt int
	elEnd int
	txt string
	txtTyp int
}

type errEl struct {
	fch byte
	line int
	errmsg string
	txt string
}

type bkEl struct {
	parEl *	parEl
	nest int
}

type comEl struct {
	txt string
}

type tblEl struct {
	rows int
	cols int
	trows []tblRow
	caption string
}

type tblRow struct {
	tblRow []tblCel
}

type tblCel struct {
	parEl *parEl
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

// line attributes
const (
	EL = iota
	HR
	PAR
	UL
	OL
	IMG
	BLK
	ERR
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
	ftnote
	sup
	sub
)

func dispState(num int)(str string) {
// function that converts state constants to strings

	var stateDisp [7]string

	stateDisp[0] = "EL"
	stateDisp[1] = "HR"
	stateDisp[2] = "PAR"
	stateDisp[3] = "UL"
	stateDisp[4] = "OL"
	stateDisp[5] = "IMG"
	stateDisp[6] = "BL"
//	stateDisp[7] = "UL2"
//	stateDisp[8] = "OL2"

	if num > len(stateDisp)-1 {return ""}
	return stateDisp[num]
}

func dispHtmlEl(num int)(str string) {
// function that converts html const to strings

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
// function that initialises the mdParse object

	mdp = new(mdParseObj)
	return mdp
}

func (mdP *mdParseObj) ParseMdFile(inpfilnam string) (err error) {
// method that opens an md file and reads it into a buffer

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
	inpfil.Close()

	fmt.Println("\n******* parsing md file! ************")

	mdP.parseMdOne()
	return nil
}

func (mdP *mdParseObj) parseMdOne()(err error) {
// method that conducts the first pass of the parser.
// the first pass creates and slice of lines.

	var lin mdLin

	buf := *(mdP.inBuf)
	ilin := 0
	ist := 0
	for i:=0; i< len(buf) ; i++ {

		if buf[i] == '\n' || buf[i] == '\r' {
			lin.linSt = ist
			lin.linEnd = i
			mdP.linList = append(mdP.linList, lin)
			ist = i+1
			ilin++
			if buf[i] == '\r' {
				i++
				ist++
			}
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
// method that reads the linList, parses the list and creates a List of elements

	var fch, sch byte

	mdP.istate = EL
//	for lin:=0; lin<len(mdP.linList); lin++ {
	maxLin := 500
	if len(mdP.linList) < maxLin {maxLin = len(mdP.linList)}

	fmt.Println("\n*************** lines **************************")
	for lin:=0; lin<maxLin; lin++ {

		linSt := mdP.linList[lin].linSt
		linEnd := mdP.linList[lin].linEnd
		linLen := linEnd - linSt

		fch = (*mdP.inBuf)[mdP.linList[lin].linSt]
		sch = 0
		if linLen > 1 {sch = (*mdP.inBuf)[mdP.linList[lin].linSt + 1]}

		if fch == '\n'||fch == '\r' {
			fmt.Printf("*** line %d: fch: CR sch: %q ", lin, sch)
		} else {
			fmt.Printf("*** line %d: fch: %q sch: %q ", lin, fch, sch)
		}
		switch fch {
			case '\n', '\r':
				// end of par?
				// end of header?
				// is cr only char?
				if linEnd > linSt + 1 {
					fmt.Printf("line %d: text after cr", lin)
					break
				}
				switch mdP.istate {
				case PAR:
					err = mdP.checkParEnd(lin)
					if err != nil {fmt.Printf("line %d ParEnd %v\n", lin, err)}
					err = mdP.checkBR()
					if err != nil {fmt.Printf("line %d checkBR %v\n", lin, err)}

				case EL, HR:
					err = mdP.checkBR()
					if err != nil {fmt.Printf("line %d checkBR %v\n", lin, err)}

				case UL, OL:
					err = mdP.checkBR()
					if err != nil {fmt.Printf("line %d checkBR %v\n", lin, err)}

				default:
					fmt.Printf("line %d error istate %s cr", lin, dispState(mdP.istate))
				}

			case '#':
				// headings
				err = mdP.checkHeading(lin)
				if err != nil {fmt.Printf("line %d: heading error %v\n", lin, err)}

			case '_':
				// horizontal ruler ___wsp/ret
				if linLen < 4 {
					fmt.Printf("line %d: linLen %d too short for HR!", lin, linLen)
					break
				}
				err = mdP.checkHr(lin)
				if err != nil {
					err1 := mdP.checkPar(lin)
					if err1 != nil { fmt.Printf("line %d: neither HR %v nor Par %v\n", err, err1)}
				}
			case '*','+','-':
				// check whether next char is whitespace
				switch {
				case sch == ' ':
					// unordered list *wsp
		fmt.Printf("un List NL:0 ")
					err = mdP.checkUnList(lin)
					if err != nil {fmt.Printf("line %d: unordered list error %v\n", lin, err)}

				case sch == fch:
					err = mdP.checkHr(lin)
					if err != nil {
						err1 := mdP.checkPar(lin)
						if err1 != nil { fmt.Printf("line %d: neither HR %v nor Par %v\n", err, err1)}
					}

				case utilLib.IsAlpha(sch) :
					fmt.Println(string((*mdP.inBuf)[linSt:linEnd]))
					err = mdP.checkPar(lin)
					if err != nil {fmt.Printf("line %d: par error %v\n", lin, err)}

				default:
				// error
					if err != nil {fmt.Printf("line %d unsuitable char %q %d after *+-\n", sch, sch)}
//					mdP.checkError(lin, fch, errStr)
				}

			case '>':
				fmt.Println("*** start blockquote")

				// block quotes
				err = mdP.checkBlock(lin)
				if err != nil {fmt.Printf("line %d: quote block error %v\n", lin, err)}

			case '`':
				fmt.Println("*** start code")

				// block quotes
				err = mdP.checkCode(lin)
				if err != nil {fmt.Printf("line %d: code block error %v\n", lin, err)}

			case ' ':
				// ws
				// nested list
//				nestLev, ord, err := mdP.checkIndent(lin)
				ch, wsNum, err := mdP.checkIndent(lin)
				if err != nil {
					fmt.Printf("line %d: indent error %v\n", lin, err)
					break
				}

				if utilLib.IsAlpha(ch) {
					err = mdP.checkBlock(lin)
					if err != nil {fmt.Printf("line %d: ind block error %v\n", lin, err)}
					break
				}

				if (ch == '*' || ch == '+') || ch == '-' {
					err := mdP.checkUnList(lin)
					if err != nil {fmt.Printf("line %d: unordered list error %v\n", lin, err)}
		fmt.Printf(" un List Nl: %d", wsNum/4)
					break
				}

				if utilLib.IsNumeric(ch) {
					err = mdP.checkOrList(lin)
					if err != nil {fmt.Printf("line %d: ordered list error %v\n", lin, err)}
		fmt.Printf(" or List Nl: %d", wsNum/4)
					break
				}
				fmt.Printf("line %d: indent error start char %q in pos: %d\n", lin, ch, wsNum +1)

			case '!':
				// image
				err = mdP.checkImage(lin)
				if err != nil { return fmt.Errorf("lin %d image: %v\n",lin, err)}

			case '[':
				// check comment
				err = mdP.checkComment(lin)
				if err != nil {
					err1 := mdP.checkPar(lin)
					if err1 != nil {fmt.Printf("lin %d comment: %v\n        par: %v\n",lin, err, err1) }
				}

			case '|':
				// table |text|
				endLin, err := mdP.checkTable(lin)
				if err != nil {fmt.Printf("line %d: table error %v\n", lin, err)}
				lin = endLin

			default:

				if utilLib.IsNumeric(fch) {
				// ordered list 1.
					err = mdP.checkOrList(lin)
					if err != nil {fmt.Printf("line %d: ordered list error %v\n", lin, err)}
		fmt.Printf("or List NL:0 ")
					break
				}

				if utilLib.IsAlpha(fch) {
					// paragraph
//			fmt.Println("*** par start ***")
					mdP.istate = PAR
//			fmt.Println(string((*mdP.inBuf)[linSt:linEnd]))
					err = mdP.checkPar(lin)
					if err != nil {fmt.Printf("line %d: par error %v\n", lin, err)}
					break
				}
				fmt.Printf("line %d: no fit for first char: %q\n", lin, fch)
//				mdP.checkError(lin, fch, errmsg)
		}
	fmt.Printf(" state: %s\n",dispState(mdP.istate))
	}

	mdP.printElList()
	return nil
}

//ppp
func (mdP *mdParseObj) parseMdTextEls()(err error) {
// method that parses the text fields

	for i:=0; i < len(mdP.elList); i++ {
		el := mdP.elList[i]
		fmt.Printf("el %3d: ", i)

		// parse par el
		if el.parEl != nil {

		}

		// parse par el
		if el.bkEl != nil {

		}

		//parse blkEl
		if el.ulEl != nil {

		}

		//parse blkEl
		if el.olEl != nil {

		}

	}

	return nil
}

func (mdP *mdParseObj) parseMdTxt(parEl *parEl)(err error) {

	buf := []byte(parEl.txt)

	linkSt := 0
	linkEnd := 0
	uriSt := 0
	uriEnd := 0

	ftnNumSt:=0
	ftnNumEnd:=0
	ftnStrSt:=0
	ftnStrEnd :=0

	for i:=0; i< len(buf); i++ {
		ch := buf[i]
		istate := 0
		switch istate {
		case 0:
			switch ch {
			case ' ':
				// ws
				istate = 1

			case '*':
				// *
				istate = 2

			case '_' :
				// _
				istate = 11

			case '[' :
				// [
				istate = 40

			default:
			}

		case 1:
			// ws
			if utilLib.IsAlphaNumeric(ch) {
				// wst
				istate = 0
			}
			if ch == '*' {
				// ws*
				istate = 2
			}
			if ch == '[' {
				// ws[
				istate = 40
			}

		case 2:
			// *
			if ch == '*' {
				// **
				istate = 20
				break
			}

			if utilLib.IsAlphaNumeric(ch) {
				// *t
				istate = 3
			}

		case 3:
				// *t
			if ch == '*' {
				// *t*
				istate = 4
			}

		case 4:
			// *txt*
			if ch == ' ' {
				// *txt*ws
				istate = 1
			} else {
				//error
			}

		case 11:
			// _
			if ch == '_' {
				// __
				istate = 15
				break
			}

			if utilLib.IsAlpha(ch) {
				// _t
				istate = 12
			}
		case 12:
			// _t
			if ch == '_' {
				// _t_
				istate = 13
			}
		case 13:
			// _t_
			if ch == ' ' {
				// _t_ws
				istate = 4
			}

		case 15:
			// __
			if utilLib.IsAlpha(ch) {
				// __t
				istate = 16
			}

		case 16:
			// __t
			if ch == '_' {
				// __t_
				istate = 17
			}

		case 17:
			// __t_
			if ch == '_' {
				// __t__
				istate = 18
			} else {
				//error
			}

		case 18:
			// __t__
			if ch == ' ' {
				// __t__ws
				istate = 4
			}

		case 20:
			// **
			if ch == '*' {
				// ***
				istate = 30
			}
			if utilLib.IsAlphaNumeric(ch) {
				// **t
				istate = 21
			}
			if ch == ' ' {
				//error
				istate = 50
			}

		case 21:
			// **t
			if ch == '*' {
				// **t*
				istate = 22
			}

		case 22:
			// **text*
			if ch == '*' {
				// **text**
				istate = 4
			}

		case 30:
			// ***
			if utilLib.IsAlphaNumeric(ch) {
				istate = 31
			}
			if ch == ' ' {
				//error ***ws
				istate = 50
			}

		case 31:
			// ***t
			if ch == '*' {
				// ***t*
				istate = 32
			}

		case 32:
			// ***text*
			if ch == '*' {
				// ***text**
				istate = 33
			}

		case 33:
			// **text**
			if ch == '*' {
				// ***text***
				istate = 4
			}

// links & footnotes
		case 40:
			// [
			if utilLib.IsAlphaNumeric(ch) {
				// [t
				linkSt = i
				istate = 41
			}
			if ch == '^' {
				// [^
				istate = 50
			}

		case 41:
			// [t
			if ch == ']' {
				// [t]
				linkEnd = i-1
				istate =42
			}
		case 42:
			// [t]
			if ch == '(' {
				// [t](
				istate = 43
			}
		case 43:
			// [t](
			if utilLib.IsUriCh(ch) {
				// [t](uri
				uriSt = i
				istate = 44
			}
		case 44:
			// [t](uri
			if ch == ')' {
				// [t](uri)
				uriEnd = i-1
				istate = 4
			}

		case 50:
			// [^
			if utilLib.IsNumeric(ch) {
				// [^1
				istate = 51
				ftnNumSt = i
			}
			if utilLib.IsAlpha(ch) {
				// [^t
				istate = 59
				ftnStrSt = i
			} else {
				fmt.Printf("error istate 50: ch %q is non-numeric!\n")
				istate = 4
			}
		case 51:
			// [^1
			if ch == ']' {
				// [^1]
				ftnNumEnd = i-1
				istate = 52
			}

		case 52:
			// [^1]
			if ch == '(' {
				// [^1](
				istate = 53
			}
		case 53:
			// [^1](
			if utilLib.IsAlpha(ch) {
				// [^1](a
				ftnStrSt =  i+1
				istate = 54
			}
		case 54:
			// [^1](a
			if ch == ')' {
				// [^1](text)
				istate = 4
				ftnStrEnd = i-1
			}
		case 59:
			// [^a
			if ch == ']' {
				// [^a]
				ftnStrEnd = i-1
				istate = 4
			}
		default:
		}
		if ftnStrEnd != 0 {
			// ftnStr = string(buf[ftnStrSt:ftnStrEnd+1])
			ftnStrSt = 0
			ftnStrEnd = 0
			//ftn Number
			fmt.Printf("*** ftn num %d:%d str %d:%d\n", ftnNumSt, ftnNumEnd, ftnStrSt, ftnStrEnd)
		}
		// link
		if linkEnd != 0 {
			linkEnd = 0
			fmt.Printf("link str: %d:%d uri: %d:%d\n", linkSt, linkEnd, uriSt, uriEnd)
		}
	}
	return nil
}

func (mdP *mdParseObj) checkError(lin int, fch byte,  errStr string) {
// method that checks for an error in a line

	var el structEl
	var errEl errEl


	fmt.Printf("line %d fch %q errmsg: %s\n", lin, fch, errStr)
	linSt := mdP.linList[lin].linSt
	linEnd := mdP.linList[lin].linEnd

	buf := (*mdP.inBuf)

	errEl.txt = string(buf[linSt:linEnd])
	errEl.fch = fch
	errEl.errmsg = errStr
	errEl.line = lin
	el.errEl = &errEl
	mdP.elList = append(mdP.elList, el)

}

func (mdP *mdParseObj) checkHeadingEOL(lin int, parEl *parEl)(err error) {
// method that tests a heading line to see whether the heading is completed in that line
// or continues.

	linSt := parEl.txtSt
	linEnd := mdP.linList[lin].linEnd
//fmt.Printf("checkHeadingEOL line %d: buf %d:%d\n", lin, linSt, linEnd)
	linLen := linEnd - linSt
	if linLen < 2 {return nil}

	buf := (*mdP.inBuf)

//fmt.Printf("*** heading EOL: %q %q\n", buf[linEnd-2], buf[linEnd -1])
//	if (buf[linLen-2] == ' ') && (buf[linLen-1] == ' ') { return true}

	// check where the text line acutally ends
	istate := 0
	newLinEnd:= 0
	numWs := 0
	for i:=linEnd -1; i>linSt; i-- {
/*
if buf[i] == '\r' || buf[i] == '\n' {
	fmt.Printf("char pos: %d char CR \n", i)
} else {
	fmt.Printf("char pos: %d char %q %d \n", i, buf[i], buf[i])
}
*/
		switch istate {
		case 0:
			if utilLib.IsSentence(buf[i]) {
				newLinEnd = i
				istate = 3
				break
			}
			if buf[i] == '#' {
				istate = 1
				break
			}
			// white spaces
			if buf[i] == ' ' {
				istate = 0
				numWs++
				break
			}
			istate = 3

		case 1:
		//#
			if buf[i] == '#' {
				istate = 1
				break
			}
			if buf[i] == ' ' {
				istate = 2
				break
			}
			if utilLib.IsSentence(buf[i]) {
				newLinEnd = i
				break
			}
			istate = 3

		case 2:
			if buf[i] == ' ' {
				istate = 2
				break
			}
			if utilLib.IsSentence(buf[i]) {
				newLinEnd = i
				break
			}
			istate = 3
		default:
			istate = 3
		}

		if newLinEnd > 0 {
			parEl.txtEnd = newLinEnd
			break
		}
		if istate == 3 {break}

	}

	if newLinEnd == 0 {return fmt.Errorf("checkHeadingEOL no lineEnd!")}

	parEl.fin = false
	if numWs > 2 {
		parEl.fin = true
	}
//	fmt.Printf(" header txt: %s %t\n", string((*mdP.inBuf)[parEl.txtSt:parEl.txtEnd+1]), parEl.fin)
	return nil
}

func (mdP *mdParseObj) checkParEnd(lin int)(err error) {
// method that checks terminates the previous paragraph after an empty line

	lastEl := len(mdP.elList) -1
	if lastEl < 0 {return fmt.Errorf("no elList")}

	el := mdP.elList[lastEl]
	if el.parEl != nil {
		if el.parEl.fin == false {
			el.parEl.fin = true
		}
	}
	return nil
}

func (mdP *mdParseObj) checkParEOL(lin int, parel *parEl)(res bool) {
// function to test whether a line is completed with a md cr

	linSt := mdP.linList[lin].linSt
	linEnd := mdP.linList[lin].linEnd

	buf := (*mdP.inBuf)

	if (linEnd - linSt) < 2 {return false}

//fmt.Printf("EOL: %q %q\n", buf[linEnd-1], buf[linEnd])
	if (buf[linEnd] == ' ') && (buf[linEnd-1] == ' ') { 
		parel.txtEnd = linEnd
		parel.fin = false
		return true
	}

	parel.txtEnd = linEnd + 1
	parel.fin = true
	return false
}

func (mdP *mdParseObj) checkWs(lin int)(fch byte, err error) {
// method that checks indented lines to find the first character

	linSt := mdP.linList[lin].linSt
	linEnd := mdP.linList[lin].linEnd

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
fmt.Printf("fch: %q numWs: %d\n", fch, numWs)

	return fch, nil
}

func (mdP *mdParseObj) checkComment(lin int)(err error) {
// method that checks whether line is a comment

	var el structEl
	var comEl comEl

	linSt := mdP.linList[lin].linSt
	linEnd := mdP.linList[lin].linEnd
	linLen := linEnd - linSt
	if linLen < 10 {return fmt.Errorf("line length too little!")}
	buf := (*mdP.inBuf)
	if buf[linSt + 1] != '/' {return fmt.Errorf("invalid char 1!")}
	if buf[linSt + 2] != '/' {return fmt.Errorf("invalid char 2!")}
	if buf[linSt + 3] != ']' {return fmt.Errorf("invalid char 3!")}
	if buf[linSt + 4] != ':' {return fmt.Errorf("invalid char 4!")}
	if buf[linSt + 5] != ' ' {return fmt.Errorf("invalid char 5!")}
	if buf[linSt + 6] != '*' {return fmt.Errorf("invalid char 6!")}
	if buf[linSt + 7] != ' ' {return fmt.Errorf("invalid char 7!")}
	if buf[linSt + 8] != '(' {return fmt.Errorf("lin %d comment: no '(' found!", lin)}

	closPar :=0
	for i:= linEnd; i> linSt+8; i-- {
		if buf[linSt + 8] != ')' {
			closPar = i
			break
		}
	}
	if closPar == 0 {return fmt.Errorf("lin %d no ')' found!", lin)}

	comEl.txt = string(buf[linSt+9:closPar-1])
	el.comEl = &comEl
	mdP.elList = append(mdP.elList, el)

	return nil
}

func (mdP *mdParseObj) checkPar(lin int)(err error) {
// method that parses a line  to check whether it is paragraph

	var el structEl
	var parEl parEl
	linSt := mdP.linList[lin].linSt
//	linEnd := mdP.linList[lin].linEnd
	buf := (*mdP.inBuf)

	mdP.checkParEOL(lin, &parEl)

	parEl.txtSt = linSt
	parEl.txt = string(buf[parEl.txtSt:parEl.txtEnd-1])
	el.parEl = &parEl
	mdP.istate = PAR

	// see whether the previous element of elList is a parEl
	last := len(mdP.elList) -1
	lastEl := mdP.elList[last]
	// prev element: par
	if lastEl.parEl != nil {
		// lastEl is a parEl
		// we tack the txtstring onto the parEl
		lastParEl := lastEl.parEl
		if !lastParEl.fin {
fmt.Printf(" new el par txt: %s ", parEl.txt)
			mdP.elList = append(mdP.elList, el)
			return nil
		}
		lastParEl.txt += " " + parEl.txt
		lastParEl.txtEnd = parEl.txtEnd
fmt.Printf(" ex el par txt: %s ", parEl.txt)
	}

	// prev element: blk
	if lastEl.bkEl != nil {
		lastParEl := lastEl.bkEl.parEl
		if !lastParEl.fin {
fmt.Printf(" new blk el txt: %s ", parEl.txt)
			mdP.elList = append(mdP.elList, el)
			return nil
		}
		lastParEl.txt += " " + parEl.txt
		lastParEl.txtEnd = parEl.txtEnd
fmt.Printf(" ex blk el txt: %s ", parEl.txt)
	}

	// prev elment: code

	//previous element: empty line
	if lastEl.emEl {
fmt.Printf(" new el par txt: %s ", parEl.txt)
		mdP.elList = append(mdP.elList, el)
		return nil
	}

	return fmt.Errorf("elList at %d has no valid el")
}


func (mdP *mdParseObj) checkHeading(lin int) (err error){
// method that parses a line for headings

	var el structEl
	var parEl parEl

//	listEl := mdP.elList[el]
	linSt := mdP.linList[lin].linSt
	linEnd := mdP.linList[lin].linEnd
	buf := (*mdP.inBuf)
//fmt.Printf("heading buffer: %s\n", buf[linSt:linEnd])

	hd := 0
	parSt := 0
	istate := 0
	for i:=linSt; i < linEnd; i++ {
		switch istate {
		case 0:
			if buf[i] == '#' {
				hd++
				break
			}
			if buf[i] == ' ' {
				istate = 1
				break
			}
			return fmt.Errorf("lin %d istate: %d char %q \n", lin, istate, buf[i])
		case 1:
			if buf[i] == ' ' {
				istate = 1
				break
			} else {
				parSt = i
				istate = 2
				break
			}

		default:
		}
		if istate > 1 {break}
	}

	// heading level
	if !(hd>0) {return fmt.Errorf(" heading: h0 not valid")}
	// no heading text start
	if parSt <1 {return fmt.Errorf(" heading txt start not found!")}

//	hdStr := fmt.Sprintf("h%d", hd)
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

//	txtstr := string(buf[parSt:linEnd])
	parEl.txtSt = parSt
	parEl.typ = hdtyp

	// need to check  for heading endings
	// adjusts linEnd
	err = mdP.checkHeadingEOL(lin, &parEl)
	if err != nil { return fmt.Errorf("checkHeadingEOL: %v", err)}
	parEl.txt = string(buf[parEl.txtSt:parEl.txtEnd+1])

	fmt.Printf(" heading: %s text: \"%s\"", dispHtmlEl(hdtyp), parEl.txt)

//fmt.Printf("linSt: %d linEnd: %d\n", linSt, linEnd)

	el.parEl = &parEl
	mdP.elList = append(mdP.elList, el)
	mdP.istate = PAR
	return nil
}


func (mdP *mdParseObj) checkHr(lin int) (err error) {
// method that parses a horizontal ruler line

	var el structEl

	linSt := mdP.linList[lin].linSt
	linEnd := mdP.linList[lin].linEnd
	buf := (*mdP.inBuf)
	numCh:=0
	ch := buf[linSt]
	for i:= linSt+1; i< linEnd; i++ {
		if buf[i] == ch {
			numCh++
		} else {
			break
		}
	}
//fmt.Printf("HR %q numCh: %d ", ch, numCh)

	if numCh >2 {
		el.hrEl = true
		mdP.elList = append(mdP.elList, el)
		mdP.istate = HR
		err = nil
	} else {
		err = fmt.Errorf("too insufficient chars %q", ch)
	}
	return err
}

func (mdP *mdParseObj) checkBold(){

}

func (mdP *mdParseObj) checkItalics() {

}

func (mdP *mdParseObj) checkUnList(lin int) (err error){
// method that parsese and unordered list item

	var el structEl
	var ulEl uList
	var parEl parEl

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

	if parSt == 0 {return fmt.Errorf("ul no text found!")}

	// nest lev
	nestLev := wsNum/4

	parEl.txtSt = parSt

	mdP.checkParEOL(lin, &parEl)
	parEl.txt = string(buf[parEl.txtSt:parEl.txtEnd-1])

fmt.Printf(" UL txt: %s ", parEl.txt)

	ulEl.nest = nestLev
	ulEl.parEl = &parEl
	el.ulEl = &ulEl
	mdP.elList = append(mdP.elList, el)
	mdP.istate = PAR

	return nil
}

func (mdP *mdParseObj) checkHtml() {

}

func (mdP *mdParseObj) checkIndent(lin int) (ch byte, numWs int, err error){
// method that chck indents and returns the nesting level and list type

	linSt := mdP.linList[lin].linSt
	linEnd := mdP.linList[lin].linEnd
	buf := (*mdP.inBuf)

	wsCount := 0
	for i:=linSt; i< linEnd; i++ {
		if buf[i] == ' ' {
			wsCount++
		} else {
			ch = buf[i]
			break
		}
	}

	if ch == 0 {return 0, numWs, fmt.Errorf("lin %d: checkIndent: no char found!")}

	if wsCount < 4 { return ch, numWs, fmt.Errorf("lin %d: checkIndent: insufficient ws!")}

/*
		case '*','-','+':
			if wsCount < 4 {return 0,false, fmt.Errorf("insufficient ws!")}
			return wsCount/4, false, nil

		case '1','2','3','4','5','6','7','8','9','0':
			if wsCount < 4 {return wsCount/4, true, fmt.Errorf("insufficient ws!")}
			return wsCount/4, true, nil
*/

	return ch, wsCount, nil
}


func (mdP *mdParseObj) checkBlock(lin int) (err error){
// A method that parses a blockquote

	var el structEl
	var bkEl bkEl
	var parEl parEl

	linSt := mdP.linList[lin].linSt
	linEnd := mdP.linList[lin].linEnd
	buf := (*mdP.inBuf)

	nest := 0
	istate :=0
	parSt :=0
	wsCount := 0
	for i:=linSt; i< linEnd; i++ {
		ch := buf[i]
		switch istate {
		case 0:
			if ch == ' ' {
				istate = 10
				break
			}
			if ch == '>' {
				nest++
				istate = 1
				break
			}
		// else error
		case 1:
			// >
			if ch == ' ' {
				// >ws
			}
			if ch == '>' {
				// >>
				nest++
				break
			}
			if utilLib.IsAlpha(ch) {
				// >a
				istate = 2
				parSt = i
				break
			}

		case 2:
			// >a, >>a, > a, >  a
			if ch == '>' {
				nest++
				istate = 0
				break
			}
			if utilLib.IsAlpha(ch) {
				istate = 98
				parSt = i
			}
		case 10:
			// ' '
			if ch == ' ' {
				wsCount++
				break
			}
			if utilLib.IsAlpha(ch) {
				istate = 99
				parSt = i
			}

		default:
			break
		}
	}

	if parSt == 0 {return fmt.Errorf("no text string found!")}

	parEl.txt = string(buf[parSt:linEnd])
	parEl.txtSt = parSt
	parEl.txtEnd = linEnd
	bkEl.nest = nest
	bkEl.parEl = &parEl
	el.bkEl = &bkEl
	mdP.elList = append(mdP.elList, el)
	mdP.istate = BLK

//bbb
fmt.Printf(" nest: % d par txt: %s ", nest, parEl.txt)

	return nil
}

func (mdP *mdParseObj) checkCode(lin int) (err error){
// A method that parses a code block

	var el structEl
	var bkEl bkEl
	var parEl parEl

	linSt := mdP.linList[lin].linSt
	linEnd := mdP.linList[lin].linEnd
	buf := (*mdP.inBuf)

	nest := 0
	istate :=0
	parSt :=0
	for i:=linSt; i< linEnd; i++ {
		ch := buf[i]
		switch istate {
		case 0:
			if ch == ' ' {istate = 1}
			if ch == '>' {
				nest++
				istate = 0
			}

		case 1:
			if ch == '>' {
				nest++
				istate = 0
			}
			if utilLib.IsAlpha(ch) {
				istate = 2
				parSt = i
			}

		default:
			break
		}
	}

	if parSt == 0 {return fmt.Errorf("no text string found!")}

	parEl.txt = string(buf[parSt:linEnd])
	parEl.txtSt = parSt
	parEl.txtEnd = linEnd
	bkEl.nest = nest
	bkEl.parEl = &parEl
	el.bkEl = &bkEl
	mdP.elList = append(mdP.elList, el)
	mdP.istate = BLK

	return nil
}

func (mdP *mdParseObj) checkStrike() {

}

func (mdP *mdParseObj) checkImage(lin int) (err error){
// a method that parse an image

	var el structEl
	var imgEl imgEl

//	listEl := mdP.elList[el]
	linSt := mdP.linList[lin].linSt
	linEnd := mdP.linList[lin].linEnd
	buf := (*mdP.inBuf)

	altSt := 0
	altEnd := 0
	srcSt := 0
	srcEnd := 0

	istate := 0
	imgEnd :=false
	for i:=linSt; i < linEnd; i++ {
		switch istate {
		case 0:
			if buf[i] == ']' {
				imgEl.alt = string(buf[linSt+2:i])
				istate = 1
			}

		case 1:
			if buf[i] == '(' {
				srcSt = i
				istate = 2
			}

		case 2:
			if buf[i] == ' ' {
				altSt = i
				srcEnd = i-1
				istate = 3
			}

			if buf[i] == ')' {
				srcEnd = i-1
				imgEnd = true
				istate = 5
			}

		case 3:
			if buf[i] == '"' {
				istate = 4
				altSt = i+1
			}

		case 4:
			if buf[i] == '"' {
				istate = 5
				altEnd = i-1
			}

		case 5:
			if buf[i] == ')' {
				imgEnd = true
			}
//				return fmt.Errorf("no starting double par!")
		default:

		}
		if imgEnd {break}

	}

	if altEnd - altSt < 3 {return fmt.Errorf("no viable altText!")
	} else { imgEl.alt = string(buf[altSt+1: altEnd -1]) }

	if srcEnd - srcSt < 3 {return fmt.Errorf("no viable img uri!")
	} else { imgEl.src = string(buf[srcSt+1: srcEnd -1]) }

	if istate == 0 {return fmt.Errorf("could not parse img el! no ']'")}
	if !imgEnd {return fmt.Errorf("could not parse img el!")}

	el.imgEl = &imgEl
	mdP.elList = append(mdP.elList, el)
	mdP.istate = IMG

	return nil
}

func (mdP *mdParseObj) checkLink() {

}

func (mdP *mdParseObj) checkTable(lin int)(endLin int, err error) {
// method that parses a table. The method returns the last line of the table.
// A table consists of several lines, unlike other elements

	var el structEl
	var tblEl tblEl

	linSt := mdP.linList[lin].linSt
	linEnd := mdP.linList[lin].linEnd
	buf := (*mdP.inBuf)

	// determine the number of cols
	col := 0
	for i:=linSt+1; i<linEnd; i++ {
		if buf[i] == '|' {
			col++
		}
	}
	if col == 0 {return lin, fmt.Errorf("no columns in table!")}

	// table rows
	row:=0

	for ilin := lin+1; ilin < lin + 10; ilin++ {
		tblLinSt := mdP.linList[ilin].linSt
//		tblLinEnd:= mdP.linList[ilin].linEnd
		if buf[tblLinSt] != '|' {
			endLin = ilin
			break
		}
		if buf[tblLinSt+1] == '-' {
			row++
		}
	}
	if row == 0 {return lin, fmt.Errorf("no rows in table!")}


	// table cells

	tblEl.rows = row
	tblEl.cols = col

	el.tblEl = &tblEl
	mdP.elList = append(mdP.elList, el)

	return endLin, nil
}

func (mdP *mdParseObj) checkOrList(lin int)(err error) {
// method that parses an ordered list

	var el structEl
	var orEl oList
	var parEl parEl

	linSt := mdP.linList[lin].linSt
	linEnd := mdP.linList[lin].linEnd
	buf := (*mdP.inBuf)

	parSt:=0
	markEnd := 0
	wsNum := 0
	for i:= linSt; i< linSt +4; i++ {
		ch:= buf[i]
		if ch == '.' {
			markEnd = i
			break
		}
		if !utilLib.IsNumeric(ch) {return fmt.Errorf(" nonumeric char in counter")}
	}

	if markEnd == 0 {return fmt.Errorf(" orList no period after counter!")}

	// check for par start
	for i:= markEnd+1; i< linEnd; i++ {
		ch:= buf[i]
		if utilLib.IsAlpha(ch) {
			parSt = i
			break
		}
		if ch == ' ' {
			wsNum++
			continue
		}
		return fmt.Errorf(" unacceptable char in pos %d before par start!", i)
	}

	if parSt == 0 {return fmt.Errorf("ol no text found!")}

	// nest lev
	nestLev := wsNum/4

	parEl.txtSt = parSt

	mdP.checkParEOL(lin, &parEl)
	parEl.txt = string(buf[parEl.txtSt:parEl.txtEnd-1])

fmt.Printf(" OL txt: %s ", parEl.txt)

	orEl.nest = nestLev
	orEl.parEl = &parEl

	el.olEl = &orEl
	mdP.elList = append(mdP.elList, el)
	mdP.istate = PAR

	return nil
}

func (mdP *mdParseObj) checkBR()(err error) {
// method that parses an empty line

	var el structEl

	el.emEl = true
	mdP.elList = append(mdP.elList, el)
	mdP.istate = EL
	return nil
}


func (mdP *mdParseObj) printLinList()() {

	fmt.Printf("line Start End Fch text\n")
	for el:=0; el<len(mdP.linList); el++ {
		fmt.Printf("line %3d: %3d %3d ", el, mdP.linList[el].linSt, mdP.linList[el].linEnd)
		str:=string((*mdP.inBuf)[mdP.linList[el].linSt:mdP.linList[el].linEnd])
		fmt.Printf("%q:%s\n", (*mdP.inBuf)[mdP.linList[el].linSt], str)
	}

}

func (mdP *mdParseObj) cvtMdToHtml()(err error) {
// method that converts the parsed element list of an md file inot an html file
	fmt.Printf("*** input file: %s\n", mdP.filnam + ".md")

	outfil, err := os.Create(mdP.filnam + ".html")
	if err != nil { return fmt.Errorf("os.Create: %v\n", err)}
	defer outfil.Close()

	return nil
}

func (mdP *mdParseObj) printElList () {
// method that prints out the structural element list

	fmt.Println("*********** El List ***********")
	fmt.Printf("Elements: %d\n", len(mdP.elList))
	fmt.Printf("  el nam typ  subels fin txt\n")
	for i:=0; i < len(mdP.elList); i++ {
		el := mdP.elList[i]
		fmt.Printf("el %3d: ", i)
		if el.emEl {
			fmt.Printf("eL: %t\n", el.emEl)
			continue
		}

		if el.hrEl {
			fmt.Printf("HR: %t\n", el.hrEl)
			continue
		}


		if el.parEl != nil {
			ParEl := *el.parEl
			fmt.Printf( "par %-5s: text: %s status: %t", dispHtmlEl(ParEl.typ), ParEl.txt, ParEl.fin)
/*
			subLen := len(ParEl.subEl)
			if subLen == 1 {
				fmt.Printf(" subel 0: \"%s\"\n", ParEl.subEl[0].txt)
			} else {
				fmt.Printf("\n")
				for i:=0; i< subLen; i++ {
					fmt.Printf("         subel %d: %s\n", i, ParEl.subEl[i].txt)
				}
			}
*/
			fmt.Printf("\n")
			continue
		}
		if el.ulEl != nil {
			fmt.Printf( "ul nest %d: ", el.ulEl.nest)
			if el.ulEl.parEl != nil {
				ParEl := el.ulEl.parEl
				fmt.Printf( "  par typ: %-5s text: %s stat: %t\n", dispHtmlEl(ParEl.typ), ParEl.txt, ParEl.fin)
/*
				subLen := len(ParEl.subEl)
				if subLen == 1 {
					fmt.Printf("         subel 0: %s\n", ParEl.subEl[0].txt)
				} else {
					for i:=0; i< subLen; i++ {
						fmt.Printf("         subel %d: %s\n", i, ParEl.subEl[i].txt)
					}
				}
*/
			} else {
				fmt.Printf( "ulEl par nil!\n")
			}

			continue
		}

		if el.olEl != nil {
			fmt.Printf( "ol nest %d ", el.olEl.nest)
			if el.olEl.parEl != nil {
				ParEl := el.olEl.parEl
				fmt.Printf( "par typ: %-5s subels: %d stat %t\n", dispHtmlEl(ParEl.typ), len(ParEl.subEl), ParEl.fin)
/*
				subLen := len(ParEl.subEl)
				if subLen == 1 {
					fmt.Printf("         subel 0: %s\n", ParEl.subEl[0].txt)
				} else {
					for i:=0; i< subLen; i++ {
						fmt.Printf("         subel %d: %s\n", i, ParEl.subEl[i].txt)
					}
				}
*/
			} else {
				fmt.Printf( "olEl par nil!\n")
			}
			continue
		}

		if el.tblEl != nil {
			fmt.Printf( "tbl: rows: %d  cols: %d \n", el.tblEl.rows, el.tblEl.cols)
			continue
		}

		if el.imgEl != nil {
			fmt.Printf( "img h: %d w: %d src: %s \n", el.imgEl.height, el.imgEl.width, el.imgEl.src)
			continue
		}

		if el.comEl != nil {
			fmt.Printf( "com: %s\n", el.comEl.txt)
			continue
		}

		if el.bkEl != nil {
			fmt.Printf( "blk lev: %d text: %s\n", el.bkEl.nest, el.bkEl.parEl.txt)
			continue
		}

		if el.errEl!= nil {
			fmt.Printf("error line %d: fch %q %d:: error %s\n", el.errEl.line, el.errEl.fch, el.errEl.fch, el.errEl.errmsg)
			continue
		}

	}
}
