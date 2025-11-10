# Post Articles API (Fiber + MySQL raw)

Project API untuk membangun REST API Post Articles Simple dengan Go Fiber, menggunakan MySQL via `database/sql` (tanpa ORM).

## Struktur Folder

- `cmd/server/main.go` — Entry point server Fiber.
- `internal/article/` — Domain Article (model, DTO, repository raw SQL, service, handler).
- `internal/router/router.go` — Registrasi routes.
- `pkg/config/` — Loader konfigurasi dari environment.
- `pkg/database/` — Koneksi MySQL (go-sql-driver/mysql) via `database/sql`.
- `pkg/response/` — Helper response JSON.
- `migrations/` — File migrasi SQL.
- `.env.example` — Contoh konfigurasi env.

## Menjalankan Aplikasi

1. Salin env (opsional untuk jalankan lokal):
   ```bash
   cp .env.example .env
   # sesuaikan DATABASE_URL
   ```
2. Inisialisasi module dan ambil dependensi (untuk lokal):
   ```bash
   cd github.com/ranggakrisnaa/sharing-vision-backend
   go mod init github.com/ranggakrisnaa/sharing-vision-backend
   go mod tidy
   ```
3. Jalankan migrasi (pakai CLI MySQL, untuk lokal):
   ```bash
   # Pastikan DB "sharing_vision" sudah dibuat
   mysql -h 127.0.0.1 -P 3306 -u root -p sharing_vision < migrations/0001_create_articles_table.sql
   ```
4. Jalankan server lokal:
   ```bash
   go run cmd/server/main.go
   ```

## Menjalankan dengan Docker Compose

1. Pindah ke direktori `github.com/ranggakrisnaa/sharing-vision-backendd` (tempat file compose):

   ```bash
   cd github.com/ranggakrisnaa/sharing-vision-backendd
   ```

2. Build dan jalankan semua layanan (DB, migrasi, API):

   ```bash
   docker compose up -d --build
   ```

   - Service `db` (MySQL) expose `3306` ke host.
   - Service `migrate` menjalankan berkas migrasi di `/app/migrations` (golang-migrate up) lalu selesai.
   - Service `api` expose `http://localhost:8080`.

3. Logs (opsional):

   ```bash
   docker compose logs -f api
   docker compose logs -f db
   ```

4. Stop dan hapus container:
   ```bash
   docker compose down
   # Jika ingin hapus data MySQL:
   docker compose down -v
   ```

## Migrasi Manual

Jalankan migrasi tanpa Docker atau secara on-demand:

- Via Go (membaca `.env`):

  ```bash
  cd github.com/ranggakrisnaa/sharing-vision-backendd
  go run cmd/migrate/main.go -dir migrations -action up
  ```

- Via Compose (gunakan image API yang sudah terbangun):
  ```bash
  cd github.com/ranggakrisnaa/sharing-vision-backendd
  docker compose run --rm migrate /app/migrate -dir /app/migrations -action up
  ```

Rollback:

```bash
cd github.com/ranggakrisnaa/sharing-vision-backendd
go run cmd/migrate/main.go -dir migrations -action down
# atau
docker compose run --rm migrate /app/migrate -dir /app/migrations -action down
```

## Endpoint

- `POST /article` — membuat artikel.
- `GET /article` - daftar artikel dengan pagination.
- `GET /article/:id` — detail artikel.
- `PUT /article/:id` — update artikel.
- `DELETE /article/:id` — hapus artikel.
