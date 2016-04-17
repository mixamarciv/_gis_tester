package main

import (
	"bytes"
	"net/http"
	"text/template"

	"fmt"
	mf "gofncstd3000"
)

func init() {
	rtr.HandleFunc("/renderquery", mf.LogreqF("/renderquery", ajax_renderquery)).Methods("POST")
	fmt.Printf("")
}

//генерим запрос на основе выбранного шаблона и выбранного набора данных
func ajax_renderquery(w http.ResponseWriter, r *http.Request) {
	json, err := getjson_from_datafiles(r)
	if checkErrorJSON("getjson_from_datafile error", err, w) {
		return
	}

	funcMap := template.FuncMap{
		"RandomGUID":   mf.Uuid,
		"CurDateTime1": func() string { s, _ := mf.StrReplaceRegexp(mf.CurTimeStr(), " ", "T"); return s },
		"CurDateTime2": mf.CurTimeStr,
	}

	type UserVars struct {
		RandomGUID1, RandomGUID2, RandomGUID3 string
		CurDateTime                           string
		HuisVer                               string
		Data                                  map[string]interface{}
		//name                     int
		//test2                    bool
	}

	vars := new(UserVars)
	vars.CurDateTime = mf.CurTimeStr()
	vars.RandomGUID1 = mf.Uuid()
	vars.HuisVer = r.FormValue("ver")
	vars.Data = json

	//рендерим xml
	t1, err := template.New("xml").Funcs(funcMap).Parse(r.FormValue("xml"))
	if checkErrorJSON("parse template1 error", err, w) {
		return
	}
	buff1 := new(bytes.Buffer)
	err = t1.Execute(buff1, vars)
	if checkErrorJSON("render template error", err, w) {
		return
	}

	//рендерим headers
	t2, err := template.New("data").Funcs(funcMap).Parse(r.FormValue("data"))
	if checkErrorJSON("parse template2 error", err, w) {
		return
	}
	buff2 := new(bytes.Buffer)
	err = t2.Execute(buff2, vars)
	if checkErrorJSON("render template error", err, w) {
		return
	}

	//собираем все в 1 json
	var ret []string
	ret = append(ret, string(buff1.Bytes()))
	ret = append(ret, string(buff2.Bytes()))

	json_ret, err := mf.ToJson(ret)
	if checkErrorJSON("ToJson error", err, w) {
		return
	}

	w.Header().Set("Content-Type", "text/json; charset=utf-8")
	w.Write(json_ret)
}

//загрузка json из файла
func getjson_from_datafiles(r *http.Request) (map[string]interface{}, error) {
	dir, err := mf.AppPath()
	if err != nil {
		return nil, err
	}

	datafilename := r.FormValue("datafilename_host")
	data_path := dir + "\\files\\data_host\\" + datafilename
	file_data, err := mf.ReadFileStr(data_path)
	if err != nil {
		return nil, err
	}

	if file_data[0:1] != "{" {
		file_data = "{" + file_data + "}"
	}

	json, err := mf.FromJson([]byte(file_data))
	if err != nil {
		return nil, err
	}

	datafilename = r.FormValue("datafilename_uk")
	data_path = dir + "\\files\\data_uk\\" + datafilename
	file_data, err = mf.ReadFileStr(data_path)
	if err != nil {
		return nil, err
	}

	if file_data[0:1] != "{" {
		file_data = "{" + file_data + "}"
	}

	json2, err := mf.FromJson([]byte(file_data))
	if err != nil {
		return nil, err
	}

	for k, _ := range json2 {
		json[k] = json2[k]
	}

	return json, err
}
