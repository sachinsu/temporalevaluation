package app

import (
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
)

func Test_Workflow(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()
	// Mock activity implementation
	filename := "C:\\Users\\sachi\\projects\\temporal\\app\\users.csv"
	dbconn := "user@password@/temporaldb"

	env.OnActivity(ImportUsers, mock.Anything, filename, dbconn).Return(2, nil)
	env.RegisterDelayedCallback(func() {
		env.SignalWorkflow(ApprovalSignalName, mock.Anything)
	}, time.Minute)
	env.OnActivity(ApproveUsers, mock.Anything, dbconn, mock.Anything).Return(0, nil)
	env.ExecuteWorkflow(OnboardUsers, filename, dbconn)
	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
}
