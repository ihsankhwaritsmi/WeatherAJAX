package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type WeatherData struct {
	Main struct {
		Temperature float64 `json:"temp`
	} `json:"temp`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
}

func main() {
	fmt.Println("hello")
	apiKey := "5df066624a2b37618a4e6f238f45c100"

	h1 := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.Execute(w, nil)

	}

	h2 := func(w http.ResponseWriter, r *http.Request) {
		//check whether an htmx was recieved or not
		log.Print("HTMX recieved")
		log.Print(r.Header.Get("HX-request"))

		//inp == city
		city := r.PostFormValue("inp")
		fmt.Println(city)

		url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s", city, apiKey)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("Error making the request for %s: %v\n", city, err)
			return
		}
		defer resp.Body.Close()

		var data WeatherData
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			fmt.Printf("Error decoding JSON for %s: %v\n", city, err)
			return
		}

		temp := data.Main.Temperature - 273.15
		desc := data.Weather[0].Description

		htmlStr := fmt.Sprintf("<p id='leTarget'>Weather in %s: %.2f°C, %s</p>", city, temp, desc)
		tmpl, _ := template.New("t").Parse(htmlStr)
		tmpl.Execute(w, nil)
	}
	http.HandleFunc("/", h1)
	http.HandleFunc("/show-weath/", h2)

	log.Fatal(http.ListenAndServe(":8000", nil))

}
