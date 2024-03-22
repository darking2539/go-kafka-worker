# This System Copy from below concept & adapt for using kafka

- This system is designed to read a large CSV file, process its contents, and store them into a PostgreSQL database. It's built using the Go programming language, and it's designed to be efficient and scalable.

### The system is composed of a few main components:

- CSVParser: Reads and parses the CSV file line by line

- WorkerPool: Manages multiple workers which handle the database insertions

- Workers: Handle batches of records and insert them into the PostgreSQL database

    [CSV File] --> [CSV Parser] --> [Worker Pool] --> [Workers] --> [PostgreSQL Database]


CSVParser

- The CSVParser reads the CSV file line by line and parses each line into a Go struct. If enabled, it can also check for and skip duplicate records based on specific fields in the struct.

WorkerPool

- The WorkerPool is responsible for managing multiple workers. It provides a queue to which the CSVParser submits records. Each worker continuously fetches from this queue, and stores the records in a batch.

Workers

- Each Worker stores records in a batch. When the batch size reaches a specified limit (50k records in this case), or when there are no more records to process, the worker generates a SQL INSERT query and executes it as a transaction to insert the batch of records into the PostgreSQL database.

Stats and Logging

- The system also tracks various statistics and logs them when processing is complete:

    - Total processing time
    
    - Total number of records processed
    
    - Total number of duplicate records (if duplicate checking is enabled)
    
    - Total number of records inserted into the database

System Flow

    
    +-------------+        Total and Duplicate Records
    |             |                   Counters
    |  CSV File   |                      |
    |             |                      |
    +-------------+                      V
           |
           V
    +-------------+
    |             |
    | CSV Parser  |
    |             |
    +-------------+
           |
           V
    +-------------+               Inserted Records
    |             |                   Counter
    | Worker Pool |                     |
    |             |                     |
    +-------------+                     V
           |
           V
    +-------------------+
    |                   |
    |   Workers (x10)   |
    |                   |
    +-------------------+
           |
           V
    +--------------------+
    |                    |
    | PostgreSQL Database|
    |                    |
    +--------------------+


The CSV File data is streamed into the CSV Parser. The parser reads and validates each record, checking for any duplicates if that setting is enabled. It sends the validated records to the Worker Pool. The workers within the pool pick up these records, batch them and finally write them to the PostgreSQL Database. Meanwhile, counters keep track of the total and duplicate records (at the parser) and the inserted records (at each worker).

This approach ensures that the system can efficiently process large files by streaming the data, avoiding loading the entire file into memory, and by using multiple workers to insert the data into the database in parallel.



Code


  
main.go
  
    package main
    
    import (
    	"log"
    	"os"
    	"time"
    )
    func main() {
    	file, err := os.Open("file.csv")
    	if err != nil {
    		log.Fatal(err)
    	}
    	defer file.Close()
    	checkDuplicates := true // set this to false if you don't want to check for duplicates
  
    	parser := NewCSVParser(file, checkDuplicates)
  
    	workerPool := NewWorkerPool(10, 6)
    	startTime := time.Now()
  
    	for {
    		record, err := parser.ParseNext()
    		if err != nil {
    			log.Println(err)
    			break
    		}
    
    		workerPool.Submit(record)
    	}
  
    	workerPool.Stop()
    	endTime := time.Now()
    	elapsedTime := endTime.Sub(startTime)
    	totalInserted := 0
    	for _, worker := range workerPool.workers {
    		totalInserted += worker.Inserted
    	}
  
    	log.Println("CSV file processing completed.")
    	log.Println("Time taken:", elapsedTime)
    	log.Println("Total records processed:", parser.Total)
    	log.Println("Total duplicate records skipped:", parser.Duplicates)
    	log.Println("Total records inserted:", totalInserted)
    }




parser.go

    package main
    
    import (
    	"encoding/csv"
    	"io"
    	"log"
    	"os"
    )
    
    type Record struct {
    	Field1 string
    	Field2 string
    	Field3 string
    }
    
    type UniqueFields struct {
    	Field1 string
    	Field2 string
    }
    
    type CSVParser struct {
    	reader          *csv.Reader
    	seen            map[UniqueFields]bool
    	checkDuplicates bool
    	Total           int
    	Duplicates      int
    }
    
    func NewCSVParser(file *os.File, checkDuplicates bool) *CSVParser {
    	return &CSVParser{
    		reader:          csv.NewReader(file),
    		seen:            make(map[UniqueFields]bool),
    		checkDuplicates: checkDuplicates,
    		Total:           0,
    		Duplicates:      0,
    	}
    }
    
    func (p *CSVParser) ParseNext() (*Record, error) {
    	record, err := p.reader.Read()
    	if err == io.EOF {
    		return nil, err
    	}
    	if err != nil {
    		log.Println(err)
    		return nil, err
    	}

  	p.Total++ // increment the total counter
  
  	r := &Record{
  		Field1: record[0],
  		Field2: record[1],
  		Field3: record[2],
  	}

  	if p.checkDuplicates {
  		uniqueFields := UniqueFields{
  			Field1: r.Field1,
  			Field2: r.Field2,
  		}
  
  		if p.seen[uniqueFields] {
  			p.Duplicates++       // increment the duplicates counter
  			return p.ParseNext() // skip duplicates
  		}
  
  		p.seen[uniqueFields] = true
  	}

	  return r, nil
    }



worker.go

    package main
    
    import (
    	"database/sql"
    	"log"
    	"sync"
    
    	_ "github.com/lib/pq" // postgres driver
    )
    
    type WorkerPool struct {
    	workers      []*Worker
    	queue        chan *Record
    	stop         chan bool
    	wg           sync.WaitGroup
    	workersCount int
    	batchSize    int
    }

    func NewWorkerPool(workersCount, batchSize int) *WorkerPool {
    	pool := &WorkerPool{
    		workers:      make([]*Worker, workersCount),
    		queue:        make(chan *Record),
    		stop:         make(chan bool),
    		workersCount: workersCount,
    		batchSize:    batchSize,
    	}

  	pool.wg.Add(workersCount)
  	for i := 0; i < workersCount; i++ {
  		worker := NewWorker(pool.queue, pool.stop, &pool.wg, batchSize)
  		worker.Start()
  		pool.workers[i] = worker
  	}

	    return pool
    }
    
    func (p *WorkerPool) Submit(record *Record) {
    	p.queue <- record
    }
    
    func (p *WorkerPool) Stop() {
    	close(p.stop)
    	p.wg.Wait()
    }

    type Worker struct {
    	queue     chan *Record
    	stop      chan bool
    	wg        *sync.WaitGroup
    	db        *sql.DB
    	records   []*Record
    	batchSize int
    	Inserted  int
    }

    func NewWorker(queue chan *Record, stop chan bool, wg *sync.WaitGroup, batchSize int) *Worker {
    	db, err := sql.Open("postgres", "user=pqgotest dbname=pqgotest sslmode=verify-full")
    	if err != nil {
    		log.Fatal(err)
    	}

    	return &Worker{
    		queue:     queue,
    		stop:      stop,
    		wg:        wg,
    		db:        db,
    		records:   make([]*Record, 0, batchSize),
    		batchSize: batchSize,
    	}
    }

    func (w *Worker) Start() {
    	go func() {
    		defer w.wg.Done()

  		for {
  			select {
  			case record := <-w.queue:
  				w.records = append(w.records, record)
  				if len(w.records) >= w.batchSize {
  					w.insertBatch()
  				}
  			case <-w.stop:
  				if len(w.records) > 0 {
  					w.insertBatch()
  				}
  				return
  			}
  		}
    	}()
    }

    func (w *Worker) insertBatch() {
    	tx, err := w.db.Begin()
    	if err != nil {
    		log.Println(err)
    		return
    	}

    	stmt, err := tx.Prepare("INSERT INTO your_table VALUES ($1, $2, $3)") // change your query
    	if err != nil {
    		log.Println(err)
    		return
    	}

    	for _, record := range w.records {
    		_, err := stmt.Exec(record.Field1, record.Field2, record.Field3) // change to match your record
    		if err != nil {
    			log.Println(err)
    			return
    		}
    	}

    	err = tx.Commit()
    	if err != nil {
    		log.Println(err)
    		return
    	}
    	w.Inserted += len(w.records) // increment the inserted counter
    
    	w.records = w.records[:0] // clear the batch
    }
