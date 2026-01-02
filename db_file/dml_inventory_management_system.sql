-- Insert Roles (3 roles sesuai requirement: super admin, admin, staff)
INSERT INTO roles (name, created_at) VALUES
('super_admin', CURRENT_TIMESTAMP),
('admin', CURRENT_TIMESTAMP),
('staff', CURRENT_TIMESTAMP);

-- Insert Users (1 user per role = 3 users)
-- Password untuk semua user adalah: "password123" (hash: $2a$10$GOArleP8YiH7jncpDEWNQORUlH25w9v28MB9Bylfe2pjZlYPZUGYm)
INSERT INTO users (name, email, password_hash, role_id, is_active, created_at, updated_at) VALUES
('Super Admin User', 'superadmin@inventory.com', '$2a$10$GOArleP8YiH7jncpDEWNQORUlH25w9v28MB9Bylfe2pjZlYPZUGYm', 1, true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Admin User', 'admin@inventory.com', '$2a$10$GOArleP8YiH7jncpDEWNQORUlH25w9v28MB9Bylfe2pjZlYPZUGYm', 2, true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Staff User', 'staff@inventory.com', '$2a$10$GOArleP8YiH7jncpDEWNQORUlH25w9v28MB9Bylfe2pjZlYPZUGYm', 3, true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Insert Categories (5 categories)
INSERT INTO categories (name, description, created_at, updated_at) VALUES
('Electronics', 'Electronic devices and accessories', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Furniture', 'Office and home furniture', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Office Supplies', 'Stationery and office equipment', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Hardware', 'Tools and hardware equipment', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Clothing', 'Work uniforms and clothing', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Insert Warehouses (5 warehouses)
INSERT INTO warehouses (name, location, created_at, updated_at) VALUES
('Main Warehouse', 'Jl. Industri No. 1, Jakarta Utara', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('South Warehouse', 'Jl. Raya Selatan No. 45, Jakarta Selatan', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('East Warehouse', 'Jl. Timur Raya No. 88, Bekasi', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('West Warehouse', 'Jl. Barat No. 12, Tangerang', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Central Warehouse', 'Jl. Pusat No. 77, Jakarta Pusat', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Insert Racks (5 racks)
INSERT INTO racks (warehouse_id, code, description, created_at, updated_at) VALUES
(1, 'A-01', 'Electronics section row A', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(1, 'B-01', 'Furniture section row B', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(2, 'C-01', 'Office supplies section', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(3, 'D-01', 'Hardware section', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(4, 'E-01', 'Clothing section', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Insert Items (10 items untuk testing berbagai kondisi stock)
INSERT INTO items (sku, name, category_id, rack_id, stock, minimum_stock, price, created_at, updated_at) VALUES
('ELC-001', 'Laptop Dell Inspiron 15', 1, 1, 50, 10, 8500000.00, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('ELC-002', 'Mouse Wireless Logitech', 1, 1, 150, 20, 250000.00, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('ELC-003', 'Keyboard Mechanical', 1, 1, 3, 5, 750000.00, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP), -- Below minimum
('FRN-001', 'Office Chair Ergonomic', 2, 2, 30, 5, 1500000.00, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('FRN-002', 'Office Desk Standing', 2, 2, 2, 5, 2500000.00, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP), -- Below minimum
('OFS-001', 'Printer HP LaserJet', 3, 3, 25, 5, 3200000.00, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('OFS-002', 'Paper A4 Sidu (1 Rim)', 3, 3, 200, 50, 45000.00, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('OFS-003', 'Pen Set', 3, 3, 4, 10, 25000.00, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP), -- Below minimum
('HRD-001', 'Hammer Tool Set', 4, 4, 80, 15, 350000.00, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('CLO-001', 'Work Uniform Shirt', 5, 5, 100, 20, 125000.00, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Insert Sales (5 sales transactions)
INSERT INTO sales (user_id, total_amount, created_at, updated_at) VALUES
(3, 9000000.00, '2025-12-01 10:30:00+07', '2025-12-01 10:30:00+07'),
(3, 1750000.00, '2025-12-02 14:15:00+07', '2025-12-02 14:15:00+07'),
(3, 3335000.00, '2025-12-05 09:45:00+07', '2025-12-05 09:45:00+07'),
(2, 500000.00, '2025-12-10 16:20:00+07', '2025-12-10 16:20:00+07'),
(2, 1000000.00, '2025-12-15 11:00:00+07', '2025-12-15 11:00:00+07');

-- Insert Sale Items (detailed items for each sale)
INSERT INTO sale_items (sale_id, item_id, quantity, price_at_sale, subtotal) VALUES
-- Sale 1: Laptop + Mouse
(1, 1, 1, 8500000.00, 8500000.00),
(1, 2, 2, 250000.00, 500000.00),

-- Sale 2: Office Chair + Mouse
(2, 4, 1, 1500000.00, 1500000.00),
(2, 2, 1, 250000.00, 250000.00),

-- Sale 3: Printer + Paper
(3, 6, 1, 3200000.00, 3200000.00),
(3, 7, 3, 45000.00, 135000.00),

-- Sale 4: Mouse
(4, 2, 2, 250000.00, 500000.00),

-- Sale 5: Uniforms
(5, 10, 8, 125000.00, 1000000.00);

-- Insert Sessions (3 active sessions for testing - one for each role)
INSERT INTO sessions (user_id, token, expired_at, revoked_at, created_at) VALUES
(1, 'a1b2c3d4-e5f6-47a8-b9c0-d1e2f3a4b5c6', CURRENT_TIMESTAMP + INTERVAL '24 hours', NULL, CURRENT_TIMESTAMP),
(2, 'b2c3d4e5-f6a7-48b9-c0d1-e2f3a4b5c6d7', CURRENT_TIMESTAMP + INTERVAL '24 hours', NULL, CURRENT_TIMESTAMP),
(3, 'c3d4e5f6-a7b8-49c0-d1e2-f3a4b5c6d7e8', CURRENT_TIMESTAMP + INTERVAL '24 hours', NULL, CURRENT_TIMESTAMP);
