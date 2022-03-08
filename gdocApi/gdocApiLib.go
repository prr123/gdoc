// v1 from rd_Gdocv3
// v2 added ext parameter to CreTxtOutFile
// v3
// 11 Jan 2022 added more error routines to Init

package gdocApiLib

import (
        "context"
        "encoding/json"
        "fmt"
        "io/ioutil"
//        "log"
        "net/http"
        "os"

        "golang.org/x/oauth2"
        "golang.org/x/oauth2/google"
        "google.golang.org/api/docs/v1"
        "google.golang.org/api/option"
)

type GdocApiStruct  struct {
	GdocCtx context.Context
	tocFile string
	Svc *docs.Service
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
func (gdoc *GdocApiStruct) CreTxtOutFile(filnam string, ext string)(outfil *os.File, err error) {

	// convert white spaces in file name to underscore
	bfil := []byte(filnam)
	for i:=0; i<len(filnam); i++ {
		if bfil[i] == ' ' {
			bfil[i] = '_'
		}
	}

    filinfo, err := os.Stat("output")
    if os.IsNotExist(err) {
        return nil, fmt.Errorf("error CreTxtOutFile: sub-dir \"output\" does not exist!")
    }
    if err != nil {
        return nil, fmt.Errorf("error CreTxtOutFile: %v \n", err)
    }
    if !filinfo.IsDir() {
        return nil, fmt.Errorf("error CreTxtOutFile -- file \"output\" is not a directory!")
    }

	path:= "output/" + string(bfil[:]) + "." +ext
	outfil, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return nil, fmt.Errorf("cannot open output text file: %v!", err)
	}
	return outfil, nil
}


func (gdoc *GdocApiStruct) Init() (err error) {
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
