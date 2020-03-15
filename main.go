package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TextType struct {
	Text string `json:"text"`
}

var connection *mongo.Database

func connect(ctx context.Context, db string) (*mongo.Database, error) {
	uri := "mongodb://localhost:27017"
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("Couldn't connect to mongodb")
	}
	err = client.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("Mongo client couldn't connect with background ctx")
	}
	todoDB := client.Database(db)
	return todoDB, nil
}

func find(collection string, filter bson.D) []*TextType {
	var query *mongo.Cursor
	var err error
	var res []*TextType

	ctxt := context.Background()

	if filter != nil {
		query, err = connection.Collection(collection).Find(ctxt, filter)
	} else {
		query, err = connection.Collection(collection).Find(ctxt, bson.D{{}})
	}

	if err != nil {
		log.Fatal(err)
	}

	for query.Next(context.TODO()) {
		var elem TextType
		err := query.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		res = append(res, &elem)
	}

	if query.Err(); err != nil {
		log.Fatal(err)
	}

	query.Close(ctxt)
	return res

}

func send(collection string, text TextType) {
	insertRes, err := connection.Collection(collection).InsertOne(context.TODO(), text)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertRes.InsertedID)
}

func home(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.WriteHeader(http.StatusOK)
		file := r.URL.RequestURI()
		switch {
		case file == "/":
			w.Header().Set("Content-Type", "text/html")
			http.ServeFile(w, r, "./index.html")
			break
		case strings.HasSuffix(file, "html"):
			w.Header().Set("Content-Type", "text/html")
			http.ServeFile(w, r, "."+file)
			break
		case strings.HasSuffix(file, "css"):
			w.Header().Set("Content-Type", "text/css")
			http.ServeFile(w, r, "."+file)
			break
		case strings.HasSuffix(file, "js"):
			w.Header().Set("Content-Type", "application/js")
			http.ServeFile(w, r, "."+file)
			break
		default:
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"message":"Not Found"`))
			break
		}
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message":"Not Found"`))
	}
}

func input(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		data := find("text", nil)
		arr := ""
		for i, val := range data {
			if i+1 == len(data) {
				arr = arr + `"` + val.Text + `"`
			} else {
				arr = arr + `"` + val.Text + `"` + ","
			}
		}
		w.Write([]byte(`{"data":[` + arr + `]}`))
	case "POST":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body %v", err)
			http.Error(w, "Can't read body", http.StatusBadRequest)
			return
		}
		var text TextType
		json.Unmarshal(body, &text)
		send("text", text)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Successfully sent text"}`))
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message":"Not Found"`))
	}
}

func main() {
	ctxt := context.Background()
	// clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := connect(ctxt, "new_db")

	if err != nil {
		log.Fatal(err)
	}

	connection = client

	http.HandleFunc("/", home)
	http.HandleFunc("/input", input)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
