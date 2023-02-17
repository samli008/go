package models

import (
	"database/sql"
	"log"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

type Person struct {
	Fso     string `json:"fso"`
	Machine string `json:"machine"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	City    string `json:"city"`
	Address string `json:"address"`
}

func ConnectDatabase() error {
	db, err := sql.Open("sqlite3", "./netapp.db")
	if err != nil {
		return err
	}

	netapp := `CREATE TABLE IF NOT EXISTS netapp (fso text, machine text, name text, phone text, email text,city text, address text);`

	_, err = db.Exec(netapp)
	if err != nil {
		log.Fatal(err)
	}

	DB = db
	return nil
}

func GetPersons(count int) ([]Person, error) {

	rows, err := DB.Query("SELECT fso,machine,name,phone,email,city,address from netapp LIMIT " + strconv.Itoa(count))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	people := make([]Person, 0)

	for rows.Next() {
		singlePerson := Person{}
		err = rows.Scan(&singlePerson.Fso, &singlePerson.Machine, &singlePerson.Name, &singlePerson.Phone, &singlePerson.Email, &singlePerson.City, &singlePerson.Address)

		if err != nil {
			return nil, err
		}

		people = append(people, singlePerson)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return people, err
}

func GetPersonByName(fso string) ([]Person, error) {

	rows, err := DB.Query("SELECT fso,machine,name,phone,email,city,address from netapp  WHERE fso like ?", "%"+fso+"%")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	people := make([]Person, 0)

	for rows.Next() {
		singlePerson := Person{}
		err = rows.Scan(&singlePerson.Fso, &singlePerson.Machine, &singlePerson.Name, &singlePerson.Phone, &singlePerson.Email, &singlePerson.City, &singlePerson.Address)

		if err != nil {
			return nil, err
		}

		people = append(people, singlePerson)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return people, err
}

func AddPerson(newPerson Person) (bool, error) {

	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("INSERT INTO netapp (fso,machine,name,phone,email,city,address) VALUES (?,?,?,?,?,?,?)")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newPerson.Fso, newPerson.Machine, newPerson.Name, newPerson.Phone, newPerson.Email, newPerson.City, newPerson.Address)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func UpdatePerson(ourPerson Person, fso string) (bool, error) {

	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("UPDATE netapp SET fso = ?, machine = ?, name = ?, phone = ?, email = ?, city = ?, address = ?  WHERE fso = ?")
	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(ourPerson.Fso, ourPerson.Machine, ourPerson.Name, ourPerson.Phone, ourPerson.Email, ourPerson.City, ourPerson.Address, fso)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func DeletePerson(fso string) (bool, error) {

	tx, err := DB.Begin()

	if err != nil {
		return false, err
	}

	stmt, err := DB.Prepare("DELETE from netapp where fso = ?")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(fso)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}
