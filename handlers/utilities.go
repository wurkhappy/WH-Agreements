package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/wurkhappy/WH-Config"
	"github.com/wurkhappy/mdp"
	"net/http"
)

func parseRequest(req *http.Request) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	return buf.Bytes()
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

type ServiceResp struct {
	StatusCode float64 `json:"status_code"`
	Body       []byte  `json:"body"`
}

func sendServiceRequest(method, service, path string, body []byte, userID string) (response []byte, statusCode int) {
	client := mdp.NewClient(config.MDPBroker, false)
	defer client.Close()
	m := map[string]interface{}{
		"Method": method,
		"Path":   path,
		"Body":   body,
	}
	req, _ := json.Marshal(m)
	request := [][]byte{req, []byte(userID)}
	reply := client.Send([]byte(service), request)
	if len(reply) == 0 {
		return nil, 404
	}
	resp := new(ServiceResp)
	json.Unmarshal(reply[0], &resp)
	return resp.Body, int(resp.StatusCode)
}
