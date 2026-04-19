# 🌱 EcoGuard

EcoGuard is a smart farming web application designed to help farmers make better decisions through climate awareness, crop monitoring, and guided farming practices.

---

## 🚀 Features

### 👤 Authentication

- User Signup and Login
- Persistent user storage
- Multiple users supported

### 🌦 Climate Alerts

- Real-time weather data from OpenWeatherMap API
- 7-day weather forecast with current conditions
- Weather alerts and warnings when available
- Location-based weather information

### 🌾 Farming Calendar

- Comprehensive seasonal farming guidance based on real agricultural data
- Location-specific crop calendars (East Africa, West Africa, Southern Africa)
- Planting and harvesting schedules for major crops
- Seasonal recommendations and upcoming farming tasks
- Pest monitoring and crop care requirements

### 🤖 AI Crop Detector

- Real plant identification using Plant.id API
- Disease and pest detection from crop images
- Health score assessment with confidence ratings
- Crop-specific care recommendations and pest management advice
- Supports major agricultural crops with detailed analysis

### 📅 Farming Calendar

- View seasonal farming information based on location
- Get planting recommendations and upcoming tasks
- Supports different regions and climates

---

## 🏗 Project Structure

```bash
ecoguard/
├── backend/
│ ├── cmd/
│ │ └── main.go # Entry point
│ └── internal/
│ ├── db.go # Database initialization and operations
│ ├── handlers.go # API handlers (login, signup, update)
│ └── models.go # Data models
│
└── frontend/
├── index.html # Main UI
├── styles.css # Styling
└── script.js # Frontend logic
```

---

## ⚙️ Installation & Setup

### Prerequisites

- Go 1.24.3 or later
- SQLite (automatically handled by the application)

### Setup Steps

1. Clone or download the repository:

   ```bash
   git clone https://github.com/Ashomondi/ecoguard.git
   cd ecoguard
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Run the application:

   ```bash
   cd backend/cmd
   export OPENWEATHER_API_KEY="your_openweather_api_key_here"
   export PLANT_ID_API_KEY="your_plant_id_api_key_here"
   go run main.go
   ```

   **API Keys Required:**
   - **OpenWeatherMap API Key**:
     1. Go to https://openweathermap.org/api
     2. Sign up for a free account
     3. Get your API key from the dashboard
   - **Plant.id API Key**:
     1. Go to https://plant.id/
     2. Sign up for a free account
     3. Get your API key (they offer 100 free identifications per month)

   Without these API keys, the app will still work but with fallback/mock data for crop detection.

4. Open your browser and navigate to `http://localhost:8080`

The application will start the server and automatically create the SQLite database (`ecoguard.db`) in the `backend/cmd` directory on first run. The backend serves both the API endpoints and the static frontend files.
