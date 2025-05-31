# OmSapa API

Backend API for [OmSEHAT](https://api-omsehat.sportsnow.app), built with **Go** and **Gin**, powered by **PostgreSQL**, **Docker**, and **Gemini**.

---

## 🚀 Features

- User Registration with OTP Verification
- AI-Powered Chat Sessions
- Doctor Diagnosis Support
- Appointment Queue Management
- AI Psychologist for Healthcare Worker Burnout
- Mental Health Support for General Users
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

The mental health chatbot provides specialized support for:
- Healthcare workers experiencing burnout and stress due to high workloads, especially in areas with high COVID-19 cases
- General users with mental health concerns
- Providing personalized coping strategies and self-care recommendations
- Offering evidence-based mental health support in a conversational format

---

## 📎 Links

* 🌐 [OmSEHAT Web App](https://omsehat.sportsnow.app)
* 🐙 [GitHub Repository](https://github.com/Om-SEHAT/omsehat-api)

---

## 🧠 Therapy Chatbot Endpoints

### 🧠 `POST /therapy`

Create a new therapy session for mental health support.

**Request Body:**

```json
{
  "stress_level": 8,
  "mood_rating": 5,
  "sleep_quality": 4,
  "is_health_worker": true,
  "specialization": "Nurse",
  "work_hours": 50
}
```

**Response:**

```json
{
  "id": "therapy-session-id",
  "user_id": "user-id",
  "stress_level": 8,
  "mood_rating": 5,
  "sleep_quality": 4,
  "is_health_worker": true,
  "specialization": "Nurse",
  "work_hours": 50,
  "created_at": "2025-05-30T14:28:45.123456Z",
  "updated_at": "2025-05-30T14:28:45.123456Z"
}
```

---

### 💬 `POST /therapy/:id`

Send a new message in a therapy session.

**Request Body:**

```json
{
  "new_message": "I've been feeling very stressed about work lately."
}
```

**Response:**

```json
{
  "message": "Chat history updated successfully",
  "next_action": "CONTINUE_CHAT",
  "reply": "I understand that you're feeling stressed about work. Can you tell me more about what aspects of your job are causing the most stress?",
  "recommendation": "Consider practicing mindfulness and deep breathing exercises when feeling overwhelmed.",
  "next_steps": "Try setting boundaries at work and taking short breaks throughout your shift.",
  "session_id": "therapy-session-id"
}
```

---

### 📄 `GET /therapy/:id`

Fetch active therapy session details and chat history.

**Sample Response:**

```json
{
  "id": "therapy-session-id",
  "user": {
    "id": "user-id",
    "name": "Sarah",
    "email": "sarah@example.com"
  },
  "stress_level": 8,
  "mood_rating": 5,
  "sleep_quality": 4,
  "is_health_worker": true,
  "specialization": "Nurse",
  "work_hours": 50,
  "messages": [
    {
      "role": "user",
      "content": "I've been feeling very stressed about work lately."
    },
    {
      "role": "therapist",
      "content": "I understand that you're feeling stressed about work. Can you tell me more about what aspects of your job are causing the most stress?"
    }
  ],
  "created_at": "2025-05-30T14:28:45.123456Z",
  "updated_at": "2025-05-30T14:28:45.123456Z"
}
```

---

### 📚 `GET /therapy/history`

Fetch a user's therapy session history.

**Sample Response:**

```json
{
  "sessions": [
    {
      "id": "therapy-session-id-1",
      "stress_level": 8,
      "mood_rating": 5,
      "sleep_quality": 4,
      "is_health_worker": true,
      "recommendation": "Practice mindfulness techniques and establish better work-life boundaries.",
      "next_steps": "Consider joining a healthcare worker support group and speaking with your supervisor about workload concerns.",
      "created_at": "2025-05-30T14:28:45.123456Z"
    },
    {
      "id": "therapy-session-id-2",
      "stress_level": 6,
      "mood_rating": 7,
      "sleep_quality": 6,
      "is_health_worker": true,
      "recommendation": "Continue with breathing exercises and add regular physical activity to your routine.",
      "next_steps": "Schedule at least 30 minutes of personal time each day and maintain a consistent sleep schedule.",
      "created_at": "2025-05-20T09:15:30.123456Z"
    }
  ]
}
```

---

### 📋 `GET /therapy/detail/:id`

Fetch details of a specific therapy session, including chat history.

**Sample Response:**

```json
{
  "id": "therapy-session-id",
  "user": {
    "id": "user-id",
    "name": "Sarah",
    "email": "sarah@example.com"
  },
  "stress_level": 8,
  "mood_rating": 5,
  "sleep_quality": 4,
  "is_health_worker": true,
  "specialization": "Nurse",
  "work_hours": 50,
  "messages": [
    {
      "role": "user",
      "content": "I've been feeling very stressed about work lately."
    },
    {
      "role": "therapist",
      "content": "I understand that you're feeling stressed about work. Can you tell me more about what aspects of your job are causing the most stress?"
    }
  ],
  "recommendation": "Practice mindfulness techniques and establish better work-life boundaries.",
  "next_steps": "Consider joining a healthcare worker support group and speaking with your supervisor about workload concerns.",
  "created_at": "2025-05-30T14:28:45.123456Z",
  "updated_at": "2025-05-30T14:45:12.789123Z"
}
```

---

### 📊 `GET /therapy/mental-health-summary`

Get a summary of a user's mental health trends over time based on their therapy sessions.

**Sample Response:**

```json
{
  "total_sessions": 5,
  "healthcare_worker": true,
  "avg_stress_level": 7.2,
  "avg_mood_rating": 5.6,
  "avg_sleep_quality": 4.8,
  "stress_trend": "improving",
  "mood_trend": "improving",
  "sleep_trend": "stable",
  "last_session_date": "2025-05-29",
  "completed_sessions": 5
}
```
