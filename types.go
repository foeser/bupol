package main

type IndexOutOfRange struct{}

func (m *IndexOutOfRange) Error() string {
	return "Index out of range."
}

//HTTPResponseObject Data structure for sending the API call status
type HTTPResponseObject struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	ErrorObject error  `json:"error"`
}

type AppSettings struct {
	RowCount      int `json:"rowCount"`
	TimeInSeconds int `json:"timeInSeconds"`
}

type AppConfig struct {
	Port int    `json:"port" env:"BUPOLPORT" env-default:"8888"`
	VDir string `json:"vdir" env:"BUPOLVDIR" env-default:"/bupol"`
}

type FirstExercise struct {
	FirstWord   string `json:"username"`
	SecondWord  string
	RandomWords []string
	Editable    string `json:"editable"`
}

type SecondExercise struct {
	Location string
	Phone    string
	Editable string `json:"editable"`
}

type AppData struct {
	FirstExercise  []*FirstExercise
	SecondExercise []*SecondExercise
	Locations      []string
	RandomWords    []string
}

type PageData struct {
	VDir            string
	TimerTime       int
	Exercise        string
	GridfieldName1  string
	GridfieldTitle1 string
	GridfieldName2  string
	GridfieldTitle2 string
}
