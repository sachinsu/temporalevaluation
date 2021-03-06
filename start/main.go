package main

import (
	"context"
	"fmt"
	"log"

	"go.temporal.io/sdk/client"

	"github.com/sachinsu/temporalevaluation/app"
)

// @@@SNIPSTART money-transfer-project-template-go-start-workflow
func main() {
	// Create the client object just once per process
	c, err := client.NewClient(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()
	options := client.StartWorkflowOptions{
		ID:        app.UserApprovalWorkflow,
		TaskQueue: app.UserApprovalTaskQueue,
	}

	we, err := c.ExecuteWorkflow(context.Background(), options, app.OnboardUsers, app.Userdata, app.Dbconn)
	if err != nil {
		log.Fatalln("error starting OnboardUsers workflow", err)
	} else {
		var count int
		err = we.Get(context.Background(), &count)
		fmt.Printf("Record Processed %d", count)
	}
	// printResults(transferDetails, we.GetID(), we.GetRunID())
}

// @@@SNIPEND

// func printResults(transferDetails app.TransferDetails, workflowID, runID string) {
// 	log.Printf(
// 		"\nTransfer of $%f from account %s to account %s is processing. ReferenceID: %s\n",
// 		transferDetails.Amount,
// 		transferDetails.FromAccount,
// 		transferDetails.ToAccount,
// 		transferDetails.ReferenceID,
// 	)
// 	log.Printf(
// 		"\nWorkflowID: %s RunID: %s\n",
// 		workflowID,
// 		runID,
// 	)
// }
