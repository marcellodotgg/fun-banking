package handler

type FormData struct {
	Data   map[string]string
	Errors map[string]string
}

func newFormData() FormData {
	return FormData{
		Data:   make(map[string]string),
		Errors: make(map[string]string),
	}
}
