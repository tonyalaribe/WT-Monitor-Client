package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/kbinani/screenshot"
)

// Creates a new file upload http request with optional extra params
func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}

// Creates a new file upload http request with optional extra params
func newfileUploadRequestRawImage(uri string, params map[string]string, paramName string, image image.Image) (*http.Request, error) {

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, "shot.png")
	if err != nil {
		return nil, err
	}
	err = png.Encode(part, image)
	if err != nil {
		return nil, err
	}

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}

func takeShot() *image.RGBA {
	// Capture each displays.
	n := screenshot.NumActiveDisplays()
	if n <= 0 {
		panic("Active display not found")

	}

	var all image.Rectangle = image.Rect(0, 0, 0, 0)

	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)
		all = bounds.Union(all)

		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			panic(err)

		}
		fileName := fmt.Sprintf("%d_%dx%d.png", i, bounds.Dx(), bounds.Dy())
		save(img, fileName)

		fmt.Printf("#%d : %v \"%s\"\n", i, bounds, fileName)

	}

	// Capture all desktop region into an image.
	fmt.Printf("%v\n", all)
	img, err := screenshot.Capture(all.Min.X, all.Min.Y, all.Dx(), all.Dy())
	if err != nil {
		panic(err)

	}
	//save(img, "all.png")
	return img
}

func takeShotAndUpload(username string) {
	extraParams := map[string]string{
		"id": username,
	}
	img := takeShot()

	req, err := newfileUploadRequestRawImage("http://"+ROOT+"/api/screenshots", extraParams, "image", img)
	if err != nil {
		log.Println(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(resp.Status)
}
