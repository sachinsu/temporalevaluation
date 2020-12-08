package app

import (
	"context"

	"github.com/jmoiron/sqlx"
)

const UserApprovalTaskQueue = "USER_APPROVAL_TASK_QUEUE"
const ApprovalSignalName = "APPROVAL_SIGNAL"
const UserApprovalWorkflow = "user-approval-workflow"

const DBSchema = `
			CREATE TABLE if not exists users (
				id int auto_increment primary key,
				name varchar(100),
				dob varchar(10),
				city varchar(10),
				isapproved int default 0
			);`

const Userdata = `
	name,dob,city
john,1980-01-01,mumbai
harry,1985-01-01,valsad`

const Dbconn = "root:passwd@tcp(localhost:3307)/temporaldb?multiStatements=true"

// GetSQLXConnection is a helper function to open connection to database
func GetSQLXConnection(ctx context.Context, dbConn string) (*sqlx.DB, func() error, error) {
	db, err := sqlx.ConnectContext(ctx, "mysql", dbConn)
	return db, db.Close, err
}
