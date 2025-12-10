package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// const CSV_FILE = "customers-100.csv"

const CSV_FILE = "customers-1m.csv"

type App struct {
	db *sql.DB
	mu sync.RWMutex
}

func main() {
	app := NewApp()
	app.Run()
}

func NewApp() *App {
	return &App{
		db: initDB(),
	}
}

func initDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./one-mille.db?_busy_timeout=10000")
	if err != nil {
		log.Fatalln(err)
	}

	_, err = db.Exec("PRAGMA journal_mode=WAL;")
	if err != nil {
		log.Fatal(err)
	}

	sql := `
		CREATE TABLE IF NOT EXISTS customers
		(id string, email string, name string, company string, city string, country string, birthday Date)
	`

	_, err = db.Exec(sql)
	if err != nil {
		log.Fatal("Cannot create table: ", err)
	}

	return db
}

func (app *App) TruncateDB() {
	_, err := app.db.Exec("DELETE FROM customers")
	if err != nil {
		log.Fatal("Cannot truncate ", err)
	}
}

func (app *App) Run() {

	defer app.db.Close()

	app.TruncateDB()
	app.SolutionOne()

	app.TruncateDB()
	app.SolutionTwo()

	app.TruncateDB()
	app.SolutionThree()

	app.TruncateDB()
	app.SolutionFour()

	// todo add even faster solutions

}

// 1 million: Solution one done in 16.52167625s
func (app *App) SolutionOne() {

	fmt.Println("Starting solution one ...")

	defer func(t time.Time) {
		eT := time.Since(t)
		fmt.Printf("Solution one done in %s\n", eT.String())
	}(time.Now())

	wg := sync.WaitGroup{}

	jobChan := make(chan []string, 10000)

	wg.Add(1)
	go func() {

		defer wg.Done()

		f, err := os.Open(CSV_FILE)
		if err != nil {
			log.Fatalln("Cannot open file", err)
		}
		defer f.Close()

		csvReader := csv.NewReader(f)

		if err != nil {
			log.Fatalln("Cannot read csv", err)
		}

		// Skip header row
		_, err = csvReader.Read()
		if err != nil {
			log.Fatalln("Cannot read header", err)
		}

		for {
			record, err := csvReader.Read()
			if err != nil {
				break
			}
			jobChan <- record
		}
		close(jobChan)
	}()

	wg.Add(1)
	go func() {

		defer wg.Done()
		for record := range jobChan {

			app.mu.Lock()
			tx, err := app.db.Begin()
			if err != nil {
				log.Fatalln("error ", err)
			}

			stmt, err := tx.Prepare("INSERT INTO customers (id, name, email, company, city, country, birthday ) VALUES (?,?,?,?,?,?,?)")
			if err != nil {
				log.Fatalln("error ", err)
			}

			_, err = stmt.Exec(record[0], record[1], record[2], record[3], record[4], record[5], record[6])
			if err != nil {
				stmt.Close()
				log.Fatalln("error ", err)
			}

			err = tx.Commit()
			if err != nil {
				stmt.Close()
				log.Fatalln("error ", err)
			}
			stmt.Close()
			app.mu.Unlock()

		}

	}()

	wg.Wait()

}

// 1 million: Solution two done in 16.973486334s
func (app *App) SolutionTwo() {

	fmt.Println("Starting solution two ...")

	defer func(t time.Time) {
		eT := time.Since(t)
		fmt.Printf("Solution two done in %s\n", eT.String())
	}(time.Now())

	wg := sync.WaitGroup{}

	jobChan := make(chan []string, 10000)

	wg.Add(1)
	go func() {

		defer wg.Done()

		f, err := os.Open(CSV_FILE)
		if err != nil {
			log.Fatalln("Cannot open file", err)
		}
		defer f.Close()

		csvReader := csv.NewReader(f)
		// Skip header row
		_, err = csvReader.Read()
		if err != nil {
			log.Fatalln("Cannot read header", err)
		}

		for {
			record, err := csvReader.Read()
			if err != nil {
				break
			}
			jobChan <- record
		}
		close(jobChan)
	}()

	for w := 1; w <= 50; w++ {

		wg.Add(1)
		go func() {

			defer wg.Done()
			for record := range jobChan {

				tx, err := app.db.Begin()
				if err != nil {
					log.Fatalln("error ", err)
				}

				stmt, err := tx.Prepare("INSERT INTO customers (id, name, email, company, city, country, birthday ) VALUES (?,?,?,?,?,?,?)")
				if err != nil {
					log.Fatalln("error ", err)
				}

				_, err = stmt.Exec(record[0], record[1], record[2], record[3], record[4], record[5], record[6])
				if err != nil {
					stmt.Close()
					log.Fatalln("error ", err)
				}

				err = tx.Commit()
				if err != nil {
					stmt.Close()
					log.Fatalln("error ", err)
				}
				stmt.Close()

			}

		}()

	}

	wg.Wait()

}

// 1 million: Solution three done in 1.7710695s
func (app *App) SolutionThree() {

	fmt.Println("Starting solution three ...")

	defer func(t time.Time) {
		eT := time.Since(t)
		fmt.Printf("Solution three done in %s\n", eT.String())
	}(time.Now())

	wg := sync.WaitGroup{}

	jobChan := make(chan []string, 1000000)

	wg.Add(1)
	go func() {

		defer wg.Done()

		f, err := os.Open(CSV_FILE)
		if err != nil {
			log.Fatalln("Cannot open file", err)
		}
		defer f.Close()

		csvReader := csv.NewReader(f)
		// Skip header row
		_, err = csvReader.Read()
		if err != nil {
			log.Fatalln("Cannot read header", err)
		}

		for {
			record, err := csvReader.Read()
			if err != nil {
				break
			}
			jobChan <- record
		}
		close(jobChan)
	}()

	wg.Add(1)
	go func() {

		defer wg.Done()

		batchSize := 1000
		var records [][]string

		for record := range jobChan {
			records = append(records, record)

			if len(records) >= batchSize {
				app.insertBatch(records)
				records = nil
			}
		}

		// Insert remaining records
		if len(records) > 0 {
			app.insertBatch(records)
		}

	}()

	wg.Wait()

}

// 1 million: Solution four done in 1.742708875s
func (app *App) SolutionFour() {

	fmt.Println("Starting solution four ...")

	defer func(t time.Time) {
		eT := time.Since(t)
		fmt.Printf("Solution four done in %s\n", eT.String())
	}(time.Now())

	wg := sync.WaitGroup{}

	jobChan := make(chan []string, 10000)

	wg.Add(1)
	go func() {

		defer wg.Done()

		f, err := os.Open(CSV_FILE)
		if err != nil {
			log.Fatalln("Cannot open file", err)
		}
		defer f.Close()

		csvReader := csv.NewReader(f)
		// Skip header row
		_, err = csvReader.Read()
		if err != nil {
			log.Fatalln("Cannot read header", err)
		}

		for {
			record, err := csvReader.Read()
			if err != nil {
				break
			}
			jobChan <- record
		}
		close(jobChan)
	}()

	for w := 1; w <= 10; w++ {

		wg.Add(1)
		go func() {

			defer wg.Done()

			batchSize := 10000
			var records [][]string

			for record := range jobChan {
				records = append(records, record)

				if len(records) >= batchSize {
					app.mu.Lock()
					app.insertBatch(records)
					app.mu.Unlock()
					records = nil
				}
			}

			// Insert remaining records
			if len(records) > 0 {
				app.mu.Lock()
				app.insertBatch(records)
				app.mu.Unlock()
			}

		}()

	}

	wg.Wait()

}

func (app *App) insertBatch(records [][]string) {
	tx, err := app.db.Begin()
	if err != nil {
		log.Fatalln("error ", err)
	}

	stmt, err := tx.Prepare("INSERT INTO customers (id, name, email, company, city, country, birthday) VALUES (?,?,?,?,?,?,?)")
	if err != nil {
		log.Fatalln("error ", err)
	}
	defer stmt.Close()

	for _, record := range records {
		_, err = stmt.Exec(record[0], record[1], record[2], record[3], record[4], record[5], record[6])
		if err != nil {
			tx.Rollback()
			log.Fatalln("error ", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalln("error ", err)
	}
}
