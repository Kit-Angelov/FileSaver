package main
import (
  "net/http"
  "strings"
  "fmt"
  "io"
  "os"
)
func uploadHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method == "GET" {
    fmt.Fprintf(w, `<html>
<head>
  <title>GoLang HTTP Fileserver</title>
</head>
<body>
<h2>Upload a file</h2>
<form action="/receive" method="post" enctype="multipart/form-data">
  <label for="file">Filename:</label>
  <input type="file" name="file" id="file">
  <br>
  <input type="submit" name="submit" value="Submit">
</form>
</body>
</html>`)
  }
}
func receiveHandler(w http.ResponseWriter, r *http.Request) {

  // the FormFile function takes in the POST input id file
  file, header, err := r.FormFile("file")

  if err != nil {
    fmt.Fprintln(w, err)
    return
  }

  defer file.Close()

  out, err := os.Create(header.Filename)
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
func sayHello(w http.ResponseWriter, r *http.Request) {
  message := r.URL.Path
  message = strings.TrimPrefix(message, "/")
  message = "Hello " + message
  w.Write([]byte(message))
}
func main() {
  dir, err := os.Getwd()
  if err != nil {
    fmt.Println("err=", err)
    os.Exit(1)
  }
  http.HandleFunc("/upload", uploadHandler)   // Display a form for user to upload file
  http.HandleFunc("/receive", receiveHandler) // Handle the incoming file
  http.Handle("/", http.FileServer(http.Dir(dir)))
  if err := http.ListenAndServe(":8080", nil); err != nil {
    panic(err)
  }
}