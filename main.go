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

func respondWithError(response http.ResponseWriter, code int, message string) {
	response.WriteHeader(code)
	fmt.Fprintf(response, `{"error": "%s"}`, message)
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

func updateUserByID(id string, updateUser User) (User, error) {
	filter := bson.M{"id": id}
	update := bson.M{
		"$set": bson.M{
			"name": updateUser.Name,
			"age":  updateUser.Age,
		},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return User{}, err
	}

	return updateUser, nil
}

func handleUsers(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	switch request.Method {
	case "GET":
		id := request.URL.Query().Get("id")
		if id != "" {
			user, err := getUserByID(id)
			if err != nil {
				respondWithError(response, http.StatusNotFound, "User not found")
				return
			}
			json.NewEncoder(response).Encode(user)
		} else {
			users, err := getUsers()
			if err != nil {
				respondWithError(response, http.StatusInternalServerError, "Internal server error")
				return
			}
			json.NewEncoder(response).Encode(users)
		}
	case "POST":
		var newUser User
		err := json.NewDecoder(request.Body).Decode(&newUser)
		if err != nil {
			respondWithError(response, http.StatusBadRequest, "Invalid request payload")
			return
		}

		newUser.ID = primitive.NewObjectID().Hex()

		_, err = collection.InsertOne(ctx, newUser)
		if err != nil {
			respondWithError(response, http.StatusInternalServerError, "Failed to create user")
			return
		}

		json.NewEncoder(response).Encode(newUser)
	case "PUT":
		id := request.URL.Query().Get("id")
		if id != "" {
			var updatedUser User
			err := json.NewDecoder(request.Body).Decode(&updatedUser)
			if err != nil {
				respondWithError(response, http.StatusBadRequest, "Invalid request payload")
				return
			}

			updatedUser.ID = id

			user, err := updateUserByID(id, updatedUser)
			if err != nil {
				respondWithError(response, http.StatusInternalServerError, "Failed to update user")
				return
			}

			json.NewEncoder(response).Encode(user)
		} else {
			respondWithError(response, http.StatusBadRequest, "ID is required for updating user")
			return
		}
	default:
		respondWithError(response, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func main() {
	initMongoDb()

	http.HandleFunc("/users", handleUsers)

	fmt.Println("Server is running on 8080...")
	http.ListenAndServe(":8080", nil)
}
