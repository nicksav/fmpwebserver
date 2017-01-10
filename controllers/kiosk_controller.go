package controllers

import (
	"net/http"
    "fmpwebserver/services"
	"fmpwebserver/services/models"
    "encoding/json"
    "github.com/gorilla/mux"
)

//RemoveKiosk - remove kiosk by Kiosk object
func RemoveKiosk (w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
    requestKiosk := new(models.Kiosk)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestKiosk)

    //Checking User claim
    usr, result:=services.CheckUserRequest(r,models.WebPortalMethod)
    if result==false{
        w.WriteHeader(http.StatusUnauthorized)
    }

    err :=models.RemoveKiosk(requestKiosk,usr)
    if err!=nil{
        w.WriteHeader(http.StatusInternalServerError)
    }

    //async revoke kiosk token
    go disableKiosk(requestKiosk)
    w.WriteHeader(http.StatusOK)
}

//UpdateKiosk - updating Kiosk by sending Kiosk object
func UpdateKiosk (w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
    
    requestKiosk := new(models.Kiosk)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestKiosk)

    //Checking User claim
    usr, result:=services.CheckUserRequest(r,models.WebPortalMethod)
    if result==false{
            w.WriteHeader(http.StatusUnauthorized)
    }

    k, err:=models.UpdateKiosk(requestKiosk,usr)
    if err!=nil{
       w.WriteHeader(http.StatusInternalServerError)
    }

    if requestKiosk.Status == 3 {
        go disableKiosk(requestKiosk)
    }

    w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(k)
}

func disableKiosk(k *models.Kiosk){
        token, err := models.GetTokenByKioskID(k)
        if err!=nil{
            return
        }
        services.LogoutKiosk(token)
}


//GetKiosks - return list of Kiosks for specific user
func GetKiosks(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {


    //Checking User claim
    usr, result:=services.CheckUserRequest(r,models.WebPortalMethod)
    if result==false{
            w.WriteHeader(http.StatusUnauthorized)
    }

    ks, err:=models.GetKiosks(usr)
    if err!=nil{
        w.WriteHeader(http.StatusInternalServerError)
    }

    w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ks)

}

//GetKiosk - return Kiosk object by id
func GetKiosk(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)
    id := vars["id"]

    usr, result:=services.CheckUserRequest(r,models.AllMethod)
    if result==false{
        w.WriteHeader(http.StatusUnauthorized)
    }

    k, err:=models.GetKioskByID(id,usr)
    if err!=nil{
        w.WriteHeader(http.StatusInternalServerError)
    }

    w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(k)
}
