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



                //mysqlへ接続。ドライバ名（mysql）と、ユーザー名・データソース(ここではgosample)を指定。
                db, err := sql.Open("mysql", "root@/ca_tech_dojo")
                log.Println("Connected to mysql.")

                //接続でエラーが発生した場合の処理
                if err != nil {
                    log.Fatal(err)
                }
                defer db.Close()

                //データベースへクエリを送信。引っ張ってきたデータがrowsに入る。
                rows, err := db.Query("INSERT USER VALUES (1, 'Satou', 'Kyoto');")
                defer rows.Close()
                if err != nil {
                    panic(err.Error())
                }

            default:
                // w.WriteHeader(http.StatusMethodNotAllowed)
                fmt.Fprint(w, "Method not allowed.\n")

                // ody := r.Body

        }
    })

    http.ListenAndServe(":8080", nil)
}