package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/chent03/apt-server/controllers"
	"github.com/chent03/apt-server/models"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world")
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Print(".env file missing")
	}
	postgresInfo := getDBConnectionInfo()
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		postgresInfo.host, postgresInfo.port, postgresInfo.user, postgresInfo.password, postgresInfo.name)

	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer us.Close()
	us.DestructiveReset()
	us.AutoMigrate()

	fmt.Println("connected!!")
	userC := controllers.NewUsers(us)

	r := mux.NewRouter()
	r.HandleFunc("/", helloHandler)
	r.HandleFunc("/api/register", userC.Register).Methods("POST")
	r.HandleFunc("/api/login", userC.Login).Methods("POST")
	handler := cors.Default().Handler(r)
	http.ListenAndServe(":"+GetPortNumber(), handler)
}
