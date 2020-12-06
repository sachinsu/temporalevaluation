package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
	"github.com/kelseyhightower/envconfig"
	"go.temporal.io/sdk/client"

	"github.com/sachinsu/temporalevaluation/app"
)

type server struct {
	Debug        bool
	Port         string `default:":8080"`
	DBConnection string `default:"root:passwd@tcp(localhost:3307)/temporaldb?multiStatements=true"`
}

// User holds user details
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	DOB  string `json:"dob"`
	City string `json:"city"`
}

// Index shows welcome message
func (s *server) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome to User Service")
}

// GetUsers returns list of users
func (s *server) GetUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var Users []User

	db, close, err := app.GetSQLXConnection(r.Context(), s.DBConnection)

	if err != nil {
		log.Fatal(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer close()

	if _, err := db.Exec(app.DBSchema); err != nil {
		log.Fatal(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = db.SelectContext(r.Context(), &Users, "select id,name,dob,city from users where isapproved=0")
	if err != nil {
		log.Fatal(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	if err := json.NewEncoder(w).Encode(Users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// UpdateUsers Updates approved status of Users
func (s *server) UpdateUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	creader := csv.NewReader(r.Body)
	records, err := creader.ReadAll()
	if err != nil {
		log.Fatal(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create the client object just once per process
	c, err := client.NewClient(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
		http.Error(w, "Internal Error :Temporal", http.StatusInternalServerError)
		return
	}
	defer c.Close()

	_, err = c.SignalWithStartWorkflow(r.Context(), app.UserApprovalWorkflow, app.ApprovalSignalName,
		records, client.StartWorkflowOptions{}, nil, nil)

	if err != nil {
		log.Fatal(err.Error())
		http.Error(w, "Internal Error: Workflow", http.StatusInternalServerError)
		return
	}

	// db, close, err := app.GetSQLXConnection(r.Context(), s.DBConnection)

	// if err != nil {
	// 	log.Fatal(err.Error())
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	// defer close()
	// sqlStmt := "update users set isapproved=1 where id=:1"

	// tx := db.MustBegin()

	// defer func() {
	// 	if err != nil {
	// 		tx.Rollback()
	// 	}
	// 	tx.Commit()
	// }()

	// for i, line := range records {
	// 	if i == 0 {
	// 		continue
	// 	}
	// 	_, err := tx.ExecContext(r.Context(), sqlStmt, line[0])
	// 	if err != nil {
	// 		log.Fatal(err.Error())
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		return
	// 	}
	// }

	fmt.Fprint(w, "Success")
}

func main() {

	var s server
	err := envconfig.Process("service", &s)

	if err != nil {
		log.Fatal(err.Error())
	}

	router := httprouter.New()
	router.GET("/", (&s).Index)
	router.GET("/Users", (&s).GetUsers)
	router.POST("/Users", (&s).UpdateUsers)

	fmt.Printf("Starting server at %s", s.Port)
	log.Fatal(http.ListenAndServe(s.Port, router))
}
