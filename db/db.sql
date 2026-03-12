-- ============================================================================
-- Setokin Database Schema
-- PostgreSQL Database for F&B Inventory Management System
-- 
-- Features:
-- - FEFO (First Expired First Out) inventory management
-- - Batch tracking with expiry dates
-- - Stock in/out transaction logging
-- - User authentication with JWT
-- - Normalized to BCNF
-- - Schema separation (data/private/api)
-- ============================================================================

-- Enable UUID extension (PostgreSQL 13+)
-- For PostgreSQL 17+, use gen_random_uuid() which is built-in
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ============================================================================
-- SCHEMA SETUP
-- ============================================================================

-- Create schemas for separation of concerns
CREATE SCHEMA IF NOT EXISTS data;      -- Tables and data storage
CREATE SCHEMA IF NOT EXISTS private;   -- Internal functions, triggers
CREATE SCHEMA IF NOT EXISTS api;       -- Public API functions/procedures

-- ============================================================================
-- USERS & AUTHENTICATION
-- ============================================================================

-- Users table for authentication
CREATE TABLE data.users (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    email text UNIQUE NOT NULL,
    password_hash text NOT NULL,
    full_name text NOT NULL,
    role text NOT NULL DEFAULT 'staff',
    is_active boolean DEFAULT true,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    
    CONSTRAINT users_email_check CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'),
    CONSTRAINT users_role_check CHECK (role IN ('owner', 'manager', 'staff'))
);

-- Refresh tokens for JWT authentication
CREATE TABLE data.refresh_tokens (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES data.users(id) ON DELETE CASCADE,
    token_hash text NOT NULL,
    expires_at timestamptz NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    revoked_at timestamptz
);

-- ============================================================================
-- INVENTORY MASTER DATA
-- ============================================================================

-- Categories for organizing items
CREATE TABLE data.categories (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name text UNIQUE NOT NULL,
    description text,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

-- Units of measurement
CREATE TABLE data.units (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name text UNIQUE NOT NULL,
    abbreviation text UNIQUE NOT NULL,
    type text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    
    CONSTRAINT units_type_check CHECK (type IN ('weight', 'volume', 'count'))
);

-- Items (bahan baku)
CREATE TABLE data.items (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name text NOT NULL,
    category_id uuid NOT NULL REFERENCES data.categories(id) ON DELETE RESTRICT,
    unit_id uuid NOT NULL REFERENCES data.units(id) ON DELETE RESTRICT,
    minimum_stock numeric(10, 3) DEFAULT 0,
    description text,
    is_active boolean DEFAULT true,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    
    CONSTRAINT items_minimum_stock_check CHECK (minimum_stock >= 0)
);

-- ============================================================================
-- SUPPLIERS
-- ============================================================================

-- Suppliers for tracking purchase sources
CREATE TABLE data.suppliers (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name text NOT NULL,
    contact_person text,
    phone text,
    email text,
    address text,
    is_active boolean DEFAULT true,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    
    CONSTRAINT suppliers_email_check CHECK (email IS NULL OR email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$')
);

-- ============================================================================
-- BATCH & INVENTORY TRACKING
-- ============================================================================

-- Batches for FEFO tracking
CREATE TABLE data.batches (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    item_id uuid NOT NULL REFERENCES data.items(id) ON DELETE RESTRICT,
    batch_number text NOT NULL,
    initial_quantity numeric(10, 3) NOT NULL,
    remaining_quantity numeric(10, 3) NOT NULL,
    expiry_date date NOT NULL,
    is_depleted boolean DEFAULT false,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    
    CONSTRAINT batches_initial_quantity_check CHECK (initial_quantity > 0),
    CONSTRAINT batches_remaining_quantity_check CHECK (remaining_quantity >= 0),
    CONSTRAINT batches_quantity_logic_check CHECK (remaining_quantity <= initial_quantity),
    CONSTRAINT batches_batch_number_key UNIQUE (item_id, batch_number)
);

-- ============================================================================
-- STOCK TRANSACTIONS
-- ============================================================================

-- Stock In transactions (pembelian)
CREATE TABLE data.stock_in (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    item_id uuid NOT NULL REFERENCES data.items(id) ON DELETE RESTRICT,
    batch_id uuid NOT NULL REFERENCES data.batches(id) ON DELETE RESTRICT,
    supplier_id uuid REFERENCES data.suppliers(id) ON DELETE SET NULL,
    quantity numeric(10, 3) NOT NULL,
    purchase_date date NOT NULL,
    purchase_price numeric(15, 2),
    notes text,
    created_by uuid NOT NULL REFERENCES data.users(id) ON DELETE RESTRICT,
    created_at timestamptz NOT NULL DEFAULT now(),
    
    CONSTRAINT stock_in_quantity_check CHECK (quantity > 0),
    CONSTRAINT stock_in_purchase_price_check CHECK (purchase_price IS NULL OR purchase_price >= 0)
);

-- Stock Out transactions (pemakaian)
CREATE TABLE data.stock_out (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    item_id uuid NOT NULL REFERENCES data.items(id) ON DELETE RESTRICT,
    quantity numeric(10, 3) NOT NULL,
    usage_date date NOT NULL,
    notes text,
    created_by uuid NOT NULL REFERENCES data.users(id) ON DELETE RESTRICT,
    created_at timestamptz NOT NULL DEFAULT now(),
    
    CONSTRAINT stock_out_quantity_check CHECK (quantity > 0)
);

-- Stock Out Details (FEFO deduction tracking)
-- This table tracks which batches were used for each stock out transaction
CREATE TABLE data.stock_out_details (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    stock_out_id uuid NOT NULL REFERENCES data.stock_out(id) ON DELETE CASCADE,
    batch_id uuid NOT NULL REFERENCES data.batches(id) ON DELETE RESTRICT,
    quantity_used numeric(10, 3) NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    
    CONSTRAINT stock_out_details_quantity_used_check CHECK (quantity_used > 0)
);

-- ============================================================================
-- INDEXES FOR PERFORMANCE
-- ============================================================================

-- Users indexes
CREATE INDEX users_email_idx ON data.users(email);
CREATE INDEX users_role_idx ON data.users(role);
CREATE INDEX users_is_active_idx ON data.users(is_active);

-- Refresh tokens indexes
CREATE INDEX refresh_tokens_user_id_idx ON data.refresh_tokens(user_id);
CREATE INDEX refresh_tokens_expires_at_idx ON data.refresh_tokens(expires_at);
CREATE INDEX refresh_tokens_token_hash_idx ON data.refresh_tokens(token_hash);

-- Items indexes
CREATE INDEX items_category_id_idx ON data.items(category_id);
CREATE INDEX items_unit_id_idx ON data.items(unit_id);
CREATE INDEX items_name_idx ON data.items(name);
CREATE INDEX items_is_active_idx ON data.items(is_active);

-- Batches indexes (critical for FEFO)
CREATE INDEX batches_item_id_idx ON data.batches(item_id);
CREATE INDEX batches_expiry_date_idx ON data.batches(expiry_date);
CREATE INDEX batches_is_depleted_idx ON data.batches(is_depleted);
CREATE INDEX batches_item_expiry_idx ON data.batches(item_id, expiry_date, is_depleted); -- Composite for FEFO queries

-- Stock In indexes
CREATE INDEX stock_in_item_id_idx ON data.stock_in(item_id);
CREATE INDEX stock_in_batch_id_idx ON data.stock_in(batch_id);
CREATE INDEX stock_in_supplier_id_idx ON data.stock_in(supplier_id);
CREATE INDEX stock_in_purchase_date_idx ON data.stock_in(purchase_date);
CREATE INDEX stock_in_created_by_idx ON data.stock_in(created_by);
CREATE INDEX stock_in_created_at_idx ON data.stock_in(created_at);

-- Stock Out indexes
CREATE INDEX stock_out_item_id_idx ON data.stock_out(item_id);
CREATE INDEX stock_out_usage_date_idx ON data.stock_out(usage_date);
CREATE INDEX stock_out_created_by_idx ON data.stock_out(created_by);
CREATE INDEX stock_out_created_at_idx ON data.stock_out(created_at);

-- Stock Out Details indexes
CREATE INDEX stock_out_details_stock_out_id_idx ON data.stock_out_details(stock_out_id);
CREATE INDEX stock_out_details_batch_id_idx ON data.stock_out_details(batch_id);

-- Suppliers indexes
CREATE INDEX suppliers_name_idx ON data.suppliers(name);
CREATE INDEX suppliers_is_active_idx ON data.suppliers(is_active);

-- ============================================================================
-- VIEWS FOR COMMON QUERIES
-- ============================================================================

-- Current inventory view (aggregated stock per item)
CREATE OR REPLACE VIEW api.v_current_inventory AS
SELECT 
    i.id AS item_id,
    i.name AS item_name,
    c.name AS category_name,
    u.abbreviation AS unit,
    COALESCE(SUM(b.remaining_quantity), 0) AS total_stock,
    i.minimum_stock,
    CASE 
        WHEN COALESCE(SUM(b.remaining_quantity), 0) <= i.minimum_stock THEN true
        ELSE false
    END AS is_low_stock,
    COUNT(b.id) FILTER (WHERE b.is_depleted = false) AS active_batches
FROM data.items i
LEFT JOIN data.categories c ON i.category_id = c.id
LEFT JOIN data.units u ON i.unit_id = u.id
LEFT JOIN data.batches b ON i.id = b.item_id AND b.is_depleted = false
WHERE i.is_active = true
GROUP BY i.id, i.name, c.name, u.abbreviation, i.minimum_stock;

-- Expiring soon view (batches expiring within 3 days)
CREATE OR REPLACE VIEW api.v_expiring_soon AS
SELECT 
    b.id AS batch_id,
    i.id AS item_id,
    i.name AS item_name,
    c.name AS category_name,
    b.batch_number,
    b.remaining_quantity,
    u.abbreviation AS unit,
    b.expiry_date,
    b.expiry_date - CURRENT_DATE AS days_until_expiry
FROM data.batches b
JOIN data.items i ON b.item_id = i.id
JOIN data.categories c ON i.category_id = c.id
JOIN data.units u ON i.unit_id = u.id
WHERE b.is_depleted = false
  AND b.expiry_date <= CURRENT_DATE + INTERVAL '3 days'
  AND b.expiry_date >= CURRENT_DATE
ORDER BY b.expiry_date ASC;

-- Batch details view (all batches with item info)
CREATE OR REPLACE VIEW api.v_batch_details AS
SELECT 
    b.id AS batch_id,
    b.batch_number,
    i.id AS item_id,
    i.name AS item_name,
    c.name AS category_name,
    b.initial_quantity,
    b.remaining_quantity,
    u.abbreviation AS unit,
    b.expiry_date,
    b.is_depleted,
    b.created_at,
    CASE 
        WHEN b.expiry_date < CURRENT_DATE THEN 'expired'
        WHEN b.expiry_date <= CURRENT_DATE + INTERVAL '3 days' THEN 'expiring_soon'
        ELSE 'good'
    END AS status
FROM data.batches b
JOIN data.items i ON b.item_id = i.id
JOIN data.categories c ON i.category_id = c.id
JOIN data.units u ON i.unit_id = u.id
ORDER BY i.name, b.expiry_date;

-- ============================================================================
-- SEED DATA
-- ============================================================================

-- Insert default units
INSERT INTO data.units (name, abbreviation, type) VALUES
    ('Kilogram', 'kg', 'weight'),
    ('Gram', 'g', 'weight'),
    ('Liter', 'L', 'volume'),
    ('Milliliter', 'ml', 'volume'),
    ('Pieces', 'pcs', 'count'),
    ('Pack', 'pack', 'count'),
    ('Box', 'box', 'count')
ON CONFLICT (abbreviation) DO NOTHING;

-- Insert default categories
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

-- ============================================================================
-- COMMENTS FOR DOCUMENTATION
-- ============================================================================

COMMENT ON SCHEMA data IS 'Data storage layer - tables and indexes';
COMMENT ON SCHEMA private IS 'Internal functions and triggers - not accessible to applications';
COMMENT ON SCHEMA api IS 'Public API layer - functions, procedures, and views for application access';

COMMENT ON TABLE data.users IS 'User accounts for authentication and authorization';
COMMENT ON TABLE data.refresh_tokens IS 'JWT refresh tokens for maintaining user sessions';
COMMENT ON TABLE data.categories IS 'Item categories for organizing inventory';
COMMENT ON TABLE data.units IS 'Units of measurement for items';
COMMENT ON TABLE data.items IS 'Master data for inventory items (bahan baku)';
COMMENT ON TABLE data.suppliers IS 'Supplier information for purchase tracking';
COMMENT ON TABLE data.batches IS 'Batch tracking for FEFO inventory management';
COMMENT ON TABLE data.stock_in IS 'Stock in transactions (purchases)';
COMMENT ON TABLE data.stock_out IS 'Stock out transactions (usage)';
COMMENT ON TABLE data.stock_out_details IS 'Detailed batch deductions for each stock out transaction (FEFO tracking)';

COMMENT ON COLUMN data.batches.batch_number IS 'Unique batch identifier, auto-generated';
COMMENT ON COLUMN data.batches.is_depleted IS 'Automatically set to true when remaining_quantity reaches 0';
COMMENT ON COLUMN data.items.minimum_stock IS 'Threshold for low stock alerts';
COMMENT ON COLUMN data.stock_out_details.quantity_used IS 'Quantity deducted from specific batch (FEFO logic)';

-- ============================================================================
-- END OF SCHEMA
-- ============================================================================

-- ============================================================================
-- FUNCTIONS & TRIGGERS (private schema)
-- ============================================================================

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION private.set_updated_at()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$;

-- Apply updated_at trigger to relevant tables
CREATE TRIGGER users_bu_updated_trg 
    BEFORE UPDATE ON data.users
    FOR EACH ROW EXECUTE FUNCTION private.set_updated_at();

CREATE TRIGGER categories_bu_updated_trg 
    BEFORE UPDATE ON data.categories
    FOR EACH ROW EXECUTE FUNCTION private.set_updated_at();

CREATE TRIGGER items_bu_updated_trg 
    BEFORE UPDATE ON data.items
    FOR EACH ROW EXECUTE FUNCTION private.set_updated_at();

CREATE TRIGGER suppliers_bu_updated_trg 
    BEFORE UPDATE ON data.suppliers
    FOR EACH ROW EXECUTE FUNCTION private.set_updated_at();

CREATE TRIGGER batches_bu_updated_trg 
    BEFORE UPDATE ON data.batches
    FOR EACH ROW EXECUTE FUNCTION private.set_updated_at();

-- Function to auto-generate batch number
CREATE OR REPLACE FUNCTION private.generate_batch_number(in_item_id uuid)
RETURNS text
LANGUAGE plpgsql
AS $$
DECLARE
    l_count integer;
    l_batch_number text;
BEGIN
    SELECT COUNT(*) INTO l_count
    FROM data.batches
    WHERE item_id = in_item_id;
    
    l_batch_number := 'BATCH-' || to_char(CURRENT_DATE, 'YYYYMMDD') || '-' || lpad((l_count + 1)::text, 4, '0');
    
    RETURN l_batch_number;
END;
$$;

-- Function to mark batch as depleted when remaining_quantity reaches 0
CREATE OR REPLACE FUNCTION private.check_batch_depletion()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
    IF NEW.remaining_quantity = 0 AND OLD.remaining_quantity > 0 THEN
        NEW.is_depleted = true;
    END IF;
    RETURN NEW;
END;
$$;

CREATE TRIGGER batches_bu_depletion_trg 
    BEFORE UPDATE ON data.batches
    FOR EACH ROW EXECUTE FUNCTION private.check_batch_depletion();
