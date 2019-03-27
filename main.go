package main

import (
	"fmt"
	"net/http"

	"sw.com/FirstWebWithGO/middleware"

	"github.com/gorilla/mux"
	"sw.com/FirstWebWithGO/controllers"
	"sw.com/FirstWebWithGO/models"
)

const (
	host     = "192.168.1.8"
	port     = 5432
	user     = "postgres"
	password = "310104shenwei"
	dbname   = "website"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	services, err := models.NewServices(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer services.Close()
	services.DestructiveReset()

	usersC := controllers.NewUsers(services.User)
	staticC := controllers.NewStatic()
	galleriesC := controllers.NewGalleries(services.Gallery)

	requiredUserMw := middleware.RequireUser{
		UserService: services.User,
	}
	newGallery := requiredUserMw.Apply(galleriesC.New)
	createGallery := requiredUserMw.ApplyFn(galleriesC.Create)

	r := mux.NewRouter()
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.Handle("/faq", staticC.Faq).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	r.Handle("/galleries/new", newGallery).Methods("GET")
	r.HandleFunc("/galleries", createGallery).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}", galleriesC.Show).Methods("GET")
	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")
	http.ListenAndServe(":3000", r)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
