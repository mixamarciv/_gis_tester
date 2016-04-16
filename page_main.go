package main

import (
	"bytes"
	"fmt"

	"net/http"

	"text/template"

	mf "gofncstd3000"
)

func init() {
	rtr.HandleFunc("/", mf.LogreqF("/", page_main)).Methods("GET")
	fmt.Printf("")
}

func page_main(w http.ResponseWriter, r *http.Request) {

	dir, err := mf.AppPath()
	if checkError("get cur app dir error", err, w) {
		return
	}

	template_file := dir + "\\pages\\main.html"
	//fmt.Println(mf.CurTimeStrShort()+" template_file: ", template_file)

	template_text, err := mf.ReadFileStr(template_file)
	if checkError("read template file error", err, w) {
		return
	}

	t, err := template.New("page_main").Parse(template_text)
	if checkError("parse template error", err, w) {
		return
	}

	type UserVars struct {
		RandomGUID1, CurDateTime string
		HuisVer                  string
		//name                     int
		//test2                    bool
	}
	vars := new(UserVars)
	vars.CurDateTime = mf.CurTimeStr()
	vars.RandomGUID1 = mf.Uuid() //fmt.Sprintf("%s", uuid.NewV4())
	vars.HuisVer = "8.6.0.6"

	buff := new(bytes.Buffer)
	err = t.Execute(buff, vars)
	if checkError("render template error", err, w) {
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(buff.Bytes())
	//ioutil.WriteFile("req/resp_"+cur_time_str2(), []byte(s), 0644)
}
