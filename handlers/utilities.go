package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func parseRequest(req *http.Request) ([]byte) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	return buf.Bytes()
}

func getUserInfo(email string) map[string]interface{} {
	client := &http.Client{}
	r, _ := http.NewRequest("GET", "http://localhost:3000/user/search?create=true&email="+email, nil)
	resp, err := client.Do(r)
	if err != nil {
		fmt.Printf("Error : %s", err)
	}
	clientBuf := new(bytes.Buffer)
	clientBuf.ReadFrom(resp.Body)
	var clientData []map[string]interface{}
	json.Unmarshal(clientBuf.Bytes(), &clientData)
	return clientData[0]
}

func sendRequest(r *http.Request) (response []byte, statusCode int) {
	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		fmt.Printf("Error : %s", err)
	}
	respBuf := new(bytes.Buffer)
	respBuf.ReadFrom(resp.Body)

	return respBuf.Bytes(), resp.StatusCode
}