package main

import (
	"flag"
	"fmt"
	"net/http"

	"sw.com/FirstWebWithGO/middleware"

	"github.com/gorilla/mux"
	"sw.com/FirstWebWithGO/controllers"
	"sw.com/FirstWebWithGO/models"
)

func main() {
	boolPtr := flag.Bool("prod", false, "Provide this flag "+
		"in production. This ensures that a .config file is "+
		"provided before the application starts.")
	flag.Parse()

	cfg := LoadConfig(*boolPtr)
	dbCfg := cfg.Database
	services, err := models.NewServices(models.WithGorm(dbCfg.Dialect(), dbCfg.ConnectionInfo()),
		models.WithLogMode(!cfg.IsProd()),
		models.WithUser(cfg.Pepper, cfg.HMACKey),
		models.WithGallery(),
		models.WithImage(),
	)
	if err != nil {
		panic(err)
	}
	defer services.Close()
	services.DestructiveReset()

	r := mux.NewRouter()
	usersC := controllers.NewUsers(services.User)
	staticC := controllers.NewStatic()
	galleriesC := controllers.NewGalleries(services.Gallery, services.Image, r)

	userMw := middleware.User{
		UserService: services.User,
	}
	requiredUserMw := middleware.RequireUser{}

	newGallery := requiredUserMw.Apply(galleriesC.New)
	createGallery := requiredUserMw.ApplyFn(galleriesC.Create)

	imageHandler := http.FileServer(http.Dir("./images"))
	r.PathPrefix("/images").Handler(http.StripPrefix("/images", imageHandler))
	r.HandleFunc("/galleries/{id:[0-9]+}/images/{filename}/delete",
		requiredUserMw.ApplyFn(galleriesC.ImageDelete)).Methods("POST")

	assetHandler := http.FileServer(http.Dir("./assets"))
	assetHandler = http.StripPrefix("/assets/", assetHandler)
	r.PathPrefix("/assets").Handler(assetHandler)

	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/logout", requiredUserMw.ApplyFn(usersC.Logout)).Methods("POST")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.Handle("/faq", staticC.Faq).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	r.Handle("/galleries/new", newGallery).Methods("GET")
	r.HandleFunc("/galleries", createGallery).Methods("POST")
	r.Handle("/galleries", requiredUserMw.ApplyFn(galleriesC.Index)).Methods("GET").Name(controllers.IndexGalleries)
	r.HandleFunc("/galleries/{id:[0-9]+}", galleriesC.Show).Methods("GET").Name(controllers.ShowGallery)
	r.HandleFunc("/galleries/{id:[0-9]+}/edit", requiredUserMw.ApplyFn(galleriesC.Edit)).Methods("GET").Name(controllers.EditGallery)
	r.HandleFunc("/galleries/{id:[0-9]+}/update", requiredUserMw.ApplyFn(galleriesC.Update)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/delete", requiredUserMw.ApplyFn(galleriesC.Delete)).Methods("POST")
	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")
	r.HandleFunc("/galleries/{id:[0-9]+}/images", requiredUserMw.ApplyFn(galleriesC.ImageUpload)).Methods("POST")
	fmt.Printf("Starting the server on :%d...\n", cfg.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), userMw.Apply(r))
}
