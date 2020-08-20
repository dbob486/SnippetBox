package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

//application struct holding all dependencies for accessibility across the application currently only our custom loggers
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {

	//Creating default configuration for the host port that will serve our application
	addr := flag.String("addr", ":8080", "HTTP network address")
	flag.Parse()

	//creating custom info and error loggers to easily access where errors and information are occurring
	//Displays the date and time into the terminal in addition to the INFORMATION the client logs
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate | log.Ltime)

	//Displays the date, time and filename:which line the error occured in into the terminal in addition to the ERROR message
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate | log.Ltime | log.Lshortfile)

	//instance of our application struct to pass in our custom loggers
	app := &application{
		errorLog: errLog,
		infoLog:  infoLog,
	}

	//to keep main from being convoluted we made a seperate file to manage routes that will return the our custom server mux
	mux := app.routes()
	
	srv := &http.Server {
		Addr:     *addr,
		ErrorLog: errLog,
		Handler:  mux,
	}
	infoLog.Printf("Starting server on %s", *addr)

	//serving file using custom server struct with
	err := srv.ListenAndServe()
	errLog.Fatal(err)
}