package main

import (
	"bytes"
	"fmt"
	"github.com/justinas/nosurf"
	"net/http"
	"runtime/debug"
	"time"
)

// The ServerError helper writes an error message and stack trace to the errorLog,
// then sends a generic 500 Internal Server Error response to the user.
// our server will log a stack trace to the server error log, unwind the stack for the affected goroutine
//(calling any deferred functions along the way) and close the underlying HTTP connection.
//But it won’t terminate the application, so importantly, any panic in your handlers won’t bring down your server.
func (app *application) ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// The ClientError helper sends a specific status code and corresponding description
// to the user. We'll use this later in the book to send responses like 400 "Bad
// Request" when there's a problem with the request that the user sent.
func (app *application) ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
	fmt.Fprintln(w, status)

}

// For consistency, we'll also implement a NotFound helper. This is simply a
// convenience wrapper around ClientError which sends a 404 Not Found response to
// the user.
func (app *application) NotFound(w http.ResponseWriter) {
	app.ClientError(w, http.StatusNotFound)
}

func (app *application) Render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := app.templateCache[name]
	if !ok {
		app.ServerError(w, fmt.Errorf("The template %s does not exist", name))
		return
	}

	buf := new(bytes.Buffer)

	// Execute the template set, passing the dynamic data with the current
	// year injected.
	err := ts.Execute(buf, app.AddDefaultData(td, r))
	if err != nil {
		app.ServerError(w, err)
		return
	}

	buf.WriteTo(w)
}

// Create an AddDefaultData helper. This takes a pointer to a templateData
// struct, adds the current year to the CurrentYear field, and then returns
// the pointer. Again, we're not using the *http.Request parameter at the
// moment, but we will do later in the book.
func (app *application) AddDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}

	td.CSRFToken = nosurf.Token(r)
	td.CurrentYear = time.Now().Year()
	td.Flash = app.session.PopString(r, "flash")

	// Add the authentication status to the template data.
	td.IsAuthenticated = app.isAuthenticated(r)
	return td
}



func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(contextKeyIsAuthenticated).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}
