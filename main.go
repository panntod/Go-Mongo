package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	collection *mongo.Collection
	ctx        = context.Background()
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func initMongoDb() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected MongoDB!")
	collection = client.Database("gomongo").Collection("users")
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	fmt.Fprintf(w, `{"error": "%s"}`, message)
}

func getUserByID(id string) (User, error) {
	var user User

	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func getUsers() ([]User, error) {
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var users []User
	for cur.Next(ctx) {
		var user User
		err := cur.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		id := r.URL.Query().Get("id")
		if id != "" {
			user, err := getUserByID(id)
			if err != nil {
				respondWithError(w, http.StatusNotFound, "User not found")
				return
			}
			json.NewEncoder(w).Encode(user)
		} else {
			users, err := getUsers()
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, "Internal server error")
				return
			}
			json.NewEncoder(w).Encode(users)
		}
	case "POST":
		var newUser User
		err := json.NewDecoder(r.Body).Decode(&newUser)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		newUser.ID = primitive.NewObjectID().Hex()

		_, err = collection.InsertOne(ctx, newUser)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to create user")
			return
		}

		json.NewEncoder(w).Encode(newUser)
	}
}

func main() {
	initMongoDb()

	http.HandleFunc("/users", handleUsers)

	fmt.Println("Server is running on 8080...")
	http.ListenAndServe(":8080", nil)
}
