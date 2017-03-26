package render

import (
	"io"
	"io/ioutil"
	"net/http"

	"gopkg.in/yaml.v2"
)

var statusCtxKey = &contextKey{"Status"}

func YAML(w http.ResponseWriter, r *http.Request, v interface{}) {
	b, err := yaml.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/x-yaml; charset=utf-8")
	if status, ok := r.Context().Value(statusCtxKey).(int); ok {
		w.WriteHeader(status)
	}
	w.Write(b)
}

func YAMLBind(r io.Reader, v interface{}) error {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(body, &v)
}

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation. This technique
// for defining context keys was copied from Go 1.7's new use of context in net/http.
type contextKey struct {
	name string
}
