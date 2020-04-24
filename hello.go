package main

import (
    "./page"

    "net/http"
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
type Token struct {
    Token string `json:"token"`
}

func main() {



    http.HandleFunc("/user/create",page.Create)

    http.HandleFunc("/user/get",page.Get)

    http.HandleFunc("/user/update",page.Update)



    http.ListenAndServe(":8080", nil)
}