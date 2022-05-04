// gdocUtilLib
// library for common utility functions
//
// author: prr
// date: 27/4/2022
// copyright 2022 prr, azul software
//
package gdocUtilLib

import (
	"fmt"
	"os"
	"net/http"
	"io"
    "unicode/utf8"
	"gopkg.in/yaml.v3"
	"google.golang.org/api/docs/v1"
)

type OptObj struct {
    DefLinSpacing float64
    BaseFontSize int
    CssFil bool
    ImgFold bool
    Verb bool
    Toc bool
    Sections bool
    DivBorders bool
	CreImgFolder bool
    Divisions []string
    DocMargin [4]int
    ElMargin [4]int
}

type OptYaml struct {
	BaseFont font_desc `yaml:"BaseFont"`
	Doc doc_desc `yaml:"Doc"`
	Output out_opt `yaml:"Output"`
	Files fil_opt `yaml:"Files"`
}

type font_desc struct {
	Name string `yaml:"Font"`
	Size string `yaml:"Size"`
	LinSpace float64 `yaml:"LineSpacing"`
}

type doc_desc struct {
	DocId string `yaml:"DocId"`
	DocMargin []string `yaml:"Doc Margin"`
	ElMargin []string `yaml:"El Margin"`
	Divisions []string `yaml:"Division"`
}

type out_opt struct {
	Verb bool `yaml:"Verbose"`
	Toc bool `yaml:"TOC"`
	Sec bool `yaml:"Sections"`
	ImgFold bool `yaml:"Img Folder"`
	CssFil bool `yaml:"Css File"`
}

type fil_opt struct {
	CssFil bool `yaml:"CssFile"`
	ScriptFil bool `yaml:"ScriptFile"`
	SummaryFil bool `yaml:"SummaryFile"`
	TagFil bool `yaml:"TagFile"`
}

func GetColor(color  *docs.Color)(outstr string) {
    outstr = ""
        if color != nil {
            blue := int(color.RgbColor.Blue*255.0)
            red := int(color.RgbColor.Red*255.0)
            green := int(color.RgbColor.Green*255)
            outstr += fmt.Sprintf("rgb(%d, %d, %d)", red, green, blue)
            return outstr
        }
    outstr = "/*no color*/\n"
    return outstr
}

func GetDash(dashStyle string)(outstr string) {

    switch dashStyle {
        case "SOLID":
            outstr = "solid"
        case "DOT":
            outstr = "dotted"
        case "DASH":
            outstr = "dashed"
        default:
            outstr = "none"
    }

    return outstr
}

func GetImgLayout (layout string) (ltyp int, err error) {

    switch layout {
        case "WRAP_TEXT":

        case "BREAK_LEFT":

        case "BREAK_RIGHT":

        case "BREAK_LEFT_RIGHT":

        case "IN_FRONT_OF_TEXT":

        case "BEHIND_TEXT":

        default:
            return -1, fmt.Errorf("layout %s not implemented!", layout)
    }
    return ltyp, nil
}

func Get_vert_align (alStr string) (outstr string) {
    switch alStr {
        case "TOP":
            outstr = "top"
        case "Middle":
            outstr = "middle"
        case "BOTTOM":
            outstr = "bottom"
        default:
            outstr = "baseline"
    }
    return outstr
}

func GetGlyphStr(nlev *docs.NestingLevel)(glyphTyp string) {

    // ordered list
    switch nlev.GlyphType {
        case "DECIMAL":
            glyphTyp = "decimal"
        case "ZERO_DECIMAL":
            glyphTyp = "decimal-leading-zero"
        case "ALPHA":
            glyphTyp = "lower-alpha"
        case "UPPER_ALPHA":
            glyphTyp = "upper-alpha"
        case "ROMAN":
            glyphTyp = "lower-roman"
        case "UPPER_ROMAN":
            glyphTyp = "upper-roman"
        default:
            glyphTyp = ""
    }
    if len(glyphTyp) > 0 {
//      cssStr = "  list-style-type: " + glyphTyp +";\n"
        return glyphTyp
    }

    // unordered list
//  cssStr =fmt.Sprintf("/*-Glyph Symbol:%x - */\n",nlev.GlyphSymbol)
    r, _ := utf8.DecodeRuneInString(nlev.GlyphSymbol)

    switch r {
        case 9679:
            glyphTyp = "disc"
        case 9675:
            glyphTyp = "circle"
        case 9632:
            glyphTyp = "square"
        default:
            glyphTyp = ""
    }
    if len(glyphTyp) > 0 {
//      cssStr = "  list-style-type: " + glyphTyp +";\n"
        return glyphTyp
    }
//  cssStr = fmt.Sprintf("/* unknown GlyphType: %s Symbol: %s */\n", nlev.GlyphType, nlev.GlyphSymbol)
    return glyphTyp
}


func GetGlyphOrd(nestLev *docs.NestingLevel)(bool) {

    ord := false
    glyphTyp := nestLev.GlyphType
    switch glyphTyp {
        case "DECIMAL":
            ord = true
        case "ZERO_DECIMAL":
            ord = true
        case "UPPER_ALPHA":
            ord = true
        case "ALPHA":
            ord = true
        case "UPPER_ROMAN":
            ord = true
        case "ROMAN":
            ord = true
        default:
            ord = false
    }
    return ord
}

func ReadYamlFil(filepath, filnam string)(opt *OptYaml, err error) {
	var optyaml OptYaml

	filpath, size, err := CheckFil(filepath, filnam)
	if err != nil {
		return nil, fmt.Errorf("checkFil: %v", err)
	}

	infil, err := os.Open(filpath)
	if err != nil {
		return nil, fmt.Errorf("error - os.Open: %v", err)
	}
	defer infil.Close()

	inbuf := make([]byte, size)
	_, err = infil.Read(inbuf)
	if err != nil {
		return nil, fmt.Errorf("error - infil.Read: %v", err)
	}

	err = yaml.Unmarshal(inbuf, &optyaml)
	if err != nil {
		fmt.Printf("unmarshall error: %v\n", err)
		return nil, fmt.Errorf("error - Unmarshal: %v", err)
	}
//	fmt.Printf("opt: %v", optyaml)
	return &optyaml, nil
}

func PrintOptYaml(opt *OptYaml) {

	fmt.Printf("******** Options ***********\n")
	fmt.Printf("Base Font:\n")
	fmt.Printf("  Name: %s\n", opt.BaseFont.Name)
	fmt.Printf("Output:\n")
	fmt.Printf("  Verbose: %t\n", opt.Output.Verb)
	fmt.Println("***************************\n")
}

func FatErr(fs string, msg string, err error) {
// function that displays a console error message and exits program
    if err != nil {
        fmt.Printf("error %s:: %s!%v\n", fs, msg, err)
    } else {
        fmt.Printf("error %s:: %s!\n", fs, msg)
    }
    os.Exit(2)
}


func CheckFil(folderPath, filnam string)(filPath string, size int64, err error) {

//	fmt.Printf("folderPath: %s filnam: %s\n", folderPath, filnam)
	fullFilNam := ""
	flen := len(filnam)
	if !(flen > 0) {
		return "", 0, fmt.Errorf("error - no filnam provided!")
	}

	// check whether filnam contains the yaml extension
	extExist := false
	if flen > 5 {
		if filnam[flen-5:flen] == ".yaml" {
			fullFilNam = filnam
			extExist = true
		}
	}
	if !extExist {
		fullFilNam = filnam + ".yaml"
	}

//	fmt.Printf("fullFilNam: %s\n", fullFilNam)

	// check whether Folder Path is a viable folder
    lenFP := len(folderPath)
    if lenFP > 0 {
        filinfo, err := os.Stat(folderPath)
        if err != nil {
            if os.IsNotExist(err) {
                return "", 0, fmt.Errorf("folderPath %s is not valid!", folderPath)
            }
        }
        if !filinfo.IsDir() {
            return "", 0, fmt.Errorf("folderPath %s is not a folder!", folderPath)
        }
		if lenFP > 1 {
	        if folderPath[lenFP-1] == '/' {
    	        filPath = folderPath + fullFilNam
        	} else {
            	filPath = folderPath + "/" + fullFilNam
        	}
		}
    } else {
        filPath = fullFilNam
    }
//	fmt.Printf("filPath: %s\n",filPath)
    // check whether file exists
    filinfo, err := os.Stat(filPath)
    if os.IsNotExist(err) {
		return "", 0, fmt.Errorf("file %s does not exist!", filPath)
	}
	size = filinfo.Size()
	return filPath, size, nil
}

func DownloadImages(doc *docs.Document, imgFoldPath string, opt *OptObj)(err error) {

    if !(len(imgFoldPath) >0) {
        return fmt.Errorf("error -- no imgFoldPath provided!")
    }

	_, err = os.Stat(imgFoldPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("error -- folder %s does not exist!", imgFoldPath)
		}
		return fmt.Errorf("error -- accessing folder %s! %v", imgFoldPath, err)
	}

	verb := true
	if opt != nil { verb = opt.Verb}

	// inline images
	if verb {fmt.Printf("*** Inline Imgs: %d ***\n", len(doc.InlineObjects))}
    for k, inlObj := range doc.InlineObjects {
        imgProp := inlObj.InlineObjectProperties.EmbeddedObject.ImageProperties
        if verb {
            fmt.Printf("Source: %s Obj %s\n", k, imgProp.SourceUri)
            fmt.Printf("Content: %s Obj %s\n", k, imgProp.ContentUri)
        }
        if !(len(imgProp.SourceUri) > 0) {
            return fmt.Errorf("error -- image %s has no URI\n", k)
        }
        imgNam := imgFoldPath + k[4:] + ".jpeg"
        if verb {fmt.Printf("image path: %s\n", imgNam)}
        URL := imgProp.ContentUri
        httpResp, err := http.Get(URL)
        if err != nil {
            return fmt.Errorf("error -- http.Get could not fetch %s! %v", URL, err)
        }
        defer httpResp.Body.Close()
//  fmt.Printf("http got %s!\n", URL)
        if httpResp.StatusCode != 200 {
            return fmt.Errorf("error -- httpResp: Received non-200 response code %d!", httpResp.StatusCode)
        }
//  fmt.Printf("http status: %d\n ", httpResp.StatusCode)
    //Create a empty file
        outfil, err := os.Create(imgNam)
        if err != nil {
            return fmt.Errorf("error -- os.Create cannot create img file! %v", err)
        }
        defer outfil.Close()
//  fmt.Println("created dir")
        //Write the bytes to the fiel
        _, err = io.Copy(outfil, httpResp.Body)
        if err != nil {
            return fmt.Errorf("error -- io.Copy cannot copy img file content from httpRespBody! %v", err)
        }
    }

	// positioned images
	fmt.Printf("*** Positioned Imgs: %d ***\n", len(doc.PositionedObjects))
    for k, posObj := range doc.PositionedObjects {
        imgProp := posObj.PositionedObjectProperties.EmbeddedObject.ImageProperties
        if verb {
            fmt.Printf("Source: %s Obj %s\n", k, imgProp.SourceUri)
            fmt.Printf("Content: %s Obj %s\n", k, imgProp.ContentUri)
        }
        if !(len(imgProp.SourceUri) > 0) {
            return fmt.Errorf("error -- image %s has no URI\n", k)
        }
        imgNam := imgFoldPath + k[4:] + ".jpeg"
        if verb {fmt.Printf("image path: %s\n", imgNam)}
        URL := imgProp.ContentUri
        httpResp, err := http.Get(URL)
        if err != nil {
            return fmt.Errorf("error -- httpGet: could not get %s! %v", URL, err)
        }
        defer httpResp.Body.Close()
//  fmt.Printf("http got %s!\n", URL)
        if httpResp.StatusCode != 200 {
            return fmt.Errorf("error -- hhtpResp: Received non-200 response code %d!", httpResp.StatusCode)
        }
//  fmt.Printf("http status: %d\n ", httpResp.StatusCode)
    //Create a empty file
        outfil, err := os.Create(imgNam)
        if err != nil {
            return fmt.Errorf("error -- osCreate: cannot create img file! %v", err)
        }
        defer outfil.Close()
//  fmt.Println("created dir")
        //Write the bytes to the fiel
        _, err = io.Copy(outfil, httpResp.Body)
        if err != nil {
            return fmt.Errorf("error -- ioCopy: cannot copy img file content! %v", err)
        }
    }

    return nil
}

func CreateImgFolder(folderPath, docNam string)(imgFolderPath string, err error) {

	// we assume folderPath is correct
	// any last slash has already been stripped from folderPath
	lenFP := len(folderPath)
	switch {
		case lenFP == 0:
			imgFolderPath = docNam + "_img"
		case lenFP > 0:
			imgFolderPath = folderPath + "/" + docNam + "_img"
	}

 //   fmt.Println("img folder path: ", imgFolderPath)

    // check whether dir folder exists, if not create one
    newDir := false
    _, err = os.Stat(imgFolderPath)
    if os.IsNotExist(err) {
        err1 := os.Mkdir(imgFolderPath, os.ModePerm)
        if err1 != nil {
            return "", fmt.Errorf("os.Mkdir: could not create img folder! %v", err1)
        }
        newDir = true
    } else {
        if err != nil {
            return "", fmt.Errorf("os.Stat: general error but not find! %v", err)
        }
    }

	// img folder is presumed to exist now
	// if img folder is not a new dir, all image files need to be deleted
    // removeAll also removes the folder, so we have to create a new folder
    if !newDir {
        err = os.RemoveAll(imgFolderPath)
        if err != nil {
            return "", fmt.Errorf("os.RemoveAll: could not delete files in image folder! %v", err)
        }
        err = os.Mkdir(imgFolderPath, os.ModePerm)
        if err != nil {
            return "", fmt.Errorf("os.Mkdir: could not create img folder! %v", err)
        }
    }

    return imgFolderPath, nil
}

func CreateOutFil(folderPath, filNam, filExt string) (outfil *os.File, err error) {
    var fullFilNam, filpath string

    if len(filNam) == 0 {
        return nil, fmt.Errorf("file name is empty!")
    }
    // create full file name: filnam + ext
    // check file extension
    if len(filExt) == 0 {
        ext := false
        for i:=len(filNam) -1; i>=0; i-- {
            if filNam[i] == '.' {
                ext = true
                break
            }
        }
        if !ext {return nil, fmt.Errorf("no file extension provided!")}
        fullFilNam = filNam
    } else {
        // check extension
        if filExt[0] == '.' {
            fullFilNam = filNam + filExt
        } else {
            fullFilNam = filNam + "." + filExt
        }
    }

    lenFP := len(folderPath)
    if lenFP > 0 {
        filinfo, err := os.Stat(folderPath)
        if err != nil {
            if os.IsNotExist(err) {
                return nil, fmt.Errorf("folderPath %s is not valid!", folderPath)
            }
        }
        if !filinfo.IsDir() {
            return nil, fmt.Errorf("folderPath %s is not a folder!", folderPath)
        }
        if folderPath[lenFP-1] == '/' {
            filpath = folderPath + fullFilNam
        } else {
            filpath = folderPath + "/" + fullFilNam
        }
    } else {
        filpath = fullFilNam
    }
    // check whether file exists
    _, err = os.Stat(filpath)
    if !os.IsNotExist(err) {
        err1:= os.Remove(filpath)
        if err1 != nil {
            return nil, fmt.Errorf("os.Remove: cannot remove existing file: %s! error: %v", filpath, err1)
        }
    }

    outfil, err = os.Create(filpath)
    if err != nil {
        return nil, fmt.Errorf("os.Create: cannot create file: %s! %v", filpath, err)
    }
    return outfil, nil
}

func CreateFileFolder(path, foldnam string)(fullPath string, existDir bool, err error) {

    // check if foldenam is valid -> no whitespaces
    fnamValid := true
    for i:=0; i< len(foldnam); i++ {
        if foldnam[i] == ' ' {
            fnamValid = false
            break
        }
    }

    if !fnamValid {
        return "", false, fmt.Errorf("error -- not a valid folder name %s!", foldnam)
    }

    // check whether foldnam folder exists
    fullPath =""
    switch {
        case len(path) == 0:
            fullPath = foldnam

        case path[0] == '/':
            return "", false, fmt.Errorf("error -- absolute path!")

        case path[len(path)  -1] == '/':
                fullPath = path + foldnam

        default:
                fullPath = path + "/" + foldnam
    }

//  fmt.Printf("full path1: %s\n", fullPath)

    // check path with folder name
    // add trimming wsp to left
    if _, err1 := os.Stat(fullPath); !os.IsNotExist(err1) {
        return fullPath, true, nil
    }

//  fmt.Printf("full path2: %s\n", fullPath)

    // path does not exist, we need to create path
    ist:=0
    for i:=0; i<len(fullPath); i++ {
        if fullPath[i] == '/' {
            parPath := string(fullPath[ist:i])
//  fmt.Printf("path %d: %s\n", i, parPath)
            if _, err1 := os.Stat(parPath); os.IsNotExist(err1) {
                err2 := os.Mkdir(parPath, os.ModePerm)
                if err2 != nil {
                    return "", false, fmt.Errorf("os.Mkdir: lev %d %v", err2, i)
                }
//          ist = i + 1
            }
        }
    }
    err = os.Mkdir(fullPath, os.ModePerm)
    if err != nil {
        return "", false, fmt.Errorf("full Path os.Mkdir: %v", err)
    }

    return fullPath, false, nil
}

func GetDefOption(opt *OptObj) {

    opt.BaseFontSize = 0
    opt.DivBorders = false
    opt.DefLinSpacing = 1.2
    opt.DivBorders = false
    opt.CssFil = false
    opt.ImgFold = true
    opt.Verb = true
    opt.Toc = true
    opt.Sections = true
    for i:=0; i< 4; i++ {opt.ElMargin[i] = 0}
    opt.Divisions = []string{"Main"}

    return
}

func PrintOptions (opt *OptObj) {

    fmt.Printf("\n************ Option Values ***********\n")
    fmt.Printf("  Base Font Size:       %d\n", opt.BaseFontSize)
    fmt.Printf("  Sections as <div>:    %t\n", opt.Sections)
    fmt.Printf("  Browser Line Spacing: %.1f\n",opt. DefLinSpacing)
    fmt.Printf("  <div> Borders:        %t\n", opt.DivBorders)
    fmt.Printf("  Divisions: %d\n", len(opt.Divisions))
    for i:=0; i < len(opt.Divisions); i++ {
        fmt.Printf("    div: %s\n", opt.Divisions[i])
    }
    fmt.Printf("  Separate CSS File:    %t\n", opt.CssFil)
    fmt.Printf("  Image Folder:         %t\n", opt.ImgFold)
    fmt.Printf("  Table of Content:     %t\n", opt.Toc)
    fmt.Printf("  Verbose output:       %t\n", opt.Verb)
    fmt.Printf("  Element Margin: ")
    for i:=0; i<4; i++ { fmt.Printf(" %3d",opt.ElMargin[i])}
    fmt.Printf("\n")
    fmt.Printf("***************************************\n\n")
}

