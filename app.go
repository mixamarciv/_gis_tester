package main

import (
	//"bytes"
	"fmt"
	//"io/ioutil"
	"log"
	"net/http"
	//"strconv"

	//"crypto/md5"
	//"regexp"
	//"strings"
	//"time"

	//"github.com/satori/go.uuid"

	"github.com/gorilla/mux"
	//"github.com/parnurzeal/gorequest"

	//"text/template"

	mf "gofncstd3000"

	//"github.com/palantir/stacktrace"
	//"runtime/debug"
)

var rtr *mux.Router

func init() {
	rtr = mux.NewRouter()
}

func main() {
	log.Println(mf.CurTimeStrShort() + " start app")

	fs := http.FileServer(http.Dir("./files/"))
	s := http.StripPrefix("/files/", fs)

	rtr.PathPrefix("/files/").Handler(mf.LogreqH("/files/", s))

	http.Handle("/", rtr)

	listenPort := 8080
	log.Println("Listening port: " + mf.IntToStr(listenPort))
	err := http.ListenAndServe(":"+mf.IntToStr(listenPort), nil)
	fmt.Sprintf("%+v", err)
}

func checkError(title string, err error, w http.ResponseWriter) bool {
	if err != nil {
		serr := "\n\n== ERROR: ======================================\n"
		serr += title + "\n"
		serr += mf.ErrStr(err)
		serr += "\n\n== /ERROR ======================================\n"
		log.Println(serr)

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte(serr))
		return true
	}
	return false
}

func checkErrorJSON(title string, err error, w http.ResponseWriter) bool {
	if err != nil {
		serr := "\n\n== ERROR: ======================================\n"
		serr += title + "\n"
		serr += mf.ErrStr(err)
		serr += "\n\n== /ERROR ======================================\n"
		log.Println(serr)

		w.Header().Set("Content-Type", "text/json; charset=utf-8")
		serr, _ = mf.StrReplaceRegexp(serr, "\"", "\\\"")
		w.Write([]byte("{\"err\":\"" + serr + "\"}"))
		return true
	}
	return false
}

func checkErrors(title string, err []error, w http.ResponseWriter) bool {
	if err != nil {
		serr := "\n\n== ERRORs: =====================================\n"
		//err := stacktrace.Propagate(err[0], title)
		serr += title + "\n"
		serr += fmt.Sprintf("%+v", err)
		serr += "\n\n== /ERRORs =====================================\n"
		log.Println(serr)

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte(serr))
		return true
	}
	return false
}
