package page

import (
	"encoding/json"
	"fmt"
	"net/http"
	"bytes"
	"io"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	jwt "github.com/dgrijalva/jwt-go"
	"strconv"
)

var secretKey = "himitu"

type UserName struct {
	Name string `json:"name"`
}

type NewUser struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Token struct {
	Token string `json:"token"`
}

func Create(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case http.MethodPost:
		
		body := r.Body
		defer body.Close()
		buf := new(bytes.Buffer)
		io.Copy(buf, body)
		var UserData UserName
		json.Unmarshal(buf.Bytes(), &UserData)
		fmt.Printf("POST hello! %s \n", string(UserData.Name))

		db, err := sql.Open("mysql", "root@/ca_tech_dojo")
		log.Println("Connected to mysql.")
		if err != nil {
			log.Print(err)
		}

		_, err = db.Exec("INSERT INTO USER(name,token) VALUES ('" + string(UserData.Name) + "','init');")
		if err != nil {
			log.Print(err)
		}		

		rows_out_user, err_out_user := db.Query("SELECT id,name FROM user ORDER BY id DESC LIMIT 1;")
		// defer rows_out_user.Close()
		if err_out_user != nil {
			log.Print(err_out_user)
		}
		var person NewUser
		for rows_out_user.Next() {

			err := rows_out_user.Scan(&person.ID, &person.Name)
			sample_json, _ := json.Marshal(person)

			if err != nil {
				log.Print(err_out_user)
			}
			fmt.Println(person.ID, person.Name)
			fmt.Println(string(sample_json))
		}

		//認証
		//アルゴリズムの指定
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":   person.ID,
			"name": person.Name,
		})

		//トークンに対して署名
		tokenString, err_JWT := token.SignedString([]byte(secretKey))
		if err_JWT == nil {
			fmt.Fprint(w, "json hello!\n")
			fmt.Printf("POST hello! %s \n", tokenString)

			var token Token
			token.Token = tokenString
			res_json, _ := json.Marshal(token)
			w.Header().Set("Content-Type", "application/json")
			w.Write(res_json)

		} else {
			fmt.Fprint(w, "Not Json hello!\n")
		}

	default:
		// w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "Method not allowed.\n")

		// ody := r.Body

	}
}

func Get(w http.ResponseWriter, r *http.Request) {
	// request bodyの読み取り

	//test
	fmt.Fprintf(w, "<h1>Test</h1>")

	switch r.Method {

	case http.MethodGet:
		// w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, "POST hello!\n")

		header_x_token := r.Header.Get("x-token")
		//tokenの解釈
		token, err := jwt.Parse(header_x_token, func(token *jwt.Token) (interface{}, error) {
			return []byte("himitu"), nil // CreateTokenにて指定した文字列を使います
		})
		if err != nil {
			log.Fatal(err)
		}
		claims := token.Claims.(jwt.MapClaims)

		var user NewUser
		user.ID = int(claims["id"].(float64))
		user.Name = claims["name"].(string)

		str_id := strconv.Itoa(user.ID)
		db, err := sql.Open("mysql", "root@/ca_tech_dojo")
		log.Println("Connected to mysql.")

		//接続でエラーが発生した場合の処理
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		rows, err := db.Query("SELECT name FROM user WHERE id = " + str_id + ";")
		defer rows.Close()
		if err != nil {
			panic(err.Error())
		}

		for rows.Next() {

			var res_data UserName
			err := rows.Scan(&res_data.Name)
			res_json, _ := json.Marshal(res_data)

			if err != nil {
				panic(err.Error())
			}
			fmt.Println(string(res_json))
			w.Header().Set("Content-Type", "application/json")
			w.Write(res_json)
		}

	default:
		fmt.Fprint(w, "Method not allowed.\n")
	}
}

func Update(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case http.MethodPut:
		fmt.Fprint(w, "POST hello!\n")

		body := r.Body
		defer body.Close()

		buf := new(bytes.Buffer)
		io.Copy(buf, body)

		var rename UserName
		json.Unmarshal(buf.Bytes(), &rename)

		header_x_token := r.Header.Get("x-token")

		//tokenの解釈
		token, err := jwt.Parse(header_x_token, func(token *jwt.Token) (interface{}, error) {
			return []byte("himitu"), nil
		})
		if err != nil {
			log.Fatal(err)
		}
		claims := token.Claims.(jwt.MapClaims)
		var user NewUser
		user.ID = int(claims["id"].(float64))
		user.Name = claims["name"].(string)

		str_id := strconv.Itoa(user.ID)
		db, err := sql.Open("mysql", "root@/ca_tech_dojo")
		log.Println("Connected to mysql.")

		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		rows, err := db.Query("UPDATE user SET name ='" + string(rename.Name) + "' WHERE id = " + str_id + ";")
		defer rows.Close()
		if err != nil {
			panic(err.Error())
		}

	default:
		fmt.Fprint(w, "Method not allowed.\n")

	}
}
