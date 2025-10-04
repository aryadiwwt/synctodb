# Go API Data Synchronizer

Sebuah service Go yang dirancang untuk mengambil data dari API eksternal, memprosesnya, dan menyimpannya ke dalam database PostgreSQL. Proyek ini dibangun dengan prinsip-prinsip *Clean Code* dan arsitektur modular untuk memastikan kode mudah di-maintain, diuji, dan dikembangkan lebih lanjut.

-----

## Fitur Utama

  * **Arsitektur Modular**: Logika dipisahkan ke dalam komponen-komponen independen (`Fetcher`, `Storer`, `Synchronizer`) untuk kejelasan dan *testability*.
  * **Aman**: Kredensial dan konfigurasi sensitif dikelola di luar kode menggunakan file `.env`.
  * **Penanganan Error yang Baik**: Error dibungkus dengan konteks untuk mempermudah proses *debugging*.
  * **Resilien**: Dilengkapi dengan *timeout* pada *request* HTTP dan menangani respons API yang kompleks, termasuk paginasi.
  * **Manajemen Dependensi**: Menggunakan Go Modules (`go.mod` dan `go.sum`) untuk memastikan *build* yang konsisten dan dapat direproduksi.

-----

## Cara Menjalankan

### **Prasyarat**

  * [Go](https://golang.org/dl/) versi 1.18 atau lebih baru.
  * [PostgreSQL](https://www.postgresql.org/download/) sebagai database.
  * `git` untuk meng-kloning repositori.

### **1. Kloning Repositori**

```bash
git clone <url-repositori-anda>
cd <nama-direktori-proyek>
```

### **2. Siapkan Database**

Pastikan PostgreSQL Anda berjalan. Buat database dan tabel yang diperlukan. Contoh skema tabel:

```sql
CREATE TABLE siskeudes_detail_output (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tahun TEXT NOT NULL,
    kode_desa TEXT NOT NULL,
    id_kegiatan TEXT NOT NULL,
    no_id TEXT NOT NULL,
    -- (tambahkan semua kolom lain sesuai struct domain Anda)
    jabatan_pptkd TEXT,
    
    -- Tambahkan UNIQUE constraint untuk kunci bisnis
    CONSTRAINT uq_output_detail_business_key UNIQUE (tahun, kode_desa, id_kegiatan, no_id)
);
```

### **3. Konfigurasi Environment**

Salin file `.env.example` (jika ada) atau buat file baru bernama `.env`. Isi file ini dengan konfigurasi Anda.

```env
# Kredensial untuk Login API
API_USERNAME="username_anda_disini"
API_PASSWORD="password_rahasia_anda"

# URL API
API_LOGIN_URL="https://konsolidasi-apbdesa.kemendagri.go.id/api/login"
API_URL="https://konsolidasi-apbdesa.kemendagri.go.id/api/rekap/output/detail"

# Parameter untuk request data
API_DATA_TAHUN="2025"
API_DATA_KD_PROV="51"
API_DATA_KD_KAB="03"

# Konfigurasi Database PostgreSQL
DATABASE_URL="postgres://user:password@localhost:5432/dbname?sslmode=disable"
```

### **4. Instal Dependensi**

Go akan secara otomatis mengunduh dependensi yang tercantum di `go.mod` saat Anda menjalankan aplikasi. Anda juga bisa menjalankannya secara manual:

```bash
go mod tidy
```

### **5. Jalankan Service**

```bash
go run main.go
```

Anda akan melihat output log di terminal yang menunjukkan proses sinkronisasi data.

-----

## Struktur Proyek

```
/
├── .env                  # Menyimpan konfigurasi & kredensial (diabaikan oleh Git)
├── .gitignore            # Daftar file yang diabaikan oleh Git
├── go.mod                # Definisi modul dan dependensi
├── go.sum                # Checksum untuk integritas dependensi
├── main.go               # Titik masuk aplikasi dan "wiring" dependensi
├── README.md               # Dokumentasi proyek
├── config/               # Mengelola pemuatan konfigurasi
├── domain/               # Definisi struct untuk entitas data inti
├── fetcher/              # Komponen untuk mengambil data dari API
├── storer/               # Komponen untuk menyimpan data ke database
└── synchronizer/         # Mengorkestrasi alur kerja fetch-and-store
```

-----

## Teknologi yang Digunakan

  * **Bahasa**: [Go](https://golang.org/)
  * **Database**: [PostgreSQL](https://www.postgresql.org/)
  * **Library Utama**:
      * `net/http`: Untuk melakukan request HTTP.
      * `github.com/jmoiron/sqlx`: Ekstensi untuk package `database/sql` standar Go.
      * `github.com/lib/pq`: Driver untuk PostgreSQL.
      * `github.com/joho/godotenv`: Untuk memuat file `.env`.