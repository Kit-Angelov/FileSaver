package main

import (
	"fmt"
	"errors"
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
	MAIN_SERVER_HOST    string = "http://localhost:8000"
	API_NEW_DOC         string = "/api/v1/me/documents/new_doc/"
	API_TOKEN_VERIFY    string = "/rest-auth/api-token-verify/"
	MEDIA_PATH          string
	HOST                string = "127.0.0.1"
	PORT                string = "8080"
	SENTRY_URL          string = ""
	LOG_DIR             string = "./"
)

type tokenStruct struct {
	Token string
}

func init() {
	flag.StringVar(&MAIN_SERVER_HOST, "SH", MAIN_SERVER_HOST, "Хост главного сервера")
	flag.StringVar(&MEDIA_PATH, "mp", MEDIA_PATH, "Путь до директории media")
	flag.StringVar(&HOST, "H", HOST, "Хост")
	flag.StringVar(&PORT, "P", PORT, "Порт")
	flag.StringVar(&SENTRY_URL, "s", SENTRY_URL, "sentry url")
	flag.StringVar(&LOG_DIR, "log", LOG_DIR, "Директория для логов")
}

// notificate main server about new file
func notifiMainServerNewFile(token, guid, name, pathToFile string) bool {
	values := url.Values{}
	values.Set("guid", guid)
	values.Set("name", name)
	values.Set("path_to_file", pathToFile)

	body := values.Encode()

	client := &http.Client{}
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s%s", MAIN_SERVER_HOST, API_NEW_DOC), strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	res, err := client.Do(req)
	if err != nil {
		logger.Raven.CaptureError(err, nil)
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
	r, err := http.PostForm(fmt.Sprintf("%s%s", MAIN_SERVER_HOST, API_TOKEN_VERIFY), url.Values{"token": {token}})
	if err != nil {
		logger.Raven.CaptureError(err, nil)
		logger.Error.Println(err)
		return false
	}
	fmt.Println(r.StatusCode)
	if r.StatusCode != 200 {
		logger.Raven.CaptureError(errors.New("invalid token"), nil)
		logger.Error.Println("invalid token")
		return false
	}
	defer r.Body.Close()

	body_byte, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Raven.CaptureError(err, nil)
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
		fmt.Fprintln(w, "not verify")
		logger.Raven.CaptureError(errors.New("not verify"), nil)
		logger.Error.Println("not verify")
		return
	}

	name := r.FormValue("name")

	file, header, err := r.FormFile("file")

	if err != nil {
		fmt.Fprintln(w, err)
		logger.Raven.CaptureError(err, nil)
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
		logger.Raven.CaptureError(err, nil)
		logger.Error.Println(err)
		return
	}
	os.Chmod(pathToFile, 0777)

	defer out.Close()

	// write the content from POST to the file
	_, err = io.Copy(out, file)
	if err != nil {
		logger.Raven.CaptureError(err, nil)
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
	// parse args
	flag.Parse()
	// init logger
	logger.Init(LOG_DIR, SENTRY_URL)
	logger.Info.Println("START APP")

	logger.Info.Println("MEDIA PATH: ", MEDIA_PATH)
	logger.Info.Println("HOST: ", HOST)
	logger.Info.Println("PORT: ", PORT)
	logger.Info.Println("MAIN_SERVER_HOST: ", MAIN_SERVER_HOST)

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
	logger.Fatal.Println(http.ListenAndServe(fmt.Sprintf("%s:%s", HOST, PORT), nil))
}