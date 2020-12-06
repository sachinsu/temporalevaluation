package app

import (
	"context"
	"encoding/csv"
	"io"
	"io/ioutil"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"go.temporal.io/sdk/workflow"
	"go.uber.org/zap"
)

var schema = `
			CREATE TABLE users if not exists (
				id int auto_increment primary key,
				name text,
				dob text,
				city text,
				isapproved int default 0
			);`

// ImportUsers is first activity in workflow
func ImportUsers(ctx workflow.Context, filename string, DbConnectionString string) error {

	logger := workflow.GetLogger(ctx)

	logger.Info("ImportUsers called.", zap.String("filename", filename))

	if _, err := os.Stat(filename); err == nil {
		logger.Error("File does not exists", zap.Error(err))
		return err
	}

	db, close, err := GetSQLXConnection(context.Background(), DbConnectionString)
	if err != nil {
		logger.Error("Cant open connection to database", zap.Error(err))
		return err
	}

	defer close()

	if _, err := db.Exec(schema); err != nil {
		logger.Error("Error while executing Schema", zap.Error(err))
		return err
	}

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Error("Unable to read from file", zap.Error(err))
		return err
	}

	r := csv.NewReader(strings.NewReader(string(content)))
	r.Comma = ';'
	r.Comment = '#'

	sqlStmt := "insert into users(name,dob,city) values(:1,:2,:3)"

	tx := db.MustBegin()

	defer func() {
		if err != nil {
			tx.Rollback()
		}
		tx.Commit()
	}()

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Error("Error while reading from file", zap.Error(err))
			return err
		}

		if _, err := tx.Exec(sqlStmt, record[0], record[1], record[2]); err != nil {
			logger.Error("Error while writing user record", zap.Error(err))
			return err

		}
	}

	return nil
}

// ApproveUsers waits for signal with list of approved users.
func ApproveUsers(ctx workflow.Context, DbConnectionString string) error {
	var signalVal string

	logger := workflow.GetLogger(ctx)

	signalChan := workflow.GetSignalChannel(ctx, ApprovalSignalName)

	s := workflow.NewSelector(ctx)
	s.AddReceive(signalChan, func(c workflow.ReceiveChannel, more bool) {
		c.Receive(ctx, &signalVal)
		logger.Info("Received signal!", zap.String("signal", ApprovalSignalName), zap.String("value", signalVal))
	})

	s.Select(ctx)

	// if len(signalVal) > 0 {
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

	// 	r := csv.NewReader(strings.NewReader(signalVal))

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

	return nil
}

// todo: SendWelcomeSMS
