## 🏨 Starightforward hotels list API with one route.

---

[![Tests](https://github.com/gooddog-rent/hotels_api/actions/workflows/develop_tests.yml/badge.svg?branch=develop)](https://github.com/gooddog-rent/hotels_api/actions/workflows/develop_tests.yml)

[![Build and Push Docker image to Registry](https://github.com/gooddog-rent/hotels_api/actions/workflows/push_to_regisrty.yml/badge.svg?branch=master)](https://github.com/gooddog-rent/hotels_api/actions/workflows/push_to_regisrty.yml)

## Description

Simple API with one json custom data file.

### Features

🔃 - auto update file changes (no need to restart binary)

⚙️ - one JSON flexible file to modify your data

✋ - rate limiter for highload perfomance

💵 - in-memory cache

🗜️ - gzip or brotli compression (more then body length 1400 bytes)

## Setup enviroments

- Create .env file at your root project folder and fill it with your data (see .env.example to reference)

### Run in Docker container via docker compose (example)

```bash
docker-compose up -d
```

---

## **Get hotels (API endpoint)**

Returns json data with list of region and hotels.

- **URL**

  /hotels

- **Method:**

  `GET`

- **URL Params**

  **Required:**

  `query=[string]`

  Hotel search query string

  **Required:**

  `limit=[decimal]`

  Limit hotels per one region

- **Success Response:**

  - **Code:** 200 <br />
    **Content:**
    ```json
    {
      "hotels": {
        "Bavaro": [
          "Impressive Resorts Spas",
          "Melia Punta Cana Beach Resort Adults Only",
          "Royalton Bavaro Resort Spa"
        ],
        "La Romana": [
          "Santana Beach Resort",
          "Viva Wyndham Dominicus Palace Resort"
        ],
        "Punta Cana Hotels": [
          "The Westin Punta Сana Resort Club",
          "Tortuga Bay Hotel at Punta Сana Resort Club"
        ],
        "Uvero Alto": [
          "Dreams Punta Cana Resort Hotel",
          "Nickelodeon Hotels Resorts Punta Cana",
          "Sensatori Resort Punta Cana"
        ]
      }
    }
    ```

- **Error Response:**

  - **Code:** 400 BAD REQUEST <br />
    **Content:** `HTTP page 400 Bad Request`

  - **Code:** 405 METHOD NOT ALLOWED <br />
    **Content:** `HTTP page 405 Method Not Allowed`

  - **Code:** 429 TOO MANY REQUESTS <br />
    **Content:** `HTTP page 429 Too Many Requests`

  - **Code:** 500 INTERNAL SERVER ERROR <br />
    **Content:** `HTTP page 500 Internal Server Error`

- **Sample Call:**

  ```shell
  curl "https://gooddog.rent/hotels?query=resort&limit=10"
  ```
