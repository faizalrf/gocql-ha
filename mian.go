package main

import (
	"fmt"
	"time"

	"github.com/gocql/gocql"
)

func main() {
	fmt.Println("Test Suite for gocql!")
	// Define the ScyllaDB cluster
	cluster := gocql.NewCluster("127.0.0.1:9002", "127.0.0.1:9003", "127.0.0.1:9004") // replace with your ScyllaDB IP address
	cluster.Keyspace = "testksp"
	cluster.Consistency = gocql.LocalOne

	// Create a session
	session, err := cluster.CreateSession()
	if err != nil {
		fmt.Println("Error creating session: ", err)
		return
	}
	defer session.Close()

	fmt.Println("Connected to ScyllaDB!")
	rowCount := 0
	for {
		// Insert data with retry mechanism
		if err := insertDataWithRetry(session, 10, 10*time.Second); err != nil {
			fmt.Println("Failed to insert data after retries: ", err)
			return
		} else {
			rowCount++
			fmt.Printf("Rows inserted: %d\n", rowCount)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func insertDataWithRetry(session *gocql.Session, maxRetries int, delay time.Duration) error {
	account_id := gocql.MustRandomUUID()
	for i := 0; i < maxRetries; i++ {
		err := session.Query(`UPDATE daily_account_stat SET "list" = "list" + ? where account_id = ? and date = ?`, 1, account_id.String(), getCurrentDate()).Exec()
		if err == nil {
			return nil // Success
		}
		fmt.Printf("Error inserting data, attempt %d: %v\n", i+1, err)
		time.Sleep(delay) // Wait before retrying
	}
	return fmt.Errorf("failed to insert data after %d attempts", maxRetries)
}

func getCurrentDate() time.Time {
	// Get current time and truncate to remove the time component, keeping only the date
	currentTime := time.Now()
	return time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())
}
