package main

import (
	"fmt"
	"log"
	"os"
	"time"

	gen "github.com/hieutrtr/generator"
)

type Ads struct {
	title string `gentype:"varchar"`
	ad_id uint8  `gentype:"smallint"`
}

var numEvents = 10000
var numWorkers = 100

var topic = os.Getenv("KAFKA_TOPIC")

func producer(sup <-chan string, res chan<- error) {
	pro, err := gen.NewProducer()
	if err != nil {
		panic(err.Error())
	}

	for q := range sup {
		e := &gen.UploadEvent{
			Topic:   topic,
			Payload: q,
		}

		err = pro.Produce(e)
		if err != nil {
			res <- fmt.Errorf("%s %s\n", q, err.Error())
		}
		res <- nil
	}
}

func generate(st interface{}, num int, sup chan<- string) {
	for i := 0; i < num; i++ {
		sup <- gen.GenJSON(st)
	}
	close(sup)
}

func main() {
	if topic == "" {
		panic("KAFKA_TOPIC env need to be set")
	}
	sup := make(chan string, numEvents)
	res := make(chan error, numEvents)

	for w := 0; w < numWorkers; w++ {
		go producer(sup, res)
	}

	start := time.Now()
	go generate(&Ads{}, numEvents, sup)
	for i := 0; i < numEvents; i++ {
		r := <-res
		if r != nil {
			fmt.Println(r)
		}
	}
	elapsed := time.Since(start)
	log.Printf("Pushing took %s", elapsed)
}
