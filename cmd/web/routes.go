// Members: Miguel Avila, Federico Rosado

package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/createblog", app.createBlogForm)
	mux.HandleFunc("/blog-add", app.createBlog)
	mux.HandleFunc("/blogs", app.blogs)
	mux.HandleFunc("/show-blog", app.showBlog)
	//create a fileserver to serve static files
	fs := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	return app.recoverPanicMiddleware(app.logRequestMiddleware(securityHeadersMiddleware(mux)))
}
