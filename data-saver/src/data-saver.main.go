package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gorilla/mux"
)

// Replace with your own connection parameters for the sql server
var server_sql = "127.0.0.1"
var sqlport = 11433
var user = "sa"
var password = "MyStrongPassword123"
var database = "master"
var tablename = "Articles"

var db *sql.DB
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

var rest_data_gen_api_url string = "http://127.0.0.1:8020"

type Article struct {
	Id          string `json:"Id"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
	Content     string `json:"Content"`
}

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

func main() {
	set_env()

	// Close the database connection pool after program executes

	// DATA_SAVER_EXT_PORT1 := os.Getenv("DATA_SAVER_EXT_PORT")
	port1, err1 := strconv.Atoi(myConfig.DATA_SAVER_EXT_PORT1)
	if port1 == 0 {
		port1 = 80
	}
	if err1 != nil {
		fmt.Println("WARNING: The env variable DATA_SAVER_EXT_PORT is not provided, using port 80.")
	}
	log.Println("Now listening on ", port1)
	handleRequests(port1)
	defer db.Close()

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
func save_in_sql(Articles []Article) {

	var result string

	// Use background context
	ctx := context.Background()
	_, err := db.Query("USE master;")
	_, err = db.Query("GO")

	// create another table
	qrr := `
	IF EXISTS (SELECT 1 
		FROM INFORMATION_SCHEMA.TABLES 
		WHERE TABLE_TYPE='BASE TABLE' 
		AND TABLE_NAME='` + tablename + `') 
	SELECT 1 AS res ELSE SELECT 0 AS res;
	`
	err = db.QueryRowContext(ctx, qrr).Scan(&result)
	if err != nil {
		fmt.Println(err.Error())
	}

	// Create the table if it does not exist
	if result == "0" {
		mm := `
		CREATE TABLE ` + tablename + ` (
			Id int NOT NULL,
			Title varchar(255) NOT NULL,
			Description  varchar(255),
			Content varchar(255),
			PRIMARY KEY (Id)
		);
		`
		_ = mm
		_, err = db.Query(mm)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	myqu1 := `DELETE FROM ` + tablename + `;`
	_, err = db.Query(myqu1)

	type Article struct {
		Id          string `json:"Id"`
		Title       string `json:"Title"`
		Description string `json:"Description"`
		Content     string `json:"Content"`
	}

	//fill in some data in the table
	var myArtile Article
	if err != nil {
		fmt.Println(err.Error())
		panic(err.Error())
	}

	_ = myqu1
	for i := 0; i < 100; i++ {
		myqu1 := `INSERT INTO ` + tablename + ` VALUES (` +
			Articles[i].Id + ", " +
			"'" + Articles[i].Title + "', " +
			"'" + Articles[i].Description + "', " +
			"'" + Articles[i].Content + "' " +
			");"
		fmt.Println(myqu1)
		_, err = db.Query(myqu1)
	}

	sqlStatement := `SELECT 
	 				 Id, Title, Description, Content
	  			     FROM ` + tablename + ` WHERE Id=$329;`
	row := db.QueryRow(sqlStatement, 1)

	switch err := row.Scan(&myArtile.Id, &myArtile.Title, &myArtile.Description, &myArtile.Content); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	case nil:
		fmt.Println("--------------------------------------------------")
		fmt.Println(myArtile)
	default:
		panic(err.Error())
	}

	if err != nil {
		panic(err.Error())
	}

	db.Close()
	return
}

// Gets data for the articles from the data-generator API
func get_data_from_api() []Article {

	// Set up the initial client
	//client := &http.Client{}
	var defaultTtransport http.RoundTripper = &http.Transport{Proxy: nil}
	client := &http.Client{Transport: defaultTtransport}

	// Change the Timeout in the client to 5 seconds
	client.Timeout = 5 * time.Second

	req, err := http.NewRequest(http.MethodGet, rest_data_gen_api_url+"/articles", nil)
	if err != nil {
		panic(err.Error())
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	body, readErr := ioutil.ReadAll(resp.Body)
	aa := string(body)
	fmt.Println(resp, err, readErr, aa)
	fmt.Println("+++++++++++++++++++++++++++++++++")
	fmt.Println(aa)

	Articles := make([]Article, 0)
	json.Unmarshal(body, &Articles)
	fmt.Printf("\n\n json object:::: \n %+v\n", Articles)
	return Articles
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `Welcome to the Data-Saver-API! 
	                /health will check my health
	                /saveAllArticles  will save all the articles in the SQL DB`)
	fmt.Println("Endpoint Hit: homePage")
}
func saveAllArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Going to save artiles..")
	fmt.Println("Going to save artiles..")
	var err error

	// Create connection string
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s",
		server_sql, user, password, sqlport, database)

	// Create connection pool
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: " + err.Error())
	}
	log.Println("Connected to sql server on ", server_sql, sqlport)

	check_sql_connection()

	Articles := get_data_from_api()
	save_in_sql(Articles)
	// json.NewEncoder(w).Encode(Articles)
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

	fmt.Fprintf(w,"++++++++++++++++++++++++ FROM data gen API ++++++++++++++++++++++++++++")
	fmt.Fprintf(w,aa)
	fmt.Fprintf(w,"+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
}

func handleRequests(port1 int) {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/health", health)
	myRouter.HandleFunc("/saveAllArticles", saveAllArticles)
	sport := ":" + strconv.Itoa(port1)
	log.Fatal(http.ListenAndServe(sport, myRouter))
}
