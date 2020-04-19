package main

import (
    "encoding/json"
    "fmt"
    // "io/ioutil"
    "net/http"
    // "Request.Body"

    // ?
    "bytes"
    // ?
    "io"


    "database/sql"
    "log"
    _ "github.com/go-sql-driver/mysql" 
)

// jsonのSchema
type InputJsonSchema struct {
    Name string `json:"name"`
}

func main() {
    http.HandleFunc("/user/create", func(w http.ResponseWriter, r *http.Request) {
    // request bodyの読み取り

         //test
        fmt.Fprintf(w, "<h1>Test</h1>")

        switch r.Method {
            case http.MethodGet:
                w.WriteHeader(http.StatusOK)
                fmt.Fprint(w, "GET hello!\n")


            case http.MethodPost:
                // w.WriteHeader(http.StatusCreated)
                fmt.Fprint(w, "POST hello!\n")

                body := r.Body
                defer body.Close()

                buf := new(bytes.Buffer)
                io.Copy(buf, body)


                var hello InputJsonSchema
                json.Unmarshal(buf.Bytes(), &hello)

                // w.WriteHeader(http.StatusCreated)
                // fmt.Fprint(w, "POST hello! %v \n", hello.Name)
                fmt.Printf("POST hello! %v \n", hello)

            default:
                // w.WriteHeader(http.StatusMethodNotAllowed)
                fmt.Fprint(w, "Method not allowed.\n")

                // ody := r.Body

        }
    })

    http.ListenAndServe(":8080", nil)
}