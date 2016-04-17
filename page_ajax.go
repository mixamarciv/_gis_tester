package main

import (
	"fmt"
	mf "gofncstd3000"
	"net/http"
)

func init() {
	rtr.HandleFunc("/loaddatahostfileslist", mf.LogreqF("/loaddatahostfileslist", ajax_loaddatahostfileslist)).Methods("GET")
	rtr.HandleFunc("/loaddataukfileslist", mf.LogreqF("/loaddataukfileslist", ajax_loaddataukfileslist)).Methods("GET")
	rtr.HandleFunc("/loadversionlist", mf.LogreqF("/loadversionlist", ajax_loadversionlist)).Methods("GET")
	rtr.HandleFunc("/loadfileslist", mf.LogreqF("/loadfileslist", ajax_loadfileslist)).Methods("GET")
	rtr.HandleFunc("/loadfiledata", mf.LogreqF("/loadfiledata", ajax_loadfiledata)).Methods("GET")

	fmt.Printf("")
}

//загрузка списка файлов с данными
func ajax_loaddatahostfileslist(w http.ResponseWriter, r *http.Request) {
	ajax_loaddatafileslist("data_host", w, r)
}
func ajax_loaddataukfileslist(w http.ResponseWriter, r *http.Request) {
	ajax_loaddatafileslist("data_uk", w, r)
}

func ajax_loaddatafileslist(typefiles string, w http.ResponseWriter, r *http.Request) {
	dir, err := mf.AppPath()
	if checkErrorJSON("get cur app dir error", err, w) {
		return
	}

	xml_path := dir + "\\files\\" + typefiles
	files, err := mf.ReadDir(xml_path)
	if checkErrorJSON("read dir "+xml_path+" error", err, w) {
		return
	}

	var filenames []string
	for _, file := range files {
		name := file.Name()
		filenames = append(filenames, name)
	}

	json, err := mf.ToJson(filenames)
	if checkErrorJSON("ToJson error", err, w) {
		return
	}

	w.Header().Set("Content-Type", "text/json; charset=utf-8")
	w.Write(json)
}

//список версий
func ajax_loadversionlist(w http.ResponseWriter, r *http.Request) {
	dir, err := mf.AppPath()
	if checkErrorJSON("get cur app dir error", err, w) {
		return
	}

	xml_path := dir + "\\files\\xml"
	files, err := mf.ReadDir(xml_path)
	if checkErrorJSON("read dir "+xml_path+" error", err, w) {
		return
	}

	var filenames []string
	for _, file := range files {
		name := file.Name()
		filenames = append(filenames, name)
	}

	json, err := mf.ToJson(filenames)
	if checkErrorJSON("ToJson error", err, w) {
		return
	}

	w.Header().Set("Content-Type", "text/json; charset=utf-8")
	w.Write(json)
}

//список файлов внутри версии
func ajax_loadfileslist(w http.ResponseWriter, r *http.Request) {

	dir, err := mf.AppPath()
	if checkErrorJSON("get cur app dir error", err, w) {
		return
	}

	ver := r.FormValue("ver")

	xml_path := dir + "\\files\\xml\\" + ver
	ex := mf.FileExists(xml_path)
	if ex == false {
		w.Header().Set("Content-Type", "text/json; charset=utf-8")
		w.Write([]byte("[0:\"file not found " + xml_path + "\"]"))
		return
	}

	files, err := mf.ReadDir(xml_path)
	if checkErrorJSON("read dir "+xml_path+" error", err, w) {
		return
	}

	re, _ := mf.RegexpCompile("xml$")
	var filenames []string
	for _, file := range files {
		name := file.Name()
		if re.MatchString(name) {
			filenames = append(filenames, name)
		}
	}

	json, err := mf.ToJson(filenames)
	if checkErrorJSON("ToJson error", err, w) {
		return
	}

	w.Header().Set("Content-Type", "text/json; charset=utf-8")
	w.Write(json)
}

//загрузка данных по выбранному файлу
func ajax_loadfiledata(w http.ResponseWriter, r *http.Request) {

	dir, err := mf.AppPath()
	if checkError("get cur app dir error", err, w) {
		return
	}

	ver := r.FormValue("ver")
	filename := r.FormValue("filename")

	xml_path := dir + "\\files\\xml\\" + ver + "\\" + filename
	ex := mf.FileExists(xml_path)
	if ex == false {
		w.Header().Set("Content-Type", "text/json; charset=utf-8")
		w.Write([]byte("[0:\"file not found " + xml_path + "\"]"))
		return
	}

	file_xml, err := mf.ReadFileStr(xml_path)
	if err != nil {
		w.Header().Set("Content-Type", "text/json; charset=utf-8")
		w.Write([]byte("[0:\"read file " + xml_path + " error\"]"))
		return
	}

	data_path := xml_path + ".json"
	file_data, err := mf.ReadFileStr(data_path)
	if err != nil {
		w.Header().Set("Content-Type", "text/json; charset=utf-8")
		w.Write([]byte("[0:\"read file " + data_path + " error\"]"))
		return
	}

	var t []string
	t = append(t, file_xml)
	t = append(t, file_data)

	json, err := mf.ToJson(t)
	if checkErrorJSON("ToJson error", err, w) {
		return
	}

	w.Header().Set("Content-Type", "text/json; charset=utf-8")
	w.Write(json)
}
