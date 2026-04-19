// Dynamic API URL - works for both local dev and production
const API_URL = window.location.origin;

let isGuest = true;

// Show a specific section and hide others
function showSection(id) {
  const restricted = ["detector", "calendar"];
  if (isGuest && restricted.includes(id)) {
    alert("This feature is only for registered users.");
    return;
  }

  document
    .querySelectorAll(".container")
    .forEach((c) => c.classList.add("hidden"));
  document.getElementById(id).classList.remove("hidden");

  if (id === "alerts") fetchClimateAlerts();
  if (id === "calendar") fetchCalendar();
}

// Go home from inner pages
function goHome() {
  showSection("home");
}

// Guest button
function continueAsGuest() {
  isGuest = true;
  showSection("home");
}

// Logout button
function logout() {
  isGuest = true;
  showSection("landing");
}

// Signup
function signup() {
  const username = document.getElementById("signupUser").value;
  const password = document.getElementById("signupPass").value;
  if (!username || !password) {
    alert("Enter username and password");
    return;
  }

  fetch(`${API_URL}/signup`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username, password }),
  })
    .then((r) => r.json())
    .then((d) => {
      alert(d.message || "Signed up successfully");
      showSection("login");
    })
    .catch((e) => alert("Error signing up: " + e.message));
}

// Login
function login() {
  const username = document.getElementById("loginUser").value;
  const password = document.getElementById("loginPass").value;
  if (!username || !password) {
    alert("Enter username and password");
    return;
  }

  fetch(`${API_URL}/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username, password }),
  })
    .then((r) => {
      if (!r.ok) throw new Error("Invalid username or password");
      return r.json();
    })
    .then((d) => {
      isGuest = false;
      alert(`Welcome ${d.username}`);
      showSection("home");
    })
    .catch((e) => alert(e.message));
}

// Weather alerts using backend OpenWeather proxy
function fetchClimateAlerts() {
  const location =
    document.getElementById("weather-location").value.trim() || "Nairobi";
  const content = document.getElementById("alerts-content");
  content.innerText = "Loading weather for " + location + "...";

  fetch(`${API_URL}/weather?q=${encodeURIComponent(location)}`)
    .then((r) => {
      if (!r.ok) {
        return r.text().then((text) => {
          throw new Error(text || "Unable to fetch weather");
        });
      }
      return r.json();
    })
    .then((data) => {
      if (!data.weather) {
        throw new Error("Malformed weather response");
      }
      displayWeather(data);
    })
    .catch((error) => {
      content.innerText = `Error: ${error.message}`;
    });
}

function displayWeather(data) {
  const content = document.getElementById("alerts-content");

  if (!data || !data.weather) {
    content.innerHTML = "No weather data received";
    return;
  }

  const weather = data.weather;
  const location = data.location || "Unknown location";

  const condition = weather.weather?.[0]?.description || "Unknown";

  const temp = weather.main?.temp;
  const humidity = weather.main?.humidity;
  const wind = weather.wind?.speed;

  let html = `<strong>${location}</strong><br>`;
  html += `Now: ${temp ?? "N/A"}°C, ${condition}<br>`;
  html += `Humidity: ${humidity ?? "N/A"}% • Wind: ${wind ?? "N/A"} m/s<br>`;

  content.innerHTML = html;
}

function fetchCalendar() {
  const location =
    document.getElementById("calendar-location").value.trim() || "Kenya";
  const content = document.getElementById("calendar-content");
  content.innerText = "Loading calendar for " + location + "...";

  fetch(`${API_URL}/calendar?location=${encodeURIComponent(location)}`)
    .then((r) => {
      if (!r.ok) {
        return r.text().then((text) => {
          throw new Error(text || "Unable to fetch calendar");
        });
      }
      return r.json();
    })
    .then((data) => {
      displayCalendar(data);
    })
    .catch((error) => {
      content.innerText = `Error: ${error.message}`;
    });
}

function displayCalendar(data) {
  const content = document.getElementById("calendar-content");

  let html = `<strong>${data.location} - ${data.current_season}</strong><br>`;
  html += `<em>${data.month} ${data.year}</em><br><br>`;

  html += `<strong>Seasonal Recommendations:</strong><br>`;
  data.recommendations.forEach((rec) => {
    html += `• ${rec}<br>`;
  });
  html += `<br>`;

  html += `<strong>Upcoming Tasks:</strong><br>`;
  data.upcoming_tasks.forEach((task) => {
    html += `• ${task}<br>`;
  });

  // Add crop calendar if available
  if (data.crop_calendar && Object.keys(data.crop_calendar).length > 0) {
    html += `<br><strong>Crop Calendar for ${data.month}:</strong><br>`;
    Object.entries(data.crop_calendar).forEach(([crop, activities]) => {
      html += `<strong>${crop}:</strong><br>`;
      activities.forEach((activity) => {
        html += `  ${activity}<br>`;
      });
      html += `<br>`;
    });
  }

  content.innerHTML = html;
}

// Optional: preview crop image
function previewImage(event) {
  const preview = document.getElementById("preview");
  const result = document.getElementById("result");

  if (event.target.files && event.target.files[0]) {
    preview.src = URL.createObjectURL(event.target.files[0]);
    preview.classList.remove("hidden");
    result.innerText =
      "Image ready for analysis. Click 'Analyze Crop' to get AI recommendations.";
  }
}

function analyzeCrop() {
  const form = document.getElementById("upload-form");
  const result = document.getElementById("result");
  const formData = new FormData(form);

  if (!formData.get("image")) {
    result.innerText = "Please select an image first.";
    return;
  }

  result.innerText = "Analyzing crop image with AI...";

  fetch(`${API_URL}/detect`, {
    method: "POST",
    body: formData,
  })
    .then((r) => {
      if (!r.ok) {
        return r.text().then((text) => {
          throw new Error(text || "Analysis failed");
        });
      }
      return r.json();
    })
    .then((data) => {
      displayCropAnalysis(data);
    })
    .catch((error) => {
      result.innerText = `Error: ${error.message}`;
    });
}

function displayCropAnalysis(data) {
  const result = document.getElementById("result");

  let html = `<strong>Crop Type:</strong> ${data.crop_type}<br>`;
  if (data.confidence) {
    html += `<strong>Confidence:</strong> ${data.confidence.toFixed(1)}%<br>`;
  }
  html += `<strong>Health Score:</strong> ${data.health_score}/100 (${data.status})<br><br>`;

  if (data.diseases && data.diseases.length > 0) {
    html += `<strong>Detected Diseases/Pests:</strong><br>`;
    data.diseases.forEach((disease) => {
      html += `• ${disease}<br>`;
    });
    html += `<br>`;
  }

  html += `<strong>Recommendations:</strong><br>`;
  data.suggestions.forEach((suggestion) => {
    html += `• ${suggestion}<br>`;
  });

  result.innerHTML = html;
}
