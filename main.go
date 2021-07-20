package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

func main() {
	setWorkingDir()
	loadConfig()
	setAppSettings()
	loadLocationsData()
	loadRandomWords()
	setupWebServer()
}

func setAppSettings() {
	appSettings = AppSettings{RowCount: 2, TimeInSeconds: 60}
}

func setWorkingDir() {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Error getting executables path: %s", err)
	}
	workingDir := filepath.Dir(exePath)
	if err = os.Chdir(workingDir); err != nil {
		log.Fatalf("Couldn't set working directory: %s", err)
	}
	log.Printf("Set working directory to: %s", workingDir)
}

func loadConfig() {
	if err := cleanenv.ReadConfig("config.json", &appConfig); err != nil {
		log.Fatalf("Error loading config.json file: %s", err)
	}
	log.Printf("Application config loaded from config.json successfully.")
}

func setupWebServer() {
	// Init HTTP Router - mux
	router := mux.NewRouter()

	// map directory to server static files
	router.PathPrefix(appConfig.VDir + "/static/").Handler(http.StripPrefix(appConfig.VDir+"/static/", http.FileServer(http.Dir("./static"))))
	router.PathPrefix(appConfig.VDir + "/exercises/first/static/").Handler(http.StripPrefix(appConfig.VDir+"/exercises/first/static/", http.FileServer(http.Dir("./static"))))
	router.PathPrefix(appConfig.VDir + "/exercises/second/static/").Handler(http.StripPrefix(appConfig.VDir+"/exercises/second/static/", http.FileServer(http.Dir("./static"))))
	router.PathPrefix(appConfig.VDir + "/results/first/static/").Handler(http.StripPrefix(appConfig.VDir+"/results/first/static/", http.FileServer(http.Dir("./static"))))
	router.PathPrefix(appConfig.VDir + "/results/second/static/").Handler(http.StripPrefix(appConfig.VDir+"/results/second/static/", http.FileServer(http.Dir("./static"))))
	router.PathPrefix(appConfig.VDir + "/exercises/first/appSettings/static/").Handler(http.StripPrefix(appConfig.VDir+"/exercises/first/appSettings/static/", http.FileServer(http.Dir("./static"))))

	// Define health check function
	router.HandleFunc(appConfig.VDir+"/health", checkHealth).Methods("GET")

	router.HandleFunc(appConfig.VDir, redirectToFirstExcerise).Methods("GET")

	// first exercise
	subrouter := router.PathPrefix(appConfig.VDir + "/exercises").Subrouter()
	subrouter.HandleFunc("", redirectToFirstExcerise).Methods("GET")
	subrouter.HandleFunc("/", redirectToFirstExcerise).Methods("GET")
	subrouter.HandleFunc("/first/", renderFirstExcerise).Methods("GET")
	subrouter.HandleFunc("/first", redirectToFirstExcerise).Methods("GET")
	// second exercise
	subrouter.HandleFunc("/second/", renderSecondExcerise).Methods("GET")
	subrouter.HandleFunc("/second", redirectToSecondExcerise).Methods("GET")

	// results
	subrouter = router.PathPrefix(appConfig.VDir + "/results").Subrouter()
	subrouter.HandleFunc("", redirectToResultsFirstExercise).Methods("GET")
	subrouter.HandleFunc("/", redirectToResultsFirstExercise).Methods("GET")
	subrouter.HandleFunc("/first/", renderResultsFirstExercise).Methods("GET")
	subrouter.HandleFunc("/first", redirectToResultsFirstExercise).Methods("GET")

	subrouter.HandleFunc("/second/", renderResultsSecondExercise).Methods("GET")
	subrouter.HandleFunc("/second", redirectToResultsSecondExercise).Methods("GET")

	// settings
	subrouter = router.PathPrefix(appConfig.VDir + "/appSettings").Subrouter()
	subrouter.HandleFunc("", renderSettings).Methods("GET")
	subrouter.HandleFunc("/get", getAppSettings).Methods("GET")
	//subrouter.HandleFunc("/", redirectToRenderSettings).Methods("GET")
	//subrouter.HandleFunc("/get", getSettings).Methods("GET")
	subrouter.HandleFunc("/set", saveAppSettings).Methods("POST")

	// exercise and results pages query this
	router.HandleFunc(appConfig.VDir+"/data/get/{exercise}", getData).Methods("GET")
	router.HandleFunc(appConfig.VDir+"/data/get/{exercise}/", getData).Methods("GET")

	router.HandleFunc(appConfig.VDir+"/data/getResults/{exercise}/{index}", getResults).Methods("GET")
	router.HandleFunc(appConfig.VDir+"/data/getResults/{exercise}/{index}/", getResults).Methods("GET")

	router.HandleFunc(appConfig.VDir+"/query/{exercise}/{first}/{second}", queryExercise).Methods("GET")

	// Setup webserver
	httpListen := "0.0.0.0:" + strconv.Itoa(appConfig.Port)
	log.Printf("Starting WebServer on %s", httpListen)

	srv := &http.Server{
		Handler:      handlers.RecoveryHandler(handlers.PrintRecoveryStack(true))(router),
		Addr:         httpListen,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
