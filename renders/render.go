package renders

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Render interface {
	Render(ctx *gin.Context)
}

type ErrorRender interface {
	Render
	WithErr(err error) ErrorRender
	Err() error
}

var jsonContentType = []string{"application/json; charset=utf-8"}

type JSON struct {
	Data interface{}
}

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}

func (r JSON) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	jsonBytes, err := json.MarshalIndent(r.Data, "", "    ")
	if err != nil {
		return err
	}
	_, err = w.Write(jsonBytes)
	return err
}

// WriteContentType (IndentedJSON) writes JSON ContentType.
func (r JSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}
