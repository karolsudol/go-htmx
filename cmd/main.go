package main

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Template struct {
	tmpl *template.Template
}

func newTemplate() *Template {
	return &Template{
		tmpl: template.Must(template.ParseGlob("views/*.html")),
	}
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.tmpl.ExecuteTemplate(w, name, data)
}

type Contact struct {
	Name  string
	Email string
	ID    int
}

type Data struct {
	Contacts []Contact
}

func NewData() *Data {
	return &Data{
		Contacts: []Contact{
			{
				Name:  "John Doe",
				Email: "john.doe@gmail.com",
				ID:    1,
			},
			{
				Name:  "Jane Doe",
				Email: "jain.doe@gmail.com",
				ID:    2,
			},
		},
	}
}

type FormData struct {
	Errors map[string]string
	Values map[string]string
}

func NewFormData() FormData {
	return FormData{
		Errors: map[string]string{},
		Values: map[string]string{},
	}
}

type PageData struct {
	Data Data
	Form FormData
}

func NewContact(id int, name, email string) Contact {
	return Contact{
		ID:    id,
		Name:  name,
		Email: email,
	}
}

func NewPageData(data Data, form FormData) PageData {
	return PageData{
		Data: data,
		Form: form,
	}
}

func contactExists(contacts []Contact, email string) bool {
	for _, c := range contacts {
		if c.Email == email {
			return true
		}
	}
	return false
}

func main() {

	e := echo.New()

	data := NewData()
	id := 3

	e.Renderer = newTemplate()
	e.Use(middleware.Logger())

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index.html", NewPageData(*data, NewFormData()))
	})

	e.POST("/contacts", func(c echo.Context) error {
		name := c.FormValue("name")
		email := c.FormValue("email")

		if contactExists(data.Contacts, email) {
			formData := FormData{
				Errors: map[string]string{
					"email": "Email already exists",
				},
				Values: map[string]string{
					"name":  name,
					"email": email,
				},
			}

			return c.Render(422, "contact-form", formData)
		}

		contact := NewContact(id, name, email)
		id++

		data.Contacts = append(data.Contacts, contact)

		formData := NewFormData()
		err := c.Render(200, "contact-form", formData)
		if err != nil {
			return err
		}

		return c.Render(200, "oob-contact", contact)
	})

	e.Logger.Fatal(e.Start(":42069"))
}
