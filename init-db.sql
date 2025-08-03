-- Tạo database cho user service
CREATE DATABASE user_service;

-- Tạo database cho product service
CREATE DATABASE product_service;

-- Tạo database cho order service
CREATE DATABASE order_service;

-- Cấp quyền cho user postgres
GRANT ALL PRIVILEGES ON DATABASE user_service TO postgres;
GRANT ALL PRIVILEGES ON DATABASE product_service TO postgres;
GRANT ALL PRIVILEGES ON DATABASE order_service TO postgres;

-- ========================================
-- USER SERVICE DATA
-- ========================================
\c user_service;

-- Tạo bảng roles
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Tạo bảng users
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    sdt TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Tạo constraint unique cho sdt
ALTER TABLE users ADD CONSTRAINT uni_users_sdt UNIQUE (sdt);

-- Tạo bảng accounts
CREATE TABLE IF NOT EXISTS accounts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    role_id INTEGER REFERENCES roles(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Thêm dữ liệu mẫu cho roles
INSERT INTO roles (id, name, created_at, updated_at) VALUES 
    (1, 'ADMIN', '2025-07-27 22:45:00.000 +0700', '2025-07-27 22:45:00.000 +0700'),
    (2, 'USER', '2025-07-27 22:45:00.000 +0700', '2025-07-27 22:45:00.000 +0700')
ON CONFLICT (id) DO NOTHING;

-- Thêm dữ liệu mẫu cho users
INSERT INTO users (id, name, sdt, created_at, updated_at, deleted_at) VALUES 
    (1, 'Admin', '0373286662', '2025-07-27 22:45:00.000 +0700', '2025-07-27 22:45:00.000 +0700', NULL),
    (2, 'Mạnh Dũng', '0383727282', '2025-07-27 22:45:00.000 +0700', '2025-07-27 22:45:00.000 +0700', NULL),
    (3, 'Ngọc Ngà', '0393837373', '2025-07-27 22:45:00.000 +0700', '2025-07-28 14:30:03.150 +0700', NULL),
    (4, 'Ngọc Báu', '0123456789', '2025-07-28 13:36:03.240 +0700', '2025-07-28 14:13:28.266 +0700', '2025-07-28 14:14:11.096 +0700')
ON CONFLICT (id) DO NOTHING;

-- Thêm dữ liệu mẫu cho accounts
INSERT INTO accounts (id, user_id, role_id, created_at, updated_at) VALUES 
    (1, 1, 1, '2025-07-27 22:45:00.000 +0700', '2025-07-27 22:45:00.000 +0700'),
    (2, 2, 2, '2025-07-27 22:45:00.000 +0700', '2025-07-27 22:45:00.000 +0700'),
    (3, 3, 2, '2025-07-27 22:45:00.000 +0700', '2025-07-27 22:45:00.000 +0700'),
    (4, 4, 2, '2025-07-28 13:36:03.249 +0700', '2025-07-28 13:36:03.249 +0700')
ON CONFLICT (id) DO NOTHING;

-- ========================================
-- PRODUCT SERVICE DATA
-- ========================================
\c product_service;

-- Tạo bảng products
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    image VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Tạo bảng inventories
CREATE TABLE IF NOT EXISTS inventories (
    id SERIAL PRIMARY KEY,
    product_id INTEGER REFERENCES products(id),
    quantity INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Thêm dữ liệu mẫu cho products
INSERT INTO products (id, name, description, price, image, created_at, updated_at, deleted_at) VALUES 
    (1, 'iPhone 14 Pro Max', 'Flagship Apple smartphone', 29990000, 'image.png', '2025-07-25 12:00:00.000 +0700', '2025-07-30 16:29:40.877 +0700', NULL),
    (2, 'Samsung Galaxy S24', 'Premium Android smartphone', 25990000, 'samsung.png', '2025-07-25 12:00:00.000 +0700', '2025-07-30 16:29:40.898 +0700', NULL),
    (3, 'Vivo V2025', 'Smartphone 2025', 17000000, 'vivo.png', '2025-07-29 12:56:38.838 +0700', '2025-07-29 12:56:38.838 +0700', NULL),
    (4, 'Appple Meo Meo', 'Smartphone 2025', 17000000, 'vivo.png', '2025-07-30 14:32:21.662 +0700', '2025-07-30 14:41:09.239 +0700', '2025-07-30 14:42:02.455 +0700'),
    (5, 'Huwai Pro', 'Smartphone 2025', 17000000, 'huwai.png', '2025-07-30 14:33:38.205 +0700', '2025-07-30 14:33:38.205 +0700', NULL)
ON CONFLICT (id) DO NOTHING;

-- Thêm dữ liệu mẫu cho inventories
INSERT INTO inventories (id, product_id, quantity, created_at, updated_at) VALUES 
    (1, 1, 90, '2025-07-25 12:05:00.000 +0700', '2025-08-01 09:52:52.614 +0700'),
    (2, 2, 98, '2025-07-25 12:15:00.000 +0700', '2025-07-30 21:11:25.949 +0700')
ON CONFLICT (id) DO NOTHING;

-- ========================================
-- ORDER SERVICE DATA
-- ========================================
\c order_service;

-- Tạo bảng orders
CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    total_price DECIMAL(15,2) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Tạo bảng order_details
CREATE TABLE IF NOT EXISTS order_details (
    id SERIAL PRIMARY KEY,
    order_id INTEGER REFERENCES orders(id),
    product_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL,
    unit_price DECIMAL(15,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Thêm dữ liệu mẫu cho orders
INSERT INTO orders (id, user_id, total_price, status, created_at, updated_at) VALUES 
    (1, 1, 29990000, 'pending', '2025-07-25 12:00:00.000 +0700', '2025-08-01 09:52:52.601 +0700'),
    (2, 2, 29990000, 'cancelled', '2025-07-25 12:10:00.000 +0700', '2025-07-25 12:10:00.000 +0700'),
    (3, 1, 81970000, 'pending', '2025-07-30 16:29:40.901 +0700', '2025-07-30 16:29:40.901 +0700'),
    (4, 1, 81970000, 'pending', '2025-07-30 21:11:25.925 +0700', '2025-07-30 21:11:25.925 +0700'),
    (5, 1, 449850000, 'pending', '2025-07-31 10:14:06.540 +0700', '2025-07-31 10:14:06.540 +0700')
ON CONFLICT (id) DO NOTHING;

-- Thêm dữ liệu mẫu cho order_details
INSERT INTO order_details (id, order_id, product_id, quantity, unit_price, created_at, updated_at) VALUES 
    (1, 1, 1, 2, 29990000, '2025-07-25 12:10:00.000 +0700', '2025-07-25 12:10:00.000 +0700'),
    (2, 2, 2, 1, 29990000, '2025-07-25 12:10:00.000 +0700', '2025-07-25 12:10:00.000 +0700'),
    (3, 3, 1, 1, 29990000, '2025-07-30 16:29:40.909 +0700', '2025-07-30 16:29:40.909 +0700'),
    (4, 3, 2, 2, 25990000, '2025-07-30 16:29:40.909 +0700', '2025-07-30 16:29:40.909 +0700'),
    (5, 4, 1, 1, 29990000, '2025-07-30 21:11:25.932 +0700', '2025-07-30 21:11:25.932 +0700'),
    (6, 4, 2, 2, 25990000, '2025-07-30 21:11:25.932 +0700', '2025-07-30 21:11:25.932 +0700'),
    (7, 5, 1, 15, 29990000, '2025-07-31 10:14:06.549 +0700', '2025-07-31 10:14:06.549 +0700')
ON CONFLICT (id) DO NOTHING; 