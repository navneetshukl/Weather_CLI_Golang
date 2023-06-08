package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
)

type Weather struct {
	Location struct {
		Name    string `json:"name"`
		Region  string `json:"Region"`
		Country string `json:"country"`
	} `json:"location"`

	Current struct {
		TempC float64 `json:"temp_c"`

		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`

	Forecast struct {
		Forecastday []struct {
			Hour []struct {
				TimeEpoch int64   `json:"time_epoch"`
				TempC     float64 `json:"temp_c"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`

				ChanceOfRain float64 `json:"chance_of_rain"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func main() {

	q:="Jhansi"
	if len(os.Args)>=2{
		q=os.Args[1]
	}
	resp, err := http.Get("http://api.weatherapi.com/v1/forecast.json?key=e88177b20c024f7d86c120114221006&q="+ q +"&days=1&aqi=no&alerts=no")

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		panic("Weather Api is not Available")
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	var weather Weather

	err = json.Unmarshal(body, &weather)
	if err != nil {
		panic(err)
	}
	//fmt.Println(weather)

	location, current, hours := weather.Location, weather.Current, weather.Forecast.Forecastday[0].Hour

	fmt.Printf(
		"%s,%s, %s: %.0fC, %s\n",
		location.Name,
		location.Region,
		location.Country,
		current.TempC,
		current.Condition.Text,
	)

	for _, hour := range hours {
		date := time.Unix(hour.TimeEpoch, 0)

		if date.Before(time.Now()) {
			continue
		}

		message := fmt.Sprintf(
			"%s - %.0fC, %.0f%%,%s\n",
			date.Format("15:04"),
			hour.TempC,
			hour.ChanceOfRain,
			hour.Condition.Text,
		)

		if hour.ChanceOfRain < 40 {
			fmt.Print(message)
		} else {
			color.Red(message)
		}
	}
}
