package main

import (
	"bytes"
	"fmt"

	"net/http"

	"text/template"

	mf "gofncstd3000"
)

func init() {
	rtr.HandleFunc("/edit", mf.LogreqF("/edit", page_edit))
	rtr.HandleFunc("/edit_save", mf.LogreqF("/edit_save", page_edit_save))
	fmt.Printf("")
}

//загрузка формы редактирования файлов
func page_edit(w http.ResponseWriter, r *http.Request) {

	dir, err := mf.AppPath()
	if checkError("get cur app dir error", err, w) {
		return
	}

	rtype := r.FormValue("type")
	file := r.FormValue("file")

	template_file := dir + "\\pages\\edit.html"

	template_text, err := mf.ReadFileStr(template_file)
	if checkError("read template file "+template_file+" error", err, w) {
		return
	}

	t, err := template.New("page_edit").Parse(template_text)
	if checkError("parse template error", err, w) {
		return
	}

	type UserVars struct {
		Data  string
		File  string
		Rtype string
	}
	vars := new(UserVars)
	vars.File = file
	vars.Rtype = rtype
	file_name := dir + "\\files\\" + rtype + "\\" + file
	vars.Data, err = mf.ReadFileStr(file_name)
	if checkError("read file "+file_name+" error", err, w) {
		return
	}

	buff := new(bytes.Buffer)
	err = t.Execute(buff, vars)
	if checkError("render template error", err, w) {
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(buff.Bytes())
}

//сохранение или удаление файлов
func page_edit_save(w http.ResponseWriter, r *http.Request) {
	dir, _ := mf.AppPath()

	rtype := r.FormValue("rtype")
	file := r.FormValue("file")
	path := dir + "\\files\\" + rtype + "\\" + file

	oper := r.FormValue("oper")

	if oper == "save" {
		data := r.FormValue("data")
		//json := r.FormValue("json")
		if mf.FileExists(path) {
			mf.CopyFile2(path, path+"_"+mf.CurTimeStrShort()+".back")
		}

		err := mf.WriteFileStr(path, data)
		if checkErrorJSON("WriteFile "+path+" error", err, w) {
			return
		}
	}

	if oper == "savexml" {
		xml := r.FormValue("xml")
		json := r.FormValue("json")
		if mf.FileExists(path) {
			tmp_name := mf.StrReplaceRegexp2(path, "\\.xml$", "")
			mf.CopyFile2(path, tmp_name+"_"+mf.CurTimeStrShort()+".xml")
			mf.CopyFile2(path+".json", tmp_name+"_"+mf.CurTimeStrShort()+".xml.json")
		}

		err := mf.WriteFileStr(path, xml)
		if checkErrorJSON("WriteFile "+path+" error", err, w) {
			return
		}

		err = mf.WriteFileStr(path+".json", json)
		if checkErrorJSON("WriteFile "+path+".json error", err, w) {
			return
		}
	}

	if oper == "delete" {
		err := mf.Rename(path, path+"_"+mf.CurTimeStrShort()+".deleted")
		if checkErrorJSON("Rename file "+path+" error", err, w) {
			return
		}
	}

	if oper == "delxml" {
		err := mf.Rename(path, path+"_"+mf.CurTimeStrShort()+".deleted")
		if checkErrorJSON("Rename file "+path+" error", err, w) {
			return
		}
		err = mf.Rename(path+".json", path+".json_"+mf.CurTimeStrShort()+".deleted")
		if checkErrorJSON("Rename file "+path+" error", err, w) {
			return
		}
	}

	d := make(map[string]string, 0)
	d["file"] = path
	j, err := mf.ToJson(d)
	if checkErrorJSON("ToJson error", err, w) {
		return
	}

	w.Header().Set("Content-Type", "text/json; charset=utf-8")
	w.Write(j)
}
