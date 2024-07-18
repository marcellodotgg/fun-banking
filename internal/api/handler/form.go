package handler

import (
	"strings"

	"github.com/gin-gonic/gin"
)

type FormData struct {
	Data   map[string]string
	Errors map[string]string
}

func NewFormData() FormData {
	return FormData{
		Data:   make(map[string]string),
		Errors: make(map[string]string),
	}
}

// Takes in a context which should have a form in the request,
// then parses the form, and inserts all the form data into
// form.Data and returns the FormData.
func GetForm(c *gin.Context) FormData {
	c.Request.ParseForm()
	form := NewFormData()

	for key, values := range c.Request.PostForm {
		if len(values) > 0 {
			form.Data[key] = strings.TrimSpace(values[0])
		}
	}

	return form
}
