package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func saveFlyImgsFrom(stores []Store) {
	for _, st := range stores {
		saveFlyImgsBy(st)
	}
}

func saveFlyImgsBy(st Store) {
	for i, fly := range st.flyers {
		crPath, _ := os.Getwd()
		dirPath := fmt.Sprintf("%v/files/%03d_%v", crPath, st.id, st.name)
		if !exists(dirPath) {
			err := os.MkdirAll(dirPath, 0777)
			if err != nil {
				panic(err)
			}
		}
		if i == 0 {
			err := deleteFilesUnder(dirPath)
			if err != nil {
				panic(err)
			}
		}
		filePath := fmt.Sprintf("%v/%03d%02d_%v_%v.jpg", dirPath, st.id, fly.id, fly.desc, time.Now().Format("20060102150405"))
		fmt.Println("\nStarted to save new image of flyer to new file")
		err := saveFlyImg(filePath, fly)
		if err != nil {
			panic(err)
		}
		err = os.Chmod(filePath, 0777)
		if err != nil {
			panic(err)
		}
	}
}

func saveFlyImg(filePath string, fly Flyer) error {
	// Reduce the burden on the target server
	time.Sleep(1 * time.Second)

	// Request the HTML page.
	fmt.Printf("Sent http Get request to : %v\n", fly.imgURL)
	res, err := http.Get(fly.imgURL)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	fmt.Printf("- The respose status of the Get request sent is : %v (%v)\n", res.StatusCode, res.Status)
	defer res.Body.Close()

	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)
	if err != nil {
		panic(err)
	}

	return err
}

func exists(dirPath string) bool {
	_, err := os.Stat(dirPath)
	return !os.IsNotExist(err)
}

func deleteFilesUnder(dir string) error {
	dirRead, err := os.Open(dir)
	if err != nil {
		return err
	}
	dirFiles, err := dirRead.ReadDir(0)
	if err != nil {
		return err
	}
	fmt.Println("\nStarted to delete old flyers per store")
	for i := range dirFiles {
		file := dirFiles[i]
		name := file.Name()
		fileName := fmt.Sprintf("%v/%v", dir, name)
		err := os.Remove(fileName)
		if err != nil {
			return err
		}
		fmt.Println("Deleted file:", fileName)
	}
	return err
}
