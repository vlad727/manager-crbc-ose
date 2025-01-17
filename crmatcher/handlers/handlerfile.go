// Package handlers get post request and parse it
// create dir uploads and create file and copy data from uploaded file
package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	DstDirName string
)

func HandlePost(w http.ResponseWriter, r *http.Request) {

	// logging
	log.Println("Func HandlePost started ")
	// parse post request
	// The argument to FormFile must match the name attribute
	// of the file input on the frontend
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create the uploads folder if it doesn't exist
	err = os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// logging
	log.Println("Dir ./uploads has been created")

	// Create a new file in the uploads directory
	dst, err := os.Create(fmt.Sprintf("./uploads/%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Copy the uploaded file to the filesystem at the specified destination
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// set destination path to file to string
	DstDirName = dst.Name() // Note this variable is visible for all func because it's global var and another code in the same package, declared in handlerfile.go

	// redirect to page with description
	http.Redirect(w, r, "/crmatcherresult", http.StatusSeeOther)

}

// about ParseMultipartForm
// https://www.mohitkhare.com/blog/file-upload-golang/
// https://youtrack.jetbrains.com/issue/GO-13454/Unresolved-reference-Close-for-os.File error
// defer dst.Close() Unresolved reference 'Close'
