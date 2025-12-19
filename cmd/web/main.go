package main

import (
	"database/sql" // New import
	"flag"
	"log/slog"
	"net/http"
	"os"
)

type application struct {
	logger *slog.Logger
}

func main() {

	// Define a new command-line flag with the name 'addr', a default value of ":4000"
	// and some short help text explaining what the flag controls. The value of the
	// flag will be stored in the addr variable at runtime.
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:password@/snippetbox?parseTime=true", "MySQL data source name")

	// Importantly, we use the flag.Parse() function to parse the command-line flag.
	// This reads in the command-line flag value and assigns it to the addr
	// variable. You need to call this *before* you use the addr variable
	// otherwise it will always contain the default value of ":4000". If any errors are
	// encountered during parsing the application will be terminated.
	flag.Parse()

	// Use the slog.New() function to initialize a new structured logger, which
	// writes to the standard out stream and uses the default settings.
	//logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// To keep the main() function tidy I've put the code for creating a connection
	// pool into the separate openDB() function below. We pass openDB() the DSN
	// from the command-line flag.
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// We also defer a call to db.Close(), so that the connection pool is closed
	// before the main() function exits.
	defer db.Close()

	// Initialize a new instance of our application struct, containing the
	// dependencies (for now, just the structured logger).
	app := &application{
		logger: logger,
	}

	// Register the two new handler functions and corresponding route patterns with
	// the servemux, in exactly the same way that we did before
	// mux := http.NewServeMux()

	//Sever Static Files
	// Create a file server which serves files out of the "./ui/static" directory.
	// Note that the path given to the http.Dir function is relative to the project
	// directory root.
	// fileServer := http.FileServer(http.Dir("./ui/static/"))

	// mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// mux.HandleFunc("GET /{$}", app.home) // Restrict this route to exact matches on / only.
	// //mux.HandleFunc("/snippet/view", snippetView)

	// mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)

	// mux.HandleFunc("GET /snippet/create", app.snippetCreate)

	// //The Post Request Handler
	// mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

	//log.Printf("Starting server on port on %s", *addr)
	// Use the Info() method to log the starting server message at Info severity
	// (along with the listen address as an attribute).
	logger.Info("Starting server", slog.Any("addr", *addr))

	// Call the new app.routes() method to get the servemux containing our routes,
	// and pass that to http.ListenAndServe().

	// Because the err variable is now already declared in the code above, we need
	// to use the assignment operator = here, instead of the := 'declare and assign'
	// operator.
	err = http.ListenAndServe(*addr, app.routes())

	//log.Fatal(err)
	logger.Error(err.Error())
	os.Exit(1)

}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
