package app

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// @@@SNIPSTART money-transfer-project-template-go-shared-task-queue
const UserApprovalTaskQueue = "USER_APPROVAL_TASK_QUEUE"
const ApprovalSignalName = "APPROVAL_SIGNAL"
const UserApprovalWorkflow = "user-approval-workflow"

// @@@SNIPEND

// type TransferDetails struct {
// 	Amount      float32
// 	FromAccount string
// 	ToAccount   string
// 	ReferenceID string
// }

// GetSQLXConnection is a helper function to open connection to database
func GetSQLXConnection(ctx context.Context, dbConn string) (*sqlx.DB, func() error, error) {
	db, err := sqlx.ConnectContext(ctx, "mysql", dbConn)
	return db, db.Close, err
}
