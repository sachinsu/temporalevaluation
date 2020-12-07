package app

import (
	"context"
	"encoding/csv"
	"io"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"go.temporal.io/sdk/activity"
	"go.uber.org/zap"
)

// ImportUsers is first activity in workflow
func ImportUsers(ctx context.Context, userdata string, DbConnectionString string) (int, error) {

	logger := activity.GetLogger(ctx)

	logger.Info("ImportUsers activity started.", zap.String("Dbconn", DbConnectionString))

	db, close, err := GetSQLXConnection(context.Background(), DbConnectionString)
	if err != nil {
		logger.Error("Cant open connection to database", zap.Error(err))
		return 0, err
	}

	defer close()

	if _, err := db.Exec(DBSchema); err != nil {
		logger.Error("Error while executing Schema", zap.Error(err))
		return 0, err
	}

	logger.Info("Database connection opened, now parsing user data")

	sqlStmt := "insert into users(name,dob,city) values(?,?,?)"

	tx := db.MustBegin()

	defer func() {
		if err != nil {
			tx.Rollback()
		}
		tx.Commit()
	}()

	r := csv.NewReader(strings.NewReader(string(userdata)))
	r.Comma = ','
	r.Comment = '#'

	records, err := r.ReadAll()
	if err != nil {
		logger.Error("Error while reading", zap.Error(err))
		return 0, err
	}

	i := 0

	for i, record := range records {
		if i == 0 {
			continue
		}

		logger.Info("Record read is ->", record[0])

		if _, err := tx.Exec(sqlStmt, record[0], record[1], record[2]); err != nil {
			logger.Error("Error while writing user record", zap.Error(err))
			return 0, err
		}
	}

	return i, nil
}

// ApproveUsers waits for signal with list of approved users.
func ApproveUsers(ctx context.Context, DbConnectionString string, Users string) (int, error) {

	logger := activity.GetLogger(ctx)
	logger.Info("ApprovedUsers called", zap.String("Dbconn", DbConnectionString), zap.String("Userlist", Users))

	db, close, err := GetSQLXConnection(context.Background(), DbConnectionString)
	if err != nil {
		logger.Error("Cant open connection to database", zap.Error(err))
		return 0, err
	}

	defer close()

	if _, err := db.Exec(DBSchema); err != nil {
		logger.Error("Error while executing Schema", zap.Error(err))
		return 0, err
	}

	r := csv.NewReader(strings.NewReader(Users))

	tx := db.MustBegin()

	defer func() {
		if err != nil {
			tx.Rollback()
		}
		tx.Commit()
	}()

	sqlStmt := "update users set isapproved =1 where id =:1"

	i := 0

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			logger.Error("Error while reading from file", zap.Error(err))
			return 0, err
		}

		if i == 0 {
			continue
		}
		i++

		if _, err := tx.Exec(sqlStmt, record[0]); err != nil {
			logger.Error("Error while writing user record", zap.Error(err))
			return 0, err

		}
	}
	return i, nil
}

// ComposeGreeting is test function
// func ComposeGreeting(ctx context.Context, name string) (string, error) {
// 	logger := activity.GetLogger(ctx)
// 	logger.Info("Composegreeting started")
// 	greeting := fmt.Sprintf("Hello %s!", name)
// 	return greeting, nil
// }
