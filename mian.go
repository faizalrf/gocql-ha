package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gocql/gocql"
)

/*
* 	Test code written to test the gocql driver behaviour if scylla nodes go down.
* 	The code fails to execute and terminates. Retrying the transactions does not work, have to re initialise the connection/session
* 	Once the session has been re-initialised, the code continues to work from where it left off.
* 	Faisal 19-Jun-2024
 */

func main() {
	fmt.Println("Test code for gocql!")
	// Define the ScyllaDB cluster
	cluster := gocql.NewCluster("127.0.0.1:9002", "127.0.0.1:9003", "127.0.0.1:9004") // replace with your ScyllaDB IP address
	cluster.Keyspace = "testksp"
	cluster.Consistency = gocql.LocalOne

	// Create a session to ScyllaDB
	session, err := cluster.CreateSession()
	if err != nil {
		fmt.Println("Error Connecting to the ScyllaDB: ", err)
		return
	}
	defer session.Close()

	fmt.Println("Connected to ScyllaDB!")

	rowCount := 0
	numRetry := 5
	numSecondWait := 10 * time.Second

	for {
		// Insert data with retry mechanism by sending the session as a pointer
		if err := writeDataWithRetry(cluster, &session, numRetry, numSecondWait); err != nil {
			fmt.Printf("Failed to insert data after %d retries\n", err)
			return
		} else {
			rowCount++
			fmt.Printf("Inserted Row Count: %d\n", rowCount)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// Write data to the ScyllaDB table with retry mechanism using a session pointer
func writeDataWithRetry(cluster *gocql.ClusterConfig, session **gocql.Session, maxRetries int, delay time.Duration) error {
	accountID := gocql.MustRandomUUID()
	for i := 0; i < maxRetries; i++ {
		err := (*session).Query(`UPDATE tab1 SET "c3" = "c3" + ? WHERE c1 = ? AND c2 = ?`, 1, accountID.String(), getCurrentDate()).Exec()
		if err == nil {
			return nil // Success
		}
		fmt.Printf("Error inserting data, attempt %d: %v\n", i+1, err)

		// Check if the error is related to connection then attempt to reconnect and retry
		if isConnectionError(err) {
			log.Println("Attempting to reconnect to the ScyllaDB...")
			time.Sleep(delay) // Wait before retrying
			if err := reConnect(cluster, session); err != nil {
				fmt.Println("Error re-connecting to the ScyllaDB: ", err)
				continue
			}
		}
	}
	return fmt.Errorf("failed to insert data after %d attempts", maxRetries)
}

// Reconnect using the session pointer
func reConnect(cluster *gocql.ClusterConfig, session **gocql.Session) error {
	(*session).Close()
	var err error
	*session, err = cluster.CreateSession()
	if err != nil {
		return err
	}
	return nil
}

// Check if the error is related to any of the following which need a reconnection to the DB
func isConnectionError(err error) bool {
	return err.Error() == "session has been closed" || err.Error() == "EOF" || err.Error() == "gocql: no hosts available in the pool"
}

func getCurrentDate() time.Time {
	// Get current time and truncate to remove the time component, keeping only the date
	currentTime := time.Now()
	return time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())
}
