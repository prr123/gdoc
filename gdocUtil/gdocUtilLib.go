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
	"gopkg.in/yaml.v3"
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
    Divisions []string
    DocMargin [4]int
    ElMargin [4]int
}

type OptYaml struct {
	BaseFont font_desc `yaml:"BaseFont"`
	Doc doc_desc `yaml:"Doc"`
}

type font_desc struct {
	Name string `yaml:"Font"`
	Size string `yaml:"Size"`
}
type doc_desc struct {
	DocId string `yaml:"DocId"`
}

type out_opt struct {
	Verb bool `yaml:"Verbose"`
	Toc bool `yaml:"TOC"`
	Sec bool `yaml:"Sections"`
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

    opt.Divisions = []string{"Summary", "Main"}
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

