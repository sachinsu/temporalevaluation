package app

import (
	"time"

	"go.temporal.io/sdk/workflow"
	"go.uber.org/zap"
)

// OnboardUsers is workflow definition functions
func OnboardUsers(ctx workflow.Context, importFileName string, DbConnectionString string) error {
	logger := workflow.GetLogger(ctx)

	logger.Info("Onboardusers", "filename", importFileName, "db Connection", DbConnectionString)

	options := workflow.ActivityOptions{
		// Timeout options specify when to automatically timeout Actvitivy functions.
		StartToCloseTimeout: time.Minute,
		// Optionally provide a customized RetryPolicy.
		// Temporal retries failures by default, this is just an example.
		// RetryPolicy: retrypolicy,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	var count int
	err := workflow.ExecuteActivity(ctx, ImportUsers, importFileName, DbConnectionString).Get(ctx, &count)
	if err != nil {
		logger.Error("Error with ImportUsers", zap.Error(err))
		return err
	}

	signalChan := workflow.GetSignalChannel(ctx, ApprovalSignalName)

	s := workflow.NewSelector(ctx)

	var signalVal string

	s.AddReceive(signalChan, func(c workflow.ReceiveChannel, more bool) {
		c.Receive(ctx, &signalVal)
		logger.Info("Received signal!", zap.String("signal", ApprovalSignalName), zap.String("value", signalVal))
	})

	logger.Info("Waiting for Signal on Channel" + ApprovalSignalName)

	s.Select(ctx)

	err = workflow.ExecuteActivity(ctx, ApproveUsers, DbConnectionString, signalVal).Get(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}

// @@@SNIPEND
