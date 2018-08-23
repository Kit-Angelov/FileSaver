package main

import (
	"fmt" 
	"net/http"
// "net/url"
//  "encoding/json"
	"os"
	// "io/ioutil"
 	"io"
 	"flag"
 	"./logger"
)

var (
	MEDIA_PATH string
	HOST       string = "127.0.0.1"
	PORT       string = "8080"
)

type test_struct struct {
	Test string
}


func init() {
	flag.StringVar(&MEDIA_PATH, "mp", MEDIA_PATH, "Путь до директории media")
	flag.StringVar(&HOST, "H", HOST, "Хост")
	flag.StringVar(&PORT, "P", PORT, "Порт")
}


func receiveHandler(w http.ResponseWriter, r *http.Request) {

	// the FormFile function takes in the POST input id file
	// r.ParseForm()

	name := r.FormValue("name")

	file, header, err := r.FormFile("file")

	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	defer file.Close()

	out, err := os.Create(fmt.Sprintf("%s.jpg", name))
	if err != nil {
		fmt.Fprintf(w, "Unable to create the file for writing. Check your write access privilege")
		return
	}

	defer out.Close()

	// write the content from POST to the file
	_, err = io.Copy(out, file)
	if err != nil {
		fmt.Fprintln(w, err)
	}

	fmt.Fprintf(w, "File uploaded successfully: ")
	fmt.Fprintf(w, header.Filename)
}


func main() {
	// init logger
	logger.Init(os.Stdout, os.Stdout, os.Stdout, os.Stderr, os.Stderr)
	logger.Info.Println("START APP")

	// parse args
	flag.Parse()
	logger.Info.Println("MEDIA PATH: ", MEDIA_PATH)
	logger.Info.Println("HOST: ", HOST)
	logger.Info.Println("PORT: ", PORT)

	// check exiting MEDIA PATH
	if _, err := os.Stat("/path/to/media"); os.IsNotExist(err) {
		logger.Raven.CaptureError(err, nil)
		logger.Fatal.Println(err)
		os.Exit(1)
	}
	// dir, err := os.Getwd()
	// if err != nil {
	// 	fmt.Println("err=", err)
	// 	os.Exit(1)
	// }
	// http.HandleFunc("/receive", receiveHandler) // Handle the incoming file
	// http.Handle("/", http.FileServer(http.Dir(dir)))
	logger.Error.Println(http.ListenAndServe(":8080", nil))
}