package main

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

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
			fmt.Println(storePath)
		}
		if !info.IsDir() {
			filePath := path
			scanFile(storePath, filePath)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", filesPath, err)
		return
	}
}

func scanFile(storePath, filePath string) {
	ctx := context.Background()

	// Creates a client.
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Sets the name of the image file to annotate.
	filename := filePath

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}
	defer file.Close()
	image, err := vision.NewImageFromReader(file)
	if err != nil {
		log.Fatalf("Failed to create image: %v", err)
	}

	texts, err := client.DetectTexts(ctx, image, nil, 10)
	if err != nil {
		log.Fatalf("Failed to detect labels: %v", err)
	}
	// saveTexts(texts)
	fmt.Println("Start to print descriptions of texts for scanned file")
	for _, text := range texts {
		fmt.Print(text.Description, " ")
	}
	fmt.Println("Printing descriptions of texts is finished.")
	fmt.Println()
}
