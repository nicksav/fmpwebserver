package main

import (  
    "fmpwebserver/routers"
    "fmpwebserver/settings"
    "github.com/urfave/negroni"
    "net/http"
)

import  "fmpwebserver/services/models"

func main() {  
    settings.Init()
    models.InitDB()
    router := routers.InitRoutes()
    n := negroni.Classic()
    n.UseHandler(router)

    http.ListenAndServe(":5000", n)
    
}