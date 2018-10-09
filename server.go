package main

import (
	"fmt"
	//"html/template"
	"net/http"
	"regexp"
	"time"
	igc "github.com/marni/goigc"
	"math/rand"
	"encoding/json"
    "strconv"
	//"path/filepath"
	"strings"
)


var timeStarted = time.Now()

type _url struct {
	URL string `json:"url"`
}
var igcFiles []Track

func trackLength(track igc.Track) float64 {

	totalDistance := 0.0

	for i := 0; i < len(track.Points)-1; i++ {
		totalDistance += track.Points[i].Distance(track.Points[i+1])
	}

	return totalDistance
}


type Track struct {
	Id string   `json:"id"`
	igcTrack igc.Track `json:"igc_track"`
}
type Attributes struct{
	HeaderDate string `json:"h_date"`
	Pilot string `json:"pilot"`
	Glider string `json:"glider"`
	Gl_id string 	`json :"glider_id"`
	Length float64 `json:"track_length"`
}

func handler(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, "{" + "\"uptime\": \""+FormatSince(timeStarted)+"\"," + "\"info\": \"Service for IGC tracks.\"," + "\"version\": \"v1\""+ "}")
}

func handler2(w http.ResponseWriter, r *http.Request){
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		ids := make([]string, 0, 0)

		for i := range igcFiles {
			ids = append(ids, igcFiles[i].Id)
		}

		json.NewEncoder(w).Encode(ids)

		break
	case http.MethodPost:

		pattern:=".*.igc"

		w.Header().Set("Content-Type", "application/json")
		//jsonR := make(map[string]string)
		URL := &_url{}

		var error = json.NewDecoder(r.Body).Decode(URL)
		if error != nil {
			fmt.Fprintln(w, "Error!! ", error)
			return
		}
		res,err:=regexp.MatchString(pattern,URL.URL)
		if err!=nil{
			http.Error(w,err.Error(),http.StatusInternalServerError)
			fmt.Fprintln(w, "Error!! ", error)
			return
		}
		if res {

			track, _ := igc.ParseLocation(URL.URL)

			Id := rand.Intn(1000)

			igcFile := Track{}
			igcFile.Id = strconv.Itoa(Id)
			igcFile.igcTrack = track

			igcFiles = append(igcFiles, igcFile)

			json.NewEncoder(w).Encode(igcFile.Id)
			return
		}
		break
	default:
		http.Error(w,"Not implemented",http.StatusNotImplemented)


	}


}


func handler3(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	parts := strings.Split(r.URL.Path, "/")

	//vals := r.URL.Query() // Returns a url.Values, which is a map[string][]string

	//productTypes, ok := vals["id"]

	attributes := &Attributes{}

	var rNum= regexp.MustCompile(`/igcinfo/api/igc/\d{1,}`)
	switch {
	case rNum.MatchString(r.URL.Path):

		for i := range igcFiles {

			if igcFiles[i].Id == parts[4] {
				attributes.HeaderDate = igcFiles[i].igcTrack.Header.Date.String()
				attributes.Pilot = igcFiles[i].igcTrack.Pilot
				attributes.Glider = igcFiles[i].igcTrack.GliderType
				attributes.Gl_id= igcFiles[i].igcTrack.GliderID
				attributes.Length = trackLength(igcFiles[i].igcTrack)

				json.NewEncoder(w).Encode(attributes)
			}

		}

	break
	default:
		fmt.Fprintln(w, "Error: something goes wrong!!")

	}

}

//Kete funksionin tjeter me shti ne handler3 edhe me thirr permes kushtit te len(parts)==4


func handler4(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	parts := strings.Split(r.URL.Path, "/")
	//attributes := &Attributes{}

	var rNum= regexp.MustCompile(`/igcinfo/api/igc/\d{1,}/\w{1,}`)
	switch {
	case rNum.MatchString(r.URL.Path):

		for i := range igcFiles {

			if igcFiles[i].Id == parts[4] {
				switch{
				case parts[5]=="pilot":
					json.NewEncoder(w).Encode(igcFiles[i].igcTrack.Pilot)
					break
				case parts[5]=="glider":
					json.NewEncoder(w).Encode(igcFiles[i].igcTrack.GliderType)
					break
				case parts[5]=="glider_id":
					json.NewEncoder(w).Encode(igcFiles[i].igcTrack.GliderID)
					break
				case parts[5]=="track_length":
					json.NewEncoder(w).Encode(trackLength(igcFiles[i].igcTrack))
					break
				case parts[5]=="h_date":
					json.NewEncoder(w).Encode(igcFiles[i].igcTrack.Header.Date.String())
					break
				default:
					http.Error(w, "400 - Bad Request, the field you entered is not on our database!", http.StatusBadRequest)
					break
				}

			}

		}

		break
	default:
		fmt.Fprintln(w, "Error: something goes wrong!!")

	}


}
func main() {

	http.HandleFunc("/igcinfo/api",handler)
	http.HandleFunc("/igcinfo/api/igc",handler2)
	http.HandleFunc("/",handler3)
	http.ListenAndServe(":8080",nil)
}
func FormatSince(t time.Time) string {
	const (
		Decisecond = 100 * time.Millisecond
		Day        = 24 * time.Hour
	)
	ts := time.Since(t)
	sign := time.Duration(1)
	if ts < 0 {
		sign = -1
		ts = -ts
	}
	ts += +Decisecond / 2
	d := sign * (ts / Day)
	ts = ts % Day
	h := ts / time.Hour
	ts = ts % time.Hour
	m := ts / time.Minute
	ts = ts % time.Minute
	s := ts / time.Second
	ts = ts % time.Second
	f := ts / Decisecond
	y := d / 365
	return fmt.Sprintf("P%dY%dD%dH%dM%d.%dS", y, d, h, m, s, f)
}
