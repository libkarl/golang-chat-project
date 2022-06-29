package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/libkarl/golang-chat-project/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	payload := getAllTasks() 
	// zakoduje JSON v payload a zapíše to do odpovědi
	json.NewEncoder(w).Encode(payload)

}

// posílám Task jako POST request do této funkce a používám NewDecoder
// k dekodování JSON těla uvnitř požadavku, po dekodování se tělo uloží do
// dle vytvořené struktury ToDoList
func CreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	var task models.ToDoList
	json.NewDecoder(r.Body).Decode(&task)
	insertOneTask(task)
	// zase ho Encoduje a výsledek uloží do odpovědi
	json.NewEncoder(w).Encode(task)
}

func TaskComplete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	// použití balíčku mux pro získání parametrů a uložení do proměnné
	// protože potřebujeme ID tasku, které bude předáno 
	params := mux.Vars(r)
	CompleteTask(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

func UndoTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params := mux.Vars(r) 
	UndoTaskByID(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params := mux.Vars(r)
	DeleteTaskByID(params["id"])

}
 
func deleteAllTasks(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	count := DeleteAll()
	json.NewEncoder(w).Encode(count)
}

func getAllTasks() error {
	return nil
}
 
func insertOneTask(task models.ToDoList) {

}

func CompleteTask(params map[string]string) {

}

func DeleteTaskByID(){

}

func DeleteAll() {
	
}