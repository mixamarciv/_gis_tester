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
	rtr.HandleFunc("/signquery", mf.LogreqF("/signquery", ajax_signquery)).Methods("POST")
	fmt.Printf("")
}

//отправляем запрос
func ajax_signquery(w http.ResponseWriter, r *http.Request) {
	xml := r.FormValue("xml")
	data := strings.Trim(r.FormValue("data"), " \n\r\t")
	if data[0:1] != "{" {
		data = "{" + data + "}"
	}

	json, err := mf.FromJson([]byte(data))
	if checkErrorJSON("FromJson error", err, w) {
		return
	}

	if _, ok := json["signserv"].(string); !ok {
		checkErrorJSON("json param \"signserv\" not found", errors.New("json param1 \"signserv\" not found!!"), w)
		return
	}

	url := json["signserv"].(string)
	req := gorequest.New().Post(url)

	//собираем все в 1 json
	ret := make(map[string]string)
	ret["xml"] = xml
	ret["data"] = data

	json_str, err := mf.ToJson(ret)
	if checkErrorJSON("ToJson error", err, w) {
		return
	}

	_, body, errs := req.Send(string(json_str)).End()
	//_, body, errs := req.Set("xml", xml).Set("data", data).Set("json_str", "").End()

	if checkErrors("request send error", errs, w) {
		return
	}

	w.Header().Set("Content-Type", "text/json; charset=utf-8")
	w.Write([]byte(body))
}
