package fetchForecast

import (
	"fmt"
	"io"
	"log"
	"net/http"
  "encoding/json"
  "time"
)

type ForecastTimeSeries struct {
  Time time.Time `json:"time"`
  Data struct {
    Instant struct {
      Details struct {
        AirTemperature           float64 `json:"air_temperature"`
        CloudAreaFraction        float64 `json:"cloud_area_fraction"`
        FogAreaFraction          float64 `json:"fog_area_fraction"`
        WindFromDirection        float64 `json:"wind_from_direction"`
        WindSpeed                float64 `json:"wind_speed"`
      } `json:"details"`
    } `json:"instant"`
  } `json:"data"`
}

type ForecastMetadata struct {
  UpdatedAt time.Time `json:"updated_at"`
  Units     struct {
    AirTemperature           string `json:"air_temperature"`
    CloudAreaFraction        string `json:"cloud_area_fraction"`
    FogAreaFraction          string `json:"fog_area_fraction"`
    WindFromDirection        string `json:"wind_from_direction"`
    WindSpeed                string `json:"wind_speed"`
  } `json:"units"`
}

type ForecastProperties struct {
  Meta ForecastMetadata `json:"meta"`
  Timeseries []ForecastTimeSeries `json:"timeseries"`
}


type ForecastApiResponse struct {
	Type     string `json:"type"`
	Geometry struct {
		Type        string `json:"type"`
		Coordinates []float64  `json:"coordinates"`
	} `json:"geometry"`
	Properties ForecastProperties `json:"properties"`
}

func ForecastApiCall(lat float64, lon float64) (ForecastApiResponse, error){
	apiUrl := "https://api.met.no/weatherapi/locationforecast/2.0/complete?lat=" + fmt.Sprintf("%.4f", lat) + "&lon=" + fmt.Sprintf("%.4f", lon)
	client := http.Client{}
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		log.Fatal(err)
		log.Fatal("Can not create request")
    return ForecastApiResponse{}, err
	}
	req.Header = http.Header{
		"Accept":     {"application/json"},
		"User-Agent": {"Astronomy-Weather"},
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		log.Fatal("Can not get response")
    return ForecastApiResponse{}, err
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
    log.Fatal("Can not read response body")
    return ForecastApiResponse{}, err
	}
  var result ForecastApiResponse
  if err := json.Unmarshal(bodyBytes, &result); err != nil {
    log.Fatal(err)
    fmt.Println("Can not unmarshal JSON")
    return ForecastApiResponse{}, err
  }
  return result, nil
}

func FilterForecast(forecast ForecastApiResponse, maxCloudCover float64, maxTime time.Time) (ForecastApiResponse, error) {
  var filteredForecastTimeSeries []ForecastTimeSeries
  for _, timeseries := range forecast.Properties.Timeseries {
    if timeseries.Data.Instant.Details.CloudAreaFraction < maxCloudCover && timeseries.Time.Before(maxTime){
      filteredForecastTimeSeries = append(filteredForecastTimeSeries, timeseries)
    }
  }
  filteredForecast := ForecastApiResponse{
    Type: forecast.Type,
    Geometry: forecast.Geometry,
    Properties: ForecastProperties{
      Meta: forecast.Properties.Meta,
      Timeseries: filteredForecastTimeSeries,
    },
  }
  return filteredForecast, nil
}
