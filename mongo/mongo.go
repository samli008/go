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

// var db *mongo.Database

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
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
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

	var db string
	fmt.Printf("current mongodb databases: %s\n", databases)
	fmt.Println("pls input export mongodb database name: ")
	fmt.Scanln(&db)

	collections, err := client.Database(db).ListCollectionNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	var col string
	fmt.Printf("%s database collections: %s\n", db, collections)
	fmt.Printf("pls input collection in %s database: ", db)
	fmt.Scanln(&col)

	//print document where type is linux.
	collection := client.Database(db).Collection(col)
	ctx, _ = context.WithTimeout(context.Background(), 30*time.Second)

	docTypes, _ := collection.Distinct(ctx, "type", bson.M{})
	var docType string
	fmt.Printf("%s types: %s\n", col, docTypes)
	fmt.Printf("pls input export %s type: ", col)
	fmt.Scanln(&docType)

	cur, err := collection.Find(ctx, bson.M{"type": docType})
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
