package app

import (
	"testing"
	"time"

	"github.com/sachinsu/temporalevaluation/app"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
)

func Test_Workflow(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()
	// Mock activity implementation

	env.OnActivity(ImportUsers, mock.Anything, app.Userdata, app.Dbconn).Return(2, nil)
	env.RegisterDelayedCallback(func() {
		env.SignalWorkflow(ApprovalSignalName, mock.Anything)
	}, time.Minute)
	env.OnActivity(ApproveUsers, mock.Anything, app.Dbconn, mock.Anything).Return(0, nil)
	env.ExecuteWorkflow(OnboardUsers, app.Userdata, app.Dbconn)
	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
}
