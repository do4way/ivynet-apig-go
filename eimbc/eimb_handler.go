package eimbc

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"sync"

	"github.com/do4way/ivynet-eimb-go/eimb"
	"github.com/gorilla/mux"
	"golang.org/x/net/html/charset"
)

//Request a common eimb request.
type Request struct {
	Method string
	Topic  string
	ID     string
	Data   string
}

//Response an common eimb response
type Response struct {
	Status  int
	Message string
}

type msgerMapper struct {
	messengers sync.Map
}

var topicPrefix = os.Getenv("APP_DOMAIN") + "::"

func (mapper *msgerMapper) get(topic string) *eimb.Messenger {
	qualifiedTopic := topicPrefix + topic
	msger, ok := mapper.messengers.Load(qualifiedTopic)
	if !ok {
		msger = eimb.NewMessenger(qualifiedTopic)
		mapper.messengers.Store(topic, msger)
	}
	return msger.(*eimb.Messenger)
}

var msgers = &msgerMapper{
	messengers: sync.Map{},
}

//HTTPPostHandler handle http post request.
func HTTPPostHandler(w http.ResponseWriter, r *http.Request) {

	req, err := parseRequest(r)
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse(203, err.Error()))
	}
	msger := msgers.get(req.Topic)
	msger.Send(req.Data)
}

//HTTPGetHandler handle http get request.
func HTTPGetHandler(w http.ResponseWriter, r *http.Request) {
}

func parseRequest(r *http.Request) (*Request, error) {

	params := mux.Vars(r)
	topic := params["topic"]
	id := params["id"]
	bodyReader, err := charset.NewReader(r.Body, getContentType(r))
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(bodyReader)
	if err != nil {
		return nil, err
	}
	return &Request{
		Method: r.Method,
		Topic:  topic,
		ID:     id,
		Data:   string(data),
	}, nil
}

func getContentType(r *http.Request) string {
	ct := r.Header["content-type"]
	if len(ct) > 0 {
		return ct[0]
	}
	return ""
}

//ErrorResponse generate error response.
func ErrorResponse(status int, msg string) *Response {

	return &Response{
		Status:  status,
		Message: msg,
	}
}
