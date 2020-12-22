package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"

	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/import/list", ListHandler).Methods("GET")
	r.HandleFunc("/import", SyncHandler).Methods("POST")
	log.Println("server running on port 8899")
	log.Fatal(http.ListenAndServe(":8899", r))
}

// Projects ...
type Projects struct {
	Source  string   `json:"source"`
	Dest    string   `json:"dest"`
	Project []string `json:"project"`
}

// SyncHandler ..
func SyncHandler(w http.ResponseWriter, r *http.Request) {

	var projects Projects

	err := json.NewDecoder(r.Body).Decode(&projects)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	// cmd := exec.Command("rsync", "--recursive ", "--delete ", "--relative ", "--stats --human-readable", source, dest)
	cmd := exec.Command("rsync", "-Rr", projects.Source, projects.Dest)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	fmt.Printf("%s\n", stdoutStderr)
	json.NewEncoder(w).Encode(stdoutStderr)
}

// ListHandler ..
func ListHandler(w http.ResponseWriter, r *http.Request) {
	res, err := ReadDir("./")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	json.NewEncoder(w).Encode(res)
}

// ReadDir ...
func ReadDir(root string) ([]string, error) {
	var dir []string
	fileInfo, err := ioutil.ReadDir(root)
	if err != nil {
		return dir, err
	}
	for _, d := range fileInfo {
		if d.IsDir() && d.Name() != ".git" {
			dir = append(dir, d.Name())
		}
	}
	return dir, nil
}
