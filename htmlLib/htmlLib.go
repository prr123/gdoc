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
//	"os"
)

func CreHtmlHead()(outstr string) {
// func create start of html doc
	outstr = "<!DOCTYPE html>\n"
	outstr += "<html lang=\"en\">\n"
	outstr += "<head>\n"
	outstr += "  <meta charset=\"UTF-8\">\n"
	// add more meta tags
	outstr += "  <title>Azul Conversion</title>\n"

	return outstr
}

func CreCss()(cssStr string) {

	ws := "    "
	cssStr = "  <style>\n"

	cssStr += "  * {\n"
	cssStr += ws + "margin: 0;\n"
	cssStr += ws + "padding: 0;\n"
	cssStr += ws + "font-family: calibri;\n"
	cssStr += ws + "list-style: none;\n"
	cssStr += ws + "text-decoration: none;\n"
	cssStr += "  }\n"

	cssStr += "  </style>\n"

	return cssStr
}

func CreHtmlMid()(outstr string) {
// func ot end head and start body
	outstr = "</head>\n<body>\n"
	return outstr
}

func CreHtmlEnd()(outstr string) {
// func to create end of html doc
	outstr = "</body>\n</html>\n"
	return outstr
}

func CreHtmlDivMain(nam string)(htmlStr string) {
// func to cre div
	htmlStr = fmt.Sprintf("  <div class=\"%s\">\n", nam)

	return htmlStr
}
