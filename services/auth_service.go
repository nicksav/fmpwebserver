package services

import (  
    "fmpwebserver/api/parameters"
    "fmpwebserver/core/authentication"
    "fmpwebserver/services/models"
    "encoding/json"
    jwt "github.com/dgrijalva/jwt-go"
    request "github.com/dgrijalva/jwt-go/request"
    "net/http"
    "github.com/gorilla/context"
)

func Login(requestUser *models.User) (int, []byte) {  
    authBackend := authentication.InitJWTAuthenticationBackend()
    if authBackend.Authenticate(requestUser) {
        token, err := authBackend.GenerateToken(requestUser.UserName)
        if err != nil {
            return http.StatusInternalServerError, []byte("")
        } else {
            response, _ := json.Marshal(parameters.TokenAuthentication{token})
            return http.StatusOK, response
        }
    }

    return http.StatusUnauthorized, []byte("")
}


func LoginKiosk(requestKiosk *models.Kiosk) (int, string) {  
    authBackend := authentication.InitJWTAuthenticationBackend()
        token, err := authBackend.GenerateTokenKiosk(requestKiosk.ID,requestKiosk.UserName)
        if err != nil {
            return http.StatusInternalServerError, ""
        }

        return http.StatusOK, token
}

func RefreshToken(requestUser *models.User) []byte {  
    authBackend := authentication.InitJWTAuthenticationBackend()
    token, err := authBackend.GenerateToken(requestUser.UserName)
    if err != nil {
        panic(err)
    }
    response, err := json.Marshal(parameters.TokenAuthentication{token})
    if err != nil {
        panic(err)
    }
    return response
}

func Logout(req *http.Request) error {  
    authBackend := authentication.InitJWTAuthenticationBackend()
    tokenRequest, err := request.ParseFromRequest(req,request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
        return authBackend.PublicKey, nil
    })
    if err != nil {
        return err
    }
    tokenString := req.Header.Get("Authorization")
    return authBackend.Logout(tokenString, tokenRequest)
}


func LogoutKiosk(tokenString string) error {  
   authBackend := authentication.InitJWTAuthenticationBackend()
   token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) { 
        return authBackend.PublicKey, nil
    })
    
    if err != nil {
        return err
    }

    claims := token.Claims.(jwt.MapClaims)
    return authBackend.LogoutKiosk(tokenString, claims["exp"])
}

func CheckUserRequest(r *http.Request,methodType int)(string,bool){
    swe := context.Get(r, "swe")
    token:=swe.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

    typ, ok := claims["typ"].(float64)
    if !ok {
       return "", false
    }

    if (typ != float64(methodType) && methodType!=3) {
        return "", false
    }

    if typ == 1 {
        
        usr, ok := claims["sub"].(string)
        if !ok {
            return "", false
        }

        return usr, true
    } else if typ == 2{
        
        usr, ok := claims["usr"].(string)
        if !ok {
            return "", false
        }

         return usr, true
    }

    return "", false

}

