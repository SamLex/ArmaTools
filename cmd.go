package main

import (
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/SamLex/armatoolsserver/armatools/unitcapturereduce"
	"github.com/SamLex/armatoolsserver/httphandlerutils"
)

var unitCaptureReduceTemplate = template.Must(template.ParseFiles("template/unitcapturereduce.html"))

var httpBindAddress = flag.String("bind", "127.0.0.1:80", "Bind address for HTTP server")

func init() {
	flag.Parse()
}

func main() {
	mux := http.NewServeMux()

	// HTTP handler that responds to / and only /, logs all requests, accepts only GET and POST,
	// times out after a certain time and serves the content from unitCaptureReducePageHandler
	mux.Handle("/",
		httphandlerutils.SimpleLogging(
			httphandlerutils.SupportedMethods(
				httphandlerutils.ExactPath("/",
					httphandlerutils.PanicHandler(
						http.HandlerFunc(unitCaptureReducePageHandler))), "GET", "POST")))

	server := http.Server{
		Addr:    *httpBindAddress,
		Handler: mux,
	}

	log.Fatalln(server.ListenAndServe())
}

func unitCaptureReducePageHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var errorString string

	var orginalFrames int
	var orginalKilobytes float64
	var reducedFrames int
	var reducedKilobytes float64
	var reduced string

	if r.Method == "POST" {
		var errorPercent float64

		captureData := r.FormValue("captureData")
		errorPercent, err = strconv.ParseFloat(r.FormValue("errorPercent"), 64)

		if err == nil {
			errorProbability := (100 - errorPercent) / 100
			reduced, orginalFrames, reducedFrames, err = unitcapturereduce.ReduceUnitCapture(captureData, errorProbability)
			if err == nil {
				orginalKilobytes = float64(len([]byte(captureData))) / 1024
				reducedKilobytes = float64(len([]byte(reduced))) / 1024
			}
		}
	}

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

	type Result struct {
		OrginalFrames    int
		OrginalKilobytes float64
		ReducedFrames    int
		ReducedKilobytes float64
		Reduced          string
	}

	data := make(map[string]interface{})
	data["Error"] = errorString
	data["Result"] = &Result{
		OrginalFrames:    orginalFrames,
		OrginalKilobytes: orginalKilobytes,
		ReducedFrames:    reducedFrames,
		ReducedKilobytes: reducedKilobytes,
		Reduced:          reduced,
	}

	err = unitCaptureReduceTemplate.Execute(w, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
