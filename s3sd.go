package main

import (
	"path/filepath"
	"os"
	"encoding/json"
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
	db := openDb()

	for _, file := range fileList {
		detectChange(db, file)
	}

}

func detectChange(s3sdDb *db.Col, file file) {
	queryResult := make(map[int]struct{}) // query result (document IDs) goes into map keys

	if err := db.EvalQuery(query, feeds, &queryResult); err != nil {
		panic(err)
	}

	//if err != nil {
	//	writeHash(db, file)
	//}
	//
	//fmt.Println(fileInDb)

}

func scanDir() []file {
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

func openDb() (*db.Col) {
	myDBDir := "./s3sd"

	// (Create if not exist) open a database
	myDB, err := db.OpenDB(myDBDir)
	if err != nil {
		panic(err)
	}

	// Create collection if not exist
	if len(myDB.AllCols()) > 0 {
		return myDB.Use("s3sd")
	} else {
		if err := myDB.Create("s3sd"); err != nil {
			panic(err)
		}
		return myDB.Use("s3sd")
	}

	//return fmt.Errorf("Could not connect to database")
}

func writeHash(s3sdDb *db.Col, file file) {
	docID, err := s3sdDb.Insert(map[string]interface{}{
		"md5" : file.md5,
		"path" : file.path,
		"size" : file.size,
		"name" : file.fname})
	if err != nil {
		panic(err)
	}

	fmt.Println(docID)

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}