package main

import (
	"net/http"

	"strings"

	"github.com/parnurzeal/gorequest"

	"fmt"
	mf "gofncstd3000"

	"errors"
)

func init() {
	rtr.HandleFunc("/sendquery", mf.LogreqF("/sendquery", ajax_sendquery)).Methods("POST")
	fmt.Printf("")
}

//отправляем запрос
func ajax_sendquery(w http.ResponseWriter, r *http.Request) {
	xml := r.FormValue("xml")
	data := strings.Trim(r.FormValue("data"), " \n\r\t")
	if data[0:1] != "{" {
		data = "{" + data + "}"
	}

	json, err := mf.FromJson([]byte(data))
	if checkErrorJSON("FromJson error", err, w) {
		return
	}

	if _, ok := json["url"].(string); !ok {
		checkErrorJSON("json param \"url\" not found", errors.New("json param1 \"url\" not found!!"), w)
		return
	}

	url := json["url"].(string)
	req := gorequest.New().Post(url)

	if basicAuth, ok := json["basicAuth"].(map[string]interface{}); ok {
		fmt.Printf("has basicAuth\n")
		req = req.SetBasicAuth(basicAuth["user"].(string), basicAuth["pass"].(string))
	} else {
		fmt.Printf("no basicAuth\n")
	}

	if reqtype, ok := json["type"].(string); ok {
		req = req.Type(reqtype)
	}

	if headers, ok := json["headers"].(map[string]interface{}); ok {
		for key, value := range headers {
			req.Set(key, value.(string))
		}
	}

	resp, body, errs := req.Send(xml).End()

	if checkErrors("request send error", errs, w) {
		return
	}

	var ret []string
	t1 := fmt.Sprintf("%+v\n=================================\n%+v", req.Header, resp.Header)
	ret = append(ret, body)
	ret = append(ret, t1)

	json_ret, err := mf.ToJson(ret)
	if checkErrorJSON("ToJson error", err, w) {
		return
	}

	w.Header().Set("Content-Type", "text/json; charset=utf-8")
	w.Write(json_ret)
}
