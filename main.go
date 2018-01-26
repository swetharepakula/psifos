package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type psifosServer struct {
	listener net.Listener
	db       *sql.DB
}

type VcapServices struct {
	Pmysql []ServiceInstances `json:"p-mysql"`
}

type ServiceInstances struct {
	Credentials Credentials `json:"credentials"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
	Name     string `json:"name"`
}

type Row struct {
	Color string
	Votes int
}

func main() {
	server := NewPsifosServer()
	port := os.Getenv("PORT")
	portNumber, err := strconv.Atoi(port)
	FreakOut(err)

	connBytes := os.Getenv("VCAP_SERVICES")

	myServices := &VcapServices{}
	err = json.Unmarshal([]byte(connBytes), myServices)
	FreakOut(err)
	creds := myServices.Pmysql[0].Credentials

	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", creds.Username, creds.Password, creds.Hostname, creds.Port, creds.Name)

	server.db, err = sql.Open("mysql", connString)
	FreakOut(err)
	defer server.db.Close()
	err = server.db.Ping()
	FreakOut(err)

	server.Start(portNumber)
	defer server.Stop()
}

func NewPsifosServer() *psifosServer {
	return &psifosServer{}
}

func (s *psifosServer) Start(port int) {
	l, e := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if e != nil {
		log.Fatal("listen error:", e)
	}

	http.Serve(l, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.Contains(path, "/put/blue") {
			_, err := s.db.Exec("update colors set votes = votes + 1 where color = ?", "blue")

			if err != nil {
				w.WriteHeader(500)
			}
		} else if strings.Contains(path, "/put/yellow") {
			_, err := s.db.Exec("update colors set votes = votes + 1 where color = ?", "yellow")

			if err != nil {
				w.WriteHeader(500)
			}
		} else if strings.Contains(path, "clear/database") {
			_, err := s.db.Exec("truncate table colors")

			if err != nil {
				w.WriteHeader(500)
			}
		} else if strings.Contains(path, "create/database") {

			_, err := s.db.Exec("CREATE TABLE colors ( color varchar(32), votes integer )")
			if err != nil {
				w.WriteHeader(500)
			}

			_, err = s.db.Exec("insert into colors (color, votes) values (?, ?)", "blue", 0)
			if err != nil {
				w.WriteHeader(500)
			}

			_, err = s.db.Exec("insert into colors (color, votes) values (?, ?)", "yellow", 0)
			if err != nil {
				w.WriteHeader(500)
			}
		}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(200)
		rows, err := s.getAllRows()
		FreakOut(err)
		w.Write([]byte(fmt.Sprintf("%v", rows)))
	}))
}

func (s *psifosServer) Stop() {
	s.listener.Close()
}

func FreakOut(err error) {
	if err != nil {
		panic(err)
	}
}

func (s *psifosServer) getAllRows() ([]Row, error) {
	rows, err := s.db.Query("select * from colors")
	if err != nil {
		return []Row{}, err
	}

	defer rows.Close()

	colors := []Row{}

	for rows.Next() {

		row := Row{}
		err = rows.Scan(&row.Color, &row.Votes)
		if err != nil {
			return []Row{}, err
		}
		colors = append(colors, row)
	}

	return colors, nil
}
