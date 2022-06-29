package middleware

import (
	"net/http"
	"log"
	"os"
	"fmt"
	"encoding/json"
	"context"
	"github.com/gorrila/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/joho/godotenv"
)

var collection *mongo.Collection

func init(){
	loadTheEnv()
	createDBInstance()
} 

func loadTheEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading environment")
	}
}

func createDBInstance() {
	connectionString := os.Getenv("DB_URL")
	dbName := os.Getenv("DB_NAME")
	collName := os.Getenv("DB_COLLECTION_NAME")
	// vytáhne nastavení klienta z mongo driverů a předá mu URL databáze
	clientOptions := options.Client().ApplyURI(connectionString)
	// vytvoří spojení s MongoDB databází, funkce Connect musí dostat kontext, který může definovat
	// například přerušení spojení po danném čase
	// zde je context.ToDo() bez této definice, který se
	// používá, když ještě nevíš jaký kontext tam bude + předávám
	// nastavení clienta s DB_URL mé databáze
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// testuje spojení mezi serverem a databází
	// lze do kontextu nastavit timeout jak dlouho má čekat
	// než vyhodí error
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	// pokud se vše podaří 
	fmt.Println("Spojení s databází navázáno")
	collection = client.Database(dbName).Collection(collName)
	fmt.Println("Vytvořeno spojení s kolekcí")
}

func GetAllTasks(w http.ResponseWriter, r *http.Request) { 
	// It sets header for request on data inside the database
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	payload := getAllTasks()
	// zakoduje JSON v payload a zapíše to do odpovědi
	json.NewEncoder(w).Encode(payload)

}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	return nil
}

func TaskComplete() error {
	return nil
}

func UndoTask() error {
	return nil
}

func DeleteTask() error {
	return nil
}

func deleteAllTasks() error {
	return nil
}