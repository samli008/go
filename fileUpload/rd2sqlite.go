package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

type doc struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Content string `json:"content"`
}

func main() {
	err := connectDatabase()
	checkErr(err)

	var folderPath string
	fmt.Print("pls input folder path: ")
	fmt.Scan(&folderPath)

	var fileType string
	fmt.Print("pls input file type: ")
	fmt.Scan(&fileType)

	rd2sqlite(folderPath, fileType)
}

func rd2sqlite(path string, fileType string) {
	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, info := range fileInfos {
		data, _ := ioutil.ReadFile(path + info.Name())
		var result doc
		result.Name = strings.Split(info.Name(), ".")[0]
		result.Type = fileType
		result.Content = string(data)
		addDoc(result)
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func connectDatabase() error {
	db, err := sql.Open("sqlite3", "./doc.db")
	if err != nil {
		return err
	}

	doc := `CREATE TABLE IF NOT EXISTS doc ( name text, type text, content text);`

	_, err = db.Exec(doc)
	if err != nil {
		log.Fatal(err)
	}

	DB = db
	return nil
}

func addDoc(newDoc doc) (bool, error) {

	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("INSERT INTO doc (name,type,content) VALUES (?,?,?)")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newDoc.Name, newDoc.Type, newDoc.Content)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}
