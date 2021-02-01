package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Image struct {
	Code    string
	Content int
}

var router = mux.NewRouter()

func main() {
	port := os.Getenv("PORT")
	mongouri := os.Getenv("MongoDB")
	clientOptions := options.Client().ApplyURI(mongouri)
	client, e := mongo.Connect(context.TODO(), clientOptions)
	CheckError(e)
	log.Println("SUCCESSFULLY CONNECTED")
	e = client.Ping(context.TODO(), nil)
	CheckError(e)
	log.Println("CONNECTION VALID")
	collection := client.Database("neptune").Collection("images")
	log.Println("DATABASE CONNECTED")
	router.HandleFunc("/image/{code}", func(w http.ResponseWriter, r *http.Request) {
		code := mux.Vars(r)["code"]
		filter := bson.D{{"code", code}}
		var res Image
		e := collection.FindOne(context.TODO(), filter).Decode(&res)
		CheckError(e)
		if res.Content != 0 {
			var ares = strconv.Itoa(res.Content)

			io.WriteString(w, ares)
		} else {
			io.WriteString(w, "no image found")
		}
	})
	log.Println("SERVER LISTENING")
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}
func CheckError(e error) {
	if e != nil {
		fmt.Println(e)
	}
}
