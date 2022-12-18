// gdocApiLib
// author: prr
// date: v3 11/Jan 2022
// copyright 2022 prr azul software
//
// v1 from rd_Gdocv3
// v2 added ext parameter to CreTxtOutFile
// v3
// 11 Jan 2022 added more error routines to Init
//

package gdocApiLib

import (
        "context"
        "encoding/json"
        "fmt"
        "io/ioutil"
        "net/http"
        "os"
        "golang.org/x/oauth2"
        "golang.org/x/oauth2/google"
        "google.golang.org/api/docs/v1"
        "google.golang.org/api/option"
)

type GdocApiObj  struct {
	GdocCtx context.Context
//	tocFile string
	Svc *docs.Service
	doc *docs.Document
}

// Retrieves a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
        tokFile := "tokGdoc.json"
        tok, err := tokenFromFile(tokFile)
        if err != nil {
                tok = getTokenFromWeb(config)
                saveToken(tokFile, tok)
        }
        return config.Client(context.Background(), tok)
}

// Requests a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
        authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
        fmt.Printf("Go to the following link in your browser then type the "+
                "authorization code: \n%v\n", authURL)

        var authCode string
        if _, err := fmt.Scan(&authCode); err != nil {
			fmt.Println("Unable to read authorization code: ", err)
			os.Exit(1)
        }

        tok, err := config.Exchange(oauth2.NoContext, authCode)
        if err != nil {
			fmt.Println("Unable to retrieve token from web: ", err)
			os.Exit(1)
        }
        return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
        f, err := os.Open(file)
        defer f.Close()
        if err != nil {
                return nil, err
        }
        tok := &oauth2.Token{}
        err = json.NewDecoder(f).Decode(tok)
        return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		fmt.Println("Unable to cache OAuth token: ", err)
		os.Exit(1)
	}
	json.NewEncoder(f).Encode(token)
}

// create text file to dump document file
// no need for mthod
func (gdoc *GdocApiObj) CreOutFile(filnam string, ext string)(outfil *os.File, err error) {
	// function create a file with the filname "filnam.ext" 
	// returns a file pointer
	//

	if !(len(filnam)>0) {return nil, fmt.Errorf("error CreOutFile:: no filnam!")}
	// check for dir
	bfil := []byte(filnam)
	fst := 0
	for i:=len(filnam)-1; i>0; i-- {
		if bfil[i] == '/' {
			fst = i
			break
		}
		if bfil[i] == '.' {
			return nil, fmt.Errorf("error CreOutFile:: found period in filnam!")
		}
	}

	if fst > (len(filnam)-1) {return nil, fmt.Errorf("error CreOutFile:: not a valid filenam end!")}


	// convert white spaces in file name to underscore
	outfilb := bfil[fst+1:]
	for i:=0; i<len(outfilb); i++ {
		if outfilb[i] == ' ' {
			outfilb[i] = '_'
		}
	}
	outfilnam := string(outfilb[:])
	outfildir := string(bfil[:fst])
	if !(len(outfildir)>0) {outfildir = "output"}
//	fmt.Printf("dir: %s file: %s\n", outfildir, outfilnam)

//    newDir := false
    _, err = os.Stat(outfildir)
    if os.IsNotExist(err) {
        err1 := os.Mkdir(outfildir, os.ModePerm)
        if err1 != nil {
            return nil, fmt.Errorf("error CreOutFile:: could not create folder! %v", err1)
        }
//        newDir = true
    } else {
        if err != nil {
            return nil, fmt.Errorf("error CreOutFile:: dir exists, but could not get info! %v", err)
        }
    }

	path:= fmt.Sprintf("%s/%s.%s", outfildir, outfilnam, ext)
	fmt.Println("output file path: ", path)
	_, err = os.Stat(path)
	if !os.IsNotExist(err) {
		err1 := os.Remove(path)
		if err1 != nil {
			return nil, fmt.Errorf("error CreOutFile: cannot remove old output file: %v!", err1)
		}
	}
	outfil, err = os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("error CreOutFile: cannot create output text file: %v!", err)
	}
	return outfil, nil
}


func (gdoc *GdocApiObj) InitGdocApi() (err error) {
        ctx := context.Background()
        b, err := ioutil.ReadFile("credGdoc.json")
        if err != nil {
			return fmt.Errorf("Unable to read client secret file: %v!", err)
		}

        // If modifying these scopes, delete your previously saved token.json.
        config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/documents.readonly")
        if err != nil {
			return fmt.Errorf("Unable to parse client secret file to config: %v", err)
        }

        client := getClient(config)

        svc, err := docs.NewService(ctx, option.WithHTTPClient(client))
        if err != nil {
			return fmt.Errorf("Unable to retrieve Docs client: %v", err)
        }
		gdoc.Svc = svc
	return nil
}

func InitGdocApiV2() (gdocObj *GdocApiObj ,err error) {

	var gdoc GdocApiObj

        ctx := context.Background()
        b, err := ioutil.ReadFile("credGdoc.json")
        if err != nil {
			return nil, fmt.Errorf("Unable to read client secret file: %v!", err)
		}

        // If modifying these scopes, delete your previously saved token.json.
        config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/documents.readonly")
        if err != nil {
			return nil, fmt.Errorf("Unable to parse client secret file to config: %v", err)
        }

        client := getClient(config)

        svc, err := docs.NewService(ctx, option.WithHTTPClient(client))
        if err != nil {
			return nil, fmt.Errorf("Unable to retrieve Docs client: %v", err)
        }
		gdoc.Svc = svc

	return &gdoc, nil
}

func (gdoc *GdocApiObj) GetDoc(docId string) (err error) {

    srv := gdoc.Svc

    doc, err := srv.Documents.Get(docId).Do()
    if err != nil {
        return fmt.Errorf("Unable to retrieve data from document: ", err)
    }
	gdoc.doc = doc
	return nil
}
