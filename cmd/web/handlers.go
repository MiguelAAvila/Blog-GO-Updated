// Members: Miguel Avila, Federico Rosado

package main

import (
	// "fmt"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"mavila_frosado.net/test1/pkg/models"
)

//Home Page
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	//  panic("This is not good! Oops")

	ts, err := template.ParseFiles("./ui/html/index.tmpl")
	port := app.addr

	if err != nil {
		app.serverError(w, err)
	}
	data := &templateData{
		Port: port,
	}
	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
	}
}

//Displays SingUp Form
func (app *application) createBlogForm(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./ui/html/form.tmpl")
	if err != nil {
		app.serverError(w, err)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		app.serverError(w, err)
	}

}

//blog Page
func (app *application) blogs(w http.ResponseWriter, r *http.Request) {
	blogs, err := app.Blogs.Read()
	port := app.addr

	if err != nil {
		app.serverError(w, err)
	}
	//instance of templateData
	data := &templateData{
		Blogs: blogs,
		Port:  port,
	}

	//Body part of tmpl
	ts, err := template.ParseFiles("./ui/html/blog.tmpl")
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.Execute(w, data)

	if err != nil {
		app.serverError(w, err)
		return
	}

}

//Extract, Validate and Write to the blogs table
func (app *application) createBlog(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	firstname := r.PostForm.Get("firstname")
	lastname := r.PostForm.Get("lastname")
	email := r.PostForm.Get("email")
	subject := r.PostForm.Get("subject")
	message := r.PostForm.Get("message")

	//Validate Form Fields
	errors := make(map[string]string)
	//Check each Fields
	if strings.TrimSpace(firstname) == "" {
		errors["firstname"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(firstname) > 20 {
		errors["firstname"] = "No more than 20 characters"
	}
	//Check each Fields
	if strings.TrimSpace(lastname) == "" {
		errors["lastname"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(lastname) > 20 {
		errors["lastname"] = "No more than 20 characters"
	}
	if strings.TrimSpace(email) == "" {
		errors["email"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(email) > 25 {
		errors["email"] = "No more than 25 characters"
	}
	if strings.TrimSpace(subject) == "" {
		errors["subject"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(subject) > 50 {
		errors["subject"] = "No more than 50 characters"
	}
	if strings.TrimSpace(message) == "" {
		errors["message"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(subject) > 500 {
		errors["message"] = "No more than 500 characters"
	}
	//Check if errors in the map
	if len(errors) > 0 {
		ts, err := template.ParseFiles("./ui/html/form.tmpl")
		err = ts.Execute(w, &templateData{
			ErrorsFromForm: errors,
			FormData:       r.PostForm,
		})
		if err != nil {
			app.serverError(w, err)
			return
		}

		if err != nil {
			app.serverError(w, err)
			return
		}
		return
	}

	//inser a blog
	id, err := app.Blogs.Insert(firstname, lastname, email, subject, message)
	if err != nil {
		app.serverError(w, err)
		return
	}
	// fmt.Fprintf(w, "Row with id %d has been inserted.", id)
	http.Redirect(w, r, fmt.Sprintf("/show-blog?id=%d", id), http.StatusSeeOther)
}

func (app *application) showBlog(w http.ResponseWriter, r *http.Request) {
	port := app.addr
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	// Creating an array of string type
	// Using var keyword
	blogs := []*models.Blog{}

	blog, err := app.Blogs.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	blogs = append(blogs, blog)

	//instance of templateData
	data := &templateData{
		Blogs: blogs,
		Port:  port,
	}

	//Body part of tmpl
	ts, err := template.ParseFiles("./ui/html/blog.tmpl")
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.Execute(w, data)

	if err != nil {
		app.serverError(w, err)
		return
	}

}
