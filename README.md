# Data Persistence & API Design with Intelligent Query Engine

This project is a high-performance Go API designed to process and store demographic profiles. It integrates three external APIs (Genderize, Agify, and Nationalize) to classify users by gender, age, and nationality, storing the results in a persistent PostgreSQL database and features a custom-built Natural Language Query (NLQ) engine for intuitive data retrieval..

---

## 🚀 Live API URL
**Base URL:** `https://stage1.doxantro.com`  


---

## 🛠 Tech Stack
- **Language:** Go (Golang) 1.22+
- **Router:** [Chi Router](https://github.com/go-chi/chi) (Standard library-compliant routing)
- **Database:** Neon Serverless PostgreSQL
- **ID Standard:** UUID v7 (Time-ordered unique identifiers)
- **Deployment:** Azure App Service

---

## ✨ Intelligence & Performance Features

### 🔍 Natural Language Query (NLQ) Engine
The centerpiece of Stage 2. The `/api/profiles/search` endpoint uses a rule-based parser (Regex & Keyword Mapping) to interpret human queries:
- **Geography:** Recognizes full country names (e.g., "People from Nigeria") and maps them to ISO codes.
- **Demographics:** Infers gender from keywords like "men," "women," "males," or "females."
- **Age Logic:** - Maps "young" to the 16-24 age range.
  - Recognizes age groups (child, teenager, adult, senior).
  - Uses regex to parse phrases like "above 40."

### ⚡ Concurrent Processing
Uses Go's `sync.WaitGroup` to fetch data from three external APIs simultaneously. This ensures that even with three network calls, the response time remains minimal.

### 🔢 Advanced Filtering & Pagination
The system handles large datasets (2,026+ records) with professional-grade pagination and sorting:
- **Metadata:** Responses include `total`, `page`, and `limit`.
- **Sorting:** Support for `age`, `gender_probability`, and `created_at`.
- **Strict 502 Handling:** Returns a **502 Bad Gateway** if upstream providers fail, ensuring only high-quality data persists.

---

## 📡 API Endpoints

### 1. Intelligence Search (NLQ)
**`GET /api/profiles/search?q=young+males+from+Hungary`**
- **Logic:** Interprets the string `q` and redirects to the filtered list logic.

### 2. Create/Retrieve Profile
**`POST /api/profiles`**
- **Body:** `{ "name": "olamide" }`
- **Behavior:** Concurrent fetching + UUID v7 generation. Idempotent by name.

### 3. List Profiles (Power Query)
**`GET /api/profiles`**
- **Parameters:** `gender`, `country_id`, `age_group`, `min_age`, `max_age`, `sort_by`, `page`, `limit`.

### 4. Delete/Get Single
**`GET /api/profiles/{id}`** | **`DELETE /api/profiles/{id}`**
- Full support for UUID v7 lookups.

---

## 🗄 Database Schema (Optimized)
The schema has been strictly aligned for performance and human-readable search:

```sql
CREATE TABLE profiles (
    id UUID PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    gender VARCHAR(50),
    gender_probability FLOAT,
    age INT,
    age_group VARCHAR(50),
    country_id VARCHAR(10),
    country_name VARCHAR(255), -- newly added for the  Intelligence Engine feature
    country_probability FLOAT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT (now() AT TIME ZONE 'utc')
);