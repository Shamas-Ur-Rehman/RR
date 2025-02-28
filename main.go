package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

type Registration struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title        string             `bson:"title" json:"title"`
	FirstName    string             `bson:"firstName" json:"firstName"`
	LastName     string             `bson:"lastName" json:"lastName"`
	CompanyName  string             `bson:"companyName" json:"companyName"`
	JobTitle     string             `bson:"jobTitle" json:"jobTitle"`
	Email        string             `bson:"email" json:"email"`
	ConfirmEmail string             `bson:"confirmEmail" json:"confirmEmail"`
	City         string             `bson:"city" json:"city"`
	Mobile       string             `bson:"mobile" json:"mobile"`
	Country      string             `bson:"country" json:"country"`
	Nationality  string             `bson:"nationality" json:"nationality"`
}

func main() {
	r := gin.Default()
	r.Use(corsMiddleware())

	clientOptions := options.Client().ApplyURI("mongodb+srv://shamasurrehman509:LqqXCkGoS6WNLXxP@cluster0.2ttxi.mongodb.net/")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	collection = client.Database("registrationDB").Collection("registrations")

	r.POST("/register", registerUser)
	r.GET("/registrations", getRegistrations)

	// Start server
	r.Run(":8080")
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173") // Allow frontend
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// Register user
func registerUser(c *gin.Context) {
	var user Registration
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	user.ID = primitive.NewObjectID()
	_, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save data"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Registration successful"})
}

// Get all registrations
func getRegistrations(c *gin.Context) {
	var users []Registration
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data"})
		return
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var user Registration
		if err := cursor.Decode(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding data"})
			return
		}
		users = append(users, user)
	}
	c.JSON(http.StatusOK, users)
}
