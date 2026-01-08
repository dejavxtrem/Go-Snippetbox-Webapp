package main

import (
	"errors"
	"fmt" // New import
	"net/http"
	"strconv"

	"github.com/dejavxtrem/snippetbox/internal/models"
)

// Define an application struct to hold the application-wide dependencies for the
// web application. For now we'll only include the structured logger, but we'll
// add more to this as development progresses.

// Change the signature of the home handler so it is defined as a method against
// *application.
// GET
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Use the Header().Add() method to add a 'Server: Go' header to the
	// response header map. The first parameter is the header name, and
	// the second parameter is the header value.
	w.Header().Add("Server", "Go")

	snippets, err := app.snippet.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// for _, snippet := range snippets {
	// 	fmt.Fprintf(w, "%v+\n", snippet)
	// }

	// Use the new render helper.
	app.render(w, r, http.StatusOK, "home.html", templateData{Snippets: snippets})

	// // // Initialize a slice containing the paths to the two files. It's important
	// // // to note that the file containing our base template must be the *first*
	// // // file in the slice.
	// files := []string{
	// 	"./ui/html/base.html",
	// 	"./ui/html/partials/nav.html",
	// 	"./ui/html/pages/home.html",
	// }

	// // // Use the template.ParseFiles() function to read the template file into a
	// // // template set. If there's an error, we log the detailed error message, use
	// // // the http.Error() function to send an Internal Server Error response to the
	// // // user, and then return from the handler so no subsequent code is executed.

	// // // Use the template.ParseFiles() function to read the files and store the
	// // // templates in a template set. Notice that we use ... to pass the contents
	// // // of the files slice as variadic arguments.
	// ts, err := template.ParseFiles(files...)

	// if err != nil {
	// 	// Because the home handler is now a method against the application
	// 	// struct it can access its fields, including the structured logger. We'll
	// 	// use this to create a log entry at Error level containing the error
	// 	// message, also including the request method and URI as attributes to
	// 	// assist with debugging.
	// 	//app.logger.Error(err.Error(), "Method", r.Method, "uri", r.URL.RequestURI())
	// 	//log.Print(err.Error())
	// 	//http.Error(w, "Internal Server Error", http.StatusInternalServerError)

	// 	app.serverError(w, r, err) // Use the serverError() helper
	// 	return
	// }

	// // Create an instance of a templateData struct holding the slice of
	// // snippets.
	// data := templateData{
	// 	Snippets: snippets,
	// }

	// // Then we use the Execute() method on the template set to write the
	// // template content as the response body. The last parameter to Execute()
	// // represents any dynamic data that we want to pass in, which for now we'll
	// // leave as nil.

	// // Use the ExecuteTemplate() method to write the content of the "base"
	// // template as the response body.

	// // Pass in the templateData struct when executing the template.
	// err = ts.ExecuteTemplate(w, "base", data)

	// if err != nil {

	// 	// app.logger.Error(err.Error(), "Method", r.Method, "uri", r.URL.RequestURI())
	// 	// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// 	app.serverError(w, r, err) // Use the serverError() helper
	// }
	// // w.Write([]byte("Hello World"))
}

// Add a snippetView handler function.
// GET
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {

	// Extract the value of the id wildcard from the request using r.PathValue()
	// and try to convert it to an integer using the strconv.Atoi() function. If
	// it can't be converted to an integer, or the value is less than 1, we
	// return a 404 page not found response.

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	// Use the SnippetModel's Get() method to retrieve the data for a
	// specific record based on its ID. If no matching record is found,
	// return a 404 Not Found response.
	snippet, err := app.snippet.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return

	}

	// Use the new render helper.
	app.render(w, r, http.StatusOK, "view.html", templateData{
		Snippet: snippet,
	})

	// files := []string{
	// 	"./ui/html/base.html",
	// 	"./ui/html/partials/nav.html",
	// 	"./ui/html/pages/view.html",
	// }

	// // Parse the template files...
	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	app.serverError(w, r, err)
	// 	return
	// }

	// // Create an instance of a templateData struct holding the snippet data.
	// data := templateData{
	// 	Snippet: snippet,
	// }

	// // And then execute them. Notice how we are passing in the snippet
	// // data (a models.Snippet struct) as the final parameter?
	// err = ts.ExecuteTemplate(w, "base", data)
	// if err != nil {
	// 	app.serverError(w, r, err)
	// }

	// // Use the fmt.Sprintf() function to interpolate the id value with a
	// // message, then write it as the HTTP response.
	// //msg := fmt.Sprintf("Display a specific snippet with ID %d..", id)
	// fmt.Fprintf(w, "%+v", snippet)
	// //w.Write([]byte(msg))
}

// GET
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet..."))
}

// POST
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	//use the method to customize() Method to send a 201 status code

	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	// Pass the data to the SnippetModel.Insert() method, receiving the
	// ID of the new record back.
	id, err := app.snippet.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)

	//w.WriteHeader(http.StatusCreated)
	//w.Write([]byte("Create a new snippet"))
}
