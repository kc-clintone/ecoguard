# 🌱 EcoGuard

EcoGuard is a smart farming web application designed to help farmers make better decisions through climate awareness, crop monitoring, and guided farming practices.

---

## 🚀 Features

### 👤 Authentication
- User Signup and Login
- Persistent user storage
- Multiple users supported

### 🌦 Climate Alerts
- Displays weather and climate-related updates
- Designed to integrate with real-time APIs (e.g., OpenWeatherMap)

### 🌾 Farming Calendar
- Provides guidance on planting, irrigation, and harvesting
- Can be extended with seasonal and crop-specific data

### 🤖 AI Crop Detector (Mock)
- Upload crop images
- Preview functionality (ready for AI integration)

### 👥 Guest Mode
- Users can explore the app without signing up
- Restricted access to some features (e.g., detector, calendar)

---

## 🏗 Project Structure
ecoguard/
├── backend/
│ ├── cmd/
│ │ └── main.go # Entry point
│ └── internal/
│ ├── handlers.go # API handlers (login, signup, update)
│ ├── storage.go # File storage logic
│ └── models.go # Data models
│
└── frontend/
├── index.html # Main UI
├── styles.css # Styling
└── script.js # Frontend logic


---

## ⚙️ Installation & Setup

### 1. Clone the repository
```bash
git clone https://github.com/YOUR_USERNAME/ecoguard.git
cd ecoguard

