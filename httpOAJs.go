package main

import (
	"fmt"
//	"io/ioutil"
	"net/http"
	"os"
	"encoding/json"
//	"google.golang.org/api/oauth2/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
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
//	Auth_provider_x509_cert_url string `json:"auth_provider_x509_cert_url"`
	ClientSecret string `json:"client_secret"`
	RedirectUris []string `json:"redirect_uris"`
}

type OA2 struct {
		RedirectURL 	string
		ClientId 		string
		ClientSecret 	string
		Scopes 			[]string
		Endpoint 		string
		oauth 			string
		state 			int
		authCode 		string
		token 			*oauth2.Token
}


func initCred()(oa2Obj *OA2, err error) {

	var cred cred
	var oa2 OA2

	credfil, err := os.Open("loginCred.json")
	if err != nil {fmt.Printf("error opening cred file: %v\n", err); os.Exit(-1);}

	info, err := credfil.Stat()
	if err != nil {fmt.Printf("error stat of cred file: %v\n", err); os.Exit(-1);}

	size := int(info.Size())

	credbuf := make([]byte, size)

	_, err = credfil.Read(credbuf)
	if err != nil {return nil, fmt.Errorf("error read cred file: %v\n", err)}

//fmt.Printf("cred : %s\n", credbuf)

	err = json.Unmarshal(credbuf,&cred)
	if err != nil {return nil, fmt.Errorf("error unMarshal cred: %v\n", err)}

	oa2 = OA2{
		RedirectURL:  "http://localhost:8090/callback",
		ClientId: "",
		ClientSecret: "",
		Scopes: 	[]string{"https://www.googleapis.com/auth/drive", "https://www.googleapis.com/auth/documents"},
//		Endpoint:     google.Endpoint,
		state: 0,
		oauth: "2kjd9nyuamtUAnwF34NQ12",
	}


	if len(cred.Installed.ClientId) > 0 {
		fmt.Printf("installed: \n")
		fmt.Printf("client secret: %s\n", cred.Installed.ClientSecret)
		fmt.Printf("client id:      %s\n", cred.Installed.ClientId)
		oa2.ClientId = cred.Installed.ClientId
		oa2.ClientSecret = cred.Installed.ClientSecret
	}
	if len(cred.Web.ClientId) > 0 {
		fmt.Printf("web:\n")
		fmt.Printf("client secret: %s\n", cred.Web.ClientSecret)
		fmt.Printf("client id:      %s\n", cred.Web.ClientId)
		oa2.ClientId = cred.Web.ClientId
		oa2.ClientSecret = cred.Web.ClientSecret
	}

	return &oa2, nil
}

func main() {
	fmt.Println("*** reading credentials file ***")
	oa2Obj, err := initCred()
	if err != nil {fmt.Printf("error initCred: %v\n", err); os.Exit(-1)}

	fmt.Printf("oa2Obj:\n%v\n\n", oa2Obj)
//	os.Exit(1)
//	http.HandleFunc("/", handleMain)
	http.HandleFunc("/", oa2Obj.handleGSignin)
//	http.HandleFunc("/login", oa2Obj.handleLogin)
	http.HandleFunc("/googleIdCallback", oa2Obj.handleGoogleIdCallback)
	http.HandleFunc("/callback", handleCallback)
	fmt.Println(http.ListenAndServe(":8090", nil))
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	var htmlIndex = `<html>
<body>
 <h2>Success</h2>
</body>
</html>`

	fmt.Fprintf(w, htmlIndex)
}


func (oa2 *OA2)handleGSignin(w http.ResponseWriter, r *http.Request) {
	var htmlIndex = `<html>
<head>
<meta name="referrer" content="no-referrer-when-downgrade" />
<script>
async function fetchAsync (url) {
  	let authResponse = await fetch(url);
  	let authData = await authResponse.json();
  	return authData;
}

function handleCredResp (response) {
  let cred = response.credential;
  let reslen = Object.keys(response).length;
  console.log("response: [" + reslen  + "]");
	keys = Object.keys(response);
//	for (i=0; i<reslen; i++) {
//		console.log("key [" + i + "]: " + keys[i])
//	}
	console.log("response clientId: " + response.clientId);
	console.log("response client_id: " + response.client_id);
	console.log("response select_by: " + response.select_by);

  let creList = cred.split(".");
  console.log("cred list: " + creList.length);
  if (creList.length != 3) {throw new Error('credential token not valid!');}
  let credDecoded= atob(creList[1]);

  console.log("response credential: "+ credDecoded);

//JSON.parse(base64_url_decode(token.split(".")[pos]));

  const credObj = JSON.parse(credDecoded)

  let credlen = Object.keys(credObj).length;
  let credKeys = Object.keys(response);
	let entries = Object.entries(credObj);

Object.entries(credObj).forEach(([key, value]) => {
  console.log("key: " + key + " value: " + value);
});

//	url = "https://accounts.google.com/o/oauth2/v2/auth?scope=https://www.googleapis.com/auth%2Fdrive.file+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fdocuments";
	url = "https://accounts.google.com/o/oauth2/v2/auth?scope=https://www.googleapis.com/auth/drive+https://www.googleapis.com/auth/documents";
	url1 = "&response_type=code&state=2kjd9nyuamtUAnwF34NQ12&redirect_uri=http://localhost:8090/googleIdCallback&client_id=962131291489-649afk49rj2v4mhjesvnvvf6fu43fi6m.apps.googleusercontent.com";
	url = url + url1;
	console.log("url: "+ url);
	location.replace(url)
//	data = fetchAsync(url);
//	console.log("replaced url")
}
</script>
</head>
<body>
<script src="https://accounts.google.com/gsi/client" async defer></script>
<script>
//        function handleCredentialResponse(response) {
//          console.log("Encoded JWT ID token: " + response.credential);
//        }
        window.onload = function () {
          google.accounts.id.initialize({
            client_id: "962131291489-649afk49rj2v4mhjesvnvvf6fu43fi6m.apps.googleusercontent.com",
            callback: handleCredResp
          });
          google.accounts.id.renderButton(document.getElementById("buttonDiv"),
            { theme: "outline", size: "large" }  // customization attributes
          );
          google.accounts.id.prompt(); // also display the One Tap dialog
        }
    </script>
    <div id="buttonDiv"></div>
</body>
</html>`

	enableCors(&w)

	fmt.Fprintf(w, htmlIndex)
}

func (oa2 *OA2) handleReqAuth(w http.ResponseWriter, r *http.Request) {

fmt.Printf("\nreceived login request\n")

	fmt.Printf("\nloginform:\n%v\n", r)


//	url := googleOauthConfig.AuthCodeURL(oa2.oauth)

	url := "https://accounts.google.com/o/oauth2/v2/auth?scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fdrive.file+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fdocuments"
	url1 := fmt.Sprintf("&response_type=code&state=%s&redirect_uri=http%3A%2F%2Flocalhost%3A8090&client_id=%s",oa2.oauth, oa2.ClientId)
fmt.Println("\ngoogle login: ", url1)
return
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

/*
func (oa2 *OA2)handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\ngoogle call back\n state: %s code: %s\n", r.FormValue("state"), r.FormValue("code"))
	fmt.Printf("\ncallback form:\n%v\n", r)
	r.ParseForm()
	for key, val := range r.Form{
		fmt.Printf("%s = %s\n", key, val)
	}

	fmt.Println("entering user info")
	content, err := oa2.getUserInfo(r.FormValue("state"), r.FormValue("code"))
	if err != nil {
		fmt.Printf("error getUserInfo: %s, %v\n", string(content), err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	fmt.Fprintf(w, "Content: %s\n", content)
}
*/

func (oa2 *OA2)handleGoogleIdCallback(w http.ResponseWriter, r *http.Request) {
//	fmt.Printf("\ngoogle call back\n state: %s code: %s\n", r.FormValue("state"), r.FormValue("code"))
	fmt.Printf("\nSign-in callback form:\n%v\n",r)
	r.ParseForm()
	fmt.Printf("form key value\n")
	for key, val := range r.Form{
		fmt.Printf("%s = %s\n", key, val)
	}
//	fmt.Fprintf(w, "Content: %s\n", content)
	oa2.authCode = r.Form["code"][0]

	state := r.Form["state"][0]
	fmt.Printf("auth code: %s state: %s\n", oa2.authCode, state)

//	if state != oa2.oauth {return fmt.Errorf("invalid oauth state")}

	err := oa2.getToken()
	if err != nil {
		fmt.Printf("error getToken: %v\n", err)
		fmt.Fprintf(w, "cannot get token\n")
	}
	fmt.Printf("received token:\n%v\n",oa2.token)
}

func (oa2 *OA2) getToken() (err error) {

	var config oauth2.Config

	config.ClientID = oa2.ClientId
	config.ClientSecret = oa2.ClientSecret
	config.Endpoint = google.Endpoint
//	config.RedirectURL = "http://localhost:8090/"
	config.RedirectURL = "http://localhost:8090/googleIdCallback"
	config.Scopes = make([]string, len(oa2.Scopes))
	for i:=0; i< len(oa2.Scopes); i++ {
		config.Scopes[i] = oa2.Scopes[i]
	}

	fmt.Printf("config:\n%v\n", config)
	fmt.Println("code exchange")

	code:= oa2.authCode
	token, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Println("failed exchange!")
		return fmt.Errorf("code exchange failed: %s", err.Error())
	}
	fmt.Printf("token: %s\n", token)
	oa2.token = token
/*
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}
*/
//	fmt.Printf("response: %s\n",contents)
	return nil
}


func enableCors(w *http.ResponseWriter) {
(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
