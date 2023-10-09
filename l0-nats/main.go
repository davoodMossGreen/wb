package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/davoodmossgreen/wb/l0-nats/models"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

func main() {
	fileContent, err := os.Open("./model.json")

	if err != nil {
	   log.Fatal(err)
	   return
	}

	defer fileContent.Close()

	byteResult, _ := ioutil.ReadAll(fileContent)

	var NewModel models.Model
	json.Unmarshal(byteResult, &NewModel)
	
	fmt.Println(NewModel)

	newModelBytes, _ := json.Marshal(NewModel)

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
    	log.Fatal(err)
	}
	defer nc.Close()



	sc, err := stan.Connect("test-cluster", "client", stan.NatsConn(nc))
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	sc.Publish("model", newModelBytes)
}