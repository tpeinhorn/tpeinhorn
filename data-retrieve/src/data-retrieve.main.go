package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"io/ioutil"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gorilla/mux"
)

// Replace with your own connection parameters
var server_sql = "127.0.0.1"
var sqlport = 11433
var user = "sa"
var password = "MyStrongPassword123"
var database = "master"
var tablename = "Articles"

var db *sql.DB
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

type Article struct {
	Id          string `json:"Id"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
	Content     string `json:"Content"`
}

var Articles_ret []Article

type configty struct {
	DATA_GEN_API_EXT_PORT1  string
	DATA_SAVER_EXT_PORT1    string
	DATA_RETRIEVE_EXT_PORT1 string
	DB_SQL_EXT_PORT1        string

	DATA_GEN_API_EXT1  string
	DATA_SAVER_EXT1    string
	DATA_RETRIEVE_EXT1 string
	DB_SQL_EXT1        string
}

func set_env() {
	myConfig.DATA_GEN_API_EXT_PORT1 = os.Getenv("DATA_GEN_API_EXT_PORT")
	myConfig.DATA_SAVER_EXT_PORT1 = os.Getenv("DATA_SAVER_EXT_PORT")
	myConfig.DATA_RETRIEVE_EXT_PORT1 = os.Getenv("DATA_RETRIEVE_EXT_PORT")
	myConfig.DB_SQL_EXT_PORT1 = os.Getenv("DB_SQL_EXT_PORT")

	myConfig.DATA_GEN_API_EXT1 = os.Getenv("DATA_GEN_API_EXT")
	myConfig.DATA_SAVER_EXT1 = os.Getenv("DATA_SAVER_EXT")
	myConfig.DATA_RETRIEVE_EXT1 = os.Getenv("DATA_RETRIEVE_EXT")
	myConfig.DB_SQL_EXT1 = os.Getenv("DB_SQL_EXT")

	if myConfig.DATA_GEN_API_EXT1 == "" {
		myConfig.DATA_GEN_API_EXT1 = "data-gen-api"
	}
	if myConfig.DATA_GEN_API_EXT_PORT1 == "" {
		myConfig.DATA_GEN_API_EXT_PORT1 = "80"
	}

	rest_data_gen_api_url = "http://" + myConfig.DATA_GEN_API_EXT1 +
		":" + myConfig.DATA_GEN_API_EXT_PORT1

	if myConfig.DB_SQL_EXT1 == "" {
		myConfig.DB_SQL_EXT1 = "sql-server-db"
	}
	if myConfig.DB_SQL_EXT_PORT1 == "" {
		myConfig.DB_SQL_EXT_PORT1 = "1433"
	}
	server_sql = myConfig.DB_SQL_EXT1 // + ":" + myConfig.DB_SQL_EXT_PORT1
	sqlport, _ = strconv.Atoi(myConfig.DB_SQL_EXT_PORT1)


	fmt.Println(` server_sql,sqlport= `+server_sql+` `+ strconv.Itoa(sqlport) +`
				  DATA_GEN_API_EXT,DATA_GEN_API_EXT_PORT1 = `+ rest_data_gen_api_url)
}

var myConfig configty
var rest_data_gen_api_url string = "http://127.0.0.1:8020"

func main() {
	set_env()
	port1, err1 := strconv.Atoi(myConfig.DATA_RETRIEVE_EXT_PORT1)
	if port1 == 0 {
		port1 = 80
	}
	if err1 != nil {
		fmt.Println("WARNING: The env variable DATA_RETRIEVE_EXT_PORT is not provided, using port 80.")
	}
	fmt.Println("Going to listen on port  ", port1)

	handleRequests(port1)

}

func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Health Checking by connecting to the data gen API")
	fmt.Println("Health Checking by connecting to the data gen API")


	var defaultTtransport http.RoundTripper = &http.Transport{Proxy: nil}
	client := &http.Client{Transport: defaultTtransport}

	// Change the Timeout in the client to 5 seconds
	client.Timeout = 5 * time.Second

	req, err := http.NewRequest(http.MethodGet, rest_data_gen_api_url+"/", nil)
	if err != nil {
		fmt.Println(err.Error())
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}
	body, readErr := ioutil.ReadAll(resp.Body)
	aa := string(body)
	fmt.Println(resp, err, readErr, aa)
	fmt.Println("++++++++++++++++++++++++ FROM data gen API ++++++++++++++++++++++++++++")
	fmt.Println(aa)
	fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")

	fmt.Fprintf(w,`
	               ++++++++++++++++++++++++ FROM data gen API ++++++++++++++++++++++++++++
				   `+aa+
				`
				+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++`)
}


// Checks the connection to the SQL server
func check_sql_connection() {
	// Use background context
	ctx := context.Background()

	// Ping database to see if it's still alive.
	// Important for handling network issues and long queries.
	err := db.PingContext(ctx)
	if err != nil {
		log.Fatal("Error pinging database: " + err.Error())
	}

	var result string

	// Run query and scan for result
	err = db.QueryRowContext(ctx, "SELECT @@version").Scan(&result)
	if err != nil {
		log.Fatal("Scan failed:", err.Error())
	}
	fmt.Printf("%s\n", result)

	_, err = db.Query("USE master;")
	if err != nil {
		log.Fatal("Scan failed:", err.Error())
	}
	fmt.Printf("%s\n", result)
}

// Reads the articles from data-gen-api and Saves the articles in
// SQL server
func read_from_sql() []Article {

	var err error

	// Create connection string
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s",
		server_sql, user, password, sqlport, database)

	// Create connection pool
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: " + err.Error())
	}
	log.Printf("Connected!\n")
	check_sql_connection()

	// Use background context
	_, err = db.Query("USE master;")

	var myArtile Article
	var Articles []Article
	var Id = ""

	for i := 0; i < 100; i++ {
		Id = strconv.Itoa(i)
		sqlStatement := `SELECT 
		Id, Title, Description, Content
		 FROM ` + tablename + ` WHERE Id=$` + Id + `;`
		row := db.QueryRow(sqlStatement, 1)

		switch err := row.Scan(&myArtile.Id, &myArtile.Title, &myArtile.Description, &myArtile.Content); err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned!")
		case nil:
			if i == 0 {
				fmt.Println("-----------Data Retrieve Printing First Artcile----------------------")
				fmt.Println(myArtile)
			}
			Articles = append(Articles, myArtile)
		default:
			panic(err.Error())
		}
	}

	if err != nil {
		panic(err.Error())
	}

	// db.Close()
	defer db.Close()
	return Articles
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `Welcome to the Data-Retrieve-API!
	                /health  will check my health
	                /articles  will retrieve from DB`)
	fmt.Println("Endpoint Hit: homePage")
}
func returnAllArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllArticles")
	Articles := read_from_sql()
	fmt.Println("printing the first article ", Articles[4])
	json.NewEncoder(w).Encode(Articles)
}

func handleRequests(port int) {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/health", health)
	myRouter.HandleFunc("/articles", returnAllArticles)
	sport := ":" + strconv.Itoa(port)
	log.Fatal(http.ListenAndServe(sport, myRouter))
}
