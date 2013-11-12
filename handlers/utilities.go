package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

var UserService string = "http://localhost:3000"
var PaymentInfoService string = "http://localhost:3120"
var AgreementsService string = "http://localhost:4050"
var CommentsService string = "http://localhost:5050"

func parseRequest(req *http.Request) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	return buf.Bytes()
}

func getUserInfo(email string) map[string]interface{} {
	client := &http.Client{}
	r, _ := http.NewRequest("GET", UserService+"/user/search?create=true&email="+email, nil)
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

func sendRequest(method, service, path string, body []byte) (response []byte, statusCode int) {
	bodyReader := bytes.NewReader(body)
	r, _ := http.NewRequest(method, service+path, bodyReader)
	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		fmt.Printf("Error : %s", err)
	}
	respBuf := new(bytes.Buffer)
	respBuf.ReadFrom(resp.Body)

	return respBuf.Bytes(), resp.StatusCode
}

func sendServiceRequest(method, service, path string, body []byte) (response []byte, statusCode int) {
	client := mdp.NewClient("tcp://localhost:5555", false)
	defer client.Close()
	m := map[string]interface{}{
		"Method": method,
		"Path":   path,
		"Body":   body,
	}
	req, _ := json.Marshal(m)
	request := [][]byte{req}
	reply := client.Send([]byte(service), request)
	if len(reply) == 0 {
		return nil, 404
	}
	resp := new(ServiceResp)
	json.Unmarshal(reply[0], &resp)
	return resp.Body, int(resp.StatusCode)
}
