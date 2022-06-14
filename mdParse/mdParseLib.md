<p style="font-size:18pt; text-align:center;">mdParse/mdParseLib</p>

# Description
     
   mdParseLib.go   
   parse markdown file   
   usage: parse file.go   
     
   author: prr azul software   
   date: 28 May 2022   
   copyright prr azul software   
     


# Types
## mdParseObj    
type mdParseObj struct     

## mdLin    
type mdLin struct     

## structEl    
type structEl struct     

## parEl    
type parEl struct     

## parSubEl    
type parSubEl struct     

## errEl    
type errEl struct     

## bkEl    
type bkEl struct     

## comEl    
type comEl struct     

## tblEl    
type tblEl struct     

## tblRow    
type tblRow struct     

## tblCel    
type tblCel struct     

## imgEl    
type imgEl struct     

## uList    
type uList struct     

## oList    
type oList struct     


# Functions
## dispState    
func dispState(num int)(str string)     

 function that converts state constants to strings    
## dispHtmlEl    
func dispHtmlEl(num int)(str string)     

 function that converts html const to strings    
## InitMdParse    
func InitMdParse() (mdp *mdParseObj)     

 function that initialises the mdParse object    

# Methods
## mdP *mdParseOb: ParseMdFile    
func (mdP *mdParseObj) ParseMdFile(inpfilnam string) (err error)     

 method that opens an md file and reads it into a buffer    
## mdP *mdParseOb: parseMdOne    
func (mdP *mdParseObj) parseMdOne()(err error)     

 method that conducts the first pass of the parser.    
 the first pass creates and slice of lines.    
## mdP *mdParseOb: parseMdTwo    
func (mdP *mdParseObj) parseMdTwo()(err error)     

 method that reads the linList, parses the list and creates a List of elements    
## mdP *mdParseOb: checkError    
func (mdP *mdParseObj) checkError(lin int, fch byte,  errStr string)     

 method that checks for an error in a line    
## mdP *mdParseOb: checkHeadingEOL    
func (mdP *mdParseObj) checkHeadingEOL(lin int, parEl *parEl)(err error)     

 method that tests a heading line to see whether the heading is completed in that line    
 or continues.    
## mdP *mdParseOb: checkParEnd    
func (mdP *mdParseObj) checkParEnd(lin int)(err error)     

 method that checks terminates the previous paragraph after an empty line    
## mdP *mdParseOb: checkParEOL    
func (mdP *mdParseObj) checkParEOL(lin int, parel *parEl)(res bool)     

 function to test whether a line is completed with a md cr    
## mdP *mdParseOb: checkWs    
func (mdP *mdParseObj) checkWs(lin int)(fch byte, err error)     

 method that checks indented lines to find the first character    
## mdP *mdParseOb: checkComment    
func (mdP *mdParseObj) checkComment(lin int)(err error)     

 method that checks whether line is a comment    
## mdP *mdParseOb: checkPar    
func (mdP *mdParseObj) checkPar(lin int)(err error)     

 method that parses a line  to check whether it is paragraph    
## mdP *mdParseOb: checkHeading    
func (mdP *mdParseObj) checkHeading(lin int) (err error)    

 method that parses a line for headings    
## mdP *mdParseOb: checkHr    
func (mdP *mdParseObj) checkHr(lin int) (err error)     

 method that parses a horizontal ruler line    
## mdP *mdParseOb: checkBold    
func (mdP *mdParseObj) checkBold()    

## mdP *mdParseOb: checkItalics    
func (mdP *mdParseObj) checkItalics()     

## mdP *mdParseOb: checkUnList    
func (mdP *mdParseObj) checkUnList(lin int) (err error)    

 method that parsese and unordered list item    
## mdP *mdParseOb: checkHtml    
func (mdP *mdParseObj) checkHtml()     

## mdP *mdParseOb: checkIndent    
func (mdP *mdParseObj) checkIndent(lin int) (nest int, ord bool, err error)    

 method that chck indents and returns the nesting level and list type    
## mdP *mdParseOb: checkBlock    
func (mdP *mdParseObj) checkBlock(lin int) (err error)    

 A method that parses a blockquote    
## mdP *mdParseOb: checkCode    
func (mdP *mdParseObj) checkCode(lin int) (err error)    

 A method that parses a code block    
## mdP *mdParseOb: checkStrike    
func (mdP *mdParseObj) checkStrike()     

## mdP *mdParseOb: checkImage    
func (mdP *mdParseObj) checkImage(lin int) (err error)    

 a method that parse an image    
## mdP *mdParseOb: checkLink    
func (mdP *mdParseObj) checkLink()     

## mdP *mdParseOb: checkTable    
func (mdP *mdParseObj) checkTable(lin int)(endLin int, err error)     

 method that parses a table. The method returns the last line of the table.    
 A table consists of several lines, unlike other elements    
## mdP *mdParseOb: checkOrList    
func (mdP *mdParseObj) checkOrList(lin int)(err error)     

 method that parses an ordered list    
## mdP *mdParseOb: checkBR    
func (mdP *mdParseObj) checkBR()(err error)     

 method that parses an empty line    
## mdP *mdParseOb: printLinList    
func (mdP *mdParseObj) printLinList()()     

## mdP *mdParseOb: cvtMdToHtml    
func (mdP *mdParseObj) cvtMdToHtml()(err error)     

 method that converts the parsed element list of an md file inot an html file    
## mdP *mdParseOb: printElList    
func (mdP *mdParseObj) printElList ()     

 method that prints out the structural element list    
