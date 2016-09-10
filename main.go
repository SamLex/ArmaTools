package main

import (
	"armaUnitCaptureReduceCmd/arma/unitcapturereduce"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

var indexTemplate = template.Must(template.ParseFiles("template/index.html"))

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.TimeoutHandler(http.HandlerFunc(indexPageHandler), 10*time.Second, ""))
	mux.Handle("/favicon.ico", http.NotFoundHandler())

	server := http.Server{
		Addr:    "127.0.0.1:4821",
		Handler: mux,
	}

	log.Println(server.ListenAndServe())
}

func indexPageHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.Proto, r.ContentLength, r.Host, r.RemoteAddr, r.RequestURI)

	var err error
	var orginalFrames int
	var orginalKilobytes float64
	var reducedFrames int
	var reducedKilobytes float64
	var reduced string

	if r.Method == "POST" {
		captureData := r.FormValue("captureData")
		stringErrorPercent := r.FormValue("errorPercent")

		var errorPercent float64
		errorPercent, err = strconv.ParseFloat(stringErrorPercent, 64)
		if err == nil {
			reduced, orginalFrames, reducedFrames, err = unitcapturereduce.ReduceUnitCapture(captureData, errorPercent)
			if err == nil {
				orginalKilobytes = float64(len([]byte(captureData))) / 1024
				reducedKilobytes = float64(len([]byte(reduced))) / 1024
			}
		}
	}

	errorString := ""
	if err != nil {
		_, isSyntaxError := err.(*json.SyntaxError)
		_, isNumError := err.(*strconv.NumError)
		if isSyntaxError {
			errorString = "Invalid UnitCapture Output"
		} else if isNumError {
			errorString = "Invalid Percentage Error Threshold"
		} else {
			errorString = err.Error()
		}
	}

	data := struct {
		Error  string
		Result struct {
			OrginalFrames    int
			OrginalKilobytes float64
			ReducedFrames    int
			ReducedKilobytes float64
			Reduced          string
		}
	}{
		Error: errorString,
		Result: struct {
			OrginalFrames    int
			OrginalKilobytes float64
			ReducedFrames    int
			ReducedKilobytes float64
			Reduced          string
		}{
			OrginalFrames:    orginalFrames,
			OrginalKilobytes: orginalKilobytes,
			ReducedFrames:    reducedFrames,
			ReducedKilobytes: reducedKilobytes,
			Reduced:          reduced,
		},
	}

	indexTemplate.Execute(w, data)
}
