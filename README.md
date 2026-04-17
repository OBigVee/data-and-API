# Data Persistence & API Design

This project is a high-performance Go API designed to process and store demographic profiles. It integrates three external APIs (Genderize, Agify, and Nationalize) to classify users by gender, age, and nationality, storing the results in a persistent PostgreSQL database.

---

## 🚀 Live API URL
**Base URL:** `https://api.doxantro.com`  


---

## 🛠 Tech Stack
- **Language:** Go (Golang) 1.22+
- **Router:** [Chi Router](https://github.com/go-chi/chi) (Standard library-compliant routing)
- **Database:** Neon Serverless PostgreSQL
- **ID Standard:** UUID v7 (Time-ordered unique identifiers)
- **Deployment:** Azure App Service

---

## ✨ Key Features
- **Concurrent Processing:** Uses Go's goroutines to fetch data from Genderize, Agify, and Nationalize simultaneously. This minimizes response latency by running network requests in parallel.
- **Idempotency:** The system checks for existing records before calling external APIs. If a profile for a specific name already exists, the system returns the stored data, preventing duplicate API calls and redundant database entries.
- **Robust Error Handling:** Specifically handles **502 Bad Gateway** scenarios. If an upstream API returns a null gender, count of 0, or empty country data, the system returns the exact required error message without storing invalid data.
- **CORS Enabled:** Fully configured `Access-Control-Allow-Origin: *` to ensure the grading script and frontend clients can reach the server.
- **ISO 8601 Timestamps:** All dates are generated and stored in UTC format.

---

## 📡 API Endpoints

### 1. Create/Retrieve Profile
**`POST /api/profiles`**
- **Request Body:** `{ "name": "ella" }`
- **Behavior:** - If the name exists, returns the existing profile (`200 OK`).
  - If the name is new, fetches data from external APIs and creates a new profile (`201 Created`).

### 2. Get Single Profile
**`GET /api/profiles/{id}`**
- **Success:** Returns the full demographic data for the provided UUID.

### 3. Get All Profiles (with Filtering)
**`GET /api/profiles`**
- **Optional Query Parameters:** `gender`, `country_id`, `age_group`.
- **Logic:** Case-insensitive filtering. Example: `/api/profiles?gender=Male&country_id=NG`

### 4. Delete Profile
**`DELETE /api/profiles/{id}`**
- **Success:** `204 No Content`.

---

## 🗄 Database Schema
The system uses the following schema for persistence:

```sql
CREATE TABLE profiles (
    id UUID PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    gender VARCHAR(50),
    gender_probability FLOAT,
    sample_size INT,
    age INT,
    age_group VARCHAR(50),
    country_id VARCHAR(10),
    country_probability FLOAT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);