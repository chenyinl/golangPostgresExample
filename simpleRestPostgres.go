package main
 
import (
    "fmt"
    "log"
    "net/http"
    "database/sql"
    "encoding/json"
    "io/ioutil"
    //"os"
    //"time"
    
    _ "github.com/lib/pq"
    "github.com/gorilla/mux"
)

var dbInfo string

const (
    DB_USER     = "USERNAME"
    DB_PASSWORD = "PASSWORD"
    DB_NAME     = "db"
)

var dbh *sql.DB

func main() {
    dbInfo = fmt.Sprintf(
        "host=localhost port=5432 user=%s password=%s dbname=%s sslmode=disable",
        DB_USER, 
        DB_PASSWORD, 
        DB_NAME)

    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/template/{templateid}/{method}", PutTemplate).Methods("PUT")
    router.HandleFunc("/template/{templateid}/{method}", PostTemplate).Methods("POST")
    router.HandleFunc("/template/{templateid}/{method}", GetTemplate).Methods("GET")
    router.HandleFunc("/template/all", GetAllTemplates).Methods("GET")
    router.HandleFunc("/ping", GetPing).Methods("GET")
    router.HandleFunc("/test", GetTest).Methods("GET")
    log.Fatal(http.ListenAndServe("dev-qa.chenl.sbx.4over.com:82", router))
}

func InitDBHandler(){
    err := *new(error)
    dbh, err = sql.Open("postgres",dbInfo)
    if err != nil {
        panic(err)
    }
   // defer dbh.Close()
}


func GetTest(w http.ResponseWriter, r *http.Request){
    w.Write([]byte("Test"))
    InitDBHandler()
    /*
    err := *new(error)
    dbh, err = sql.Open("postgres",dbInfo)
    if err != nil {
        panic(err)
    }
    defer dbh.Close()
    */
    Part2(w, r)
}
    
    
    
func Part2(w http.ResponseWriter, r *http.Request){
    query := "SELECT environment_uuid, environment, port FROM public.where_am_i()"
    rows, err := dbh.Query(query)
    if err != nil {
        panic(err)
    }
    var uuid string
    var environment string
    var port string
    rows.Next();
    err = rows.Scan(&uuid, &environment, &port)

    m:=  map[string]string{
        "environment_uuid": string(uuid),
        "environment": string(environment),
        "port": string(port)}

    b,err:=json.Marshal(m);
    if err != nil {
        panic(err)
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(b)
    
}
func GetPing(w http.ResponseWriter, r *http.Request) {
    /*
    type Pingdata struct{
        Environment_uuid string
        Environment string
        Port string
    }*/
    query := "SELECT environment_uuid, environment, port FROM public.where_am_i()"
        
     /* connect to DB */
    dbhandler, err := sql.Open("postgres", dbInfo)
    err = dbhandler.Ping()
    if err != nil {
        panic(err)
    }
    defer dbhandler.Close()
    
    /* run the query */
    rows, err := dbhandler.Query(query)
    
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
    }

    var uuid string
    var environment string
    var port string
    rows.Next();
    err = rows.Scan(&uuid, &environment, &port)

   /* different approach: use struct
    m := &Pingdata{
        Environment_uuid: string(uuid),
        Environment: string(environment),
        Port: string(port)}*/

    m:=  map[string]string{
        "environment_uuid": string(uuid),
        "environment": string(environment),
        "port": string(port)}

    b,err:=json.Marshal(m);
    if err != nil {
        panic(err)
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(b)
}

func PostTemplate(w http.ResponseWriter, r *http.Request) {
    var query string;
    vars := mux.Vars(r)
    templateId := vars["templateid"]
    method := vars["method"]

    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        panic(err)
    }
    bodystring := string(body)
    query = "INSERT INTO  automatedtesting.response_templates (id, methods, template_data)"+
        "values('"+templateId+"','"+method+"','"+bodystring+"')" 
        
     /* connect to DB */
    dbhandler, err := sql.Open("postgres", dbInfo)
    err = dbhandler.Ping()
    if err != nil {
        panic(err)
    }
    defer dbhandler.Close()
    
    /* run the query */
    _, err = dbhandler.Query(query)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
    }
}
 
func GetTemplate(w http.ResponseWriter, r *http.Request) {

    var query string;
    vars := mux.Vars(r)
    templateId := vars["templateid"] 
    
    query = "select template_data from automatedtesting.response_templates where id='"+templateId+"'"

    /* connect to DB */
    dbhandler, err := sql.Open("postgres", dbInfo)
    err = dbhandler.Ping()
    if err != nil {
        panic(err)
    }
    defer dbhandler.Close()
    
    /* run the query */
    rows, err := dbhandler.Query(query)
    // fmt.Printf("Type: %T Value: %v\n",rows, rows)
    
    for rows.Next(){
      var temp_data string
      err = rows.Scan(&temp_data)
      fmt.Fprintln(w,temp_data)
    }
}

func PutTemplate(w http.ResponseWriter, r *http.Request) {

    var query string;
    vars := mux.Vars(r)
    templateId := vars["templateid"]
    method := vars["method"]

    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        panic(err)
    }
    bodystring := string(body)
    query = "UPDATE  automatedtesting.response_templates "+
        "SET template_data='"+bodystring+"'"+
        " methods='"+method+"'"+
        "' WHERE id='"+templateId+"'"
        
     /* connect to DB */
    dbhandler, err := sql.Open("postgres", dbInfo)
    err = dbhandler.Ping()
    if err != nil {
        panic(err)
    }
    defer dbhandler.Close()
    
    /* run the query */
    _, err = dbhandler.Query(query)
    
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
    }
}

func GetAllTemplates(w http.ResponseWriter, r *http.Request) {

    var query string;
    //vars := mux.Vars(r)
    
    query = "select * from automatedtesting.response_templates"

    /* connect to DB */
    dbhandler, err := sql.Open("postgres", dbInfo)
    err = dbhandler.Ping()
    if err != nil {
        panic(err)
    }
    defer dbhandler.Close()
    
    /* run the query */
    rows, err := dbhandler.Query(query)
    // fmt.Printf("Type: %T Value: %v\n",rows, rows)
    

    if err != nil {
      panic(err)
    }
    defer rows.Close()
    columns, err := rows.Columns()
    if err != nil {
      panic( err)
    }
    count := len(columns)
    tableData := make([]map[string]interface{}, 0)
    values := make([]interface{}, count)
    valuePtrs := make([]interface{}, count)
    for rows.Next() {
      for i := 0; i < count; i++ {
          valuePtrs[i] = &values[i]
      }
      rows.Scan(valuePtrs...)
      entry := make(map[string]interface{})
      for i, col := range columns {
          var v interface{}
          val := values[i]
          b, ok := val.([]byte)
          if ok {
              v = string(b)
          } else {
              v = val
          }
          entry[col] = v
      }
      tableData = append(tableData, entry)
    }
    jsonData, err := json.Marshal(tableData)
        if err != nil {
                fmt.Println(w, err)
        }
    //fmt.Println(w,string(jsonData))
    w.Header().Set("Content-Type", "application/json")
    w.Write(jsonData)
   //fmt.Fprintln(string(jsonData), nil 
}
/*
func getJSONIsql(String string)(string,error){
rows, err := db.Query(sqlString)
  if err != nil {
      return "", err
  }
  defer rows.Close()
  columns, err := rows.Columns()
  if err != nil {
      return "", err
  }
  count := len(columns)
  tableData := make([]map[string]interface{}, 0)
  values := make([]interface{}, count)
  valuePtrs := make([]interface{}, count)
  for rows.Next() {
      for i := 0; i < count; i++ {
          valuePtrs[i] = &values[i]
      }
      rows.Scan(valuePtrs...)
      entry := make(map[string]interface{})
      for i, col := range columns {
          var v interface{}
          val := values[i]
          b, ok := val.([]byte)
          if ok {
              v = string(b)
          } else {
              v = val
          }
          entry[col] = v
      }
      tableData = append(tableData, entry)
  }
  jsonData, err := json.Marshal(tableData)
  if err != nil {
      return "", err
  }
  fmt.Println(string(jsonData))
  return string(jsonData), nil 
}
*/

func respondWithError(w http.ResponseWriter, code int, message string) {
    respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    response, _ := json.Marshal(payload)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(response)
}
