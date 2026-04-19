package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type geoLocation struct {
    Name    string  `json:"name"`
    Lat     float64 `json:"lat"`
    Lon     float64 `json:"lon"`
    Country string  `json:"country"`
    State   string  `json:"state"`
}

type WeatherResponse struct {
    Current struct {
        Temp      float64 `json:"temp"`
        Humidity  int     `json:"humidity"`
        WindSpeed float64 `json:"wind_speed"`
        Weather   []struct {
            Description string `json:"description"`
            Main        string `json:"main"`
        } `json:"weather"`
    } `json:"current"`
    Daily []struct {
        Dt   int64   `json:"dt"`
        Temp struct {
            Min float64 `json:"min"`
            Max float64 `json:"max"`
        } `json:"temp"`
        Weather []struct {
            Main string `json:"main"`
        } `json:"weather"`
        Pop     float64 `json:"pop"` // Probability of precipitation
    } `json:"daily"`
    Alerts []struct {
        Event     string `json:"event"`
        Start     int64  `json:"start"`
        End       int64  `json:"end"`
        Description string `json:"description"`
    } `json:"alerts"`
}

func WeatherHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

    apiKey := os.Getenv("OPENWEATHER_API_KEY")
    if apiKey == "" {
        http.Error(w, "OpenWeather API key not configured. Please set OPENWEATHER_API_KEY environment variable with your OpenWeatherMap API key. Get one at https://openweathermap.org/api", http.StatusInternalServerError)
        return
    }

    locationQuery := r.URL.Query().Get("q")
    if locationQuery == "" {
        locationQuery = "Nairobi"
    }

    geo, err := fetchGeolocation(locationQuery, apiKey)
    if err != nil {
        // Provide more helpful error message
        if strings.Contains(err.Error(), "401") {
            http.Error(w, "Invalid OpenWeather API key. Please check your OPENWEATHER_API_KEY environment variable. Get a free API key at https://openweathermap.org/api", http.StatusUnauthorized)
            return
        }
        http.Error(w, fmt.Sprintf("Failed to resolve location: %v", err), http.StatusBadRequest)
        return
    }

	weatherData, err := fetchOpenWeather(geo.Lat, geo.Lon, apiKey)
    if err != nil {
        if strings.Contains(err.Error(), "401") {
            http.Error(w, "Invalid OpenWeather API key. Please check your OPENWEATHER_API_KEY environment variable. Get a free API key at https://openweathermap.org/api", http.StatusUnauthorized)
            return
        }
        http.Error(w, fmt.Sprintf("Failed to fetch weather: %v", err), http.StatusInternalServerError)
        return
    }

    response := map[string]any{
        "location": fmt.Sprintf("%s, %s", geo.Name, geo.Country),
        "weather":  weatherData,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func fetchGeolocation(query, apiKey string) (geoLocation, error) {
    encoded := url.QueryEscape(query)
    url := fmt.Sprintf("https://api.openweathermap.org/geo/1.0/direct?q=%s&limit=1&appid=%s", encoded, apiKey)

    resp, err := http.Get(url)
    if err != nil {
        return geoLocation{}, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return geoLocation{}, fmt.Errorf("geocoding returned status %d", resp.StatusCode)
    }

    var locations []geoLocation
    if err := json.NewDecoder(resp.Body).Decode(&locations); err != nil {
        return geoLocation{}, err
    }
    if len(locations) == 0 {
        return geoLocation{}, fmt.Errorf("location not found")
    }
    return locations[0], nil
}

func fetchOpenWeather(lat, lon float64, apiKey string) (map[string]any, error) {
	url := fmt.Sprintf(
		"https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&units=metric&appid=%s",
		lat, lon, apiKey,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weather API error %d: %s", resp.StatusCode, string(body))
	}

	var data map[string]any
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return data, nil
}