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
	Uri string `json:"uri"`
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

	// json.Unmarshal([]byte(os.Getenv("VCAP_SERVICES")), &vcapServices)

	// postgresCredentials := vcapServices["postgres"].([]interface{})[0].(map[string]interface{})["credentials"].(map[string]interface{})
	// jdbcUri := postgresCredentials["jdbc_uri"].(string)

	connBytes := os.Getenv("VCAP_SERVICES")
	// myServices := &vcapServices{
	// 	pmysql: []serviceInstances{
	// 		serviceInstances{
	// 			credentials: credentials{},
	// 		},
	// 	},
	// }

	myServices := &VcapServices{}
	err = json.Unmarshal([]byte(connBytes), myServices)
	FreakOut(err)

	server.db, err = sql.Open("mysql", myServices.Pmysql[0].Credentials.Uri)
	FreakOut(err)
	defer server.db.Close()
	err = server.db.Ping()
	FreakOut(err)

	_, err = server.db.Exec("create database if not exists colors;", Row{})
	FreakOut(err)

	_, err = server.db.Exec("insert into colors (color, votes) values (blue, 0)  ")
	FreakOut(err)
	_, err = server.db.Exec("insert into colors (color, votes) values (yellow, 0)  ")
	FreakOut(err)

	fmt.Println("CONNECTION STRING: ", myServices.Pmysql[0].Credentials.Uri)

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
			_, err := s.db.Exec("update colors set votes = votes + 1 where color = blue")

			if err != nil {
				w.WriteHeader(500)
			}
			w.Write([]byte("blue"))
		} else if strings.Contains(path, "/put/yellow") {
			w.Write([]byte("yellow"))
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
