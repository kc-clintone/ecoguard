const API_URL = "http://localhost:8080"; // backend API

let isGuest = true;

// Show a specific section and hide others
function showSection(id) {
  const restricted = ["detector", "calendar"];
  if (isGuest && restricted.includes(id)) {
    alert("This feature is only for registered users.");
    return;
  }

  document.querySelectorAll(".container").forEach(c => c.classList.add("hidden"));
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
  if(!username || !password) { alert("Enter username and password"); return; }

  fetch(`${API_URL}/signup`, {
    method: "POST",
    headers: {"Content-Type": "application/json"},
    body: JSON.stringify({username, password})
  })
  .then(r => r.json())
  .then(d => {
    alert(d.message || "Signed up successfully");
    showSection("login");
  })
  .catch(e => alert("Error signing up: " + e.message));
}

// Login
function login() {
  const username = document.getElementById("loginUser").value;
  const password = document.getElementById("loginPass").value;
  if(!username || !password) { alert("Enter username and password"); return; }

  fetch(`${API_URL}/login`, {
    method: "POST",
    headers: {"Content-Type": "application/json"},
    body: JSON.stringify({username, password})
  })
  .then(r => {
    if(!r.ok) throw new Error("Invalid username or password");
    return r.json();
  })
  .then(d => {
    isGuest = false;
    alert(`Welcome ${d.Username}`);
    showSection("home");
  })
  .catch(e => alert(e.message));
}

// Placeholder functions for alerts and calendar
function fetchClimateAlerts() {
  document.getElementById("alerts-content").innerText =
    "Heavy rainfall expected.\nHigh temperatures next week.";
}

function fetchCalendar() {
  document.getElementById("calendar-content").innerText =
    "Plant maize this week.\nIrrigate twice weekly.";
}

// Optional: preview crop image
function previewImage(event) {
  const preview = document.getElementById("preview");
  preview.src = URL.createObjectURL(event.target.files[0]);
  preview.classList.remove("hidden");
  document.getElementById("result").innerText = "Image ready for AI detection (mock)";
}