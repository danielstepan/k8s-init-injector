package handler

import (
	"net/http"
)

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("HandleRoot!"))
}
