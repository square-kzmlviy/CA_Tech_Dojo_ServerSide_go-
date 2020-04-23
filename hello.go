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
    // "io/ioutil"


    "database/sql"
    "log"
    _ "github.com/go-sql-driver/mysql" 

    jwt "github.com/dgrijalva/jwt-go"
    "strconv"

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

type Response struct {
    TOKEN   string `json:"token"`
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
                rows_out_user, err_out_user := db.Query("SELECT id,name FROM user ORDER BY id DESC LIMIT "+"1;")
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

                var res Response
                //アルゴリズムの指定
                token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
                    "id":person.ID,
                    "name":person.Name,
                })

                
                //トークンに対して署名
                tokenString, err_JWT := token.SignedString([]byte(secretKey))
                if err_JWT == nil {

                    fmt.Fprint(w, "json hello!\n")

                    res.TOKEN = tokenString
                    res_json, _ := json.Marshal(res)

                    w.Write(res_json)

                    fmt.Printf("POST hello! %s \n", string(res_json))
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
                // w.WriteHeader(http.StatusOK)

                fmt.Fprint(w, "GET hello!\n")
                //レスポンスヘッダの解釈
                tokenString := r.Header.Get("x-token")
                fmt.Printf("%s\n", tokenString)



                token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
                    return []byte("himitu"), nil // CreateTokenにて指定した文字列を使います
                })

                if err != nil {
                    fmt.Fprint(w, "Not Json Token\n")
                }

                claims := token.Claims.(jwt.MapClaims)
                fmt.Println(claims["name"])

                //ここヤバイ理解できてない
                //https://medium.com/since-i-want-to-start-blog-that-looks-like-men-do/%E5%88%9D%E5%BF%83%E8%80%85%E3%81%AB%E9%80%81%E3%82%8A%E3%81%9F%E3%81%84interface%E3%81%AE%E4%BD%BF%E3%81%84%E6%96%B9-golang-48eba361c3b4
                var user InputJsonSchema
                // user.Name = claims["name"].(string)

                //interface型の元々float64を直して、それを文字列型に変更する値をチェック
                str_test := strconv.FormatFloat(claims["id"].(float64), 'G', 4, 64)

                fmt.Println(claims["id"].(float64))
                fmt.Println(str_test)



                 //mysqlへ接続。ドライバ名（mysql）と、ユーザー名・データソース(ここではgosample)を指定。
                db, err := sql.Open("mysql", "root@/ca_tech_dojo")
                log.Println("Connected to mysql.")

                 //接続でエラーが発生した場合の処理
                if err != nil {
                    log.Fatal(err)
                }
                defer db.Close()

                 //データベースへクエリを送信。引っ張ってきたデータがrowsに入る。
                rows, err := db.Query("SELECT name FROM user WHERE id = " + str_test + ";")
                defer rows.Close()
                if err != nil {
                    panic(err.Error())
                }


                for rows.Next() {
                    err = rows.Scan(&user.Name)

                    if err != nil {
                        panic(err.Error())
                    }
                    fmt.Printf("SQL確認 \n")
                    fmt.Println(user.Name)
                    //構造体をjsonへ
                    res_json, _ := json.Marshal(user)
                    fmt.Printf("レスポンス確認 \n")
                    fmt.Println(string(res_json))
                    w.Write(res_json)
                }
                
                




            default:
                // w.WriteHeader(http.StatusMethodNotAllowed)
                fmt.Fprint(w, "Method not allowed.\n")

            }
    })

    http.ListenAndServe(":8080", nil)
}