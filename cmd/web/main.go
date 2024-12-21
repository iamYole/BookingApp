package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/iamYole/BookingApp/internal/config"
	"github.com/iamYole/BookingApp/internal/handlers"
	"github.com/iamYole/BookingApp/internal/helpers"
	"github.com/iamYole/BookingApp/internal/models"
	"github.com/iamYole/BookingApp/internal/render"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {

	err := run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Starting application on %s ...", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {

	//change this to true when in production
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	gob.Register(models.Reservation{})

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Error creating template cache")
		return err
	}

	app.TemplateCache = tc
	app.UseCache = false

	render.NewTemplate(&app)

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	helpers.NewHelpers(&app)

	return nil
}
