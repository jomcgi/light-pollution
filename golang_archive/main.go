package main

import (
  "astronomy-weather/fetchForecast"
  "fmt"
  "encoding/json"
  "log"
  "time"
)

func prettyPrint(i interface{}) string {
    s, _ := json.MarshalIndent(i, "", "  ")
    return string(s)
}

func main() {
  forecast, err:= fetchForecast.ForecastApiCall(55.885981, -4.279137)
  if err != nil {
    log.Fatal(err)
  }
  numberOfDaystoForecast:= 3
  numberOfHoursToForecast:= time.Duration(24*numberOfDaystoForecast)
  maxForecastTime:= time.Now().Add(time.Hour*numberOfHoursToForecast)
  filteredForecast, err := fetchForecast.FilterForecast(forecast, 5, maxForecastTime)
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println(prettyPrint(filteredForecast))
}