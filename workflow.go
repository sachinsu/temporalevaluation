package app

import (
	"time"

	"go.temporal.io/sdk/workflow"
	"go.uber.org/zap"
)

// OnboardUsers is workflow definition functions
func OnboardUsers(ctx workflow.Context, userdata string, DbConnectionString string) error {
	logger := workflow.GetLogger(ctx)

	logger.Info("Onboardusers called", "db Connection", DbConnectionString)

	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	var count int
	err := workflow.ExecuteActivity(ctx, ImportUsers, userdata, DbConnectionString).Get(ctx, &count)
	if err != nil {
		logger.Error("Error with ImportUsers", zap.Error(err))
		return err
	}

	// Configure to wait on channel for signal
	signalChan := workflow.GetSignalChannel(ctx, ApprovalSignalName)

	s := workflow.NewSelector(ctx)

	var signalVal string

	s.AddReceive(signalChan, func(c workflow.ReceiveChannel, more bool) {
		c.Receive(ctx, &signalVal)
		logger.Info("Received signal!", zap.String("signal", ApprovalSignalName), zap.String("value", signalVal))
	})

	logger.Info("Waiting for Signal on Channel" + ApprovalSignalName)

	s.Select(ctx)

	// Call ApproveUsers activity with data received in signal
	err = workflow.ExecuteActivity(ctx, ApproveUsers, DbConnectionString, signalVal).Get(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}
