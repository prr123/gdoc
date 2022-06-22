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
	"google/gdoc/htmlLib"
)

type mdParseObj struct {
	istate int
	errCount int
	cnest int
	imgCount int
	ftnCount int
	cliTyp bool
	cliCount [10]int
	filnam string
	inBuf *[]byte
	linList []mdLin
	elList []structEl
	errList [10]errObj
}

type errObj struct {
	lin int
	msg string
}

type mdLin struct {
	linSt int
	linEnd int
	typ byte
	fchar byte
}

type structEl struct {
	elTyp int
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
	nest int
	txtSt int
	txtEnd int
	txt string
	subEl []parSubEl
}

type parSubEl struct {
	elSt int
	elEnd int
	bold bool
	italic bool
	sup bool
	sub bool
	link bool
	strike bool
	ftn	int
	img int
	txtTyp int
	txt string
	lkUri string
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
	count [6]int
	parEl *parEl
}

// line attributes
const (
	UK = iota 	// unknown
	EP			// empty line after par
	EUL			// empty line after ulist
	EOL			// empty line after olist
	EB			// empty line after block
	HR			// horizintal ruler
	COM			// comment
	PAR			// paragraph
	UL			// unordered list
	OL			// ordered list
	IMG			// image
	BLK			// block
	COD			// code
	TBL			// table
	ERR			// error ???
)

//html elements
const (
	br = iota
	par
	hr
	sp
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
	cm
	cod
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
	link
)

func dispTxtAtt(num int)(str string) {

	var txtAtt [8]string
	txtAtt[bold] = "bold"
	txtAtt[italic] = "italic"
	txtAtt[strike] = "strike"
	txtAtt[html] = "html"
	txtAtt[ftnote] = "ftnote"
	txtAtt[sup] = "sup"
	txtAtt[sub] = "sub"
	txtAtt[link] = "link"

	if num > len(txtAtt)-1 {return ""}
	return txtAtt[num]
}

func dispState(num int)(str string) {
// function that converts state constants to strings

	var stateDisp [16]string

	stateDisp[UK] = "UK"   //unknown
	stateDisp[EP] = "EP"
	stateDisp[EUL] = "EUL"
	stateDisp[EOL] = "EOL"
	stateDisp[EB] = "EB"
	stateDisp[HR] = "HR"
	stateDisp[COM] = "COM"
	stateDisp[PAR] = "PAR"
	stateDisp[UL] = "UL"
	stateDisp[OL] = "OL"
	stateDisp[IMG] = "IMG"
	stateDisp[BLK] = "BLK"
	stateDisp[COD] = "COD"
	stateDisp[TBL] = "TBL"
	stateDisp[ERR] = "ERR"

	if num > len(stateDisp)-1 {return ""}
	return stateDisp[num]
}

func dispHtmlEl(num int)(str string) {
// function that converts html const to strings

	var htmlDisp [18]string

	htmlDisp[br] = "br"
	htmlDisp[par] = "p"
	htmlDisp[hr] = "hr"
	htmlDisp[sp] = "span"
	htmlDisp[ul] = "ul"
	htmlDisp[ol] = "ol"
	htmlDisp[li] = "li"
	htmlDisp[img] = "img"
	htmlDisp[ftn] = "ftn"
	htmlDisp[h1] = "h1"
	htmlDisp[h2] = "h2"
	htmlDisp[h3] = "h3"
	htmlDisp[h4] = "h4"
	htmlDisp[h5] = "h5"
	htmlDisp[h6] = "h6"
	htmlDisp[cm] = "cm"
	htmlDisp[cod] = "cod"

	if num > len(htmlDisp)-1 {return ""}
	return htmlDisp[num]
}

func dispElTyp(el structEl) {

	fmt.Printf("el Type: ")
	if el.emEl {fmt.Printf("empty line EL!\n")}
	if el.hrEl {fmt.Printf("hor ruler HR!\n")}
	if el.comEl != nil {fmt.Printf("comment CM!\n")}
	if el.parEl != nil {fmt.Printf("paragraph PAR!\n")}
	if el.tblEl != nil {fmt.Printf("table TBL!\n")}
	if el.imgEl != nil {fmt.Printf("image IMG!\n")}
	if el.ulEl != nil {fmt.Printf("unord List UL\n")}
	if el.olEl != nil {fmt.Printf("ord List OL!\n")}
	if el.bkEl != nil {fmt.Printf("block BLK!\n")}
	if el.errEl != nil {fmt.Printf("error ERR!\n")}

}

func getElTyp(el structEl) (outstr string){

	outstr = "el Type: "
	if el.emEl {outstr += "empty line EL! "}
	if el.hrEl {outstr += "hor ruler HR! "}
	if el.comEl != nil {outstr += "comment CM! "}
	if el.parEl != nil {outstr += "paragraph PAR! "}
	if el.tblEl != nil {outstr += "table TBL! "}
	if el.imgEl != nil {outstr += "image IMG! "}
	if el.ulEl != nil {outstr += "unord List UL "}
	if el.olEl != nil {outstr += "ord List OL! "}
	if el.bkEl != nil {outstr += "block BLK! "}
	if el.errEl != nil {outstr += "error ERR! "}
	return outstr
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

	inpfil, err := os.Open(inpfilnam)
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

//	mdP.istate = EP

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
					mdP.istate = ERR
					break
				}
				switch mdP.istate {
				case PAR:
					err = mdP.checkParEnd(lin)
					if err != nil {fmt.Printf("line %d ParEnd %v\n", lin, err); mdP.istate = PAR;}
					err = mdP.checkBR()
					mdP.istate = EP
					if err != nil {fmt.Printf("line %d checkBR %v\n", lin, err)}

				case BLK:
					err = mdP.checkBR()
					mdP.istate = EB
					if err != nil {fmt.Printf("line %d checkBR %v\n", lin, err)}

				case EP, HR:
					err = mdP.checkBR()
					mdP.istate = EP
					if err != nil {fmt.Printf("line %d checkBR %v\n", lin, err)}

				case UL:
					err = mdP.checkBR()
					mdP.istate = EUL
					if err != nil {fmt.Printf("line %d checkBR %v\n", lin, err)}

				case OL:
					err = mdP.closeOrList(lin)
					if err != nil {fmt.Printf("line %d closeOL %v\n", lin, err)}
					err = mdP.checkBR()
					mdP.istate = EOL
					if err != nil {fmt.Printf("line %d checkBR %v\n", lin, err)}

				default:
					fmt.Printf("line %d unknown istate %s cr", lin, dispState(mdP.istate))
				}

			case '#':
				// headings
				err = mdP.checkHeading(lin)
				mdP.istate = PAR
				if err != nil {fmt.Printf("line %d: heading error %v\n", lin, err)}

			case '_':
				// horizontal ruler ___wsp/ret
				if linLen < 4 {
					fmt.Printf("line %d: linLen %d too short for HR!", lin, linLen)
					break
				}
				err = mdP.checkHr(lin)
				mdP.istate = HR
				if err != nil {
					err1 := mdP.checkPar(lin)
					mdP.istate = PAR
					if err1 != nil { fmt.Printf("line %d: neither HR %v nor Par %v\n", err, err1); mdP.istate = ERR}
				}
			case '*','+','-':
				// check whether next char is whitespace
				switch {
				case sch == ' ':
					// unordered list *wsp
		fmt.Printf("un List NL:0 ")
					err = mdP.checkUnList(lin)
					if err != nil {fmt.Printf("line %d: unordered list error %v\n", lin, err)}
					mdP.istate = UL

				case sch == fch:
					err = mdP.checkHr(lin)
					if err != nil {
						err1 := mdP.checkPar(lin)
						if err1 != nil {
							fmt.Printf("line %d: neither HR %v nor Par %v\n", err, err1)
							mdP.istate = ERR
						}
					} else {mdP.istate = HR}

				case utilLib.IsAlpha(sch):
//					fmt.Println(string((*mdP.inBuf)[linSt:linEnd]))
					err = mdP.checkPar(lin)
					mdP.istate = PAR
					if err != nil {fmt.Printf("line %d: par error %v\n", lin, err)}

				default:
				// error
					fmt.Printf("line %d unsuitable char %q %d after *+-\n", sch, sch)
					mdP.istate = ERR
				}

			case '>':
				// block quotes
				err = mdP.checkBlock(lin, 0)
				mdP.istate = BLK
				if err != nil {fmt.Printf("line %d: quote block error %v\n", lin, err)}

			case '`':
				// block quotes
				endLin, err := mdP.checkCode(lin)
				mdP.istate = COD
				if err != nil {fmt.Printf("line %d: code block error %v\n", lin, err)}
				lin = endLin + 1

			case ' ':
				// ws
				// nested list
				ch, wsNum, err := mdP.checkIndent(lin)
				if err != nil {
					mdP.addErr(lin,  fmt.Sprintf("indent error %v", err))
					fmt.Printf("line %d: indent error %v\n", lin, err)
					break
				}

				if ch == 0 {
					mdP.addErr(lin, "indent error no char")
					mdP.istate = EP
					break
				}

				if ch == '>' {
					// blockquote within list
					err = mdP.checkBlock(lin, wsNum)
					break
				}
				if utilLib.IsAlpha(ch) {
					switch mdP.istate {
					case OL, EOL:
						err = mdP.checkOrList(lin)
						if err != nil {fmt.Printf("line %d: OL error %v\n", lin, err)}
						mdP.istate = OL

					case UL, EUL:
						err = mdP.checkUnList(lin)
						if err != nil {fmt.Printf("line %d: UL error %v\n", lin, err)}
						mdP.istate = UL

					case EB, BLK, PAR, EP:
						err = mdP.checkBlock(lin, wsNum)
						mdP.istate = BLK
						if err != nil {fmt.Printf("line %d: block error %v\n", lin, err)}

					default:
						mdP.addErr(lin,fmt.Sprintf("indent istate: %s !", dispState(mdP.istate)))

					}
					break
				}
				if (ch == '*' || ch == '+') || ch == '-' {
					err := mdP.checkUnList(lin)
					mdP.istate = UL
					if err != nil {fmt.Printf("line %d: unordered list error %v\n", lin, err)}
		fmt.Printf(" unList Nl: %d", wsNum/4)
					break
				}

				if utilLib.IsNumeric(ch) {
					err = mdP.checkOrList(lin)
					mdP.istate = OL
					if err != nil {fmt.Printf("line %d: ordered list error %v\n", lin, err)}
		fmt.Printf(" orList Nl: %d", wsNum/4)
					break
				}
				fmt.Printf("line %d: indent error start char %q in pos: %d\n", lin, ch, wsNum +1)

			case '!':
				// image
				err = mdP.checkImage(lin)
				mdP.istate = IMG
				if err != nil { return fmt.Errorf("lin %d image: %v\n",lin, err)}

			case '[':
				// check comment
				err = mdP.checkComment(lin)
				mdP.istate = COM
				if err != nil {
					err1 := mdP.checkPar(lin)
					if err1 != nil {fmt.Printf("lin %d comment: %v\n        par: %v\n",lin, err, err1) }
				}

			case '|':
				// table |text|
				endLin, err := mdP.checkTable(lin)
				mdP.istate = TBL
				if err != nil {fmt.Printf("line %d: table error %v\n", lin, err)}
				lin = endLin + 1

			default:
				if utilLib.IsNumeric(fch) {
				// ordered list 1.
					err = mdP.checkOrList(lin)
					mdP.istate = OL
					if err != nil {fmt.Printf("line %d: ordered list error %v\n", lin, err)}
					break
				}
				if utilLib.IsAlpha(fch) {
					// paragraph, block continuation
fmt.Printf("alpha: el state: %s ", dispState(mdP.istate))

					switch mdP.istate {
					case UL:
						err = mdP.checkUnList(lin)
						if err != nil {fmt.Printf("line %d: UL par err: %v\n", lin, err)}
						mdP.istate = UL
					case OL:
						err = mdP.checkOrList(lin)
						if err != nil {fmt.Printf("line %d: OL par err: %v\n", lin, err)}
						mdP.istate = OL

					case EP,EOL, EUL, PAR:
						err = mdP.checkPar(lin)
						mdP.istate = PAR
						if err != nil {fmt.Printf("line %d: par error %v\n", lin, err)}

					default:
						fmt.Printf("line %d: alpha state %s\n", lin, dispState(mdP.istate))
					}
					break
				}
				fmt.Printf("line %d: no fit for first char: %q\n", lin, fch)
				mdP.istate = ERR
//				mdP.checkError(lin, fch, errmsg)
		}
	fmt.Printf(" state: %s\n",dispState(mdP.istate))
	}

	mdP.printElList()
//	mdP.printErrList()
	fmt.Println("/n******** parse sub el ***********")
	err = mdP.parseMdSubEl()
	if err != nil { return fmt.Errorf("error parseMdSubEl: %v", err)}

	mdP.printElList()

	return nil
}

//ppp
func (mdP *mdParseObj) parseMdSubEl()(err error) {
// method that parses the text fields

	for i:=0; i < len(mdP.elList); i++ {
		el := mdP.elList[i]
		fmt.Printf("el %3d: ", i)

		// parse par if el.parEl != nil
		if el.parEl != nil {
			err = mdP.parsePar(el.parEl)
		}

		// parse par el
		if el.bkEl != nil {

		}

		//parse blkEl
		if el.ulEl != nil {
			err = mdP.parsePar(el.ulEl.parEl)
		}

		//parse blkEl
		if el.olEl != nil {
			err = mdP.parsePar(el.olEl.parEl)
		}

	}

	return nil
}



//ppar
func (mdP *mdParseObj) parsePar(parEl *parEl)(err error) {
// method which parses the text of of parEl

	var subEl parSubEl

	txtbuf := []byte(parEl.txt)

fmt.Printf("parsePar %d: %s\n", len(txtbuf), parEl.txt)

	linkSt := 0
	linkEnd := 0
	uriSt := 0
	uriEnd := 0

	ftnNumSt:=0
	ftnNumEnd:=0
	ftnStrSt:=0
	ftnStrEnd :=0

	txtSt := 0
	txtEnd := 0
	istate := 0

	last := len(txtbuf) -1

	for i:=0; i< len(txtbuf); i++ {
		ch := txtbuf[i]
fmt.Printf("i %d: %q istate: %d %s\n", i, ch, istate, string(txtbuf[txtSt:i+1]))

		switch istate {
		case 0:
			txtSt = i
			switch ch {

			case '*':
				// *
				istate = 3

			case '_' :
				// _
				istate = 11

			case '[' :
				// [
				istate = 40

			default:
				if utilLib.IsAlpha(ch) {
				// t
					istate = 1
				}
			}

		case 1:
			// t
			switch ch {
			case ' ':
				// Tws
				istate = 2

			default:
				if i == last {
					txtEnd = i
					subEl.txt = string(txtbuf[txtSt:txtEnd+1]) + "\n"
					parEl.subEl = append(parEl.subEl, subEl)
				}
			}

		case 2:
			// tws
			switch ch {
			case ' ':

			case '*':
				// tws*
				txtEnd = i-1
				subEl.txt = string(txtbuf[txtSt:txtEnd])
				parEl.subEl = append(parEl.subEl, subEl)

				istate = 3

			case '_':
				// tws_

				txtEnd = i-1
				subEl.txt = string(txtbuf[txtSt:txtEnd])
				parEl.subEl = append(parEl.subEl, subEl)
				istate =11

			default:
				// wsT
				if utilLib.IsAlphaNumeric(ch) {

				}
			}
		case 3:
			// *
			if ch == '*' {
				// **
				istate = 20
				break
			}

			if utilLib.IsAlphaNumeric(ch) {
				// *t
				txtSt = i
				istate = 4
			}

		case 4:
				// *t
			switch ch {
			case '*':
				// *t*
				istate = 5

			case ' ':
				// error

			case '_','[',']':
				// error
			}

		case 5:
			// *txt*
			switch ch {
			case ' ':
				// *txt*ws
				txtEnd = i-1
				subEl.italic = true
				subEl.txt = string(txtbuf[txtSt:txtEnd+1])
				istate = 1

			default:
				// *txt*t
			}

		case 9:
			// "/r/n"
			switch ch {
			case '\n':
				subEl.txt += "\n"
				parEl.subEl = append(parEl.subEl, subEl)
				istate = 0
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
				txtSt = i
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
			switch ch {
			case ' ':
				// _t_ws
				txtEnd := i-1
				subEl.italic = true
				subEl.txt = string(txtbuf[txtSt:txtEnd+1])
				istate = 0

			default:
				// _t_a

			}

		case 15:
			// __
			if utilLib.IsAlpha(ch) {
				// __t
				txtSt = 1
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
			switch ch {
			case ' ':
				// __t__ws
				txtEnd := i-1
				subEl.italic = true
				subEl.txt = string(txtbuf[txtSt:txtEnd+1])
				istate = 0

			default:
				// __t__a

			}

		case 20:
			// **
			switch ch {
			case '*':
				// ***
				istate = 30

			case ' ': 
				//error
				istate = 50

			default:
				if utilLib.IsAlphaNumeric(ch) {
					// **t
					istate = 21
					txtSt = i
				}
				// error
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
				txtSt = i
			}
			if ch == ' ' {
				//error ***ws
				istate = 50
			}

		case 31:
			// ***t
			switch ch {
			case '*':
				// ***t*
				istate = 32
			default:
			}

		case 32:
			// ***text*
			if ch == '*' {
				// ***text**
				istate = 33
			} else {
				// error
			}

		case 33:
			// ***text***
			if ch == '*' {
				// ***text***
				istate = 34
			}

		case 34:
			switch ch {
			case ' ':
				// __t__ws
				txtEnd := i-1
				subEl.italic = true
				subEl.bold = true
				subEl.txt = string(txtbuf[txtSt:txtEnd+1])
				istate = 0

			case '\n':
				txtEnd := i
				subEl.italic = true
				subEl.bold = true
				subEl.txt = string(txtbuf[txtSt:txtEnd+1])
				istate = 0

			default:
				// __t__a
				// error
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
				subEl.link = true
				subEl.txt = string(txtbuf[linkSt:linkEnd +1])
fmt.Printf("link txt: %s\n",subEl.txt)
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
				subEl.lkUri = string(txtbuf[uriSt:uriEnd +1])
				parEl.subEl = append(parEl.subEl, subEl)
//ppsub
fmt.Printf("link uri: %s\n",subEl.lkUri)
fmt.Printf(" **** subel %d bold %t italic %t link: %s txt: %s\n", i, subEl.bold, subEl.italic, subEl.lkUri, subEl.txt)

				istate = 0
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

	//debug
	sublen := len(parEl.subEl)
	fmt.Printf("parsePar sub %d\n", sublen)
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

//	if mdP.istate != EL {return fmt.Errorf("checkParEnd: current el not EL!")}
	lastEl := len(mdP.elList) -1
	if lastEl < 0 {return fmt.Errorf("no elList")}

	el := mdP.elList[lastEl]
	if el.parEl != nil {
		el.parEl.fin = true
	}
	return nil
}

func (mdP *mdParseObj) checkParEOL(lin int, parel *parEl)(err error) {
// function to test whether a line is completed with a md cr

	linSt := mdP.linList[lin].linSt
	linEnd := mdP.linList[lin].linEnd

	buf := (*mdP.inBuf)

	if (linEnd - linSt) < 2 {return fmt.Errorf("checkParEOL: empty line")}

//fmt.Printf("EOL: %q %q\n", buf[linEnd-1], buf[linEnd])
	if (buf[linEnd-2] == ' ') && (buf[linEnd-1] == ' ') {
		parel.txtEnd = linEnd -3
		parel.fin = true
		return nil
	}

	parel.txtEnd = linEnd -1
	parel.fin = false
	return nil
}

func (mdP *mdParseObj) checkPar(lin int)(err error) {
// method that parses a line  to check whether it is paragraph

	var el structEl
	var parEl parEl
	linSt := mdP.linList[lin].linSt
//	linEnd := mdP.linList[lin].linEnd
	buf := (*mdP.inBuf)

	err = mdP.checkParEOL(lin, &parEl)
	if err != nil {fmt.Printf(" error checkPar: empty line!")}

	parEl.txtSt = linSt
	parEl.txt = string(buf[parEl.txtSt:parEl.txtEnd+1])
	parEl.typ = par
	el.parEl = &parEl

	// see whether the previous element of elList is a parEl
	last := len(mdP.elList) -1
	lastEl := mdP.elList[last]
fmt.Printf("\n last el: %d %s\n", last, getElTyp(lastEl))
	// prev element: par
	if lastEl.parEl != nil {
		// lastEl is a parEl
		// we tack the txtstring onto the parEl
		lastParEl := lastEl.parEl

		if lastParEl.fin {
			// last paragraph was ended -> create new parel
fmt.Printf(" new el par txt: %s ", parEl.txt)
			mdP.elList = append(mdP.elList, el)
			mdP.istate = PAR
			return nil
		}
		// last par was not ended
		lastParEl.txt += " " + parEl.txt
		lastParEl.txtEnd = parEl.txtEnd
fmt.Printf(" ex el par txt: %s ", parEl.txt)
		mdP.istate = PAR
		return nil
	}


	if lastEl.emEl {
fmt.Printf(" new el par txt: %s ", parEl.txt)
		mdP.elList = append(mdP.elList, el)
		mdP.istate = PAR
		return nil
	}

	return fmt.Errorf("state: %s elList at %d has no valid el", dispState(mdP.istate), last)
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

func (mdP *mdParseObj) checkUnList(lin int) (err error){
// method that parsese and unordered list item

	var el structEl
	var ulEl uList
	var parEl parEl
	var ulCh byte

//	listEl := mdP.elList[el]
	linSt := mdP.linList[lin].linSt
	linEnd := mdP.linList[lin].linEnd
	buf := (*mdP.inBuf)
	parSt:=0
	istate := 0
	wsNum := 0
	ulCh = 0

	for i:= linSt; i< linEnd; i++ {
		switch istate {
			case 0:
				switch buf[i] {
				case ' ':
					wsNum++
				case '*','+','-':
					istate = 1
					ulCh = buf[i]
				default:
					if utilLib.IsAlpha(buf[i]) {
						parSt = i
					}
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

	if parSt == 0 {return fmt.Errorf("UL no text found!")}

	// nest lev
	nestLev := wsNum/4

	parEl.txtSt = parSt
	parEl.typ = ul
	parEl.nest = nestLev
	err = mdP.checkParEOL(lin, &parEl)
	parEl.txt = string(buf[parEl.txtSt:parEl.txtEnd+1])

fmt.Printf(" UL txt: %s ", parEl.txt)

	last := len(mdP.elList) - 1

	if last > -1 {
		// if there is a lastEl
		lastEl := mdP.elList[last]
		if lastEl.ulEl != nil {
			// if prev el is ul
			if ulCh == 0 {
				// if no mark char, par
				lastEl.ulEl.parEl.txt += parEl.txt
				lastEl.ulEl.parEl.txtEnd = parEl.txtEnd
				lastEl.ulEl.parEl.fin = false
				return nil
			} else {
				// new element
				lastEl.ulEl.parEl.fin = true
			}
		} else {
			// if not there is nothing to do
		}
	}

	// new element
	ulEl.nest = nestLev
	ulEl.parEl = &parEl
	el.ulEl = &ulEl
	mdP.elList = append(mdP.elList, el)
	mdP.istate = UL

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

	return ch, wsCount, nil
}


func (mdP *mdParseObj) checkBlock(lin int, ws int) (err error){
// A method that parses a blockquote

	var el structEl
	var bkEl bkEl
	var parEl parEl

	linSt := mdP.linList[lin].linSt + ws
	linEnd := mdP.linList[lin].linEnd
	buf := (*mdP.inBuf)
	nest := -1

	fch := buf[linSt]
//fmt.Printf("BLK first char: %q\n", fch)

	if mdP.istate == BLK && fch != '>' {
		last := len(mdP.elList) -1
		lastEl := mdP.elList[last]
		if lastEl.bkEl != nil {
			nest = lastEl.bkEl.nest
		} else {
			return fmt.Errorf("lin %d: istate: %d no blkEl found!", lin, mdP.istate)
		}

		parEl.txt = string(buf[linSt:linEnd])
		parEl.txtSt = linSt
		parEl.txtEnd = linEnd

		bkEl.nest = nest
		bkEl.parEl = &parEl
		el.bkEl = &bkEl
		mdP.elList = append(mdP.elList, el)
		mdP.istate = BLK

fmt.Printf(" implied blk nest: %d par txt: %s ", nest, parEl.txt)
		return nil
	}

	istate :=0
	wsCount := 0
	parSt := 0

	for i:=linSt; i< linEnd; i++ {
		ch := buf[i]
		switch istate {
		case 0:
			switch ch {
				case ' ':
				istate = 10
				case '>':
				nest++
				istate = 1
				default:
			}
		// else error
		case 1:
			// >
			switch ch {
				case ' ':
				// >ws

				case '>':
				// >>
				nest++
				default:
				if utilLib.IsAlpha(ch) {
				// >a
				istate = 98
				parSt = i
				}
			}
		case 2:
			// >a, >>a, > a, >  a
			switch ch {
				case '>':
				nest++
				istate = 0

				default:
				if utilLib.IsAlpha(ch) {
					istate = 98
					parSt = i
				}
			}
		case 10:
			// ' '
			switch ch {
				case ' ':
				wsCount++
				default:
				if utilLib.IsAlpha(ch) {
					istate = 99
					parSt = i
				}
			}
		default:
			break
		}
	}

	if parSt == 0 {
		parEl.txt = ""
		parEl.txtSt = linEnd -1
		parEl.txtEnd = linEnd -1
	} else {
		parEl.txt = string(buf[parSt:linEnd])
		parEl.txtSt = parSt
		parEl.txtEnd = linEnd -1
	}


	bkEl.nest = nest
	bkEl.parEl = &parEl
	el.bkEl = &bkEl
	mdP.elList = append(mdP.elList, el)
	mdP.istate = BLK

fmt.Printf(" nest: %d par %d:%d txt: %s ", nest, parEl.txtSt, parEl.txtEnd, parEl.txt)

	return nil
}

func (mdP *mdParseObj) checkCode(lin int) (endLin int, err error){
// A method that parses a code block

	var el structEl
	var parEl parEl

	linSt := mdP.linList[lin].linSt
	linEnd := mdP.linList[lin].linEnd
	buf := (*mdP.inBuf)

	// find code start
	numBq :=0
	for i:= linSt; i< linEnd; i++ {
		switch buf[i] {
			case '`':
				numBq++
			case '\r', '\n':
				if numBq < 3 {
					return lin, fmt.Errorf("line %d insufficient backquotes: %d", lin, numBq)
				}
			default:
					return lin, fmt.Errorf("line %d wrong char %q!", lin, buf[i])
		}
		if numBq > 2 {break}
	}

	endLin = lin
	endPos := 0
	stPos :=mdP.linList[lin+1].linSt
	for ilin:= lin+1; ilin< len(mdP.linList); ilin++ {
		if buf[stPos] == '`' {
			endLin = ilin
			endPos = mdP.linList[ilin-1].linEnd-1
			break
		}
	}
	if endLin < lin+2 {
		return endLin, fmt.Errorf("el COD lines %d:%d no content!", lin, endLin)
	}

	parEl.txt = string(buf[stPos:endPos])
	parEl.txtSt = stPos
	parEl.txtEnd = endPos
	parEl.typ = cod
	el.parEl = &parEl
	mdP.elList = append(mdP.elList, el)
	mdP.istate = COD

	return endLin, nil
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

func (mdP *mdParseObj) closeOrList(lin int)(err error) {

	last := len(mdP.elList) -1
	lastEl := mdP.elList[last] 
	if lastEl.olEl == nil {return fmt.Errorf("last el is not ol!")}
	nest := lastEl.olEl.nest
	lastEl.olEl.count[nest] = 0
	return nil
}


func (mdP *mdParseObj) checkOrList(lin int)(err error) {
// method that parses an ordered list

// codes
// 0 orlist item
// 1 not an orlist item
//
	var el structEl
	var orEl oList
	var parEl parEl

	linSt := mdP.linList[lin].linSt
	linEnd := mdP.linList[lin].linEnd
	buf := (*mdP.inBuf)

	parSt:=0

	markSt:= 0
	markEnd := 0
	wsNum := 0
	istate := 0
	mNum :=0

	for i:= linSt; i< linEnd; i++ {
		ch:= buf[i]
		switch istate {
		case 0:
			if ch == ' ' {
				wsNum++
			}
			if utilLib.IsNumeric(ch) {
				// new list item
				markSt = i
				istate = 1
				mNum = int(ch) - 48 + mNum*10
			}
			if utilLib.IsAlpha(ch) {
				parSt = i
			}

		case 1:
			if ch == '.' {
				markEnd = i
				istate = 2
			}

		case 2:
			if ch == ' ' {
				istate = 3
			}

		case 3:
			if utilLib.IsAlpha(ch) {
				parSt = i
			}
		default:

		}
		if parSt > 0 {break}
	}

//	if markSt == 0 {return fmt.Errorf(" orList: no counter digit!")}
	if markEnd == 0 {return fmt.Errorf(" orList: no period after counter!")}

	// convert count into number
//	numStr := string(buf[markSt:markEnd])

	if parSt == 0 {return fmt.Errorf("olList text start not found!")}

	// nest lev
	nest := wsNum/4

	if mNum > 0 {fmt.Printf("orlist marker: %d nest: %d\n", mNum, nest)}

	parEl.txtSt = parSt
	parEl.typ = ol
	parEl.nest = nest
	mdP.checkParEOL(lin, &parEl)
	parEl.txt = string(buf[parEl.txtSt:parEl.txtEnd+1])

fmt.Printf(" OL txt: %s ", parEl.txt)

	last := len(mdP.elList) -1
	if last > -1 {
		// if there is a last el
		lastEl := mdP.elList[last]
		if lastEl.olEl != nil {
			// if the lastEL is ol item, then
			// if there is no mark, add the text to the previous ol element
			if markSt == 0 {
				lastEl.olEl.parEl.txt += parEl.txt
				lastEl.olEl.parEl.txtEnd = parEl.txtEnd
				lastEl.olEl.parEl.fin = false
				mdP.istate = OL
				return nil
			} else {
				// new element
				lastEl.olEl.parEl.fin = true
			}
		} else {
			// if lastEl is not a ol element, there is nothing to do
		}
	}

	// a new list item
	mdP.cliCount[nest]++
	orEl.nest = nest
	orEl.count[nest] = 	mdP.cliCount[nest]
	orEl.parEl = &parEl
	el.olEl = &orEl
	mdP.elList = append(mdP.elList, el)
	mdP.istate = OL

	return nil
}

func (mdP *mdParseObj) closeOL()() {

	for inest:=0; inest< mdP.cnest; inest++ {
		mdP.cliCount[inest] = 0
	}

}


func (mdP *mdParseObj) checkBR()(err error) {
// method that parses an empty line

	var el structEl

	last := len(mdP.elList) -1
	if last > -1 {
		lastEl := mdP.elList[last]
		switch mdP.istate {
			case PAR:
				lastEl.parEl.fin = true
			case UL:
				lastEl.ulEl.parEl.fin = true
			case OL:
				lastEl.olEl.parEl.fin = true
			default:
		}
	}

	el.emEl = true
	mdP.elList = append(mdP.elList, el)
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

func cvtParHtml (parel *parEl) (htmlStr, cssStr string, err error) {


	return htmlStr, cssStr, err
}

//html
func parseEl (el structEl) (htmlStr, cssStr string, err error) {

	var eltyp int

	typStr := ""
	switch {
	case el.parEl != nil:
		eltyp = PAR
		typStr = fmt.Sprintf("<!--- %s --->\n", dispState(eltyp))
		htmlStr, cssStr, err = cvtParHtml(el.parEl)
	case el.emEl:
		eltyp = EP
	case el.hrEl:
		eltyp = HR
	case el.ulEl != nil:
		eltyp = UL
	case el.olEl !=nil:
		eltyp = OL
	case el.comEl != nil:
		eltyp = COM
	case el.tblEl != nil:
		eltyp = TBL
	case el.imgEl != nil:
		eltyp = IMG
	case el.bkEl != nil:
		eltyp = BLK
	default:
		eltyp = UK

	}

	htmlStr = typStr + htmlStr
	return htmlStr, cssStr, err
}


func (mdP *mdParseObj) cvtElListHtml()(htmlStr string, cssStr string, err error) {

	var el structEl

	for elIdx:=0; elIdx<len(mdP.elList); elIdx++ {
		el = mdP.elList[elIdx]
		errStr := ""
		thtmlStr, tcssStr, err := parseEl(el)
		if err != nil {errStr = fmt.Sprintf("<!--- error el %d: %v --->\n", elIdx, err)}
		htmlStr += thtmlStr
		if len(errStr) > 0 {htmlStr += errStr}
		cssStr += tcssStr
	}
	return htmlStr, cssStr, nil
}

func (mdP *mdParseObj) CvtMdToHtml(outfil *os.File)(err error) {
// method that converts the parsed element list of an md file inot an html file

	nam := outfil.Name()
	fmt.Printf("out file name: %s\n", nam)

	htmlStr, cssStr,_ := mdP.cvtElListHtml()

    outstr := htmlLib.CreHtmlHead()

	outstr += "<style>\n"
	outstr += cssStr
	outstr += "</style>\n"

    outstr += htmlLib.CreHtmlMid()

	outstr += htmlLib.CreHtmlDivMain("main")

	outstr += htmlStr

	outstr += "  </div>\n"

    outstr += htmlLib.CreHtmlEnd()

	_,err = outfil.WriteString(outstr)
	if err != nil {return fmt.Errorf("cannot write outstr: %v", err)}
	return nil
}

func (mdP *mdParseObj) addErr (lin int, msg string) {

	if mdP.errCount > 9 {return}

	mdP.errCount++
	idx := mdP.errCount - 1
	err := mdP.errList[idx]
	err.lin = lin
	err.msg = msg
	return
}

func (mdP *mdParseObj) printErrList () {
// method that prints the error list

	if mdP.errCount == 0 {
		fmt.Println("*** no Parsing Errors ***")
		return
	}

	fmt.Println("*********** Error List ***********")
	max_errors := 10
	if mdP.errCount < max_errors {max_errors = mdP.errCount}
	for i:=0; i< max_errors; i++ {
		fmt.Printf("err %d line %d: %s\n", mdP.errList[i].lin, mdP.errList[i].msg)
	}
	if max_errors == 10 {fmt.Println("too many errors only 10 are displayed!")}
	return
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
			fmt.Printf( "par typ %-5s %t: text: %s ", dispHtmlEl(ParEl.typ), ParEl.fin, ParEl.txt)

			subLen := len(ParEl.subEl)
			if subLen == 1 {
				fmt.Printf(" subel 0: \"%s\"\n", ParEl.subEl[0].txt)
			} else {
				fmt.Printf("\n")
				for i:=0; i< subLen; i++ {
					subEl := ParEl.subEl[i]
					if subEl.link {
						fmt.Printf("         subel %d bold %t italic %t link: %s: %s\n", i, subEl.bold, subEl. italic, subEl.lkUri, subEl.txt)
					} else {
						fmt.Printf("         subel %d bold %t italic %t link: %s: %s\n", i, subEl.bold, subEl. italic, subEl.txt)
					}
				}
			}

			fmt.Printf("\n")
			continue
		}
		if el.ulEl != nil {
			fmt.Printf( "ul nest %d: ", el.ulEl.nest)
			if el.ulEl.parEl != nil {
				ParEl := el.ulEl.parEl
				subLen := len(ParEl.subEl)
				fmt.Printf( "  par typ: %-5s %t len: %d text: %s\n", dispHtmlEl(ParEl.typ), ParEl.fin, subLen, ParEl.txt)

				for i:=0; i< subLen; i++ {
					subEl := ParEl.subEl[i]
					if subEl.link {
						fmt.Printf("         subel %d bold %t italic %t link: %s txt: %s\n", i, subEl.bold, subEl.italic, subEl.lkUri, subEl.txt)
					} else {
						fmt.Printf("         subel %d bold %t italic %t txt: %s: \n", i, subEl.bold, subEl.italic, subEl.txt)
					}
				}

			} else {
				fmt.Printf( "ulEl error par is nil!\n")
			}

			continue
		}

		if el.olEl != nil {
			nest := el.olEl.nest
			fmt.Printf( "ol nest %d counter %d: ", el.olEl.nest, el.olEl.count[nest])
			if el.olEl.parEl != nil {
				ParEl := el.olEl.parEl
				subLen := len(ParEl.subEl)
				fmt.Printf( "  par typ: %-5s %t len: %d text: %s\n", dispHtmlEl(ParEl.typ), ParEl.fin, subLen, ParEl.txt)

				for i:=0; i< subLen; i++ {
					subEl := ParEl.subEl[i]
					if subEl.link {
						fmt.Printf("         subel %d bold %t italic %t link: %s txt: %s\n", i, subEl.bold, subEl.italic, subEl.lkUri, subEl.txt)
					} else {
						fmt.Printf("         subel %d bold %t italic %t txt: %s: \n", i, subEl.bold, subEl.italic, subEl.txt)
					}
				}

			} else {
				fmt.Printf( "olEl error par is nil!\n")
			}

			continue
		}


		if el.tblEl != nil {
			fmt.Printf( "tbl: rows: %d  cols: %d \n", el.tblEl.rows, el.tblEl.cols)
			continue
		}

		if el.imgEl != nil {
			fmt.Printf( "img h: %d w: %d src: %s alt: %s\n", el.imgEl.height, el.imgEl.width, el.imgEl.src, el.imgEl.alt)
			continue
		}

		if el.comEl != nil {
			fmt.Printf( "com: %s\n", el.comEl.txt)
			continue
		}

		if el.bkEl != nil {
			txtSt := el.bkEl.parEl.txtSt
			txtEnd := el.bkEl.parEl.txtEnd
			txtLen := txtEnd - txtSt
			fmt.Printf( "blk lev: %d txtlen %d: ", el.bkEl.nest, txtLen)
			if txtLen > 100 {
				fmt.Printf("st: %d:%d: ", txtSt, txtEnd)
				fmt.Printf(" text 100: %s\n", string(el.bkEl.parEl.txt[:100]))
			} else {
				fmt.Printf(" text %s\n", el.bkEl.parEl.txt)
			}
			continue
		}

		if el.errEl!= nil {
			fmt.Printf("error line %d: fch %q %d:: error %s\n", el.errEl.line, el.errEl.fch, el.errEl.fch, el.errEl.errmsg)
			continue
		}

	}
}
