package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func songUpload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("songUpload called.")
	var Buf bytes.Buffer
	// in your case file would be fileupload
	// fmt.Println(r)
	file, header, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	fmt.Println("step 1")
	fmt.Printf("File name %s\n", header.Filename)
	// Copy the file data to my buffer
	io.Copy(&Buf, file)
	fmt.Println("step 2")
	// do something with the contents...
	// I normally have a struct defined and unmarshal into a struct, but this will
	// work as an example
	// contents := Buf.String()
	err = ioutil.WriteFile(header.Filename, Buf.Bytes(), 0644)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(contents)
	// I reset the buffer in case I want to use it again
	// reduces memory allocations in more intense projects
	Buf.Reset()
	// do something else
	// etc write header
	return

}

// func downloadFile(filepath string, url string) error {

// 	// Get the data
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	// Create the file
// 	out, err := os.Create(filepath)
// 	if err != nil {
// 		return err
// 	}
// 	defer out.Close()

// 	// Write the body to file
// 	_, err = io.Copy(out, resp.Body)
// 	return err
// }
