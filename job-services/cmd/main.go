package main

import (
	"context"
	"encoding/json"
	"job-services/internal/app"
	"job-services/internal/models"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup

	//init consumer
	concurrency := 1
	wg.Add(concurrency)
	go func() {
		defer wg.Done()

		for app.ConsumerRunning {
			msg, err := app.ConsumerReader.Reader.ReadMessage(context.Background())
			if err == nil {
				app.ConsumerMsgChan <- &msg
			} else if !app.ConsumerRunning {
				log.Println("Stop Consumer")
				return
			} else {
				log.Println("Error: ", err.Error())
			}

			if err := app.ConsumerReader.Reader.CommitMessages(context.Background(), msg); err != nil {
				log.Printf("Error committing offset: %v\n", err)
			}

		}

	}()

	// (consumer --> worker)
	wg.Add(1)
	go func() {
		defer wg.Done()

		for msg := range app.ConsumerMsgChan {
			//TODO
			var record models.TransactionModel
			err := json.Unmarshal(msg.Value, &record)
			if err == nil {
				//fmt.Println(string(msg.Value))
				app.WorkerPool.Submit(&record)
			}
		}
	}()

	log.Println("Application Start Sucessful")

	<-sigterm // Await a sigterm signal

	//stop consumer first
	app.ConsumerRunning = false
	app.ConsumerReader.Reader.Close()
	close(app.ConsumerMsgChan)
	wg.Wait()

	//stop workerpool
	app.WorkerPool.Stop()

	log.Println("Wait 5s before shut down!!!")
	time.Sleep(time.Second * 5)
	os.Exit(0)
}
