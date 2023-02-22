package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	//print databases
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)

	//print document where type is linux.
	collection := client.Database("login").Collection("docs")
	ctx, _ = context.WithTimeout(context.Background(), 30*time.Second)

	cur, err := collection.Find(ctx, bson.M{"type": "linux"})
	if err != nil {
		log.Fatal("find:" + err.Error())
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var result doc
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Printf("%s\n", result.Name)
		addDoc(result)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
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
