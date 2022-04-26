// Members: Miguel Avila, Federico Rosado

package main

import (
	"net/http"

	"github.com/bmizerany/pat"
)

func (app *application) routes() http.Handler {
	//Thrird party router/multiplexer
	mux := pat.New()
	mux.Get("/", http.HandlerFunc(app.home))
	mux.Get("/blog/create", http.HandlerFunc(app.createBlogForm))
	mux.Post("/blog/create", http.HandlerFunc(app.createBlog))
	mux.Get("/blog/:id", http.HandlerFunc(app.showBlog))
	// mux.HandleFunc("/blogs", app.blogs)
	//create a fileserver to serve static files
	fs := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static/", fs))

	return app.recoverPanicMiddleware(app.logRequestMiddleware(securityHeadersMiddleware(mux)))
}
