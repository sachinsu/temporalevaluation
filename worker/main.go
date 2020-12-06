package main

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"github.com/sachinsu/temporalevaluation/app"
)

// @@@SNIPSTART money-transfer-project-template-go-worker
func main() {
	// Create the client object just once per process
	c, err := client.NewClient(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()
	// This worker hosts both Worker and Activity functions
	w := worker.New(c, app.UserApprovalTaskQueue, worker.Options{})
	w.RegisterWorkflow(app.OnboardUsers)
	w.RegisterActivity(app.ImportUsers)
	w.RegisterActivity(app.ApproveUsers)
	// Start listening to the Task Queue
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}

// @@@SNIPEND
