package models

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

var DB *sql.DB

type Doc struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Content string `json:"content"`
}

func Md(content string) string {
	input := []byte(content)
	unsafe := blackfriday.MarkdownCommon(input)
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	return string(html)
}

func ConnectDatabase() error {
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

func GetDocs(docType string) ([]Doc, error) {

	rows, err := DB.Query("SELECT name,type,content from doc WHERE type = ?", docType)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	docs := make([]Doc, 0)

	for rows.Next() {
		doc := Doc{}
		err = rows.Scan(&doc.Name, &doc.Type, &doc.Content)

		if err != nil {
			return nil, err
		}

		docs = append(docs, doc)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return docs, err
}

func GetDocByName(name string, docType string) ([]Doc, error) {

	rows, err := DB.Query("SELECT name,type,content from doc  WHERE type =  ? and name like ?", docType, "%"+name+"%")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	docs := make([]Doc, 0)

	for rows.Next() {
		doc := Doc{}
		err = rows.Scan(&doc.Name, &doc.Type, &doc.Content)

		if err != nil {
			return nil, err
		}

		docs = append(docs, doc)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return docs, err
}

func GetDocByCon(content string, docType string) ([]Doc, error) {

	rows, err := DB.Query("SELECT name,type,content from doc  WHERE type = ? and content like ?", docType, "%"+content+"%")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	docs := make([]Doc, 0)

	for rows.Next() {
		doc := Doc{}
		err = rows.Scan(&doc.Name, &doc.Type, &doc.Content)

		if err != nil {
			return nil, err
		}

		docs = append(docs, doc)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return docs, err
}

func AddDoc(newDoc Doc) (bool, error) {

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

func UpdateDoc(ourDoc Doc, name string) (bool, error) {

	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("UPDATE doc SET name = ?, type = ?, content = ?  WHERE name = ?")
	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(ourDoc.Name, ourDoc.Type, ourDoc.Content, name)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func DeleteDoc(name string, docType string) (bool, error) {

	tx, err := DB.Begin()

	if err != nil {
		return false, err
	}

	stmt, err := DB.Prepare("DELETE from doc where name = ? and type = ?")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(name, docType)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}
