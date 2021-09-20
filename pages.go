package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
)

var templates = template.Must(template.ParseGlob("static/templates/*"))

func renderResultsSecondExercise(w http.ResponseWriter, r *http.Request) {
	pageData := PageData{
		VDir:            appConfig.VDir,
		RowCount:        appSettings.RowCount,
		Exercise:        "second",
		GridfieldName1:  "Location",
		GridfieldTitle1: "Location",
		GridfieldName2:  "Phone",
		GridfieldTitle2: "Phone number",
	}
	if err := templates.ExecuteTemplate(w, "resultsSecondExercise", pageData); err != nil {
		log.Fatalf("Error rendering template: %s", err)
	}
	log.Println("Rendered the results second exercise.")
}

func redirectToResultsSecondExercise(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, appConfig.VDir+"/results/second/", http.StatusFound)
}

func renderResultsFirstExercise(w http.ResponseWriter, r *http.Request) {
	pageData := PageData{
		VDir:            appConfig.VDir,
		RowCount:        appSettings.RowCount,
		Exercise:        "first",
		GridfieldName1:  "FirstWord",
		GridfieldTitle1: "First Word",
		GridfieldName2:  "SecondWord",
		GridfieldTitle2: "Second Word",
	}
	if err := templates.ExecuteTemplate(w, "resultsFirstExercise", pageData); err != nil {
		log.Fatalf("Error rendering template: %s", err)
	}
	log.Println("Rendered the results first exercise.")
}

func redirectToResultsFirstExercise(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, appConfig.VDir+"/results/first/", http.StatusFound)
}

func renderSecondExcerise(w http.ResponseWriter, r *http.Request) {
	pageData := PageData{
		VDir:            appConfig.VDir,
		TimerTime:       appSettings.TimeInSeconds,
		Exercise:        "second",
		GridfieldName1:  "Location",
		GridfieldTitle1: "Location",
		GridfieldName2:  "Phone",
		GridfieldTitle2: "Phone number",
	}
	if err := templates.ExecuteTemplate(w, "secondExercise", pageData); err != nil {
		log.Fatalf("Error rendering template: %s", err)
	}
	log.Println("Rendered the second exercise.")
}

func redirectToSecondExcerise(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, appConfig.VDir+"/exercises/second/", http.StatusFound)
}

func renderFirstExcerise(w http.ResponseWriter, r *http.Request) {
	pageData := PageData{
		VDir:            appConfig.VDir,
		TimerTime:       appSettings.TimeInSeconds,
		Exercise:        "first",
		GridfieldName1:  "FirstWord",
		GridfieldTitle1: "First Word",
		GridfieldName2:  "SecondWord",
		GridfieldTitle2: "Second Word",
	}
	if err := templates.ExecuteTemplate(w, "firstExercise", pageData); err != nil {
		log.Fatalf("Error rendering template: %s", err)
	}
	log.Println("Rendered the first exercise.")
}

func redirectToFirstExcerise(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, appConfig.VDir+"/exercises/first/", http.StatusFound)
}

func renderSettings(w http.ResponseWriter, r *http.Request) {
	pageData := PageData{
		VDir: appConfig.VDir,
	}
	if err := templates.ExecuteTemplate(w, "settings", pageData); err != nil {
		log.Fatalf("Error rendering template: %s", err)
	}
	log.Println("Rendered settings.")
}

func redirectToRenderSettings(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, appConfig.VDir+"/settings", http.StatusFound)
}

func checkHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	io.WriteString(w, "alive")
}
