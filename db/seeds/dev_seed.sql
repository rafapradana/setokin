-- ==========================================================================
-- Setokin Development Seed Data
-- ==========================================================================

-- Default units
INSERT INTO data.units (name, abbreviation, type) VALUES
    ('Kilogram', 'kg', 'weight'),
    ('Gram', 'g', 'weight'),
    ('Liter', 'L', 'volume'),
    ('Milliliter', 'ml', 'volume'),
    ('Pieces', 'pcs', 'count'),
    ('Pack', 'pack', 'count'),
    ('Box', 'box', 'count')
ON CONFLICT (abbreviation) DO NOTHING;

-- Default categories
INSERT INTO data.categories (name, description) VALUES
    ('Daging', 'Daging sapi, ayam, ikan, dll'),
    ('Sayuran', 'Sayuran segar'),
    ('Bumbu', 'Bumbu dapur dan rempah'),
    ('Dairy', 'Susu, keju, mentega, dll'),
    ('Tepung & Biji-bijian', 'Tepung, beras, pasta, dll'),
    ('Minyak & Lemak', 'Minyak goreng, margarin, dll'),
    ('Minuman', 'Minuman kemasan dan bahan minuman'),
    ('Lain-lain', 'Bahan lainnya')
ON CONFLICT (name) DO NOTHING;

-- Test admin user (password: Admin123!)
-- bcrypt hash of "Admin123!"
INSERT INTO data.users (email, password_hash, full_name, role) VALUES
    ('admin@setokin.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Admin Setokin', 'owner'),
    ('manager@setokin.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Manager Setokin', 'manager'),
    ('staff@setokin.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Staff Setokin', 'staff')
ON CONFLICT (email) DO NOTHING;
