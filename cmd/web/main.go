// Members: Miguel Avila, Federico Rosado

package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq" //Third party package
	"mavila_frosado.net/test1/pkg/models/postgresql"
)

// Database Function
func setUpDB(dsn string) (*sql.DB, error) {

	//Establish connection to the dababase
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	//test connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

//Dependencies Injection (passing)
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	Blogs    *postgresql.BlogModel
	addr     string
}

func main() {
	// create a command line flag for the port
	addr := flag.String("addr", ":8080", "HTTP network address")
	dsn := flag.String("dsn", os.Getenv("DBTEST_DB_DSN"), "PostgreSQL DSN (Data Source Name)")

	flag.Parse()

	//Create a logger - for loggin anything that is not an error
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	//Create a logger - for loggin  error
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	//create and initialize
	var db, err = setUpDB(*dsn)

	//Check if errors
	if err != nil {
		log.Fatal(err)
	}

	//A must before exiting
	defer db.Close()
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		Blogs:    &postgresql.BlogModel{DB: db},
		addr:     *addr,
	}

	//Create a custom server
	srv := &http.Server{
		Addr:     *addr,
		Handler:  app.routes(),
		ErrorLog: errorLog,
	}

	//Instead of the default logger
	infoLog.Printf("Starting Server on port %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)

}
