package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	vision "cloud.google.com/go/vision/apiv1"
)

func visionScan() {
	crPath, _ := os.Getwd()
	filesPath := fmt.Sprintf("%v/files/", crPath)
	err := os.Chdir(filesPath)
	if err != nil {
		fmt.Printf("unable to change tmpDir path : %v\n", err)
	}
	var storePath string
	err = filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() && path != "." {
			storePath = path
			fmt.Printf("\nStarted scanning images of store: %v\n", storePath)
			return nil
		}
		if !info.IsDir() {
			imgfPath := path
			scanFile(storePath, imgfPath)
			return nil
		}
		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", filesPath, err)
		return
	}
}

func scanFile(storePath, imgfPath string) {
	ctx := context.Background()

	// Creates a client.
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Sets the name of the image file to annotate.
	file, err := os.Open(imgfPath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}
	defer file.Close()
	image, err := vision.NewImageFromReader(file)
	// _, err = vision.NewImageFromReader(file)
	if err != nil {
		log.Fatalf("Failed to create image: %v", err)
	}
	// TODO: #5 Vison APIをリクエスト時に、許可拒否エラーが発生する

	texts, err := client.DetectTexts(ctx, image, nil, 10)
	if err != nil {
		log.Fatalf("Failed to detect labels: %v", err)
	}
	var records []string
	for i, text := range texts {
		records[i] = text.Description
	}

	tmpSl := strings.Split(imgfPath, "/")
	tmpSl = strings.Split(tmpSl[len(tmpSl)-1], ".")
	newfPath := fmt.Sprintf("%v/%v.csv", storePath, tmpSl[0])
	csvFile, err := os.Create(newfPath)
	if err != nil {
		log.Fatalln("failed to create file", err)
	}
	defer csvFile.Close()
	w := csv.NewWriter(csvFile)
	defer w.Flush()
	// Using Write
	if err := w.Write(records); err != nil {
		log.Fatalln("error writing record to file", err)
	}
	fmt.Println("Created file:", newfPath)
}
