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

	fmt.Println("YamlToGo success")
}

func parseYaml(buf *[]byte) (outstr string, err error) {

	var linByt [](*[]byte)
	bufLen := len(*buf)

//	fmt.Printf("buf len: %d\n", bufLen)

	linst:=0
	linend:= linst+100
	if linend > bufLen {linend = bufLen}
	linCount:=0
	for {
//fmt.Printf("line [%d:%d] %d\n", linst, linend, linCount)
		if linend<linst {
			break
		}

		for i:=linst; i<linend; i++ {
			if (*buf)[i] == '\n' {
				linend = i
				break
			}
		}
		linbuf := (*buf)[linst:linend]
//fmt.Printf("line [%d:%d][%d]: %s\n", linst, linend,linCount, string(linbuf))
		idx := bytes.Index(linbuf, []byte(":"))
		if idx > -1 {
			linByt = append(linByt, &linbuf)
			linCount++
		}
		linst = linend+1
		linend = linst+100
		if linend > bufLen {linend = bufLen}
	}

// parse lines
	outstr = "type tplObj struct {\n"

	fmt.Printf("lines: %d\n", linCount)
	ilin := -1
	errstr := ""
	for ilin=0; ilin< len(linByt); ilin++ {
		lin := *linByt[ilin]
//fmt.Printf("line [%d]: %s\n", ilin, string(lin))
		col:=-1
		for j:=0; j<len(lin); j++ {
			if lin[j] == ':' {
				col = j
				break
			}
		} //j
		if col == -1 {
			errstr += fmt.Sprintf("error line [%d] no colon\n", ilin)
			continue
		}

fmt.Printf("key: %s val: %s\n", string(lin[:col]), string(lin[col+1:]))

		istate:=0
		typSt:= -1
		typEnd:= -1
		for j:=0; j<col; j++ {
			switch istate {
			case 0:
				if lin[j] == '<' {
					typSt = j
					istate = 1
				}
			case 1:
				if lin[j] == '>' {
					typEnd = j
					istate = 2
				}


			default:
			errstr += fmt.Sprintf("error line [%d] istate %d \n", ilin, istate)

			}
			if istate > 1 {break}
		}// j

		fmt.Printf("colon %d typSt %d typEnd %d\n", col, typSt, typEnd)


		typStr := ""
		namStr := ""
		if typSt < 0 {
			// type is string
			typStr = "string"
			namStr = string(lin[:col])
		fmt.Printf("name0: %s type: %s\n",namStr, typStr)
			continue
		} else {
			if typEnd < 0 {errstr += fmt.Sprintf("error line [%d] no end bracker \">\"\n", ilin, istate); continue;}
		}
		typStr = string(lin[typSt+1:typEnd])
		namStr = string(lin[:typSt])
		fmt.Printf("name1: %s type: %s\n",namStr, typStr)

	} //ilin

	return outstr, nil
}

