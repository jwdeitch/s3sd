package main

import (
	"path/filepath"
	"os"
	"fmt"
	"github.com/emirozer/go-helpers"
	"github.com/HouzuoGuo/tiedot/db"
	_ "github.com/HouzuoGuo/tiedot/dberr"
)

type file struct {
	fname string
	path  string
	size  int64
	md5   string
}

func main() {
	fileList := scanDir()
	db, _ := openDb()

	for _, file := range fileList {
		writeHash(db, file)
	}

}

func scanDir() *[]file {
	fileList := []file{}
	err := filepath.Walk("./", func(path string, f os.FileInfo, err error) error {

		if !f.IsDir() {
			md5Str, _ := helpers.Md5Hash(path)
			fileList = append(fileList, file{
				fname: f.Name(),
				path:  path,
				size:  f.Size(),
				md5:   md5Str})
		}
		return nil
	})
	check(err)

	return fileList
}

func openDb() (*db.Col, error) {
	myDBDir := "./s3sd"
	os.RemoveAll(myDBDir)
	defer os.RemoveAll(myDBDir)

	// (Create if not exist) open a database
	myDB, err := db.OpenDB(myDBDir)
	if err != nil {
		panic(err)
	}

	// Create collection if not exist
	if len(myDB.AllCols()) {
		return myDB.Use("s3sd")
	} else {
		if err := myDB.Create("s3sd"); err != nil {
			panic(err)
		}
		return myDB.Use("s3sd")
	}

	return err
}

func writeHash(db *db.Col, file file) {
	f, err := os.Create("./s3sd")
	check(err)

	defer f.Close()

	n3, err := f.WriteString("wrdwqdqwdqwdq\n")
	fmt.Printf("wrote %d bytes\n", n3)

	f.Sync()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}