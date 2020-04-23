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

    jwt "github.com/dgrijalva/jwt-go"
)
var secretKey = "himitu"

// jsonのSchema
type InputJsonSchema struct {
    Name string `json:"name"`
}

type NewUser struct {
    ID   int `json:"id"`
    Name string `json:"name"`
}

func main() {



    http.HandleFunc("/user/create", func(w http.ResponseWriter, r *http.Request) {
    // request bodyの読み取り

         //test
        fmt.Fprintf(w, "<h1>Test</h1>")

        switch r.Method {


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
                fmt.Printf("POST hello! %s \n", string(hello.Name))


                //mysqlへ接続。ドライバ名（mysql）と、ユーザー名・データソース(ここではgosample)を指定。
                db, err := sql.Open("mysql", "root@/ca_tech_dojo")
                log.Println("Connected to mysql.")

                //接続でエラーが発生した場合の処理
                if err != nil {
                    log.Fatal(err)
                }
                defer db.Close()

                //データベースへクエリを送信。引っ張ってきたデータがrowsに入る。
                // rows, err := db.Query("INSERT USER VALUES (1, " + string(hello.Name) + ",'Kyoto');")
                rows, err := db.Query("INSERT INTO USER(name,token) VALUES ('" + string(hello.Name) + "','init');")
                // rows, err := db.Query(fmt.Sprintf("INSERT USER VALUES (1, " + "Satou" + ", 'Kyoto');"))
                defer rows.Close()
                if err != nil {
                    panic(err.Error())
                }


                // // INSERTしたレコードを読めるかテスト
                // LIMIT := 1
                rows_out_user, err_out_user := db.Query("SELECT id,name FROM user ORDER BY id DESC LIMIT 1;")
                defer rows_out_user.Close()
                if err_out_user != nil {
                    panic(err_out_user.Error())
                }
                var person NewUser //構造体Person型の変数personを定義
                //レコード一件一件をあらかじめ用意しておいた構造体に当てはめていく。
                for rows_out_user.Next() {
                    
                    err := rows_out_user.Scan(&person.ID, &person.Name)
                    sample_json, _ := json.Marshal(person)
            
                    if err != nil {
                        panic(err_out_user.Error())
                    }
                    fmt.Println(person.ID, person.Name)
                    fmt.Println(string(sample_json))
                }
            
                


                //認証
                //アルゴリズムの指定
                token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
                    "id":person.ID,
                    "name":person.Name,
                })

                
                //トークンに対して署名
                tokenString, err_JWT := token.SignedString([]byte(secretKey))
                if err_JWT == nil {
                    fmt.Fprint(w, "json hello!\n")
                    fmt.Printf("POST hello! %s \n", tokenString)
                } else {
                    fmt.Fprint(w, "Not Json hello!\n")
                }


            default:
                // w.WriteHeader(http.StatusMethodNotAllowed)
                fmt.Fprint(w, "Method not allowed.\n")

                // ody := r.Body

        }
    })

    http.HandleFunc("/user/get", func(w http.ResponseWriter, r *http.Request) {
        // request bodyの読み取り
    
             //test
            fmt.Fprintf(w, "<h1>Test</h1>")
    
            switch r.Method {
    
    
                case http.MethodGet:
                    // w.WriteHeader(http.StatusCreated)
                    fmt.Fprint(w, "POST hello!\n")
    
                    header_x_token := r.Header.Get("x-token")
                    // //headerが読み込まれている確認
                    // fmt.Println(header_x_token)


                    //tokenの解釈
                    token, err := jwt.Parse(header_x_token, func(token *jwt.Token) (interface{}, error) {
                        return []byte("himitu"), nil // CreateTokenにて指定した文字列を使います
                    })
                    if err != nil {
                        log.Fatal(err)
                    }
                    claims := token.Claims.(jwt.MapClaims)
                    //解釈したPAYLOAD:DATA ["name"]の確認
                    fmt.Println(claims["name"].(string))


    
    
                default:
                    // w.WriteHeader(http.StatusMethodNotAllowed)
                    fmt.Fprint(w, "Method not allowed.\n")
    
                    // ody := r.Body
    
            }
        })



    http.ListenAndServe(":8080", nil)
}