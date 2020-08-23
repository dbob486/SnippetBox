package main

import (
	"danielgarcia.net/snippetbox/pkg/models/mysql"
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

//application struct holding all dependencies for accessibility across the application currently only our custom loggers
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *mysql.SnippetModel
}

func main() {

	//Creating command line flag with default configuration for the host port that will serve our application
	addr := flag.String("addr", ":8080", "HTTP network address")

	//Defining a new command line flag for the MySQL DSN(Data Source Name) string.
	dsn := flag.String("dns", "web:another@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	//creating custom info and error loggers to easily access where errors and information are occurring
	//Displays the date and time into the terminal in addition to the INFORMATION the client logs
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate | log.Ltime)

	//Displays the date, time and filename:which line the error occured in into the terminal in addition to the ERROR message
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate | log.Ltime | log.Lshortfile)

	//pass dsn into our openDB(dns string) function that will create a db connection pool
	db, err := openDB(*dsn)
	if err != nil {
		errLog.Fatal(err)
	}
	defer db.Close()

	//instance of our application struct to pass in our custom loggers
	app := &application{
		errorLog: errLog,
		infoLog:  infoLog,
		snippets: &mysql.SnippetModel{DB: db} ,
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
	err = srv.ListenAndServe()
	errLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}