// googe drive library
// author prr
// 30/1/2022
// copywrite 2022 prr

package gdApiLib

import (
        "context"
        "encoding/json"
        "fmt"
//		"io"
        "io/ioutil"
//		"strings"
        "net/http"
        "os"
        "golang.org/x/oauth2"
        "golang.org/x/oauth2/google"
        "google.golang.org/api/drive/v3"
        "google.golang.org/api/option"
//		"google.golang.org/api/googleapi"
)

type gdApiObj  struct {
	Ctx context.Context
	Svc *drive.Service
}

type FileInfo struct {
	Id string
	MimeType string
	Name string
	Ext string
	ParentName string
	ParentId string
	SingleParent bool
	ModTime string
	Size int64
}

var Gapp = map[string]string {
	"gdoc": "application/vnd.google-apps.document",
	"gsheet": "application/vnd.google-apps.spreadsheet",
	"gdraw": "application/vnd.google-apps.drawing",
	"gscript": "application/vnd.google-apps.script",
	"photo": "application/vnd.google-apps.photo",
	"gslide": "application/vnd.google-apps.presentation",
	"gmap": "application/vnd.google-apps.map",
	"gform": "application/vnd.google-apps.form",
	"folder": "application/vnd.google-apps.folder",
	"file": "application/vnd.google-apps.file",
	"jpg": "image/jpeg",
	"png": "image/png",
	"svg": "image/svg+xml",
	"pdf": "application/pdf",
	"html": "text/html",
	"text": "text/plain",
	"rich": "application/rtf",
	"word": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	"excel": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	"csv": "text/csv",
	"ppt": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
}

// Retrieves a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	tokFile := "/home/peter/go/src/google/gdrive/tokGdrive.json"
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
func (gdrive *gdApiObj) CreDumpFile(fid string, filnam string)(err error) {

//	old_filenam := filnam

//	filExt := -1
//	var nfilnam []byte
	nfilnam := make([]byte,len(filnam), len(filnam)+5)
	for i:= len(filnam) -1; i > -1; i-- {
		nfilnam[i] = filnam[i]
		if filnam[i] == '.' {
			nfilnam[i] = '_'
		}
	}
	ext := "txt"

/*
	if filExt > 0 {
		ext = filnam[filExt+1:]
		filnam = filnam[:filExt]
	} else {
		if !(len(ext) >0) {
//			return nil, fmt.Errorf("error CreDumpFile: no file extension!")
		}
	}
*/

	for i:=0; i<len(nfilnam); i++ {
		if nfilnam[i] == ' ' {
			nfilnam[i] = '_'
		}
	}

	// check whether output directory exists
	filinfo, err := os.Stat("output")
	if os.IsNotExist(err) {
		return fmt.Errorf("error gdrive::CreDumpFile: sub-dir \"output\" does not exist!")
	}
	if err != nil {
		return fmt.Errorf("error gdrive::CreDumpFile: %v \n", err)
	}
	if !filinfo.IsDir() {
		return fmt.Errorf("error gdrive::CreDumpFile -- file \"output\" is not a directory!")
	}
	path:= "output/" + string(nfilnam) + "." + ext
//	fmt.Println("path: ",path)
	outfil, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("error gdrive::CreDumpFile: cannot open output file: %v \n", err)
	}

	// get file attributes
	svc := gdrive.Svc
	gfil, err := svc.Files.Get(fid).Do()
	if err != nil {
		return fmt.Errorf("error gdrive::CreDumpFile: cannot get file with id: %s! %v", fid, err)
	}

	outstr := fmt.Sprintf("File Name: %s Extension: %s Full Ext: %s\n", gfil.Name, gfil.FileExtension, gfil.FullFileExtension)
	outstr += fmt.Sprintf("Mime Type: %s Size: %d\n", gfil.MimeType, gfil.Size)
	outstr += fmt.Sprintf("File Id: %s Version: %d\n", gfil.Id, gfil.Version)
	outstr += fmt.Sprintf("Created: %s\n", gfil.CreatedTime)
	outstr += fmt.Sprintf("Modified: %s\n", gfil.ModifiedTime)
	outstr += fmt.Sprintf("Description: %s\n", gfil.Description)
	outstr += fmt.Sprintf("Original Name: %s \n", gfil.OriginalFilename)
	outstr += fmt.Sprintf("Parents: %d\n", len(gfil.Parents))
	outstr += fmt.Sprintf("Thumbnail: %s\n", gfil.ThumbnailLink)
	outstr += fmt.Sprintf("Web Content Link: %s\n", gfil.WebContentLink)
	outstr += fmt.Sprintf("Web View Link: %s\n", gfil.WebViewLink)

	outfil.WriteString(outstr)

	return nil
}


func InitDriveApi() (svc *drive.Service, err error) {
        ctx := context.Background()
//		gdObj.Ctx = ctx
        b, err := ioutil.ReadFile("/home/peter/go/src/google/gdrive/credGdrive.json")
        if err != nil {
			return nil, fmt.Errorf("Unable to read client secret file: %v!", err)
		}

        // If modifying these scopes, delete your previously saved token.json.
        config, err := google.ConfigFromJSON(b, drive.DriveScope)
        if err != nil {
			return nil, fmt.Errorf("Unable to parse client secret file to config: %v!", err)
        }

        client := getClient(config)

        svc, err = drive.NewService(ctx, option.WithHTTPClient(client))
        if err != nil {
			return nil, fmt.Errorf("Unable to retrieve Drive client: %v !", err)
        }
	return svc, nil
}

