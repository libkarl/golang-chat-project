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
 
func DeleteAllTasks(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	count := DeleteAll()
	json.NewEncoder(w).Encode(count)
}

func getAllTasks() []primitive.M {
	// do cur se uloží vše v kolekci pomocí funkce Find 
	// pokud dám funkci Find prázdný dotaz jak to zapsal níže,
	// znamená to v mongoDB, že chci vše co je v kolekci 
	// říká se tomu empty query
	// pokud bych chtěl konkrétní task, psal bych name = něco např.
	// decodovat = odhalit informaci skrytou uvnitř
	cur, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	// proměnná results je slice formátu M z bson primitiv
	// to dělá to to je neuspořídaná reprezentace Bson formatu
	var results []primitive.M
	// next posunuje mezi jednotlivými položkami v kolekci při
	// iteraci skrz cursor
	for cur.Next(context.Background()){
		var result bson.M
		// provede unmarshal předané položky do result
		// protože golang nemůže pracovat přímo s formatem jaký je
		// v mongoDB 
		// v případě problému uloží do e error
		e := cur.Decode(&result)
		if e != nil {
			log.Fatal(e)
		}
		// decodovaný výsledek přidá do results
		results = append(results, result)

		
	}

	// what the fuck is this.. ? 
	if err := cur.Err(); err != nil {
		log.Fatal(err)

	}

	// ukončení spojení s čím ? 
	cur.Close(context.Background())
	return results


}
 


func CompleteTask(task string) {
	// úkol který chci označit jako dokončený má nějaké ID
	// tohle id předávám jako filter do databáze 
	// vytáhnu ho pomocí var z requestu, je ve formátu Hex
	// takže ho překlopím do structury ObjectID
	id, _ := primitive.ObjectIDFromHex(task)
	// získané id typu ObjectID překlopím do typu primitive.M
	filter := bson.M{"_id":id}
	// změním stav u tohoto Id na true
	// napsaná kriteria dle kterých se upgraduje stav položky
	// v databazi
	update := bson.M{"$set":bson.M{"status":true}}
	// UpdateOne je funkce z MongoDB, dám jí kontext, podle čeho má vybrat co upgraduje
	// dám jí parametry změny a dostanu výsledek 
	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err !=nil {
		log.Fatal(err)
	}

	fmt.Println("modified count:", result)
}


func insertOneTask(task models.ToDoList) {
	insertResult, err := collection.InsertOne(context.Background(), task)

	if err !=nil {
		log.Fatal(err)
	}

	fmt.Println("Insert a single task: ", insertResult.InsertedID)
} 

func UndoTaskByID(task string){
	id,_ := primitive.ObjectIDFromHex(task)
	filter := bson.M{"_id":id}
	update := bson.M{"$set":bson.M{"status":false}}
	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err !=nil {
		log.Fatal(err)
	}
	// modified count na result vrátí počet modifikovaných dokumentů
	fmt.Println("The number of modified document", result.ModifiedCount)

}

func DeleteTaskByID(task string){
	id, _ := primitive.ObjectIDFromHex(task)
	filter := bson.M{"_id":id}
	d, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
	// vypíše počet odstraněných dokumentů
	fmt.Println("Deleted Document", d.DeletedCount)
}

func DeleteAll() int64{
	// pro smazání mnoha záznamů najednou je nutné použí funkci delete many
	d, err := collection.DeleteMany(context.Background(), bson.D{{}})
	if err !=nil {
		log.Fatal(err)
	}
	fmt.Println("Number of deleted document", d.DeletedCount)
	// funkce z mongo DB DeletedCount vrací výstup jako int64 
	// proto je výstup z funkce nastaven na int64
	return d.DeletedCount
}