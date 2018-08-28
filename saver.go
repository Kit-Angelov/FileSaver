package main

import (
	"fmt"
	"strings"
	"net/http"
	"net/url"
	"encoding/json"
	"os"
 	"io"
 	"io/ioutil"
 	"flag"
 	"./logger"
 	"./utils"
 	"github.com/gorilla/mux"
)

var (
	SERVER_NEW_DOC      string = "http://192.168.2.243:8000/api/v1/me/documents/new_doc/"
	SERVER_TOKEN_VERIFY string = "http://192.168.2.243:8000/rest-auth/api-token-verify/"
	MEDIA_PATH          string
	HOST                string = "127.0.0.1"
	PORT                string = "8080"
	SENTRY_URL          string = ""
)

type tokenStruct struct {
	Token string
}

func init() {
	flag.StringVar(&MEDIA_PATH, "mp", MEDIA_PATH, "Путь до директории media")
	flag.StringVar(&HOST, "H", HOST, "Хост")
	flag.StringVar(&PORT, "P", PORT, "Порт")
	flag.StringVar(&SENTRY_URL, "s", "sentry url")
}

// notificate main server about new file
func notifiMainServerNewFile(token, guid, name, pathToFile string) bool {
	values := url.Values{}
	values.Set("guid", guid)
	values.Set("name", name)
	values.Set("path_to_file", pathToFile)

	body := values.Encode()

	client := &http.Client{}
	req, _ := http.NewRequest("POST", SERVER_NEW_DOC, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	res, err := client.Do(req)
	if err != nil {
		logger.Error.Println(err)
		return false
	}
	data, _ := ioutil.ReadAll(res.Body)
	logger.Info.Println(res.Status)
	logger.Info.Println(string(data))
	return true
}

// verification token with main server
func verifyToken(token string) bool {
	r, err := http.PostForm(SERVER_TOKEN_VERIFY, url.Values{"token": {token}})
	if err != nil {
		logger.Error.Println(err)
		return false
	}
	defer r.Body.Close()

	body_byte, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Error.Println(err)
		return false
	}
	resToken := tokenStruct{}
	json.Unmarshal(body_byte, &resToken)
	if resToken.Token == token {
		return true
	} else {
		return false
	}
}

func fileSaver(w http.ResponseWriter, r *http.Request) {
	token := strings.Split(r.Header.Get("Authorization"), " ")[1]
	verify := verifyToken(token)

	if verify == false {
		fmt.Fprintln(w, "invalid token")
		logger.Error.Println("invalid token")
		return
	}

	name := r.FormValue("name")

	file, header, err := r.FormFile("file")

	if err != nil {
		fmt.Fprintln(w, err)
		logger.Error.Println(err)
		return
	}

	if name == "" {
		name = header.Filename
	}

	defer file.Close()

	pathToFile, relPathToFile, guid := utils.GenPath(MEDIA_PATH, header.Filename)

	out, err := os.Create(pathToFile)
	if err != nil {
		fmt.Fprintln(w, err)
		logger.Error.Println(err)
		return
	}

	defer out.Close()

	// write the content from POST to the file
	_, err = io.Copy(out, file)
	if err != nil {
		logger.Error.Println(err)
		fmt.Fprintln(w, err)
		return
	}

	fmt.Fprintf(w, "File uploaded successfully: ")
	fmt.Fprintf(w, header.Filename)
	logger.Info.Printf("File uploaded successfully: %s", header.Filename)
	notifiResult := notifiMainServerNewFile(token, fmt.Sprintf("%s",guid), name, relPathToFile)
	logger.Info.Println(notifiResult)

}


func main() {
	// init logger
	logger.Init(SENTRY_URL, os.Stdout, os.Stdout, os.Stdout, os.Stderr, os.Stderr)
	logger.Info.Println("START APP")

	// parse args
	flag.Parse()
	logger.Info.Println("MEDIA PATH: ", MEDIA_PATH)
	logger.Info.Println("HOST: ", HOST)
	logger.Info.Println("PORT: ", PORT)

	// check exiting MEDIA PATH
	if _, err := os.Stat(MEDIA_PATH); os.IsNotExist(err) {
		logger.Raven.CaptureError(err, nil)
		logger.Fatal.Println(err)
		os.Exit(1)
	}

	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/upload/", fileSaver).Methods("POST") // Handle the incoming file 
	http.Handle("/", r)
	// http.Handle("/", http.FileServer(http.Dir(dir)))
	logger.Fatal.Println(http.ListenAndServe(fmt.Sprintf(":%s", PORT), nil))
}