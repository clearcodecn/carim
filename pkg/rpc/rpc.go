package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
)

type handler func(req interface{}) (v interface{}, err error)

type server struct {
	handlers map[string]handler
	params   map[string]interface{}
}

func main() {
	s := new(server)
	s.handlers = make(map[string]handler)
	s.params = make(map[string]interface{})
	s.handlers["com.demo.queryUser"] = queryUser
	s.params["com.demo.queryUser"] = QueryUserParams{}

	s.start()
}

type QueryUserParams struct {
	Username string `json:"username"`
}

type QueryUserResponse struct {
	Data []Users `json:"data"`
}

type Users struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func queryUser(req interface{}) (v interface{}, err error) {
	// 查询用户.
	// TODO:: ...
	param := req.(*QueryUserParams)

	log.Println("receive:", param.Username)
	_ = param
	resp := &QueryUserResponse{
		Data: []Users{
			{
				Username: "foo",
				Password: "bar",
			},
		},
	}

	return resp, nil
}

func (s *server) start() {
	http.ListenAndServe(":6000", s)
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	var param = make(map[string]string)
	err = json.Unmarshal(data, &param)
	if err != nil {
		return
	}
	// 找方法
	method, ok := param["method"]
	if !ok {
		return
	}
	// 找handler
	h, ok := s.handlers[method]
	if !ok {
		return
	}
	p, ok := s.params[method]
	if !ok {
		return
	}
	pi := reflect.New(reflect.TypeOf(p)).Interface()
	// 反序列化参数
	err = json.Unmarshal([]byte(param["data"]), &pi)
	if err != nil {
		return
	}
	// 调用函数
	resp, err := h(pi)
	if err != nil {
		return
	}
	// 序列化响应
	data, err = json.Marshal(resp)
	if err != nil {
		return
	}
	// 写参数
	w.Write(data)
}
