package controllers

import (
	"fmpwebserver/services"
	"fmpwebserver/services/models"
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
)




func SignUp(w http.ResponseWriter, r *http.Request){
	requestUser := new(models.User)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestUser)

	userExists, err:= models.CheckUserExists(requestUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		if(userExists==nil){
			u, err:= models.SignUpUser(requestUser)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				code, err :=models.RequestConfirmationCode(u)
				if err != nil{
					w.WriteHeader(http.StatusInternalServerError)
				}else{
					//Sending email in goroutine
				 	go services.SendConfirmEmail(u,code)
				}
				json.NewEncoder(w).Encode(u)
			}
		}else{
				json.NewEncoder(w).Encode(userExists)
		}
	}
}


func SignUpKiosk(w http.ResponseWriter, r *http.Request){
	requestKiosk := new(models.Kiosk)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestKiosk)

	k,confirmCode,err:= models.SignUpKiosk(requestKiosk.UserName)
	if err!= nil {
		w.WriteHeader(http.StatusInternalServerError)
	}else{
		go services.SendConfirmEmailKiosk(k,confirmCode)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(k)
	}
}

func LinkKioskToAccount(w http.ResponseWriter, r *http.Request){
	requestKiosk := new(models.Kiosk)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestKiosk)

	responseStatus, token := services.LoginKiosk(requestKiosk)
	requestKiosk.Token = token
	k,err:= models.LinkKioskToAccount(requestKiosk)
	if err!= nil {
		w.WriteHeader(responseStatus)
	}else{
		if k!=nil{
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(k)
		} else {
			w.WriteHeader(http.StatusUnauthorized)			
		}		
	}
}



func ValidateConfirmCode(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
    code := vars["code"]


    usr, err :=models.VlidateConfirmCode(code)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
	}else{
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(usr)
	}
}



func Login(w http.ResponseWriter, r *http.Request) {
	requestUser := new(models.User)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestUser)

	responseStatus, token := services.Login(requestUser)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseStatus)
	w.Write(token)
}

func RefreshToken(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	requestUser := new(models.User)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestUser)

	w.Header().Set("Content-Type", "application/json")
	w.Write(services.RefreshToken(requestUser))
}

func Logout(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	err := services.Logout(r)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
