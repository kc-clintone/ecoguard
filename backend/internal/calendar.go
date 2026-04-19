package internal

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type CalendarData struct {
	Location    string   `json:"location"`
	CurrentSeason string  `json:"current_season"`
	Month       string   `json:"month"`
	Year        int      `json:"year"`
	Recommendations []string `json:"recommendations"`
	UpcomingTasks   []string `json:"upcoming_tasks"`
	CropCalendar    map[string][]string `json:"crop_calendar"`
}

type CropInfo struct {
	Name         string
	PlantingMonths []time.Month
	HarvestMonths  []time.Month
	Requirements  []string
	Pests         []string
}

func CalendarHandler(w http.ResponseWriter, r *http.Request) {
	location := r.URL.Query().Get("location")
	if location == "" {
		location = "Kenya" // Default location
	}

	calendar := generateCalendarData(location)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(calendar)
}

func generateCalendarData(location string) CalendarData {
	now := time.Now()
	month := now.Month()
	year := now.Year()

	// Determine season based on location and month
	season := getSeason(location, month)

	// Get recommendations based on season and location
	recommendations := getSeasonalRecommendations(season, location)
	upcomingTasks := getUpcomingTasks(season, month)

	// Generate crop calendar for the location
	cropCalendar := generateCropCalendar(location, month)

	return CalendarData{
		Location:    location,
		CurrentSeason: season,
		Month:       month.String(),
		Year:        year,
		Recommendations: recommendations,
		UpcomingTasks:   upcomingTasks,
		CropCalendar:    cropCalendar,
	}
}

func getSeason(location string, month time.Month) string {
	// Enhanced season determination based on real agricultural calendars
	location = strings.ToLower(location)

	// East Africa (Kenya, Tanzania, Uganda, Ethiopia)
	if strings.Contains(location, "kenya") || strings.Contains(location, "tanzania") ||
	   strings.Contains(location, "uganda") || strings.Contains(location, "ethiopia") ||
	   strings.Contains(location, "rwanda") || strings.Contains(location, "burundi") {

		switch month {
		case time.December, time.January, time.February:
			return "Dry Season (Short Rains Ending)"
		case time.March, time.April, time.May:
			return "Long Rains"
		case time.June, time.July, time.August:
			return "Dry Season"
		case time.September, time.October, time.November:
			return "Short Rains"
		}
	}

	// West Africa
	if strings.Contains(location, "nigeria") || strings.Contains(location, "ghana") ||
	   strings.Contains(location, "senegal") || strings.Contains(location, "mali") {

		switch month {
		case time.November, time.December, time.January, time.February, time.March:
			return "Dry Season"
		case time.April, time.May, time.June, time.July, time.August, time.September, time.October:
			return "Wet Season"
		}
	}

	// Southern Africa
	if strings.Contains(location, "south africa") || strings.Contains(location, "zimbabwe") ||
	   strings.Contains(location, "mozambique") || strings.Contains(location, "zambia") {

		switch month {
		case time.April, time.May, time.June, time.July, time.August, time.September:
			return "Dry Winter"
		case time.October, time.November, time.December, time.January, time.February, time.March:
			return "Wet Summer"
		}
	}

	// Default Northern Hemisphere seasons
	switch month {
	case time.December, time.January, time.February:
		return "Winter"
	case time.March, time.April, time.May:
		return "Spring"
	case time.June, time.July, time.August:
		return "Summer"
	case time.September, time.October, time.November:
		return "Fall"
	}

	return "Unknown"
}

func getSeasonalRecommendations(season, location string) []string {
	recommendations := []string{}

	switch strings.ToLower(season) {
	case "long rains", "short rains", "wet season", "wet summer":
		recommendations = append(recommendations,
			"Plant drought-resistant crop varieties",
			"Prepare for potential flooding and waterlogging",
			"Apply organic fertilizers before planting",
			"Monitor soil moisture levels regularly",
			"Use raised beds in flood-prone areas",
			"Implement proper drainage systems",
		)
	case "dry season", "dry winter":
		recommendations = append(recommendations,
			"Implement irrigation systems (drip or sprinkler)",
			"Use mulch to retain soil moisture",
			"Plant drought-tolerant crops like sorghum, millet, or cowpeas",
			"Conserve water resources through efficient irrigation",
			"Test soil for nutrient deficiencies",
			"Prepare land for next rainy season",
		)
	case "spring":
		recommendations = append(recommendations,
			"Start planting warm-season crops",
			"Prepare soil with compost and organic matter",
			"Set up pest monitoring systems",
			"Plan crop rotation to prevent soil depletion",
			"Test soil pH and nutrient levels",
			"Install trellises for climbing crops",
		)
	case "summer":
		recommendations = append(recommendations,
			"Increase irrigation frequency during heat waves",
			"Apply foliar fertilizers for quick nutrient uptake",
			"Monitor for heat stress in crops",
			"Harvest early maturing varieties first",
			"Provide shade for sensitive crops",
			"Control weeds aggressively",
		)
	case "fall":
		recommendations = append(recommendations,
			"Plant cool-season crops like cabbage, carrots, and onions",
			"Clean up garden debris to prevent disease",
			"Test soil for next season planning",
			"Order seeds for spring planting",
			"Amend soil with organic matter",
			"Plan cover crops for winter",
		)
	case "winter":
		recommendations = append(recommendations,
			"Protect plants from frost with row covers",
			"Plan indoor seed starting for early spring",
			"Maintain greenhouse temperatures",
			"Prune dormant fruit trees",
			"Research new crop varieties",
			"Attend agricultural workshops",
		)
	}

	// Location-specific additions
	if strings.Contains(strings.ToLower(location), "kenya") {
		recommendations = append(recommendations,
			"Consider altitude for crop selection (highland vs lowland crops)",
			"Monitor for coffee berry disease in coffee-growing areas",
			"Use integrated pest management for common pests",
			"Participate in government agricultural extension programs",
		)
	}

	if strings.Contains(strings.ToLower(location), "nigeria") {
		recommendations = append(recommendations,
			"Plant early maturing varieties to avoid terminal drought",
			"Use neem-based pesticides for pest control",
			"Implement agroforestry practices",
			"Join farmers' cooperatives for better market access",
		)
	}

	return recommendations
}

func getUpcomingTasks(season string, month time.Month) []string {
	tasks := []string{}

	switch strings.ToLower(season) {
	case "long rains", "short rains", "wet season":
		tasks = append(tasks,
			"Prepare seed beds and till soil",
			"Treat seeds with appropriate fungicides",
			"Apply pre-emergent herbicides",
			"Set up rain gauges for monitoring",
			"Prepare irrigation backup systems",
			"Stock up on fertilizers and pesticides",
		)
	case "dry season":
		tasks = append(tasks,
			"Install or repair irrigation systems",
			"Monitor crops for drought stress symptoms",
			"Apply drought-tolerant fertilizers",
			"Harvest mature crops promptly",
			"Prepare storage facilities",
			"Plan next season's crop rotation",
		)
	case "spring":
		tasks = append(tasks,
			"Start tomato and pepper seedlings indoors",
			"Plant potatoes and onions",
			"Till and amend soil with compost",
			"Install trellises and supports",
			"Set up drip irrigation lines",
			"Order beneficial insects for pest control",
		)
	case "summer":
		tasks = append(tasks,
			"Harvest vegetables frequently to maintain production",
			"Control weeds with mulching and hand weeding",
			"Mulch around plants to conserve moisture",
			"Stake tall plants like tomatoes and peppers",
			"Monitor for and treat common pests",
			"Side-dress plants with additional fertilizer",
		)
	case "fall":
		tasks = append(tasks,
			"Plant garlic, shallots, and multiplier onions",
			"Direct sow carrots, beets, and radishes",
			"Clean up spent plant debris",
			"Test soil pH and nutrient levels",
			"Order seeds for spring planting",
			"Plan winter cover crops",
		)
	case "winter":
		tasks = append(tasks,
			"Plan garden layout for next season",
			"Start seeds indoors for early spring transplanting",
			"Maintain greenhouse and cold frames",
			"Prune fruit trees and berry bushes",
			"Research new varieties and techniques",
			"Attend winter agricultural workshops",
		)
	}

	return tasks
}

func generateCropCalendar(location string, currentMonth time.Month) map[string][]string {
	cropCalendar := make(map[string][]string)

	// Define crops and their planting/harvesting times based on location
	crops := getCropsForLocation(location)

	for _, crop := range crops {
		activities := []string{}

		// Check if current month is planting time
		isPlantingTime := false
		for _, plantMonth := range crop.PlantingMonths {
			if plantMonth == currentMonth {
				isPlantingTime = true
				break
			}
		}

		// Check if current month is harvesting time
		isHarvestingTime := false
		for _, harvestMonth := range crop.HarvestMonths {
			if harvestMonth == currentMonth {
				isHarvestingTime = true
				break
			}
		}

		if isPlantingTime {
			activities = append(activities, "🌱 Plant "+crop.Name)
		}

		if isHarvestingTime {
			activities = append(activities, "🌾 Harvest "+crop.Name)
		}

		// Add general care activities
		if isPlantingTime || isHarvestingTime {
			activities = append(activities, "💧 Water regularly")
			activities = append(activities, "🐛 Monitor for pests: "+strings.Join(crop.Pests, ", "))
			activities = append(activities, "🌿 Requirements: "+strings.Join(crop.Requirements, ", "))
		}

		if len(activities) > 0 {
			cropCalendar[crop.Name] = activities
		}
	}

	return cropCalendar
}

func getCropsForLocation(location string) []CropInfo {
	location = strings.ToLower(location)

	// East African crops (Kenya, Tanzania, Uganda)
	if strings.Contains(location, "kenya") || strings.Contains(location, "tanzania") ||
	   strings.Contains(location, "uganda") || strings.Contains(location, "ethiopia") {

		return []CropInfo{
			{
				Name: "Maize",
				PlantingMonths: []time.Month{time.March, time.April, time.August, time.September, time.October},
				HarvestMonths: []time.Month{time.July, time.August, time.December, time.January},
				Requirements: []string{"Well-drained soil", "Full sun", "Regular watering"},
				Pests: []string{"Maize stalk borer", "Army worm", "Aphids"},
			},
			{
				Name: "Beans",
				PlantingMonths: []time.Month{time.March, time.April, time.August, time.September},
				HarvestMonths: []time.Month{time.June, time.July, time.November, time.December},
				Requirements: []string{"Loamy soil", "Partial shade tolerance", "Nitrogen fixation"},
				Pests: []string{"Bean fly", "Aphids", "Bean bruchid"},
			},
			{
				Name: "Tomatoes",
				PlantingMonths: []time.Month{time.January, time.February, time.June, time.July, time.August},
				HarvestMonths: []time.Month{time.April, time.May, time.October, time.November},
				Requirements: []string{"Well-drained soil", "Full sun", "Support staking"},
				Pests: []string{"Tomato hornworm", "Aphids", "Whiteflies"},
			},
			{
				Name: "Potatoes",
				PlantingMonths: []time.Month{time.March, time.April, time.August, time.September},
				HarvestMonths: []time.Month{time.June, time.July, time.November, time.December},
				Requirements: []string{"Loose, well-drained soil", "Full sun", "Hilling up soil"},
				Pests: []string{"Potato beetle", "Aphids", "Potato cyst nematode"},
			},
		}
	}

	// Default crops for other locations
	return []CropInfo{
		{
			Name: "Tomatoes",
			PlantingMonths: []time.Month{time.March, time.April, time.May},
			HarvestMonths: []time.Month{time.June, time.July, time.August, time.September},
			Requirements: []string{"Full sun", "Well-drained soil", "Regular watering"},
			Pests: []string{"Hornworms", "Aphids", "Blossom end rot"},
		},
		{
			Name: "Lettuce",
			PlantingMonths: []time.Month{time.March, time.April, time.August, time.September},
			HarvestMonths: []time.Month{time.May, time.June, time.October, time.November},
			Requirements: []string{"Partial shade", "Rich soil", "Consistent moisture"},
			Pests: []string{"Slugs", "Aphids", "Cabbage loopers"},
		},
		{
			Name: "Carrots",
			PlantingMonths: []time.Month{time.March, time.April, time.August, time.September},
			HarvestMonths: []time.Month{time.June, time.July, time.November, time.December},
			Requirements: []string{"Loose soil", "Full sun", "Regular thinning"},
			Pests: []string{"Carrot fly", "Aphids", "Wireworms"},
		},
	}
}