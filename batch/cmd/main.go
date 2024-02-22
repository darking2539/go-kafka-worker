package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	klcsv "kafka-batch/csv"
	"kafka-batch/internal/app"
	"kafka-batch/internal/services"
	"kafka-batch/internal/worker"

	"github.com/joho/godotenv"
	_ "go.uber.org/ratelimit"
)

func init() {
	godotenv.Load()
}

func main() {

	startTime := time.Now()
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	//create produce Worker
	var wgProduceWorker sync.WaitGroup
	concurrency := 40
	wgProduceWorker.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			defer wgProduceWorker.Done()
			for msg := range app.ProducerQueue {
				app.ProducerWriter.PublishMessage(msg.Message)
			}
		}()
	}

	//init work Group
	workerCount := 10
	workerPool := worker.NewWorkerPool(workerCount)
	log.Printf("worker: %d", workerCount)

	//prepare csv file data
	csvFilePath := "boss.csv"
	countCSV := 0
	callback := func(record []string) error {
		return services.ProcessRecord(record, workerPool, &countCSV)
	}

	//read csv file and process
	go func() {
		err := klcsv.ReadCSV(csvFilePath, callback, nil, 1)
		if err != nil {
			log.Fatal("Read CSV Error:", err.Error())
		}
		sigterm <- syscall.SIGINT
	}()

	<-sigterm // Await a sigterm signal

	//stop worker pool
	workerPool.Stop()

	//stop producerQueue
	close(app.ProducerQueue)

	//wait all workers is shutdown
	wgProduceWorker.Wait()

	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	totalInserted := 0

	for _, worker := range workerPool.Workers {
		totalInserted += worker.Inserted
	}

	log.Println("ReadFrom CSV: ", countCSV, " row")
	log.Println("Produce Data to Kafka Record:", totalInserted)
	log.Println("TimeTaken:", elapsedTime)
	log.Println("Wait 5s before shut down!!!")
	time.Sleep(time.Second * 5)
	os.Exit(0)

}
