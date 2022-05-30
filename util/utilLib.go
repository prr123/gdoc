package utilLib
//
//
// library that parses text for chars
//
// author prr
// created 12/4/2022
//
// copyright 2022 prr
//

import (
	"os"
	"fmt"
)

func IsAlpha(let byte)(res bool) {
// function that test whether byte is alpha
    res = false
    if (let >= 'a' && let <= 'z') || (let >= 'A' && let <= 'Z') { res = true}
    return res
}

func IsAlphaNumeric(let byte)(res bool) {
// function that test whether byte is aphanumeric
    res = false
    tbool := (let >= 'a' && let <= 'z') || (let >= 'A' && let <= 'Z')
    if tbool || (let >= '0' && let <= '9') { res = true }
    return res
}

func IsNumeric(let byte)(res bool) {
// function that test whether byte is aphanumeric
    res = false
    if (let >= '0') && (let <= '9') { res = true }
    return res
}

func CvtBytToNum(let byte)(res int) {
// function that converts a ascii byte into an integer
	res = -1
    if (let >= '0') && (let <= '9') { 
		res = int(let) - 49
	}
	return res
}

func IsWsp(let byte)(res bool) {
// fuction that tests a white space
	res = false
	if (let ==' ')||(let == '\t') { res = true}
	return res
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

