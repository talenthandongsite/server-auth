// response package is
package response

import (
	"encoding/json"
	"net/http"
)

type StandardOutput struct {
	Status bool            `json:"status"`
	Data   json.RawMessage `json:"data"`
}

func JSONResponse(w http.ResponseWriter, body []byte, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	out := StandardOutput{
		Status: true,
		Data:   body,
	}
	b, err := json.Marshal(out)
	if err != nil {
		w.Write(body)
	}
	w.Write(b)
}

func JSONError(w http.ResponseWriter, err error, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	e, _ := json.Marshal(err.Error())

	errOut := StandardOutput{
		Status: false,
		Data:   e,
	}

	b, err := json.Marshal(errOut)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	w.Write(b)
}
