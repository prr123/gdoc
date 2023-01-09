// gdocApiV2Lib
// author: prr
// date: 9 Jan 2023
//
// copyright 2023 prr azul software
//
// for license description and documentation:
// see github.com/prr123/gdoc/gdocApiV2
//
//

package gdocApiV2Lib

import (
        "context"
        "encoding/json"
        "fmt"
//        "net/http"
        "os"
        "golang.org/x/oauth2"
        "golang.org/x/oauth2/google"
        "google.golang.org/api/docs/v1"
//        "google.golang.org/api/drive/v3"
        "google.golang.org/api/option"
)

type cred struct {
    Installed credItems `json:"installed"`
    Web credItems `json:"web"`
}

type credItems struct {
    ClientId string `json:"client_id"`
    ProjectId string `json:"project_id"`
    AuthUri string `json:"auth_uri"`
    TokenUri string `json:"token_uri"`
//  Auth_provider_x509_cert_url string `json:"auth_provider_x509_cert_url"`
    ClientSecret string `json:"client_secret"`
    RedirectUris []string `json:"redirect_uris"`
}

type GdocApiObj  struct {
	GdocCtx context.Context
	tocFile string
	Svc *docs.Service
//	DrivSvc *drive.Service
	doc *docs.Document
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

// create text file to dump document file
// no need for mthod
func CreOutFile(filnam string, ext string)(outfil *os.File, err error) {
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
            return nil, fmt.Errorf("os.MkDir: %v", err1)
        }
//        newDir = true
    } else {
        if err != nil {
            return nil, fmt.Errorf("os.Stat %s: %v", outfildir, err)
        }
    }

	path:= fmt.Sprintf("%s/%s.%s", outfildir, outfilnam, ext)
//	fmt.Println("output file path: ", path)
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

// need to read add token file ref and scopes

func InitGdocApi(credFilNam string) (api *GdocApiObj, err error) {
// method that initializes the Gdoc Api by creating a service

	var gdoc GdocApiObj
	var svc *docs.Service

	var cred cred
    var config oauth2.Config

	ctx := context.Background()

	if len(credFilNam) < 1 {credFilNam = "/home/peter/go/src/google/gdoc/loginCred.json"}
	credbuf, err := os.ReadFile(credFilNam)
	if err != nil {return nil, fmt.Errorf("ReadFile %s: %v!", err)}

    err = json.Unmarshal(credbuf,&cred)
    if err != nil {return nil, fmt.Errorf("error unMarshal credbuf: %v\n", err)}

    if len(cred.Installed.ClientId) > 0 {
        config.ClientID = cred.Installed.ClientId
        config.ClientSecret = cred.Installed.ClientSecret
    }
    if len(cred.Web.ClientId) > 0 {
        config.ClientID = cred.Web.ClientId
        config.ClientSecret = cred.Web.ClientSecret
    }

    config.Scopes = make([]string,2)
    config.Scopes[0] = "https://www.googleapis.com/auth/drive"
    config.Scopes[1] = "https://www.googleapis.com/auth/documents"

    config.Endpoint = google.Endpoint

    tokFile := "tokNew.json"
    tok, err := tokenFromFile(tokFile)
    if err != nil {return nil, fmt.Errorf("error retrieving token: %v", err)}

    client := config.Client(context.Background(), tok)
	svc, err = docs.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {return nil, fmt.Errorf("Unable to retrieve Docs client: %v", err)}
	gdoc.Svc = svc

	return &gdoc, nil
}

func (gdoc *GdocApiObj) GetDoc(docId string) (err error) {

    svc := gdoc.Svc

    doc, err := svc.Documents.Get(docId).Do()
    if err != nil {return fmt.Errorf("document.Get: ", err)}
    gdoc.doc = doc
    return nil
}
