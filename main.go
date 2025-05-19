package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Teneieiza/Golang-API/middleware"
	"github.com/Teneieiza/Golang-API/config"
	"github.com/Teneieiza/Golang-API/routes"
	// "github.com/Teneieiza/Golang-API/00-marshal"
	// webservice "github.com/Teneieiza/Golang-API/01-webservice"
	// urlpath "github.com/Teneieiza/Golang-API/02-urlpath"
	// middleware "github.com/Teneieiza/Golang-API/03-middleware"
	// corsorigin "github.com/Teneieiza/Golang-API/04-corsorigin"
	// connectDB "github.com/Teneieiza/Golang-API/05-connectDB"
	// gowithmysql "github.com/Teneieiza/Golang-API/06-gowithmysql"
	// apiwithdb "github.com/Teneieiza/Golang-API/07-APIwithDB"
)

func main() {
	fmt.Println("Hello, World! with GoLang eieiza")

	// marshal.Marshal()
	// marshal.UnMarshal()
	// webservice.WorkRequest()
	// urlpath.UrlPath()
	// middleware.Middleware()
	// corsorigin.Corsorigin()
	// connectDB.ConnectDB()
	// gowithmysql.Gowithmysql()
	// apiwithdb.ApiWithDB()
	config.SetupDB()
	routes.SetupRoutes()

	log.Println("Server listening on http://localhost:5000")
	http.ListenAndServe(":5000", middleware.CorsMiddleware(http.DefaultServeMux))

}
