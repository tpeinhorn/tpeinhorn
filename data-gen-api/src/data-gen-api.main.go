// data-gen-api.go
package main

import (
	"encoding/json"
	"fmt"
	_ "fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	// "path/filepath"
	"reflect"
	"strconv"
	"time"
	"bufio"
    "flag"

	"github.com/gorilla/mux"
)

// Article - Our struct for all articles
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

var rest_data_retrieve_api_url string = "http://127.0.0.1:8021"
var server_sql = "127.0.0.1"
var sqlport = 11433

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func set_env() {
	myConfig.DATA_GEN_API_EXT_PORT1 = os.Getenv("DATA_GEN_API_EXT_PORT")
	myConfig.DATA_SAVER_EXT_PORT1 = os.Getenv("DATA_SAVER_EXT_PORT")
	myConfig.DATA_RETRIEVE_EXT_PORT1 = os.Getenv("DATA_RETRIEVE_EXT_PORT")
	myConfig.DB_SQL_EXT_PORT1 = os.Getenv("DB_SQL_EXT_PORT")

	myConfig.DATA_GEN_API_EXT1 = os.Getenv("DATA_GEN_API_EXT")
	myConfig.DATA_SAVER_EXT1 = os.Getenv("DATA_SAVER_EXT")
	myConfig.DATA_RETRIEVE_EXT1 = os.Getenv("DATA_RETRIEVE_EXT")
	myConfig.DB_SQL_EXT1 = os.Getenv("DB_SQL_EXT")


	if myConfig.DATA_RETRIEVE_EXT1 == "" {
		myConfig.DATA_RETRIEVE_EXT1 = "data-retrieve"
	}
	if myConfig.DATA_RETRIEVE_EXT_PORT1 == "" {
		myConfig.DATA_RETRIEVE_EXT_PORT1 = "80"
	}

	rest_data_retrieve_api_url = "http://" + myConfig.DATA_RETRIEVE_EXT1 +
		":" + myConfig.DATA_RETRIEVE_EXT_PORT1

	server_sql = myConfig.DB_SQL_EXT1 // + ":" + myConfig.DB_SQL_EXT_PORT1
	sqlport, _ = strconv.Atoi(myConfig.DB_SQL_EXT_PORT1)

	
	fmt.Println(` server_sql,sqlport= `+server_sql+` `+ strconv.Itoa(sqlport) +`
	              rest_data_retrieve_api_url = `+ rest_data_retrieve_api_url)
}

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

var Articles []Article

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `Welcome to the Data-Gen-API!
	                /articles   will show all articles
					/article/3  will show article number 3
					/articles_ret  will get articles from the retrive api
					/articles_compare  will compare what got from the retrieve api
					
					The data Retrieve URL is:`+ rest_data_retrieve_api_url +
					`
					`)
	fmt.Println("Endpoint Hit: homePage")
}

func returnAllArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllArticles")
	json.NewEncoder(w).Encode(Articles)
	fmt.Println("finished json encode")
}

func returnAllArticlesRet(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllArticles")
	Articles_ret := get_data_from_data_retrieve()
	json.NewEncoder(w).Encode(Articles_ret)
	fmt.Println("finished json encode")
}

func returnAllArticlesCompare(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllArticlesCompare")
	Articles_ret := get_data_from_data_retrieve()
	var if_same = true
	a := len(Articles)
	b := len(Articles_ret)
	//Articles_ret[20].Content = "My content"
	var i = 0
	if a != b {
		if_same = false
	}
	for i = 0; i < a; i++ {
		if_same = (reflect.DeepEqual(Articles[i], Articles_ret[i])) && if_same
		if !if_same {
			break
		}
	}

	var same_st = "GOOD NEWS! The Articles and Articles_ret are the same."
	if !if_same {
		same_st = "ERROR: The Articles and Articles_ret are not the same at Article " + strconv.Itoa(i)
	}
	fmt.Println(same_st)
	fmt.Fprintf(w, same_st)

	_ = Articles_ret
}

func returnSingleArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnSingleArticle")
	vars := mux.Vars(r)
	key := vars["Id"]

	for _, article := range Articles {
		if article.Id == key {
			json.NewEncoder(w).Encode(article)
		}
	}
}

func createNewArticle(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// unmarshal this into a new Article struct
	// append this to our Articles array.
	fmt.Println("Endpoint Hit: createNewArticle")
	reqBody, _ := ioutil.ReadAll(r.Body)
	var article Article
	json.Unmarshal(reqBody, &article)
	// update our global Articles array to include
	// our new Article
	Articles = append(Articles, article)

	json.NewEncoder(w).Encode(article)
}

func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Health Checking by connecting to the data retrieve API")
	fmt.Println("Health Checking by connecting to the data retrieve API")


	var defaultTtransport http.RoundTripper = &http.Transport{Proxy: nil}
	client := &http.Client{Transport: defaultTtransport}

	// Change the Timeout in the client to 5 seconds
	client.Timeout = 5 * time.Second

	req, err := http.NewRequest(http.MethodGet, rest_data_retrieve_api_url+"/", nil)
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
	fmt.Println("++++++++++++++++++++++++ FROM data Retrieve API ++++++++++++++++++++++++++++")
	fmt.Println(aa)
	fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")

	fmt.Fprintf(w,`
	               ++++++++++++++++++++++++ FROM data Retrieve API ++++++++++++++++++++++++++++
				   `+aa+
				`
				+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++`)
}

func writeFile(fSize int64, root string) (string,error) {
    
	speed_st := " This speed was not measured"
	fName := `/home/diskio` // test file
    defer os.Remove(fName)
    f, err := os.Create(fName)
    if err != nil {
        return speed_st,err
    }
    const defaultBufSize = 4096
    buf := make([]byte, defaultBufSize)
    buf[len(buf)-1] = '\n'
    w := bufio.NewWriterSize(f, len(buf))

    start := time.Now()
    written := int64(0)
    for i := int64(0); i < fSize; i += int64(len(buf)) {
        nn, err := w.Write(buf)
        written += int64(nn)
        if err != nil {
            return speed_st,err
        }
    }
    err = w.Flush()
    if err != nil {
        return speed_st,err
    }
    err = f.Sync()
    if err != nil {
        return speed_st,err
    }
    since := time.Since(start)

    err = f.Close()
    if err != nil {
        return speed_st,err
    }
    fmt.Printf("written: %dB %dns %.2fGB %.2fs %.2fMB/s\n",
        written, since,
        float64(written)/1000000000, float64(since)/float64(time.Second),
        (float64(written)/1000000)/(float64(since)/float64(time.Second)),
    )

	speed_st = fmt.Sprintf("written: %dB %dns %.2fGB %.2fs %.2fMB/s\n",
	written, since,
	float64(written)/1000000000, float64(since)/float64(time.Second),
	(float64(written)/1000000)/(float64(since)/float64(time.Second)),
    )

    return speed_st, nil
}

var size = flag.Int("size", 2, "file size in GiB")

func disk_speed(w http.ResponseWriter, r *http.Request) {

	root := "/"
	vars := mux.Vars(r)
	root = root + vars["Path"]

	fmt.Println("Going to check the speed in ",root)
	fmt.Fprintf(w,"Going to check the speed in " + root)
    flag.Parse()
    fSize := int64(*size) * (1024 * 1024 * 1024)
    speed_st,err := writeFile(fSize,root)
    if err != nil {
        fmt.Fprintln(os.Stderr, fSize, err)
		fmt.Fprintf(w,err.Error())	
    }
	fmt.Fprintf(w,"Speed test:====" + speed_st)
}


// Simply lists files in the mounted file system
func list_files(w http.ResponseWriter, r *http.Request)  {
	root := "/"
	vars := mux.Vars(r)
	root = root + vars["Path"]

	fmt.Fprintf(w, `Listing files in the directory: ` + root+ `  
	`)
	fmt.Println   (`Listing files in the directory: ` + root + ` 
	`)

    // var files []string

    
    files, err := ioutil.ReadDir(root)
    if err != nil {
        fmt.Println (err.Error())
    }
 
    for _, f := range files {
            fmt.Println(f.Name())
			fmt.Fprintf(w,f.Name()+" \n ")
    }
}

func deleteArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: deleteArticle")
	vars := mux.Vars(r)
	Id := vars["Id"]

	for index, article := range Articles {
		if article.Id == Id {
			Articles = append(Articles[:index], Articles[index+1:]...)
		}
	}

}

func handleRequests(port int) {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/health", health)
	myRouter.HandleFunc("/disk_speed/{Path}", disk_speed)
	myRouter.HandleFunc("/list_files/{Path}", list_files)
	myRouter.HandleFunc("/articles", returnAllArticles)
	myRouter.HandleFunc("/articles_ret", returnAllArticlesRet)
	myRouter.HandleFunc("/articles_compare", returnAllArticlesCompare)
	myRouter.HandleFunc("/article", createNewArticle).Methods("POST")
	myRouter.HandleFunc("/article/{Id}", deleteArticle).Methods("DELETE")
	myRouter.HandleFunc("/article/{Id}", returnSingleArticle)

	sport := ":" + strconv.Itoa(port)
	log.Fatal(http.ListenAndServe(sport, myRouter))
}

var myConfig configty

func main() {

	// var myConfig configty
	//Get the enviroment variables first
	set_env()

	Articles = []Article{}
	var aa = ""
	var Id = ""
	var myArticle Article
	for i := 0; i < 100; i++ {
		aa = RandStringRunes(20)
		Id = strconv.Itoa(i)
		myArticle.Id = Id
		myArticle.Title = "my article Title " + Id + aa
		myArticle.Description = "Description for article " + Id + aa
		myArticle.Content = "Content for article " + Id + "  " + aa
		Articles = append(Articles, myArticle)
	}

	port1, err1 := strconv.Atoi(myConfig.DATA_GEN_API_EXT_PORT1)

	if port1 == 0 {
		port1 = 80
	}
	if err1 != nil {
		fmt.Println("WARNING: The env variable DATA_GEN_API_EXT_PORT is not provided, using port 80.")
	}
	fmt.Println(`Welcome to the Data-Gen-API!
	/health  Checks my health
	/list_files/{Path} List files in the directory
	/disk_speed/{Path}  will run the disk speed checks
	/articles   will show all articles
	/article/3  will show article number 3
	/articles_ret  will get articles from the retrive api
	/articles_compare  will compare what got from the retrieve api

	
	The data Retrieve URL is: `+ rest_data_retrieve_api_url +
	`
	`)
	handleRequests(port1)

}

// Gets data for the articles from the data-generator API
func get_data_from_data_retrieve() []Article {

	// Set up the initial client
	//client := &http.Client{}
	var defaultTtransport http.RoundTripper = &http.Transport{Proxy: nil}
	client := &http.Client{Transport: defaultTtransport}

	// Change the Timeout in the client to 5 seconds
	client.Timeout = 5 * time.Second

	req, err := http.NewRequest(http.MethodGet, rest_data_retrieve_api_url+"/articles", nil)
	if err != nil {
		panic(err.Error())
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	body, readErr := ioutil.ReadAll(resp.Body)
	_ = readErr
	Articles_ret := make([]Article, 0)
	json.Unmarshal(body, &Articles_ret)
	fmt.Printf("\n\n first artcile retrieve json object:::: \n %+v\n", Articles_ret[0])
	return Articles_ret
}
