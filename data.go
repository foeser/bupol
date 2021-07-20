package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

var appConfig AppConfig
var appData AppData
var appSettings AppSettings

func loadDataFirstExercise() {
	appData.FirstExercise = nil
	getRandomPersonData()
	addRandomWords()
}

func loadDataSecondExercise() {
	appData.SecondExercise = nil
	rand.Seed(time.Now().UnixNano())

	/*for index := range appData.PersistentPersons {
		appData.PersistentPersons[index].Location = locations[rand.Intn(len(locations))]
	}*/

	for i := 0; i < appSettings.RowCount; i++ {
		appData.SecondExercise = append(appData.SecondExercise, &SecondExercise{Location: appData.Locations[rand.Intn(len(appData.Locations))], Phone: fmt.Sprintf("%03d", rand.Intn(1000))})
	}
}

func loadRandomWords() {
	fileBytes, err := ioutil.ReadFile("final")
	if err != nil {
		log.Fatalf("Error reading random words file: %s", err)
	}

	appData.RandomWords = strings.Split(string(fileBytes), ":")
	if appData.RandomWords == nil {
		log.Fatal("Error reading random words!")
	}
	/*fmt.Printf("Count before removing duplicates: %d\n", len(appData.RandomWords))
	keys := make(map[string]bool)
	var finalList []string
	for _, word := range appData.RandomWords {
		if _, value := keys[word]; !value {
			keys[word] = true
			finalList = append(finalList, word)
		}
	}
	fmt.Printf("Count after removing duplicates: %d\n", len(finalList))

	file, err := os.OpenFile("final", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	datawriter := bufio.NewWriter(file)

	for _, data := range finalList {
		_, _ = datawriter.WriteString(data + ":")
	}

	datawriter.Flush()
	file.Close()*/
}

func loadLocationsData() {
	/*fileBytes, err := ioutil.ReadFile("locations")
	if err != nil {
		log.Fatalf("Error reading locations file: %s", err)
	}
	appData.Locations = strings.Split(string(fileBytes), "\n")*/
	file, err := os.Open("locations")
	if err != nil {
		//handle error
		return
	}
	defer file.Close()

	s := bufio.NewScanner(file)
	for s.Scan() {
		appData.Locations = append(appData.Locations, s.Text())

		// other code what work with parsed line...
	}
}

func getRandomPersonData() {
	url := fmt.Sprintf("https://randomname.de/?format=json&count=%d&phone=a", appSettings.RowCount)
	res, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error while querying random person data: %s", err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Error while reading response body of query to random person data: %s", err)
	}
	if err = json.Unmarshal(body, &appData.FirstExercise); err != nil {
		log.Fatalf("Error reading valid json back from random person data: %s", err)
	}
}

func addRandomWords() {
	rand.Seed(time.Now().UnixNano())
	for _, entry := range appData.FirstExercise {
		// get five random words for each entry
		for count := 0; count < 5; count++ {
			entry.RandomWords = append(entry.RandomWords, appData.RandomWords[rand.Intn(len(appData.RandomWords))])
		}
		// select one randomly as the second word
		entry.SecondWord = entry.RandomWords[rand.Intn(len(entry.RandomWords))]
	}
}

/*func queryRandomWords() {
	url := fmt.Sprintf("https://www.palabrasaleatorias.com/zufallige-worter.php?fs=%d&fs2=0&Submit=Neues+Wort", RowCount)
	for _, data := range appData.FirstExercise {
		resp, err := http.Get(url)
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

		data.RandomWords = result
		// now select one final randomly
		rand.Seed(time.Now().UnixNano())
		data.SecondWord = result[rand.Intn(len(result))]

	}
}*/
