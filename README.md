# 🌱 EcoGuard

EcoGuard is a production-ready smart farming web application designed to help farmers make better decisions through climate awareness, crop monitoring, and guided farming practices.

---

## 🚀 Features

### 👤 Authentication

- User Signup and Login with password validation
- Persistent user storage with SQLite
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

---

## 📋 Prerequisites

- Go 1.24.3 or later
- OpenWeather API key (free at https://openweathermap.org/api)
- Plant.id API key (optional, for enhanced crop detection)

## 🛠 Local Development

### Setup

```bash
# Clone repository
git clone <repository-url>
cd ecoguard

# Install dependencies
go mod download

# Create environment file
cp .env.example .env

# Edit .env with your API keys
# OPENWEATHER_API_KEY=your_key_here
# PLANT_ID_API_KEY=your_key_here
```

### Running Locally

```bash
# From project root
cd backend/cmd
go run main.go
```

Server runs on `http://localhost:8080`

## 🚀 Production Deployment

See [DEPLOYMENT.md](DEPLOYMENT.md) for comprehensive deployment guide including:

- Environment configuration
- Building optimized binaries
- Running with systemd
- Docker containerization
- Nginx reverse proxy setup
- Security best practices
- Troubleshooting

Quick start:

```bash
# Build production binary
go build -ldflags="-w -s" -o ecoguard ./backend/cmd/

# Run with environment variables
DB_PATH=/var/lib/ecoguard/ecoguard.db \
PORT=8080 \
OPENWEATHER_API_KEY=your_key \
./ecoguard
```

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

The application will start the server and automatically create the SQLite database (`ecoguard.db`) on first run. The backend serves both the API endpoints and the static frontend files.

## 📡 API Endpoints

### Authentication

#### POST `/signup`

Create a new user account.

**Request:**

```json
{
  "username": "farmer",
  "password": "secure_password"
}
```

**Response (200):**

```json
{
  "message": "Account created, please login"
}
```

#### POST `/login`

Authenticate and login user.

**Request:**

```json
{
  "username": "farmer",
  "password": "secure_password"
}
```

**Response (200):**

```json
{
  "message": "Login successful",
  "username": "farmer"
}
```

#### POST `/update`

Update user profile.

**Request:**

```json
{
  "username": "farmer",
  "password": "new_password"
}
```

### Weather

#### GET `/weather`

Get weather data and 7-day forecast.

**Query Parameters:**

- `latitude` (required): Latitude coordinate
- `longitude` (required): Longitude coordinate

**Example:** `/weather?latitude=0.35&longitude=35.83`

**Response (200):**

```json
{
  "current": {
    "temperature": 25.5,
    "condition": "Sunny",
    "humidity": 65
  },
  "forecast": [
    {
      "date": "2024-01-15",
      "high": 28,
      "low": 18,
      "condition": "Sunny"
    }
  ]
}
```

### Crop Detection

#### POST `/detect`

Analyze crop health and detect diseases.

**Request (multipart/form-data):**

- `image`: Image file (JPEG/PNG)

**Response (200):**

```json
{
  "plant_name": "Zea mays",
  "common_name": "Maize",
  "is_healthy": {
    "probability": 0.92,
    "confidence": 0.95
  },
  "diseases_detected": [],
  "care_tips": ["Water regularly", "Monitor for pests"]
}
```

### Farming Calendar

#### GET `/calendar`

Get seasonal farming guidance.

**Query Parameters:**

- `region`: east_africa, west_africa, or southern_africa
- `crop`: crop name (maize, beans, rice, sorghum, etc.)

**Example:** `/calendar?region=east_africa&crop=maize`

**Response (200):**

```json
{
  "crop": "Maize",
  "region": "East Africa",
  "planting_months": [3, 4, 5],
  "harvest_months": [8, 9, 10],
  "soil_requirements": "Well-drained fertile soil",
  "pests": ["Fall armyworm", "Stem borers"],
  "diseases": ["Leaf rust", "Maize streak virus"],
  "seasonal_tasks": ["Prepare land", "Plant seeds", "Weed"]
}
```

## 🔐 Security

- CORS headers configured for production
- Database connection pooling for performance
- Environment-based configuration for API keys
- Input validation on all endpoints
- SQLite with parameterized queries prevents SQL injection

## 🐛 Troubleshooting

**Issue: Database locked error**

```
Solution: Ensure database file is on local storage, not NFS
```

**Issue: Weather API returns 401**

```
Solution: Check OPENWEATHER_API_KEY is set and valid
Get free key at: https://openweathermap.org/api
```

**Issue: Crop detector returns mock data**

```
Solution: Set PLANT_ID_API_KEY environment variable
Get free key at: https://plant.id/api
```

**Issue: Frontend shows blank page**

```
Solution: Clear browser cache (Ctrl+Shift+Delete)
Check browser console for errors (F12)
```

## 💡 Tips for Production

1. **Use absolute database paths** - Avoid relative paths
2. **Enable HTTPS** - Use reverse proxy with SSL certificates
3. **Monitor logs** - Set up log aggregation
4. **Regular backups** - Backup SQLite database daily
5. **Rate limiting** - Implement rate limiting for API endpoints
6. **Monitoring** - Track API response times and errors

## 📝 License

Open source - feel free to use and modify

## 🤝 Contributing

Contributions welcome! Areas to enhance:

- Add PostgreSQL support for scaling
- Implement JWT authentication
- Add multi-language support
- Expand crop database
- Mobile app version

---

Made with 🌱 for farmers
