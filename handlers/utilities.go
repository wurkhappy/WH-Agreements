package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func parseRequest(req *http.Request) (map[string]interface{}, []byte) {
	var m map[string]interface{}
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	reqBytes := buf.Bytes()
	json.Unmarshal(reqBytes, &m)
	return m, reqBytes
}
