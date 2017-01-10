package routers

import (
	"fmpwebserver/controllers"
	"fmpwebserver/core/authentication"
	"github.com/urfave/negroni"
	"github.com/gorilla/mux"
)

func SetAuthenticationRoutes(router *mux.Router) *mux.Router {
	router.HandleFunc("/token-auth", controllers.Login).Methods("POST")
	router.HandleFunc("/signup", controllers.SignUp).Methods("POST")
	router.HandleFunc("/signupkiosk", controllers.SignUpKiosk).Methods("POST")
	router.HandleFunc("/linkkiosk", controllers.LinkKioskToAccount).Methods("POST")
	router.HandleFunc("/confirmation/{code}", controllers.ValidateConfirmCode).Methods("GET")

	router.Handle("/refresh-token-auth",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(controllers.RefreshToken),
		)).Methods("GET")
	router.Handle("/logout",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(controllers.Logout),
		)).Methods("GET")
	return router
}
