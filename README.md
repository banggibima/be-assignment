# be-assignment

Proyek ini merupakan backend service berbasis Go Gin + PostgreSQL

## how to run (using docker compose)

1. Clone repository

   ```sh
   git clone https://github.com/banggibima/be-assignment.git
   cd be-assignment
   ```

2. Build & start containers

   ```sh
   docker compose up --build
   ```

3. Jalankan container di background

   ```sh
   docker compose up -d
   ```

4. Tunggu hingga log menunjukkan:

   ```sh
   [GIN-debug] listening and serving http on :8080
   ```

   Aplikasi akan otomatis membuat database, menjalankan migrasi, dan seeding data.

5. Akses server di:
   - base url: http://localhost:8081
   - health check: http://localhost:8081/health

## run tests (manual)

Jika ingin menjalankan test secara manual:

1.  Pastikan container aplikasi dan database sudah berjalan:

    ```sh
    docker compose up -d
    ```

2.  Jalankan container test secara terpisah:

    ```sh
    docker compose run --rm test
    ```

## swagger ui

Untuk melihat dokumentasi, buka:  
[http://localhost:8081/swagger/index.html](http://localhost:8081/swagger/index.html)

## api endpoints

| method   | endpoint             | description                             |
| -------- | -------------------- | --------------------------------------- |
| **GET**  | `/health`            | mengecek status server                  |
| **POST** | `/orders`            | membuat order baru                      |
| **GET**  | `/orders/:id`        | mendapatkan detail order berdasarkan ID |
| **GET**  | `/jobs/:id`          | mendapatkan status job tertentu         |
| **POST** | `/jobs/:id/cancel`   | membatalkan job yang sedang berjalan    |
| **POST** | `/jobs/settlement`   | menjalankan proses settlement job       |
| **GET**  | `/downloads/:job_id` | mengunduh hasil job berdasarkan ID      |

## notes

- default backend container berjalan di port 3000, di PC bisa diakses 3001
- default database container berjalan di port 5432, di PC bisa diakses 5433
- environment variables diatur otomatis lewat docker-compose.yaml
