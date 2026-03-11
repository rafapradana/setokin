# Product Requirements Document (PRD)

## Product Name

Setokin

## Product Overview

Web application untuk membantu bisnis **F&B di Indonesia** mengelola stok bahan baku dengan sistem sederhana dan praktis.

Produk ini fokus pada kebutuhan operasional dapur dan gudang:

* mencatat barang masuk (pembelian)
* mencatat barang keluar (dipakai dapur)
* melacak expiry bahan
* menggunakan sistem **FEFO (First Expired First Out)** - prioritas pakai bahan yang paling cepat expired
* menyediakan laporan penggunaan stok

Targetnya bukan enterprise warehouse, tapi **resto, cafe, cloud kitchen, bakery, UMKM F&B**.

> **Note:** Sistem menggunakan FEFO (First Expired First Out), bukan FIFO tradisional. Artinya saat stock out, sistem otomatis mengambil dari batch dengan expiry date paling dekat, bukan berdasarkan tanggal pembelian.

---

# Problem

Banyak bisnis F&B di Indonesia mengalami masalah:

1. Stok bahan tidak tercatat dengan baik
2. Banyak bahan **expired atau terbuang**
3. Tidak tahu **berapa bahan yang benar-benar dipakai**
4. Pencatatan masih manual (Excel / buku)
5. Sulit membuat laporan stok

Akibatnya:

* food cost tidak terkontrol
* sering kehabisan bahan
* sering overstock

---

# Goals

### Business Goals

* Membuat sistem inventory sederhana untuk F&B
* Mengurangi food waste
* Memberikan visibilitas stok bahan

### User Goals

User bisa:

* tahu stok bahan saat ini
* tahu bahan yang hampir expired
* mencatat bahan masuk dan keluar dengan cepat
* melihat laporan penggunaan bahan

---

# Target Users

### Primary Users

**Owner restoran**

* ingin tahu stok dan penggunaan bahan

**Kitchen manager**

* mengelola bahan dapur

**Staff gudang**

* mencatat barang masuk dan keluar

---

# Core Features (MVP)

## 1. Item Management

User dapat membuat dan mengelola daftar bahan baku.

### Fields

* Item Name
* Category
* Unit (kg, gram, liter, pcs)
* Minimum Stock (optional)

### Example

```
Item: Ayam Fillet
Category: Daging
Unit: kg
Minimum stock: 5
```

### Actions

User dapat:

* create item
* edit item
* delete item
* view item list

---

# 2. Stock In (Barang Masuk)

Digunakan saat membeli bahan dari supplier.

Setiap barang masuk akan membuat **batch baru**.

### Fields

* Item
* Quantity
* Purchase date
* Expiry date
* Supplier (optional)
* Purchase price (optional)

### Example

```
Item: Ayam Fillet
Quantity: 10kg
Purchase Date: 10 Mar
Expiry Date: 18 Mar
```

### System Behavior

System akan:

* membuat **batch baru**
* menambahkan stok ke inventory

---

# 3. Stock Out (Barang Dipakai)

Digunakan saat dapur menggunakan bahan.

### Fields

* Item
* Quantity
* Date
* Notes (optional)

### Example

```
Item: Ayam Fillet
Quantity: 2kg
Notes: dipakai dapur
```

### FEFO Logic (First Expired First Out)

System otomatis mengambil stok dari batch dengan **expiry date paling dekat**.

Example:

```
Batch A (dibeli 1 Mar)
5kg — exp 20 Mar

Batch B (dibeli 5 Mar)
5kg — exp 15 Mar
```

Jika user pakai 3kg, system akan ambil dari **Batch B** (exp lebih dekat):

```
Batch A → 5kg (exp 20 Mar)
Batch B → 2kg sisa (exp 15 Mar)
```

Ini memastikan bahan yang lebih cepat expired dipakai duluan, mengurangi waste.

---

# 4. Inventory Dashboard

Halaman utama yang menampilkan stok saat ini.

### Information

* Item name
* Total stock
* Unit
* Low stock indicator

### Example

```
Ayam Fillet — 8kg
Tepung — 25kg
Susu — 12L
```

User dapat klik item untuk melihat batch.

---

# 5. Batch & Expiry Tracking

User dapat melihat batch per item.

### Information

* Batch quantity
* Expiry date
* Remaining stock

### Example

```
Ayam Fillet

Batch A
3kg — exp 15 Mar

Batch B
5kg — exp 20 Mar
```

---

# 6. Expiry Alerts

System memberi peringatan bahan yang akan expired.

### Rules

Default alert:

* 3 hari sebelum expiry

### Example

```
Expiring Soon

Ayam Fillet — exp in 2 days
Susu — exp tomorrow
```

---

# 7. Reports

User dapat melihat laporan stok.

## Daily Report

Menampilkan aktivitas hari tersebut.

Example:

```
10 Mar

Stock In
+ Ayam Fillet 10kg

Stock Out
- Ayam Fillet 2kg
- Tepung 1kg
```

---

## Weekly Usage Report

Menampilkan total penggunaan bahan selama minggu tersebut.

Example:

```
Weekly Usage

Ayam Fillet — 18kg
Tepung — 9kg
Susu — 12L
```

---

## Monthly Usage Report

Menampilkan penggunaan bahan selama bulan tersebut.

Example:

```
Monthly Usage

Ayam Fillet — 72kg
Tepung — 40kg
Minyak — 28L
```

---

# User Flow

## Flow 1 — Barang Datang

```
User buka Stock In
↓
Input barang
↓
Isi quantity
↓
Isi expiry date
↓
Save
```

System membuat batch baru.

---

## Flow 2 — Bahan Dipakai

```
User buka Stock Out
↓
Pilih item
↓
Masukkan quantity
↓
Save
```

System otomatis:

* menggunakan batch FEFO (ambil dari expiry paling dekat)
* mengurangi stok

---

## Flow 3 — Cek Expiry

```
User buka Dashboard
↓
Lihat section Expiring Soon
↓
Klik item
↓
Lihat batch detail
```

---

# Data Model (Simplified)

## items

```
id
name
category
unit
minimum_stock
created_at
```

---

## batches

```
id
item_id
quantity
remaining_quantity
expiry_date
created_at
```

---

## stock_in

```
id
item_id
batch_id
quantity
purchase_date
supplier
purchase_price
created_at
```

---

## stock_out

```
id
item_id
quantity
notes
created_at
```

---

# Tech Stack & Architecture

## Architecture Overview

Setokin menggunakan **monorepo architecture** dengan struktur:

```
setokin/
├── frontend/          # Next.js app
├── backend/           # Go Fiber API
├── docker-compose.yml
├── .env              # Single env file (root level)
├── Makefile
└── scripts/          # Bash & PowerShell scripts
```

---

## Tech Stack

### Frontend

* **Next.js** - React framework dengan SSR/SSG
* **TypeScript** - Type safety
* **React** - UI library
* **Shadcn UI** - Component library
* **Tailwind CSS** - Styling
* **Phosphor Icons** - Icon set

### Backend

* **Go (Golang)** - Backend language
* **Fiber** - Web framework (Express-like untuk Go)
* **GORM** - ORM untuk database operations
* **JWT** - Authentication dengan access & refresh token

### Database & Storage

* **PostgreSQL** - Primary database
* **MinIO** - Object storage (S3-compatible)

### DevOps & Infrastructure

* **Docker** - Containerization
* **Docker Compose** - Multi-container orchestration
* **Nginx** - Reverse proxy & static file serving
* **Makefile** - Task automation
* **Bash & PowerShell** - Cross-platform scripting

---

## Architecture Details

### Monorepo Structure

Semua service dalam satu repository untuk:

* Easier dependency management
* Shared configuration
* Simplified deployment

### Environment Management

* **Single `.env` file** di root directory
* Environment variables di-inject via Docker Compose
* Shared across frontend & backend

### Authentication Flow

```
User login
↓
Backend generates:
- Access token (short-lived, 15 min)
- Refresh token (long-lived, 7 days)
↓
Frontend stores tokens
↓
Access token expired?
→ Use refresh token to get new access token
```

### API Architecture

```
Client (Next.js)
↓
Nginx (reverse proxy)
↓
Backend API (Go Fiber)
↓
PostgreSQL / MinIO
```

### Container Setup

Docker Compose menjalankan:

* `frontend` - Next.js app
* `backend` - Go Fiber API
* `postgres` - Database
* `minio` - Object storage
* `nginx` - Reverse proxy

---

## Development Workflow

### Local Development

```bash
# Start all services
make dev

# Run frontend only
make frontend

# Run backend only
make backend
```

### Build & Deploy

```bash
# Build Docker images
make build

# Deploy with Docker Compose
make deploy
```

---

# Non Functional Requirements

### Performance

* Dashboard load < 2 seconds

### Usability

* Form input harus sederhana
* Mobile friendly (banyak dipakai di dapur)

### Security

* Authentication required
* Role-based access optional (future)

---

# Future Features (Post MVP)

Tidak termasuk dalam MVP tapi bisa ditambahkan nanti.

* POS integration
* Recipe → auto stock deduction
* Supplier management
* Purchase order
* Multi outlet inventory
* Stock opname mode
* WhatsApp restock alerts

---

# MVP Scope Summary

Fitur yang **HARUS ada di V1**:

1. Item management
2. Stock In
3. Stock Out
4. FEFO stock deduction (prioritas expiry terdekat)
5. Expiry tracking
6. Inventory dashboard
7. Daily / Weekly / Monthly reports

---