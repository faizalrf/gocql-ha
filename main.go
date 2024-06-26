// Automatic Retry mechanism with gocql driver
package main

import (
	"fmt"
	"time"

	"github.com/gocql/gocql"
)

/*
* 	Test code written to test the gocql driver behaviour if scylla nodes go down.
* 	The code using the built-in Simple Retry policy to automatically rerty the faild execution
*   The retry is seamless and the application does not see any failures and executes smoothly
* 	Faisal 25-Jun-2024
 */

func main() {
	fmt.Println("Automatic SimpleRetry GOCQL Policy Test!")
	// Define the ScyllaDB cluster
	cluster := gocql.NewCluster("172.18.0.2", "172.18.0.3", "172.18.0.4")
	cluster.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(gocql.RoundRobinHostPolicy())
	cluster.Logger = gocql.Logger // Enable logging

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

	for {
		if err := writeDataWithSimpleRetryPolicy(session, numRetry); err != nil {
			fmt.Printf("Failed to insert data after %d retries\n", err)
			return
		} else {
			rowCount++
			fmt.Printf("Inserted Row Count: %d\n", rowCount)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func writeDataWithSimpleRetryPolicy(session *gocql.Session, maxRetries int) error {
	accountID := gocql.TimeUUID()

	// The simple retry policy and Speculative Execution
	retryPolicy := &gocql.SimpleRetryPolicy{
		NumRetries: maxRetries,
	}

	// Speculative Execution Policy for preventing lags in query execution by
	// automatically executing the query on another node if the original one times out
	sp := &gocql.SimpleSpeculativeExecution{
		NumAttempts:  maxRetries,
		TimeoutDelay: 10 * time.Millisecond,
	}
	//Counters Update
	cql := `UPDATE tab1 SET "c3" = "c3" + ? WHERE c1 = ? AND c2 = ?`

	//Build the query with Simple Retry and Speculative Execution Policy
	qry := session.Query(cql, 1, accountID.String(), getCurrentDate()).
		RetryPolicy(retryPolicy).
		SetSpeculativeExecutionPolicy(sp).
		Idempotent(true)

	//Execute the query, the query will now automatically retry on a different node
	// if the original node went down or failed execution
	err := qry.Exec()
	if err != nil {
		fmt.Printf("The query failed with '%v'!\n", err)
	}
	return err
}

// Return a date for writing into the table
func getCurrentDate() time.Time {
	// Get current time and truncate to remove the time component, keeping only the date
	currentTime := time.Now()
	return time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())
}
