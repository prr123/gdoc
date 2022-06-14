# golang libraries for google document and markdown files

The libraries allow the text conversion from one file format into another.

Input files are:  
*  gdoc
*  txt
*  md

Output files are:
*  gdoc
*  txt
*  md
*  html (integrated css with html file and seperate html and css files)
*  dom (js and css files; base html file)

## documentation
There will be md files for each library

## examples
to come 

## gdocTxtLib.go
This library is mostly for debugging. It generates a text output of a gdoc sutructure.

## gdocMdLib.go
This library contains methods to convert a google document file into a markdown file.

## gdocApiLib.go
This library contains the subroutines to access the google docs api.

## gdocHtmlLib.go
library that converts a Gdoc document file to a html file with a css section or a html file and a css file

## gdocDocLib
library that converts a Gdoc document to a base html file with a a js section and a css section.
The js builds the dom to convert the gdoc file into a web file.

## gdocDomLib.go
library that converts a Gdoc document to a base html file with a a js section and a css section.
The js builds the dom to convert the gdoc file into a web file.

## gdocUtilLib.go  
utility library for the options object and file/ file folder creation

## mdParseLib.go  
library that parses a markdown file.
-- work in progress

## mdDomLib.go  
library that will convert a markdown file into a javascript Dom script file.
-- not functional yet!

## mdGdocLib.go  
library that will convert a markdown file into a google document.
-- not functional yet!

## mdHtmlLib.go  
library that will convert a markdown file into a html file.
-- not functional yet!

## txtGdocLib.go
library that creates a new google document from a plain text file..

---

