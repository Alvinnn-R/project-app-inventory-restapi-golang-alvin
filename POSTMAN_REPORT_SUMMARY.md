# Postman Testing - Report Summary

Dokumentasi testing untuk endpoint Report Barang dengan laporan total barang, penjualan, dan pendapatan.

---

## Setup Awal

### 1. Login untuk mendapatkan Token

**Endpoint:** `POST http://localhost:8080/api/v1/auth/login`

**Body (JSON):**

```json
{
  "email": "superadmin@inventory.com",
  "password": "password123"
}
```

**Response:**

```json
{
  "status": true,
  "message": "login successful",
  "data": {
    "token": "a1b2c3d4-e5f6-47a8-b9c0-d1e2f3a4b5c6",
    "user": {
      "id": 1,
      "name": "Super Admin User",
      "email": "superadmin@inventory.com",
      "role_id": 1,
      "role_name": "super_admin",
      "is_active": true
    }
  }
}
```

**Simpan token** dari response untuk digunakan di header `Authorization: Bearer <token>` pada request berikutnya.

---

## GET REPORT SUMMARY

### Endpoint

```
GET http://localhost:8080/api/v1/reports/summary
```

### Headers

```
Authorization: Bearer a1b2c3d4-e5f6-47a8-b9c0-d1e2f3a4b5c6
```

### Expected Response (200 OK)

```json
{
  "status": true,
  "message": "success get report summary",
  "data": {
    "total_items": 10,
    "low_stock_items": 3,
    "total_sales": 5,
    "total_revenue": 15585000,
    "active_users": 3,
    "total_categories": 5,
    "total_warehouses": 5
  }
}
```

### Response Field Explanation

| Field              | Description                                        | Query Source                                                                |
| ------------------ | -------------------------------------------------- | --------------------------------------------------------------------------- |
| `total_items`      | Total semua barang di inventory                    | `SELECT COUNT(*) FROM items`                                                |
| `low_stock_items`  | Jumlah barang dengan stock di bawah minimum        | `SELECT COUNT(*) FROM items WHERE stock < minimum_stock`                    |
| `total_sales`      | Total transaksi penjualan (exclude yang di-delete) | `SELECT COUNT(*) FROM sales WHERE deleted_at IS NULL`                       |
| `total_revenue`    | Total pendapatan dari semua penjualan              | `SELECT COALESCE(SUM(total_amount), 0) FROM sales WHERE deleted_at IS NULL` |
| `active_users`     | Jumlah user aktif                                  | `SELECT COUNT(*) FROM users WHERE is_active = true`                         |
| `total_categories` | Total kategori barang                              | `SELECT COUNT(*) FROM categories`                                           |
| `total_warehouses` | Total gudang                                       | `SELECT COUNT(*) FROM warehouses`                                           |

---

## Test Cases

### TC-1: Get Report as Super Admin

**Endpoint:**

```
GET http://localhost:8080/api/v1/reports/summary
```

**Headers:**

```
Authorization: Bearer a1b2c3d4-e5f6-47a8-b9c0-d1e2f3a4b5c6
```

**Expected:** 200 OK dengan data lengkap

```json
{
  "status": true,
  "message": "success get report summary",
  "data": {
    "total_items": 10,
    "low_stock_items": 3,
    "total_sales": 5,
    "total_revenue": 15585000,
    "active_users": 3,
    "total_categories": 5,
    "total_warehouses": 5
  }
}
```

**Verifikasi:**

- ✅ Status code: 200 OK
- ✅ `total_items` = jumlah barang di database (10 items dari DML)
- ✅ `low_stock_items` = 3 (Keyboard Mechanical stock=3, Office Desk Standing stock=2, Pen Set stock=4)
- ✅ `total_sales` = 5 transaksi
- ✅ `total_revenue` = sum dari semua total_amount
- ✅ `active_users` = 3 (super_admin, admin, staff)
- ✅ `total_categories` = 5 kategori
- ✅ `total_warehouses` = 5 gudang

---

### TC-2: Get Report as Admin

**Login sebagai Admin:**

```json
POST http://localhost:8080/api/v1/auth/login
Body: {
    "email": "admin@inventory.com",
    "password": "password123"
}
```

**Get Report:**

```
GET http://localhost:8080/api/v1/reports/summary
Authorization: Bearer <admin_token>
```

**Expected:** 200 OK dengan data lengkap

```json
{
  "status": true,
  "message": "success get report summary",
  "data": {
    "total_items": 10,
    "low_stock_items": 3,
    "total_sales": 5,
    "total_revenue": 15585000,
    "active_users": 3,
    "total_categories": 5,
    "total_warehouses": 5
  }
}
```

**Note:** Admin juga bisa akses report (sesuai requirements)

---

### TC-3: Get Report as Staff (Optional - Tergantung Authorization Rule)

**Login sebagai Staff:**

```json
POST http://localhost:8080/api/v1/auth/login
Body: {
    "email": "staff@inventory.com",
    "password": "password123"
}
```

**Get Report:**

```
GET http://localhost:8080/api/v1/reports/summary
Authorization: Bearer <staff_token>
```

**Expected (Current):** 200 OK - Staff bisa akses report
**Expected (If Authorization Applied):** 403 Forbidden - Staff tidak boleh akses revenue report

**Note:** Sesuai requirements:

- ✅ Super Admin: Bisa akses report
- ✅ Admin: Bisa akses report
- ❌ Staff: Tidak boleh akses report revenue (opsional)

---

### TC-4: Unauthorized Access - No Token

**Endpoint:**

```
GET http://localhost:8080/api/v1/reports/summary
```

**Headers:** (No Authorization header)

**Expected:** 401 Unauthorized

```json
{
  "status": false,
  "message": "unauthorized",
  "data": null
}
```

---

### TC-5: Invalid Token

**Endpoint:**

```
GET http://localhost:8080/api/v1/reports/summary
```

**Headers:**

```
Authorization: Bearer invalid-token-12345
```

**Expected:** 401 Unauthorized

```json
{
  "status": false,
  "message": "invalid token",
  "data": null
}
```

---

### TC-6: Expired Token

**Endpoint:**

```
GET http://localhost:8080/api/v1/reports/summary
```

**Headers:**

```
Authorization: Bearer <expired_token>
```

**Expected:** 401 Unauthorized

```json
{
  "status": false,
  "message": "token expired",
  "data": null
}
```

---

## Data Verification Flow

### Verifikasi Low Stock Items

**1. Get Report:**

```
GET http://localhost:8080/api/v1/reports/summary
```

Response: `"low_stock_items": 3`

**2. List All Items:**

```
GET http://localhost:8080/api/v1/items?page=1
```

**3. Check Items dengan stock < minimum_stock:**

- Keyboard Mechanical: stock=3, minimum_stock=5 ❌ Low stock
- Office Desk Standing: stock=2, minimum_stock=5 ❌ Low stock
- Pen Set: stock=4, minimum_stock=10 ❌ Low stock
- **Total: 3 items** ✅ Sesuai dengan report

---

### Verifikasi Total Sales & Revenue

**1. Get Report:**

```
GET http://localhost:8080/api/v1/reports/summary
```

Response:

```json
{
  "total_sales": 5,
  "total_revenue": 15585000
}
```

**2. List All Sales:**

```
GET http://localhost:8080/api/v1/sales?page=1
```

**3. Manual Calculation:**

```
Sale 1: Rp 9,000,000
Sale 2: Rp 1,750,000
Sale 3: Rp 3,335,000
Sale 4: Rp   500,000
Sale 5: Rp 1,000,000
------------------------
Total: Rp 15,585,000 ✅ Sesuai dengan report
```

---

### Verifikasi Total Items

**1. Get Report:**

```
GET http://localhost:8080/api/v1/reports/summary
```

Response: `"total_items": 10`

**2. List All Items:**

```
GET http://localhost:8080/api/v1/items?page=1
```

**3. Check Total Records:**

- Pagination response: `"total_records": 10` ✅ Sesuai

---

## Real-Time Testing Scenario

### Scenario 1: Impact of Creating New Item

**1. Get Initial Report:**

```
GET http://localhost:8080/api/v1/reports/summary
Authorization: Bearer <token>
```

Response: `"total_items": 10`

**2. Create New Item:**

```
POST http://localhost:8080/api/v1/items
Authorization: Bearer <token>
Body: {
    "sku": "ELC-004",
    "name": "Monitor LED 24 inch",
    "category_id": 1,
    "rack_id": 1,
    "stock": 25,
    "minimum_stock": 5,
    "price": 2500000
}
```

**3. Get Report Again:**

```
GET http://localhost:8080/api/v1/reports/summary
Authorization: Bearer <token>
```

**Expected:** `"total_items": 11` ✅ Bertambah 1

---

### Scenario 2: Impact of Creating Sale

**1. Get Initial Report:**

```
GET http://localhost:8080/api/v1/reports/summary
Authorization: Bearer <token>
```

Response:

```json
{
  "total_sales": 5,
  "total_revenue": 15585000
}
```

**2. Create New Sale:**

```
POST http://localhost:8080/api/v1/sales
Authorization: Bearer <token>
Body: {
    "items": [
        {
            "item_id": 2,
            "quantity": 5
        }
    ]
}
```

Response:

```json
{
  "total_amount": 1250000
}
```

**3. Get Report Again:**

```
GET http://localhost:8080/api/v1/reports/summary
Authorization: Bearer <token>
```

**Expected:**

```json
{
  "total_sales": 6,
  "total_revenue": 16835000
}
```

- Total sales: 5 + 1 = 6 ✅
- Total revenue: 15585000 + 1250000 = 16835000 ✅

---

### Scenario 3: Impact of Deleting Sale (Soft Delete)

**1. Get Initial Report:**

```
GET http://localhost:8080/api/v1/reports/summary
Authorization: Bearer <token>
```

Response:

```json
{
  "total_sales": 6,
  "total_revenue": 16835000
}
```

**2. Delete Sale #6:**

```
DELETE http://localhost:8080/api/v1/sales/6
Authorization: Bearer <token>
```

**3. Get Report Again:**

```
GET http://localhost:8080/api/v1/reports/summary
Authorization: Bearer <token>
```

**Expected:**

```json
{
  "total_sales": 5,
  "total_revenue": 15585000
}
```

- Sale yang di-soft delete tidak masuk hitungan ✅
- Revenue kembali ke nilai sebelumnya ✅

---

### Scenario 4: Impact of Deactivating User

**1. Get Initial Report:**

```
GET http://localhost:8080/api/v1/reports/summary
Authorization: Bearer <token>
```

Response: `"active_users": 3`

**2. Deactivate Staff User:**

```
PUT http://localhost:8080/api/v1/users/3
Authorization: Bearer <token>
Body: {
    "name": "Staff User",
    "email": "staff@inventory.com",
    "password": "password123",
    "role_id": 3,
    "is_active": false
}
```

**3. Get Report Again:**

```
GET http://localhost:8080/api/v1/reports/summary
Authorization: Bearer <token>
```

**Expected:** `"active_users": 2` ✅ Berkurang 1

---

## Complete Testing Flow

### Full Report Testing Workflow

```bash
# 1. Login sebagai Super Admin
POST /api/v1/auth/login
Body: {"email": "superadmin@inventory.com", "password": "password123"}
→ Simpan token

# 2. Get initial report snapshot
GET /api/v1/reports/summary
Header: Authorization: Bearer <token>
→ Catat semua values

# 3. Create new item (testing total_items)
POST /api/v1/items
Header: Authorization: Bearer <token>
Body: {new item data}
→ total_items should increase

# 4. Get report - verify total_items increased
GET /api/v1/reports/summary
Header: Authorization: Bearer <token>
→ Verify total_items += 1

# 5. Create new sale (testing total_sales & total_revenue)
POST /api/v1/sales
Header: Authorization: Bearer <token>
Body: {items: [...]}
→ Note the total_amount

# 6. Get report - verify sales metrics
GET /api/v1/reports/summary
Header: Authorization: Bearer <token>
→ Verify total_sales += 1
→ Verify total_revenue += total_amount

# 7. Create item with low stock (testing low_stock_items)
POST /api/v1/items
Header: Authorization: Bearer <token>
Body: {stock: 2, minimum_stock: 5}
→ Should be counted as low stock

# 8. Get report - verify low_stock_items increased
GET /api/v1/reports/summary
Header: Authorization: Bearer <token>
→ Verify low_stock_items += 1

# 9. Delete sale (testing soft delete exclusion)
DELETE /api/v1/sales/{sale_id}
Header: Authorization: Bearer <token>

# 10. Get report - verify deleted sale not counted
GET /api/v1/reports/summary
Header: Authorization: Bearer <token>
→ Verify total_sales -= 1
→ Verify total_revenue -= deleted_amount
```

---

## Error Responses

### 401 Unauthorized - Missing Token

```json
{
  "status": false,
  "message": "unauthorized",
  "data": null
}
```

### 401 Unauthorized - Invalid Token

```json
{
  "status": false,
  "message": "invalid token",
  "data": null
}
```

### 401 Unauthorized - Expired Token

```json
{
  "status": false,
  "message": "token expired",
  "data": null
}
```

### 403 Forbidden - No Permission (If authorization applied)

```json
{
  "status": false,
  "message": "forbidden: insufficient permissions",
  "data": null
}
```

### 500 Internal Server Error

```json
{
  "status": false,
  "message": "failed to fetch report: <error details>",
  "data": null
}
```

---

## Role-Based Access Control (Future Implementation)

Sesuai requirements, ideal authorization untuk report:

| Role        | Access                                         |
| ----------- | ---------------------------------------------- |
| Super Admin | ✅ Full access ke semua report                 |
| Admin       | ✅ Full access ke semua report                 |
| Staff       | ❌ Tidak boleh akses revenue report (opsional) |

### Implementation Example (Future):

```go
// middleware/permission.go
func ReportPermission(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        user := r.Context().Value("user").(*model.User)

        // Only super_admin and admin can access reports
        if user.RoleName != "super_admin" && user.RoleName != "admin" {
            utils.ResponseBadRequest(w, http.StatusForbidden,
                "forbidden: only admin and super_admin can access reports", nil)
            return
        }

        next.ServeHTTP(w, r)
    })
}
```

---

## Notes

1. **Authentication Required:** Endpoint report memerlukan Bearer token yang valid
2. **Real-Time Data:** Report selalu menampilkan data terkini dari database
3. **Soft Delete Awareness:** Total sales dan revenue hanya menghitung yang `deleted_at IS NULL`
4. **Zero Handling:** Revenue menggunakan `COALESCE(SUM(total_amount), 0)` untuk handle case tanpa sales
5. **Performance:** Query aggregation efisien dengan index pada:
   - `items(stock)`
   - `sales(deleted_at)`
   - `users(is_active)`
6. **Consistency:** Data report konsisten dengan data yang ditampilkan di list endpoints
7. **Low Stock Definition:** Items dengan `stock < minimum_stock`
8. **Active Users:** Only users dengan `is_active = true`

---

## Postman Environment Variables

Setup environment di Postman untuk testing lebih mudah:

**Variables:**

```
base_url: http://localhost:8080/api/v1
token: (auto-set dari login response)
super_admin_token: a1b2c3d4-e5f6-47a8-b9c0-d1e2f3a4b5c6
admin_token: b2c3d4e5-f6a7-48b9-c0d1-e2f3a4b5c6d7
staff_token: c3d4e5f6-a7b8-49c0-d1e2-f3a4b5c6d7e8
```

**Usage in Postman:**

```
URL: {{base_url}}/reports/summary
Header: Authorization: Bearer {{token}}
```

---

## Checklist Testing Report

- [ ] Login sebagai super admin
- [ ] Get report summary
- [ ] Verify all fields present dan valid
- [ ] Verify total_items matches actual count
- [ ] Verify low_stock_items dengan manual check
- [ ] Verify total_sales matches sales count
- [ ] Verify total_revenue calculation correct
- [ ] Create new item, verify total_items increases
- [ ] Create new sale, verify metrics update
- [ ] Delete sale, verify metrics decrease
- [ ] Test with admin role
- [ ] Test with staff role (verify access control)
- [ ] Test without token (401 error)
- [ ] Test with invalid token (401 error)
- [ ] Verify real-time data accuracy
