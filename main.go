package main

import (  
    "fmpwebserver/routers"
    "fmpwebserver/settings"
    "fmpwebserver/services/models"
    "github.com/urfave/negroni"
    "net/http"
)

func main() {  
    settings.Init()
    models.InitDB()
    router := routers.InitRoutes()
    n := negroni.Classic()
    n.UseHandler(router)
    http.ListenAndServe(":5000", n)
}