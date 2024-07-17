package param

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func Int(r *http.Request, param string) (int, error) {
	val, err := strconv.Atoi(chi.URLParam(r, param))
	if err != nil {
		return 0, err
	}

	return val, nil
}
