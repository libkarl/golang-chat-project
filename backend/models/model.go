package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	
)


// definice vlastního datového typu do kterého se překlopí informace z databáze
// při načtení, je nutné co z json souboru bude co pro zařazení do datové struktury

type ToDoList struct {
	ID		primitive.ObjectID	`json:"_id,omitempty" bson: "_id,omitempty"`
	Task 	string				`json:"task,omitempty"`
	Status	bool				`json: "status,omitempty"`
}



