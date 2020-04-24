package main

import (
    "./page"

    "net/http"
)


func main() {



    http.HandleFunc("/user/create",page.Create)

    http.HandleFunc("/user/get",page.Get)

    http.HandleFunc("/user/update",page.Update)



    http.ListenAndServe(":8080", nil)
}