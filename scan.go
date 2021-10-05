package main

import (
	"bufio"
	"context"
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
			return nil
		}
		if !info.IsDir() {
			imgfPath := path
			fmt.Printf("\nStarted ocr-scanning image: %v\n", imgfPath)
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
	annotations, err := client.DetectTexts(ctx, image, nil, 10)
	if err != nil {
		log.Fatalf("Failed to detect labels: %v", err)
	}

	tmpSl := strings.Split(imgfPath, "/")
	tmpSl = strings.Split(tmpSl[len(tmpSl)-1], ".")
	newfPath := fmt.Sprintf("%v/%v.txt", storePath, tmpSl[0])
	f, err := os.Create(newfPath)
	if err != nil {
		log.Fatalf("Failed to create file for saving sanned contents : %v", err)
	}
	defer f.Close()

	for _, annotation := range annotations {
		fw := bufio.NewWriter(f)
		_, err = fw.Write([]byte(annotation.Description))
		if err != nil {
			log.Fatalf("Failed to write sanned contents : %v", err)
		}
		err = fw.Flush()
		if err != nil {
			log.Fatalf("Failed to flush sanned contents to file : %v", err)
		}
	}
	fmt.Println("Created file:", newfPath)
}
