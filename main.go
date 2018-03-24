package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

// FileMetadata to create metadata json response
type FileMetadata struct {
	Size int64  `json:"size"`
	Name string `json"fileName"`
}

func main() {
	router := http.NewServeMux()

	css := http.FileServer(http.Dir("css"))
	router.Handle("/css/", http.StripPrefix("/css/", css))

	router.HandleFunc("/", homePage)
	router.HandleFunc("/submitForm", form)

	server := http.Server{
		Addr:    getPort(),
		Handler: router,
	}

	server.ListenAndServe()
}

// homePage function diplay index.html page
func homePage(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("index.html"))
	t.Execute(w, nil)
}

// form function get a file size and name from formFile and display it in a JSON format in response writer
func form(w http.ResponseWriter, r *http.Request) {
	var metaData FileMetadata
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	metaData.Size = fileHeader.Size
	metaData.Name = fileHeader.Filename
	fileMetadaJSON, err := json.Marshal(metaData)
	fmt.Fprint(w, string(fileMetadaJSON))
}

// get the rigth port based on the environment
func getPort() string {
	port := os.Getenv("PORT")

	if port == "" {
		log.Println("Running on port 8080 on localhost...")
		return ":8080"
	}
	log.Println("Running on port: " + port + "...")
	return ":" + port
}
