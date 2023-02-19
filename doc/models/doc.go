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

	doc := `CREATE TABLE IF NOT EXISTS doc ( name text, content text);`

	_, err = db.Exec(doc)
	if err != nil {
		log.Fatal(err)
	}

	DB = db
	return nil
}

func GetDocs() ([]Doc, error) {

	rows, err := DB.Query("SELECT name,content from doc")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	docs := make([]Doc, 0)

	for rows.Next() {
		doc := Doc{}
		err = rows.Scan(&doc.Name, &doc.Content)

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

func GetDocByName(name string) ([]Doc, error) {

	rows, err := DB.Query("SELECT name,content from doc  WHERE name like ?", "%"+name+"%")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	docs := make([]Doc, 0)

	for rows.Next() {
		doc := Doc{}
		err = rows.Scan(&doc.Name, &doc.Content)

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

	stmt, err := tx.Prepare("INSERT INTO doc (name,content) VALUES (?,?)")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newDoc.Name, newDoc.Content)

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

	stmt, err := tx.Prepare("UPDATE doc SET name = ?, content = ?  WHERE name = ?")
	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(ourDoc.Name, ourDoc.Content, name)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func DeleteDoc(name string) (bool, error) {

	tx, err := DB.Begin()

	if err != nil {
		return false, err
	}

	stmt, err := DB.Prepare("DELETE from doc where name = ?")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(name)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}
