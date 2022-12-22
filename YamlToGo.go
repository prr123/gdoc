// golang program that parses a yaml file and generates a go struct
// author: prr, azul software
// created: 20/12/2022
// copyright 2022 prr, Peter Riemenschneider, Azul Software
//
// start: YamlToGo
//

package main

import (
	"os"
	"fmt"
	"bytes"
    util "google/gdoc/utilLib"
	)

type linItem struct {
	nam string
	namYaml string
	typ string
	val string
}

func main() {

    numArgs:= len(os.Args)

    if numArgs < 2 {
        fmt.Printf("error - exit: insufficient command line arguments\n")
        fmt.Printf("usage is: YamlToGo \"input file\" [\\out=] [\\dbg]\n")
        os.Exit(-1)
    }

    inFilNam := os.Args[1]
	idx := bytes.Index([]byte(inFilNam),[]byte(".tpl"))
	if idx == -1 {fmt.Printf("invalid input file name \"%s\": no ext \".tpl\"!\n", inFilNam); os.Exit(-1);}

	infil, err := os.Open(inFilNam)
	if err != nil {fmt.Printf("cannot open input file \"%s\": %v!\n", inFilNam, err); os.Exit(-1);}
	defer infil.Close()

    flags := [] string {"out", "dbg"}

    argmap, err := util.ParseFlagsStart(os.Args, flags,2)
    if err != nil {fmt.Printf("error ParseFlags: %v\n", err); os.Exit(-1);}

    outFilNam, ok := argmap["out"]
	outFilNamStr := ""
	tplFilNamStr := string(inFilNam[:idx]) + ".go"
    if !ok {
		outFilNamStr += tplFilNamStr
		fmt.Printf("no output filename provided! using \"%s\"\n", outFilNamStr)

	} else {
    	outFilNamStr := outFilNam.(string)
		outIdx := bytes.Index([]byte(outFilNamStr),[]byte(".go"))
		if outIdx == -1 {fmt.Printf("invalid output file name \"%s\": no ext \".go\"!\n", outFilNamStr); os.Exit(-1);}
    	if outFilNamStr == "none" {outFilNamStr += tplFilNamStr}
	}

fmt.Printf("out file: %s\n",outFilNamStr)

	outfil, err := os.Create(outFilNamStr)
	if err != nil {fmt.Printf("error creating output file: %v",err); os.Exit(-1)}
	defer outfil.Close()
	outstr := fmt.Sprintf("// go struct for yaml file: %s\n", inFilNam)
	outfil.WriteString(outstr)

	info, err := infil.Stat()
	if err != nil {fmt.Printf("error reading stat of input file: %v",err); os.Exit(-1)}

//	fmt.Printf("infile size: %d\n", info.Size())

	inbuf := make([]byte, info.Size())

	_, err = infil.Read(inbuf)
	if err != nil {fmt.Printf("error reading input file: %v",err); os.Exit(-1)}

//	fmt.Println("inbuf: ", string(inbuf[:10]))

	outstr, err = parseYaml(&inbuf)
	if err != nil {fmt.Printf("error parsing inbuf: %v",err); os.Exit(-1)}
	outfil.WriteString(outstr)
	outfil.WriteString("}")
	fmt.Println("YamlToGo success")
}

func parseYaml(buf *[]byte) (outstr string, err error) {

	bufLen := len(*buf)

//	fmt.Printf("buf len: %d\n", bufLen)

	linst:=0
	linend:= linst+100
	if linend > bufLen {linend = bufLen}

	outstr = "type tplObj struct {\n"

//	exeStr := "\nfunc initYaml () (tpl tplObj, err error) {\n"

	linCount:=0
	for {
//fmt.Printf("line [%d:%d] %d\n", linst, linend, linCount)
		if linend<linst {
			break
		}

		for i:=linst; i<linend; i++ {
			if (*buf)[i] == '\n' {
				linend = i
				linCount++
				break
			}
		}


		linByt := (*buf)[linst:linend]
//fmt.Printf("line [%d:%d][%d]: %s\n", linst, linend,linCount, string(linByt))

		item, errp := parseLin(linByt)
//fmt.Printf("item: %v\n", item)
		if errp != nil {
			outstr += fmt.Sprintf("//line [%d]: %s\n", linCount, string(linByt))
			outstr += fmt.Sprintf("//error line %d: %v\n", linCount, errp)
		} else {
			outstr += fmt.Sprintf("  %s %s `yaml:\"%s\"`\n", (*item).nam, (*item).typ, (*item).namYaml)
		}

		if linend == bufLen {break}
		linst = linend+1
		linend = linst+100
		if linend > bufLen {linend = bufLen}
		if linst >= bufLen {break}
	}

	fmt.Printf("lines: %d\n", linCount)

	return outstr, nil
}

// parse lines

func parseLin (lin []byte) (itemP *linItem, err error) {

	var item linItem

//	errstr := ""
//fmt.Printf("line %s\n",string(lin))

	col:=-1
	com := -1
	for j:=0; j<len(lin); j++ {
		switch lin[j] {
			case '#':
				com = j
			case ':':
				col = j
		}
	} //j
	if col == -1 && com == -1 {
		return nil, fmt.Errorf("no colon\n")
	}
	// comment only
	if col == -1 { return nil, nil}

	if col> com && com > 0 {return nil, fmt.Errorf("comment before colon")}
// fmt.Printf("key: %s val: %s\n", string(lin[:col]), string(lin[col+1:]))

	item.namYaml = string(lin[:col])
	// crude method of capitalising the first letter
	if lin[0] > 96 {lin[0] = lin[0] - 32}
	item.nam = string(lin[:col])

	if com == -1 {
		// no comment thus no type
		item.val = string(lin[col:])
		item.typ = "string"
		return &item, nil
	}

	// we have a comment, so v
	item.val = string(lin[col:com])


	// checking for typ
	typSt:= -1
	typEnd:= -1
	for j:=com; j< len(lin); j++ {
		switch lin[j] {
			case '<':
				typSt = j
			case '>':
				typEnd = j
		}
	}// j

	if typSt > typEnd {return nil, fmt.Errorf("\">\" before \"<\"")}
	if typSt == -1 && typEnd> 0 {return nil, fmt.Errorf("no \"<\" before \">\"")}
	if typSt < 0 {item.typ = "string"; return &item, nil;}

	item.typ = string(lin[typSt+1:typEnd])
	return &item, nil

}

