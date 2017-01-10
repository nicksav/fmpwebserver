package routers

import (
	"fmpwebserver/controllers"
	"fmpwebserver/core/authentication"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

func SetKioskRoutes(router *mux.Router) *mux.Router {
	router.Handle("/kiosks/{id}",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(controllers.GetKiosk),
		)).Methods("GET")

	router.Handle("/kiosks",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(controllers.GetKiosks),
		)).Methods("GET")

    router.Handle("/kiosks",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(controllers.UpdateKiosk),
		)).Methods("PUT")

    router.Handle("/kiosks",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(controllers.RemoveKiosk),
		)).Methods("DELETE")

	return router
}
