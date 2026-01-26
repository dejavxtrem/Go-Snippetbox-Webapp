package main

import (
	"errors"
	"fmt" // New import
	"net/http"
	"strconv"
	"strings" // New import
	"unicode/utf8"

	"github.com/dejavxtrem/snippetbox/internal/models"
)

// Define a snippetCreateForm struct to represent the form data and validation
// errors for the form fields. Note that all the struct fields are deliberately
// exported (i.e. start with a capital letter). This is because struct fields
// must be exported in order to be read by the html/template package when
// rendering the template.
type snippetCreateForm struct {
	Title       string
	Content     string
	Expires     int
	FieldErrors map[string]string
}

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
	//w.Header().Add("Server", "Go")

	snippets, err := app.snippet.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// for _, snippet := range snippets {
	// 	fmt.Fprintf(w, "%v+\n", snippet)
	// }

	// Call the newTemplateData() helper to get a templateData struct containing
	// the 'default' data (which for now is just the current year), and add the
	// snippets slice to it.
	data := app.newTemplateData(r)
	data.Snippets = snippets

	// Pass the data to the render() helper as normal.
	app.render(w, r, http.StatusOK, "home.html", data)

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

	// And do the same thing again here...
	data := app.newTemplateData(r)
	data.Snippet = snippet

	// Use the new render helper.
	app.render(w, r, http.StatusOK, "view.html", data)

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

	data := app.newTemplateData(r)

	// Initialize a new snippetCreateForm instance and pass it to the template.
	// Notice how this is also a great opportunity to set any default or
	// 'initial' values for the form --- here we set the initial value for the
	// snippet expiry to 365 days.

	data.Form = snippetCreateForm{
		Expires: 365,
	}

	//data := app.newTemplateData(r)

	app.render(w, r, http.StatusOK, "create.html", data)

	//w.Write([]byte("Display a form for creating a new snippet..."))
}

// POST
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	//use the method to customize() Method to send a 201 status code

	// title := "O snail"
	// content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	// expires := 7
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Use the r.PostForm.Get() method to retrieve the title and content
	// from the r.PostForm map.
	//title := r.PostForm.Get("title")
	//content := r.PostForm.Get("content")

	// The r.PostForm.Get() method always returns the form data as a *string*.
	// However, we're expecting our expires value to be a number, and want to
	// represent it in our Go code as an integer. So we need to manually convert
	// the form data to an integer using strconv.Atoi(), and send a 400 Bad
	// Request response if the conversion fails.
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Create an instance of the snippetCreateForm struct containing the values
	// from the form and an empty map for any validation errors.
	// form := snippetCreateForm{
	// 	Title:       r.PostForm.Get("title"),
	// 	Content:     r.PostForm.Get("content"),
	// 	Expires:     expires,
	// 	FieldErrors: map[string]string{},
	// }

	form := snippetCreateForm{
		Title:       r.PostForm.Get("title"),
		Content:     r.PostForm.Get("content"),
		Expires:     expires,
		FieldErrors: map[string]string{},
	}

	// Initialize a map to hold any validation errors for the form fields.
	//fieldErrors := make(map[string]string)

	// Check that the title value is not blank and is not more than 100
	// characters long. If it fails either of those checks, add a message to the
	// errors map using the field name as the key.
	if strings.TrimSpace(form.Title) == "" {
		form.FieldErrors["title"] = "The field cannot be blank"

	} else if utf8.RuneCountInString(form.Title) > 100 {
		form.FieldErrors["title"] = "This field cannot be more than 100 characters long"
	}

	// Check that the content value isn't blank.
	if strings.TrimSpace(form.Content) == "" {
		form.FieldErrors["content"] = "This field cannot be blank"
	}

	// Check the expires value matches one of the permitted values (1, 7 or
	// 365).
	if form.Expires != 1 && form.Expires != 7 && form.Expires != 365 {
		form.FieldErrors["expires"] = "This field must equal 1, 7 or 365"
	}

	// If there are any errors, dump them in a plain-text HTTP response and
	// return from the handler.

	// If there are any validation errors, then  the create.tmpl template,
	// passing in the snippetCreateForm instance as dynamic data in the Form
	// field. Note that we use the HTTP status code 422 Unprocessable Entity
	// when sending the response to indicate that there was a validation error.

	//fmt.Printf("this is the form error title %+v\n", form.FieldErrors)

	if len(form.FieldErrors) > 0 {
		data := app.newTemplateData(r)
		data.Form = form

		app.render(w, r, http.StatusUnprocessableEntity, "create.html", data)
		//fmt.Fprint(w, FieldErrors)
		return
	}

	// Pass the data to the SnippetModel.Insert() method, receiving the
	// ID of the new record back.

	// We also need to update this line to pass the data from the
	// snippetCreateForm instance to our Insert() method.
	id, err := app.snippet.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)

	//w.WriteHeader(http.StatusCreated)
	//w.Write([]byte("Create a new snippet"))
}
