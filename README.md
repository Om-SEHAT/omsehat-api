# OmSapa API

Backend API for [OmSEHAT](https://api-omsehat.sportsnow.app), built with **Go** and **Gin**, powered by **PostgreSQL**, **Docker**, and **Gemini**.

---

## 🚀 Features

- User Registration with OTP Verification
- AI-Powered Chat Sessions
- Doctor Diagnosis Support
- Appointment Queue Management
- PostgreSQL for Persistence
- Dockerized Setup
- Gemini Model Integration
- SMTP Support for Email Notifications

---

## 🧱 Tech Stack

- Go + Gin (API Framework)
- PostgreSQL (Database)
- Docker (Containerization)
- Gemini (AI integration)

---

## 🔧 Setup Instructions

### 1. Clone the Repository

```bash
git clone git@github.com:Om-SEHAT/omsehat-api.git
cd omsehat-api
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Configure Environment

Create a `.env` file in the root directory based on `.env.example`.

### 4. Run the Application

```bash
go run main.go
```

Or use [air](https://github.com/cosmtrek/air) for live reload:

```bash
air
```

### 5. Run with Docker

Make sure Docker is installed, then run:

```bash
docker compose up --build
```

> The API will be available at [http://localhost:8080](http://localhost:8080)

---

## ⚠️ Error Body Format

Error responses follow this JSON structure:

```json
{
  "message": "Error message",
  "details": {
    "FieldName": "Error description"
  }
}
```

---

## 📚 API Endpoints

### 🔐 `POST /register`

Register a new user.

**Request Body:**

```json
{
  "name": "Mario",
  "email": "new@gmail.com",
  "nationality": "Indonesian",
  "dob": "2004-04-01",
  "weight": 40.2,
  "gender": "male",
  "height": 165.6,
  "heartrate": 98.6,
  "bodytemp": 35.5
}
```

**Response:**

```json
{
  "message": "User registered successfully",
  "user": {
    "id": "some-uuid",
    "name": "Mario",
    "email": "new@gmail.com"
  }
}
```

---

### ✅ `POST /verify-otp`

Verify OTP and create a session.

**Request Body:**

```json
{
  "email": "new@gmail.com",
  "OTP": "132267",
  ...
}
```

---

### 💬 `POST /session/:id`

Send a new message in a session. It also determines the next action (continue chat or schedule an appointment).

**Request Body:**

```json
{
  "new_message": "What is my name and age?"
}
```

**Response:**

```json
{
  "message": "Chat history updated successfully",
  "reply": "...",
  "next_action": "CONTINUE_CHAT", // or "APPOINTMENT
  "session_id": "uuid"
}
```

---

### 🩺 `POST /session/:id/diagnose`

Add diagnosis to a session.

**Request Body:**

```json
{
  "diagnosis": "Patient has mild fever."
}
```

**Response:**

```json
{
  "message": "Diagnosis saved successfully"
}
```

---

### 📄 `GET /session/:id`

Fetch session details and chat history.

**Sample Response:**

```json
{
  "id": "session-id",
  "user": {
    "id": "user-id",
    "name": "Mario",
    "email": "new@gmail.com"
  },
  "messages": [
    {
      "role": "user",
      "content": "Hello!"
    },
    {
      "role": "omsapa",
      "content": "Hello Mario! ..."
    }
  ]
}
```

---

### 📅 `GET /queue/:doctor_id/`

Fetch current appointment queue for a doctor.

---

### 🧑‍⚕️ `GET /doctor/:id`

Fetch doctor details for home page.

**Sample Response:**

```json
{
  "appointment_count_all_time": 2,
  "appointment_count_daily": 2,
  "current_queue": {
    "id": "861ae8de-4a35-4640-9302-20d82f97e3f6",
    "doctor_id": "f186afd5-a175-420e-b06e-d35a713d3616",
    "session_id": "4cc39394-f2b8-4133-9ea1-c03215a58a72",
    "session": {
      "id": "4cc39394-f2b8-4133-9ea1-c03215a58a72",
      "user_id": "22f94f7d-e96f-4da5-ad2b-6539d6f9543d",
      "weight": 52,
      "height": 170,
      "heartrate": 90,
      "bodytemp": 36,
      "prediagnosis": "Possible influenza (flu)",
      "created_at": "2025-05-16T11:17:53.123726Z",
      "updated_at": "2025-05-16T11:29:38.812007Z"
    },
    "number": 1,
    "created_at": "2025-05-16T11:29:31.672677Z",
    "updated_at": "2025-05-16T11:29:31.672677Z"
  },
  "doctor": {
    "id": "f186afd5-a175-420e-b06e-d35a713d3616",
    "name": "dr. Udin",
    "email": "udin@example.com",
    "specialty": "General Practitioner",
    "roomno": "A2"
  }
}
```

---

### 📄 `GET /user/:id`

Fetch user details, current session, and session history.

**Sample Response:**

```json
{
  "user": {
    "id": "22f94f7d-e96f-4da5-ad2b-6539d6f9543d",
    "name": "Mario",
    "email": "new@gmail.com",
    "gender": "Male",
    "nationality": "Italian",
    "age": "23"
  },
  "current_session": {
    "queue": {
      "id": "861ae8de-4a35-4640-9302-20d82f97e3f6",
      "number": 1
    },
    "session_id": "4cc39394-f2b8-4133-9ea1-c03215a58a72",
    "bodytemp": 36,
    "doctor_diagnosis": "",
    "heartrate": 90,
    "height": 170,
    "prediagnosis": "Possible influenza (flu)",
    "weight": 52,
    "created_at": "2025-05-16T11:29:38.812007Z"
  },
  "history_sessions": [
    {
      "session_id": "18e92407-992c-4f47-9126-0b9eac8bd520",
      "bodytemp": 37,
      "doctor_diagnosis": "Common cold",
      "heartrate": 85,
      "height": 170,
      "prediagnosis": "Mild cold symptoms",
      "weight": 51,
      "created_at": "2025-01-10T08:45:23.123456Z"
    }
  ]
}
```

---

## 🧠 Gemini Integration

OmSEHAT leverages [Gemini](https://deepmind.google/technologies/gemini/) for contextual and medical-like conversational intelligence. The AI uses your user metrics (age, weight, vitals, etc.) to provide personalized replies.

---

## 📎 Links

* 🌐 [OmSEHAT Web App](https://omsehat.sportsnow.app)
* 🐙 [GitHub Repository](https://github.com/Om-SEHAT/omsehat-api)
