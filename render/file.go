package render

import "net/http"

type File struct {
	Path    string
	Request *http.Request
}

func (r File) Write(w http.ResponseWriter) error {
	http.ServeFile(w, r.Request, r.Path)
	return nil
}
