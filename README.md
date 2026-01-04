# Inventory Management System - RESTful API

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)
![Chi Router](https://img.shields.io/badge/Chi_Router-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![UUID](https://img.shields.io/badge/UUID-Token-orange?style=for-the-badge)

Aplikasi Inventory Management System berbasis RESTful API untuk mengelola data barang, kategori, rak, gudang, dan proses penjualan barang. Dibuat menggunakan Go (Golang), Chi Router, dan PostgreSQL sebagai project **Golang Intermediate Daytime Class - Lumoshive Academy Bootcamp**.

## Video Demo

[![Watch Demo](https://img.shields.io/badge/YouTube-FF0000?style=for-the-badge&logo=youtube&logoColor=white)](https://youtu.be/cz4xMg-xHPE)

**[Tonton Video Penjelasan Sistem](https://youtu.be/cz4xMg-xHPE)**

---

## Fitur Utama

- **User Authentication** - Login system dengan UUID Token & Bcrypt password hashing
- **Session Management** - Token dengan masa berlaku (expired_at) dan pencabutan sesi (revoked_at)
- **Role-Based Authorization** - Super Admin, Admin, dan Staff dengan permission berbeda
- **CRUD Barang (Items)** - Manajemen data barang dengan SKU, stock, dan harga
- **CRUD Kategori** - Pengelompokan barang berdasarkan kategori
- **CRUD Rak** - Manajemen rak penyimpanan barang
- **CRUD Gudang** - Manajemen lokasi gudang
- **CRUD User** - Manajemen akun pengguna dengan pengaturan role
- **CRUD Penjualan (Sales)** - Pencatatan transaksi penjualan barang
- **Report Summary** - Laporan total barang, penjualan, dan pendapatan
- **Cek Stok Minimum** - Alert barang dengan stok di bawah threshold (default: 5)
- **Pagination** - Pagination untuk semua list endpoint
- **Input Validation** - Validasi data menggunakan go-playground/validator
- **Logging System** - Zap Logger dengan log rotation
- **Unit Testing** - Testing dengan mock pattern (Coverage >50%)

---

## Role & Permissions

### Super Admin

- ✅ Full access semua fitur
- ✅ CRUD semua entity
- ✅ Report & Cek stok minimum
- ✅ Manage role user

### Admin

- ✅ CRUD barang, kategori, rak, gudang
- ✅ CRUD users
- ✅ Create sale, list sales, detail sale
- ✅ Report & Cek stok minimum

### Staff

- ✅ Read barang/kategori/rak/gudang
- ✅ Create sale, list sales, detail sale
- ✅ Cek stok minimum
- ❌ Delete master data
- ❌ Manage users/roles
- ❌ Akses report revenue

---

## Struktur Project

```
project-app-inventory-restapi-golang-alvin/
├── database/
│   └── database.go        # Database connection dengan pgx/v5
├── db_file/
│   ├── ddl_inventory_management_system.sql    # Database schema
│   ├── dml_inventory_management_system.sql    # Sample data
│   └── Project App Inventory.postman_collection.json
├── dto/
│   ├── auth.go            # Auth request/response DTOs
│   ├── pagination.go      # Pagination DTO
│   └── *.go               # Other DTOs
├── handler/
│   ├── auth.go            # Login/Logout handlers
│   ├── item.go            # Item CRUD handlers
│   ├── category.go        # Category CRUD handlers
│   ├── rack.go            # Rack CRUD handlers
│   ├── warehouse.go       # Warehouse CRUD handlers
│   ├── user.go            # User CRUD handlers
│   ├── sale.go            # Sale CRUD handlers
│   └── report.go          # Report handler
├── middleware/
│   ├── auth.go            # Authentication & Role middleware
│   ├── logging.go         # Request logging middleware
│   └── middleware.go      # Middleware setup
├── model/
│   ├── item.go            # Item model
│   ├── category.go        # Category model
│   ├── user.go            # User & Session models
│   └── *.go               # Other models
├── repository/
│   ├── item.go            # Item repository
│   ├── item_test.go       # Item repository tests
│   ├── category.go        # Category repository
│   ├── session.go         # Session repository
│   ├── session_test.go    # Session repository tests
│   └── *.go               # Other repositories
├── router/
│   └── router.go          # Route definitions dengan Chi
├── service/
│   ├── auth.go            # Auth service
│   ├── auth_test.go       # Auth service tests
│   ├── item.go            # Item service
│   └── *.go               # Other services
├── utils/
│   ├── config.go          # Viper configuration
│   ├── logger.go          # Zap logger setup
│   ├── password_hash.go   # Bcrypt utilities
│   ├── token.go           # UUID token generation
│   ├── validator.go       # Input validation
│   └── response.go        # HTTP response helpers
├── logs/                  # Application logs
├── .env                   # Environment variables
├── main.go                # Entry point
├── go.mod                 # Go modules
└── go.sum                 # Dependencies
```

## API Endpoints

### Auth Endpoints

| Method | Endpoint         | Description | Auth Required |
| ------ | ---------------- | ----------- | ------------- |
| POST   | `/api/v1/login`  | User login  | ❌            |
| POST   | `/api/v1/logout` | User logout | ✅            |

### Items Endpoints

| Method | Endpoint                  | Description         | Role Required      |
| ------ | ------------------------- | ------------------- | ------------------ |
| GET    | `/api/v1/items`           | Get all items       | All authenticated  |
| GET    | `/api/v1/items/{id}`      | Get item by ID      | All authenticated  |
| GET    | `/api/v1/items/low-stock` | Get low stock items | All authenticated  |
| POST   | `/api/v1/items`           | Create new item     | Super Admin, Admin |
| PUT    | `/api/v1/items/{id}`      | Update item         | Super Admin, Admin |
| DELETE | `/api/v1/items/{id}`      | Delete item         | Super Admin, Admin |

### Categories Endpoints

| Method | Endpoint                  | Description         | Role Required      |
| ------ | ------------------------- | ------------------- | ------------------ |
| GET    | `/api/v1/categories`      | Get all categories  | All authenticated  |
| GET    | `/api/v1/categories/{id}` | Get category by ID  | All authenticated  |
| POST   | `/api/v1/categories`      | Create new category | Super Admin, Admin |
| PUT    | `/api/v1/categories/{id}` | Update category     | Super Admin, Admin |
| DELETE | `/api/v1/categories/{id}` | Delete category     | Super Admin, Admin |

### Racks Endpoints

| Method | Endpoint             | Description     | Role Required      |
| ------ | -------------------- | --------------- | ------------------ |
| GET    | `/api/v1/racks`      | Get all racks   | All authenticated  |
| GET    | `/api/v1/racks/{id}` | Get rack by ID  | All authenticated  |
| POST   | `/api/v1/racks`      | Create new rack | Super Admin, Admin |
| PUT    | `/api/v1/racks/{id}` | Update rack     | Super Admin, Admin |
| DELETE | `/api/v1/racks/{id}` | Delete rack     | Super Admin, Admin |

### Warehouses Endpoints

| Method | Endpoint                  | Description          | Role Required      |
| ------ | ------------------------- | -------------------- | ------------------ |
| GET    | `/api/v1/warehouses`      | Get all warehouses   | All authenticated  |
| GET    | `/api/v1/warehouses/{id}` | Get warehouse by ID  | All authenticated  |
| POST   | `/api/v1/warehouses`      | Create new warehouse | Super Admin, Admin |
| PUT    | `/api/v1/warehouses/{id}` | Update warehouse     | Super Admin, Admin |
| DELETE | `/api/v1/warehouses/{id}` | Delete warehouse     | Super Admin, Admin |

### Users Endpoints

| Method | Endpoint             | Description     | Role Required      |
| ------ | -------------------- | --------------- | ------------------ |
| GET    | `/api/v1/users`      | Get all users   | Super Admin, Admin |
| GET    | `/api/v1/users/{id}` | Get user by ID  | Super Admin, Admin |
| POST   | `/api/v1/users`      | Create new user | Super Admin, Admin |
| PUT    | `/api/v1/users/{id}` | Update user     | Super Admin, Admin |
| DELETE | `/api/v1/users/{id}` | Delete user     | Super Admin, Admin |

### Sales Endpoints

| Method | Endpoint             | Description     | Role Required      |
| ------ | -------------------- | --------------- | ------------------ |
| GET    | `/api/v1/sales`      | Get all sales   | All authenticated  |
| GET    | `/api/v1/sales/{id}` | Get sale by ID  | All authenticated  |
| POST   | `/api/v1/sales`      | Create new sale | All authenticated  |
| PUT    | `/api/v1/sales/{id}` | Update sale     | Super Admin, Admin |
| DELETE | `/api/v1/sales/{id}` | Delete sale     | Super Admin, Admin |

### Report Endpoints

| Method | Endpoint                  | Description        | Role Required      |
| ------ | ------------------------- | ------------------ | ------------------ |
| GET    | `/api/v1/reports/summary` | Get summary report | Super Admin, Admin |

---

## Author

**Alvin Rama S**

- GitHub: [@Alvinnn-R](https://github.com/Alvinnn-R)
- Bootcamp: Golang Intermediate Daytime Class - Lumoshive Academy

---

## License

This project is for educational purposes as part of Lumoshive Academy Bootcamp.

---

## Acknowledgments

- Lumoshive Academy - Golang Bootcamp
- Instructor & Mentors
- Fellow bootcamp participants
