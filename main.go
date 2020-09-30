package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/chent03/apt-server/controllers"
	"github.com/chent03/apt-server/middleware"
	"github.com/chent03/apt-server/models"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world updated")
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Print(".env file missing")
	}
	postgresInfo := getDBConnectionInfo()
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		postgresInfo.host, postgresInfo.port, postgresInfo.user, postgresInfo.password, postgresInfo.name)

	services, err := models.NewServices(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer services.Close()
	services.DestructiveReset()
	userC := controllers.NewUsers(services.User)
	reviewC := controllers.NewReviews(services.Review)

	requireUserMw := middleware.RequireUser{
		UserService: services.User,
	}

	getUserInfo := requireUserMw.ApplFn(userC.GetUserInfo)
	createReview := requireUserMw.ApplFn(reviewC.Create)

	r := mux.NewRouter()
	r.HandleFunc("/", helloHandler)
	r.HandleFunc("/api/register", userC.Register).Methods("POST")
	r.HandleFunc("/api/login", userC.Login).Methods("POST")
	r.HandleFunc("/api/getUserInfo", getUserInfo).Methods("GET")
	r.HandleFunc("/api/review", createReview).Methods("POST")

	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "https://renters-review.herokuapp.com"},
		AllowedMethods:   []string{"POST", "GET", "OPTIONS", "PUT"},
		AllowedHeaders:   []string{"Accept", "Accept-Language", "Content-Type"},
		AllowCredentials: true,
	}).Handler(r)

	http.ListenAndServe(":"+GetPortNumber(), handler)
}
