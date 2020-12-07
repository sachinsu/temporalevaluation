package app

import (
	"context"
	"encoding/csv"
	"io"
	"io/ioutil"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"go.temporal.io/sdk/activity"
	"go.uber.org/zap"
)

// ImportUsers is first activity in workflow
func ImportUsers(ctx context.Context, filename string, DbConnectionString string) (int, error) {

	logger := activity.GetLogger(ctx)

	logger.Info("ImportUsers activity started.", zap.String("filename", filename),
		zap.String("Dbconn", DbConnectionString))

	if _, err := os.Stat(filename); err == nil {
		logger.Error("File does not exists", zap.Error(err))
		return 0, err
	}

	// db, close, err := GetSQLXConnection(context.Background(), DbConnectionString)
	// if err != nil {
	// 	logger.Error("Cant open connection to database", zap.Error(err))
	// 	return 0, err
	// }

	// defer close()

	// if _, err := db.Exec(DBSchema); err != nil {
	// 	logger.Error("Error while executing Schema", zap.Error(err))
	// 	return 0, err
	// }

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Error("Unable to read from file", zap.Error(err))
		return 0, err
	}

	r := csv.NewReader(strings.NewReader(string(content)))
	r.Comma = ','
	r.Comment = '#'

	// sqlStmt := "insert into users(name,dob,city) values(:1,:2,:3)"

	// tx := db.MustBegin()

	// defer func() {
	// 	if err != nil {
	// 		tx.Rollback()
	// 	}
	// 	tx.Commit()
	// }()

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

		logger.Info("Record read is ->", len(record))

		// if _, err := tx.Exec(sqlStmt, record[0], record[1], record[2]); err != nil {
		// 	logger.Error("Error while writing user record", zap.Error(err))
		// 	return 0, err
		// }
	}

	return i, nil
}

// ApproveUsers waits for signal with list of approved users.
func ApproveUsers(ctx context.Context, DbConnectionString string, Users string) (int, error) {

	logger := activity.GetLogger(ctx)
	logger.Info("ApprovedUsers called", zap.String("Dbconn", DbConnectionString), zap.String("Userlist", Users))

	// 	db, close, err := GetSQLXConnection(context.Background(), DbConnectionString)
	// 	if err != nil {
	// 		logger.Error("Cant open connection to database", zap.Error(err))
	// 		return err
	// 	}

	// 	defer close()

	// 	if _, err := db.Exec(schema); err != nil {
	// 		logger.Error("Error while executing Schema", zap.Error(err))
	// 		return err
	// 	}

	r := csv.NewReader(strings.NewReader(Users))

	userList, err := r.ReadAll()

	if err != nil {
		logger.Error("Error reading user list", zap.Error(err))
		return 0, err
	}

	return len(userList), nil

	// 	tx := db.MustBegin()

	// 	defer func() {
	// 		if err != nil {
	// 			tx.Rollback()
	// 		}
	// 		tx.Commit()
	// 	}()

	// 	sqlStmt := "update users set isapproved =1 where id =:1"

	// 	for {
	// 		record, err := r.Read()
	// 		if err == io.EOF {
	// 			break
	// 		}
	// 		if err != nil {
	// 			logger.Error("Error while reading from file", zap.Error(err))
	// 			return err
	// 		}

	// 		if _, err := tx.Exec(sqlStmt, record[0]); err != nil {
	// 			logger.Error("Error while writing user record", zap.Error(err))
	// 			return err

	// 		}
	// 	}

	// }
}

// ComposeGreeting is test function
// func ComposeGreeting(ctx context.Context, name string) (string, error) {
// 	logger := activity.GetLogger(ctx)
// 	logger.Info("Composegreeting started")
// 	greeting := fmt.Sprintf("Hello %s!", name)
// 	return greeting, nil
// }
