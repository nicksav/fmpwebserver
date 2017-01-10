package controllers

import (
	"net/http"
	"fmpwebserver/services/models"
)

func HelloController(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	bks, err := models.AllUsers()
    if err != nil {
        http.Error(w, http.StatusText(500), 500)
        return
    }
 	for _, bk := range bks {
		w.Write([]byte(bk.UserName +" "+ bk.Password))
    }


}
