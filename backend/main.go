package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/libkarl/golang-chat-project/backend/router"
)

func main(){
	r := router.Router()
	fmt.Println("starting the server on port 9000..") 

	log.Fatal(htpp.ListenAndServe(":9000", r))
	
}