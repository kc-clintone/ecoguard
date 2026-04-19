package internal

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type PlantIDRequest struct {
	Images        []string `json:"images"`
	Latitude      float64  `json:"latitude,omitempty"`
	Longitude     float64  `json:"longitude,omitempty"`
	SimilarImages bool     `json:"similar_images,omitempty"`
	Symptoms      bool     `json:"symptoms,omitempty"`
}

type PlantIDResponse struct {
	Result struct {
		IsHealthy struct {
			Binary      bool    `json:"binary"`
			Probability float64 `json:"probability"`
		} `json:"is_healthy"`

		Disease struct {
			Suggestions []struct {
				Name        string  `json:"name"`
				Probability float64 `json:"probability"`
			} `json:"suggestions"`
		} `json:"disease"`
	} `json:"result"`
}

type CropAnalysis struct {
	HealthScore int      `json:"health_score"`
	Status      string   `json:"status"`
	Suggestions []string `json:"suggestions"`
	CropType    string   `json:"crop_type"`
	Confidence  float64  `json:"confidence"`
	Diseases    []string `json:"diseases,omitempty"`
}

func CropDetectorHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "No image file provided", http.StatusBadRequest)
		return
	}
	defer file.Close()

	imageData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read image", http.StatusInternalServerError)
		return
	}

	analysis, err := analyzeCropWithPlantID(imageData)
	if err != nil {
		fmt.Println("Plant.id error:", err)
		analysis = analyzeCropMock()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(analysis)
}

func analyzeCropWithPlantID(imageData []byte) (CropAnalysis, error) {
	apiKey := os.Getenv("PLANT_ID_API_KEY")
	if apiKey == "" {
		return CropAnalysis{}, fmt.Errorf("PLANT_ID_API_KEY not set")
	}

	base64Image := "data:image/jpeg;base64," +
		base64.StdEncoding.EncodeToString(imageData)

	requestBody := PlantIDRequest{
		Images:        []string{base64Image},
		Latitude:      -1.286389,
		Longitude:     36.817223,
		SimilarImages: true,
		Symptoms:      true,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return CropAnalysis{}, err
	}

	req, err := http.NewRequest(
		"POST",
		"https://plant.id/api/v3/health_assessment",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return CropAnalysis{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return CropAnalysis{}, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// Accept 200–299 as success
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return CropAnalysis{}, fmt.Errorf(
			"Plant.id error %d: %s",
			resp.StatusCode,
			string(body),
		)
	}

	var plantResponse PlantIDResponse
	err = json.Unmarshal(body, &plantResponse)
	if err != nil {
		return CropAnalysis{}, err
	}

	return processPlantIDResponse(plantResponse), nil
}

func processPlantIDResponse(response PlantIDResponse) CropAnalysis {
	healthScore := 90
	status := "Healthy"
	diseases := []string{}
	suggestions := []string{}

	if !response.Result.IsHealthy.Binary {
		healthScore = 60
		status = "Unhealthy"
	}

	for _, disease := range response.Result.Disease.Suggestions {
		if disease.Probability > 0.15 {
			diseases = append(diseases, disease.Name)

			healthScore -= int(disease.Probability * 40)

			suggestions = append(
				suggestions,
				fmt.Sprintf(
					"%s detected (%.0f%% probability)",
					disease.Name,
					disease.Probability*100,
				),
			)
		}
	}

	if healthScore < 0 {
		healthScore = 0
	}

	if healthScore >= 80 {
		status = "Healthy"
	} else if healthScore >= 50 {
		status = "Fair"
	} else {
		status = "Poor"
	}

	if len(suggestions) == 0 {
		suggestions = []string{
			"No visible disease detected",
			"Continue regular watering",
			"Monitor leaves for changes",
		}
	}

	return CropAnalysis{
		HealthScore: healthScore,
		Status:      status,
		Suggestions: suggestions,
		CropType:    "Detected Crop",
		Confidence:  response.Result.IsHealthy.Probability * 100,
		Diseases:    diseases,
	}
}

func analyzeCropMock() CropAnalysis {
	return CropAnalysis{
		HealthScore: 75,
		Status:      "Fallback",
		Suggestions: []string{
			"Plant.id API unavailable",
			"Check API key",
			"Upload clearer image",
		},
		CropType:   "Unknown",
		Confidence: 0,
	}
}