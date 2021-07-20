package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func getData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	exercise := vars["exercise"]

	if exercise == "first" {
		loadDataFirstExercise()
		json.NewEncoder(w).Encode(appData.FirstExercise)
	} else if exercise == "second" {
		loadDataSecondExercise()
		json.NewEncoder(w).Encode(appData.SecondExercise)
	}
	log.Printf("Request for Application data served")

}

func getResults(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	exercise := vars["exercise"]
	index := vars["index"]

	var result HTTPResponseObject
	result.Success = false

	if exercise == "first" {
		i, _ := strconv.Atoi(index)
		if i >= len(appData.FirstExercise) {
			w.WriteHeader(http.StatusInternalServerError)
			result.Success = false
			result.Message = "Index out of range."
			result.ErrorObject = &IndexOutOfRange{}
			json.NewEncoder(w).Encode(result)
		} else {
			var toReturn = FirstExercise{FirstWord: appData.FirstExercise[i].FirstWord, SecondWord: strings.Join(appData.FirstExercise[i].RandomWords, ","), Editable: "true", RandomWords: appData.FirstExercise[i].RandomWords}
			log.Printf("Second word: %s\n", appData.FirstExercise[i].SecondWord)
			json.NewEncoder(w).Encode(toReturn)
		}
	} else if exercise == "second" {
		i, _ := strconv.Atoi(index)
		if i >= len(appData.FirstExercise) {
			w.WriteHeader(http.StatusInternalServerError)
			result.Success = false
			result.Message = "Index out of range."
			result.ErrorObject = &IndexOutOfRange{}
			json.NewEncoder(w).Encode(result)
		} else {
			var toReturn = SecondExercise{Location: appData.SecondExercise[i].Location, Phone: "?", Editable: "true"}
			log.Printf("Second word: %s\n", appData.SecondExercise[i].Phone)
			json.NewEncoder(w).Encode(toReturn)
		}
	}

	log.Printf("Request for Exercise data served")

}

func getAppSettings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(appSettings)
	log.Printf("Request for Settings data served")
}

func saveAppSettings(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var result HTTPResponseObject

	log.Printf("New Application settings received")
	resBody, err := ioutil.ReadAll(r.Body)
	//data := []byte(`{"rowCount":2,"timeInSeconds":60}`)
	err = json.Unmarshal(resBody, &appSettings)
	//err := json.NewDecoder(r.Body).Decode("{rowCount: 22, timeInSeconds: 99}")
	if err != nil {
		// http.Error(w, err.Error(), http.StatusBadRequest)
		result.Success = false
		result.Message = "Issues decoding application data. " + err.Error()
		result.ErrorObject = err
		log.Printf(" - Issues decoding application data")
	} else {
		result.Success = true
		result.Message = "New application settings saved."
		log.Printf(" - New application settings saved.")
	}
	json.NewEncoder(w).Encode(result)

}

func queryExercise(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)

	first := vars["first"]
	second := vars["second"]
	exercise := vars["exercise"]
	var result HTTPResponseObject
	result.Success = false

	if exercise == "first" {
		var resultFirstExercise *FirstExercise

		for _, c := range appData.FirstExercise {
			if c.FirstWord == first && c.SecondWord == second {
				result.Success = true
				result.Message = c.SecondWord
				result.ErrorObject = nil
				resultFirstExercise = c
				// this is important to disable the row editing for successful solved ones
				resultFirstExercise.Editable = "false"
			}
		}

		if result.Success == false && result.ErrorObject == nil {
			w.WriteHeader(http.StatusInternalServerError)
			result.Success = false
			json.NewEncoder(w).Encode(result)
		} else {
			json.NewEncoder(w).Encode(resultFirstExercise)
		}
	} else if exercise == "second" {
		var resultSecondExercise *SecondExercise

		for _, c := range appData.SecondExercise {
			if c.Location == first && c.Phone == second {
				result.Success = true
				result.Message = c.Phone
				result.ErrorObject = nil
				resultSecondExercise = c
				// this is important to disable the row editing for successful solved ones
				resultSecondExercise.Editable = "false"
			}
		}

		if result.Success == false && result.ErrorObject == nil {
			w.WriteHeader(http.StatusInternalServerError)
			result.Success = false
			json.NewEncoder(w).Encode(result)
		} else {
			json.NewEncoder(w).Encode(resultSecondExercise)
		}
	}
}
