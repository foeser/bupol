package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type Phone struct {
	Land   string `json:"land"`
	Mobile string `json:"mobile"`
}

type Person struct {
	Username string `json:"username"`
	Phone    Phone  `json:"phone"`
	Location string
}

type Entry struct {
	Person     Person
	RandomWord string
}

var currentEntries []Entry

var templates = template.Must(template.ParseGlob("static/templates/*"))

func getGermanWords() []string {
	resp, err := http.Get("https://www.palabrasaleatorias.com/zufallige-worter.php?fs=10&fs2=0&Submit=Neues+Wort")
	if err != nil {
		log.Fatalln(err)
	}
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//Convert the body to type string
	sb := string(body)
	//log.Printf(sb)
	scanner := bufio.NewScanner(strings.NewReader(sb))

	var result []string
	for scanner.Scan() {
		line := scanner.Text()
		//fmt.Println(line)
		if strings.Contains(line, "<br /><div style=\"font-size:3em; color:#6200C5;\">") {
			// get next line and remove div to get word
			scanner.Scan()
			//fmt.Println(scanner.Text())
			result = append(result, strings.TrimSuffix(scanner.Text(), "</div>"))
		}

	}

	if scanner.Err() != nil {
		log.Println(scanner.Err())
	}

	return result
}

func getRandomPersonData(count int) []*Person {
	url := fmt.Sprintf("https://randomname.de/?format=json&count=%d&phone=a", count)
	res, err := http.Get(url)
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}
	var persons []*Person
	json.Unmarshal([]byte(body), &persons)
	//fmt.Printf("Persons : %+v", persons)

	// add additional randomized location
	locations := readLocation("locations")
	rand.Seed(time.Now().UnixNano())

	/*for index := range persons {
		persons[index].Location = locations[rand.Intn(len(locations))]
	}*/

	for _, person := range persons {
		person.Location = locations[rand.Intn(len(locations))]
	}

	return persons
}

func readLocation(fileName string) []string {
	fileBytes, err := ioutil.ReadFile("locations")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	locations := strings.Split(string(fileBytes), "\n")

	//fmt.Println(locations)

	return locations
}

func main() {

	//fileServer := http.FileServer(http.Dir("./static")) // New code
	//http.Handle("/", fileServer)
	//http.HandleFunc("/", wordsHandler)
	http.HandleFunc("/words", wordsHandler)
	http.HandleFunc("/phoneNumbers", phoneNumbersHandler)
	http.HandleFunc("/results", resultsHandler)

	fmt.Printf("Starting server at port 8888\n")
	if err := http.ListenAndServe("0.0.0.0:8888", nil); err != nil {
		log.Fatal(err)
	}
}

func resultsHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "results", currentEntries)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
func phoneNumbersHandler(w http.ResponseWriter, r *http.Request) {
	// reset current entries and get new data
	persons := getRandomPersonData(10)

	// merge results from currentPersons (second exercise: location + phone) into currentEntries
	for i := range currentEntries {
		currentEntries[i].Person.Location = persons[i].Location
		currentEntries[i].Person.Phone.Mobile = persons[i].Phone.Mobile
	}

	err := templates.ExecuteTemplate(w, "phonenumbers", currentEntries)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
func wordsHandler(w http.ResponseWriter, r *http.Request) {
	/*if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	fmt.Fprintf(w, "POST request successful")
	name := r.FormValue("name")
	address := r.FormValue("address")

	fmt.Fprintf(w, "Name = %s\n", name)
	fmt.Fprintf(w, "Address = %s\n", address)*/

	// reset current entries and get new data
	persons := getRandomPersonData(10)
	currentEntries = nil
	for i, s := range getGermanWords() {
		currentEntries = append(currentEntries, Entry{Person: *persons[i], RandomWord: s})
	}

	err := templates.ExecuteTemplate(w, "words", currentEntries)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
