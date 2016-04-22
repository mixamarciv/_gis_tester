package main

import (
	"fmt"
	mf "gofncstd3000"
	"net/http"
	str "strings"
)

func init() {
	rtr.HandleFunc("/loaddatahostfileslist", mf.LogreqF("/loaddatahostfileslist", ajax_loaddatahostfileslist)).Methods("GET")
	rtr.HandleFunc("/loaddataukfileslist", mf.LogreqF("/loaddataukfileslist", ajax_loaddataukfileslist)).Methods("GET")
	rtr.HandleFunc("/loadversionlist", mf.LogreqF("/loadversionlist", ajax_loadversionlist)).Methods("GET")
	rtr.HandleFunc("/loadfileslist", mf.LogreqF("/loadfileslist", ajax_loadfileslist)).Methods("GET")
	rtr.HandleFunc("/loadfiledata", mf.LogreqF("/loadfiledata", ajax_loadfiledata)).Methods("GET")
	rtr.HandleFunc("/saveresult", mf.LogreqF("/saveresult", ajax_saveresult))
	rtr.HandleFunc("/loadresultlist", mf.LogreqF("/loadresultlist", ajax_loadresultlist))
	rtr.HandleFunc("/loadresult", mf.LogreqF("/loadresult", ajax_loadresult))

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

	re, _ := mf.RegexpCompile("deleted$")

	var filenames []string
	for _, file := range files {
		name := file.Name()
		if re.MatchString(name) {
			continue
		}
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
	if checkErrorJSON("get cur app dir error", err, w) {
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
	file_xml = prepare_xmlfile(file_xml)

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

func prepare_xmlfile(s string) string {

	ss := str.SplitAfter(s, "<soapenv:Header>")
	if len(ss) != 2 {
		return s
	}
	h1 := ss[0]

	ss = str.SplitAfter(ss[1], "</soapenv:Header>")
	if len(ss) != 2 {
		return s
	}

	h2 := ss[0]
	b := ss[1]

	h1 = mf.StrReplaceRegexp2(h1, "/\\d+\\.\\d+\\.\\d+\\.\\d+/", "{{.HuisVer}}")
	h2 = mf.StrReplaceRegexp2(h2, "<ns:Date>\\?</ns:Date>", "<ns:Date>{{CurDateTime1}}</ns:Date>")
	h2 = mf.StrReplaceRegexp2(h2, "<ns:MessageGUID>\\?</ns:MessageGUID>", "<ns:MessageGUID>{{RandomGUID}}</ns:MessageGUID>")
	h2 = mf.StrReplaceRegexp2(h2, "<ns:SenderID>\\?</ns:SenderID>", "<ns:SenderID>{{index .Data \"SenderID\"}}</ns:SenderID>")

	b = mf.StrReplaceRegexp2(b, "<ns:TransportGUID>\\?</ns:TransportGUID>", "<ns:TransportGUID>{{RandomGUID}}</ns:TransportGUID>")
	return h1 + h2 + b
}

//сохранение файла
func ajax_saveresult(w http.ResponseWriter, r *http.Request) {

	dir, err := mf.AppPath()
	if checkErrorJSON("get cur app dir error", err, w) {
		return
	}

	path := dir + "\\files\\result\\" + mf.CurTimeStrShort()[0:8]
	err = mf.MkdirAll(path)
	if checkErrorJSON("get cur app dir error", err, w) {
		return
	}

	d := make(map[string]string, 0)
	d["data"] = r.FormValue("data")
	d["xml"] = r.FormValue("xml")
	d["res_data"] = r.FormValue("res_data")
	d["res_xml"] = r.FormValue("res_xml")

	json, err := mf.ToJson(d)
	if checkErrorJSON("ToJson error", err, w) {
		return
	}

	file := path + "\\" + mf.CurTimeStrShort() + "_" + r.FormValue("name_prefix") + ".res"
	err = mf.WriteFile(file, json)
	if checkErrorJSON("WriteFile error", err, w) {
		return
	}

	w.Header().Set("Content-Type", "text/json; charset=utf-8")
	w.Write([]byte("{\"ok\":1}"))
}

//список файлов внутри версии
func ajax_loadresultlist(w http.ResponseWriter, r *http.Request) {

	dir, err := mf.AppPath()
	if checkErrorJSON("get cur app dir error", err, w) {
		return
	}

	path := dir + "\\files\\result\\"
	ex := mf.FileExists(path)
	if ex == false {
		w.Header().Set("Content-Type", "text/json; charset=utf-8")
		w.Write([]byte("[0:\"file not found " + path + "\"]"))
		return
	}

	files, err := mf.ReadDir(path)
	if checkErrorJSON("read dir "+path+" error", err, w) {
		return
	}

	use_filter := true
	filter := r.FormValue("filter")
	if len(filter) == 0 {
		use_filter = false
	}
	re, _ := mf.RegexpCompile(filter)
	var filenames []string
	for _, file := range files {
		name := file.Name()
		i_path := path + "\\" + file.Name()
		i_files, err := mf.ReadDir(i_path)
		if checkErrorJSON("read dir "+i_path+" error", err, w) {
			return
		}
		for _, i_file := range i_files {
			i_name := name + "/" + i_file.Name()
			if use_filter {
				if re.MatchString(i_name) {
					filenames = append(filenames, i_name)
				}
			} else {
				filenames = append(filenames, i_name)
			}
		}
	}

	json, err := mf.ToJson(filenames)
	if checkErrorJSON("ToJson error", err, w) {
		return
	}

	w.Header().Set("Content-Type", "text/json; charset=utf-8")
	w.Write(json)
}

//загрузка ранее сохраненного файла
func ajax_loadresult(w http.ResponseWriter, r *http.Request) {

	dir, err := mf.AppPath()
	if checkErrorJSON("get cur app dir error", err, w) {
		return
	}

	filename, _ := mf.StrReplaceRegexp(r.FormValue("file"), "/", "\\")
	path := dir + "\\files\\result\\" + filename

	file_data, err := mf.ReadFileStr(path)
	if checkErrorJSON("ReadFileStr error", err, w) {
		return
	}
	/*
		json, err := mf.FromJson([]byte(file_data))
		if err != nil {
			return nil, err
		}
	*/
	w.Header().Set("Content-Type", "text/json; charset=utf-8")
	w.Write([]byte(file_data))
}
