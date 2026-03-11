-- ==========================================================================
-- Setokin Database Migration 001 — Initial Schema
-- ==========================================================================

-- UP
-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create schemas
CREATE SCHEMA IF NOT EXISTS data;
CREATE SCHEMA IF NOT EXISTS private;
CREATE SCHEMA IF NOT EXISTS api;

-- Users
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

-- Refresh tokens
CREATE TABLE data.refresh_tokens (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES data.users(id) ON DELETE CASCADE,
    token_hash text NOT NULL,
    expires_at timestamptz NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    revoked_at timestamptz
);

-- Categories
CREATE TABLE data.categories (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name text UNIQUE NOT NULL,
    description text,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

-- Units
CREATE TABLE data.units (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name text UNIQUE NOT NULL,
    abbreviation text UNIQUE NOT NULL,
    type text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT units_type_check CHECK (type IN ('weight', 'volume', 'count'))
);

-- Items
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

-- Suppliers
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

-- Batches
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

-- Stock In
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

-- Stock Out
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

-- Stock Out Details
CREATE TABLE data.stock_out_details (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    stock_out_id uuid NOT NULL REFERENCES data.stock_out(id) ON DELETE CASCADE,
    batch_id uuid NOT NULL REFERENCES data.batches(id) ON DELETE RESTRICT,
    quantity_used numeric(10, 3) NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT stock_out_details_quantity_used_check CHECK (quantity_used > 0)
);

-- Upload tracking
CREATE TABLE data.uploads (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES data.users(id) ON DELETE CASCADE,
    file_name text NOT NULL,
    file_key text NOT NULL,
    file_type text NOT NULL,
    file_size bigint NOT NULL,
    purpose text NOT NULL,
    status text NOT NULL DEFAULT 'pending',
    confirmed_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT uploads_status_check CHECK (status IN ('pending', 'confirmed', 'expired')),
    CONSTRAINT uploads_purpose_check CHECK (purpose IN ('stock_in_invoice', 'item_image', 'supplier_document'))
);

-- ==========================================================================
-- INDEXES
-- ==========================================================================
CREATE INDEX users_email_idx ON data.users(email);
CREATE INDEX users_role_idx ON data.users(role);
CREATE INDEX users_is_active_idx ON data.users(is_active);
CREATE INDEX refresh_tokens_user_id_idx ON data.refresh_tokens(user_id);
CREATE INDEX refresh_tokens_expires_at_idx ON data.refresh_tokens(expires_at);
CREATE INDEX refresh_tokens_token_hash_idx ON data.refresh_tokens(token_hash);
CREATE INDEX items_category_id_idx ON data.items(category_id);
CREATE INDEX items_unit_id_idx ON data.items(unit_id);
CREATE INDEX items_name_idx ON data.items(name);
CREATE INDEX items_is_active_idx ON data.items(is_active);
CREATE INDEX batches_item_id_idx ON data.batches(item_id);
CREATE INDEX batches_expiry_date_idx ON data.batches(expiry_date);
CREATE INDEX batches_is_depleted_idx ON data.batches(is_depleted);
CREATE INDEX batches_item_expiry_idx ON data.batches(item_id, expiry_date, is_depleted);
CREATE INDEX stock_in_item_id_idx ON data.stock_in(item_id);
CREATE INDEX stock_in_batch_id_idx ON data.stock_in(batch_id);
CREATE INDEX stock_in_supplier_id_idx ON data.stock_in(supplier_id);
CREATE INDEX stock_in_purchase_date_idx ON data.stock_in(purchase_date);
CREATE INDEX stock_in_created_by_idx ON data.stock_in(created_by);
CREATE INDEX stock_in_created_at_idx ON data.stock_in(created_at);
CREATE INDEX stock_out_item_id_idx ON data.stock_out(item_id);
CREATE INDEX stock_out_usage_date_idx ON data.stock_out(usage_date);
CREATE INDEX stock_out_created_by_idx ON data.stock_out(created_by);
CREATE INDEX stock_out_created_at_idx ON data.stock_out(created_at);
CREATE INDEX stock_out_details_stock_out_id_idx ON data.stock_out_details(stock_out_id);
CREATE INDEX stock_out_details_batch_id_idx ON data.stock_out_details(batch_id);
CREATE INDEX suppliers_name_idx ON data.suppliers(name);
CREATE INDEX suppliers_is_active_idx ON data.suppliers(is_active);
CREATE INDEX uploads_user_id_idx ON data.uploads(user_id);
CREATE INDEX uploads_status_idx ON data.uploads(status);

-- ==========================================================================
-- VIEWS
-- ==========================================================================
CREATE OR REPLACE VIEW api.v_current_inventory AS
SELECT
    i.id AS item_id, i.name AS item_name, c.name AS category_name,
    u.abbreviation AS unit, COALESCE(SUM(b.remaining_quantity), 0) AS total_stock,
    i.minimum_stock,
    CASE WHEN COALESCE(SUM(b.remaining_quantity), 0) <= i.minimum_stock THEN true ELSE false END AS is_low_stock,
    COUNT(b.id) FILTER (WHERE b.is_depleted = false) AS active_batches
FROM data.items i
LEFT JOIN data.categories c ON i.category_id = c.id
LEFT JOIN data.units u ON i.unit_id = u.id
LEFT JOIN data.batches b ON i.id = b.item_id AND b.is_depleted = false
WHERE i.is_active = true
GROUP BY i.id, i.name, c.name, u.abbreviation, i.minimum_stock;

CREATE OR REPLACE VIEW api.v_expiring_soon AS
SELECT
    b.id AS batch_id, i.id AS item_id, i.name AS item_name, c.name AS category_name,
    b.batch_number, b.remaining_quantity, u.abbreviation AS unit,
    b.expiry_date, b.expiry_date - CURRENT_DATE AS days_until_expiry
FROM data.batches b
JOIN data.items i ON b.item_id = i.id
JOIN data.categories c ON i.category_id = c.id
JOIN data.units u ON i.unit_id = u.id
WHERE b.is_depleted = false AND b.expiry_date <= CURRENT_DATE + INTERVAL '3 days' AND b.expiry_date >= CURRENT_DATE
ORDER BY b.expiry_date ASC;

CREATE OR REPLACE VIEW api.v_batch_details AS
SELECT
    b.id AS batch_id, b.batch_number, i.id AS item_id, i.name AS item_name,
    c.name AS category_name, b.initial_quantity, b.remaining_quantity,
    u.abbreviation AS unit, b.expiry_date, b.is_depleted, b.created_at,
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

-- ==========================================================================
-- FUNCTIONS & TRIGGERS (private schema)
-- ==========================================================================
CREATE OR REPLACE FUNCTION private.set_updated_at()
RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$;

CREATE TRIGGER users_bu_updated_trg BEFORE UPDATE ON data.users
    FOR EACH ROW EXECUTE FUNCTION private.set_updated_at();
CREATE TRIGGER categories_bu_updated_trg BEFORE UPDATE ON data.categories
    FOR EACH ROW EXECUTE FUNCTION private.set_updated_at();
CREATE TRIGGER items_bu_updated_trg BEFORE UPDATE ON data.items
    FOR EACH ROW EXECUTE FUNCTION private.set_updated_at();
CREATE TRIGGER suppliers_bu_updated_trg BEFORE UPDATE ON data.suppliers
    FOR EACH ROW EXECUTE FUNCTION private.set_updated_at();
CREATE TRIGGER batches_bu_updated_trg BEFORE UPDATE ON data.batches
    FOR EACH ROW EXECUTE FUNCTION private.set_updated_at();

CREATE OR REPLACE FUNCTION private.generate_batch_number(in_item_id uuid)
RETURNS text LANGUAGE plpgsql AS $$
DECLARE
    l_count integer;
    l_batch_number text;
BEGIN
    SELECT COUNT(*) INTO l_count FROM data.batches WHERE item_id = in_item_id;
    l_batch_number := 'BATCH-' || to_char(CURRENT_DATE, 'YYYYMMDD') || '-' || lpad((l_count + 1)::text, 4, '0');
    RETURN l_batch_number;
END;
$$;

CREATE OR REPLACE FUNCTION private.check_batch_depletion()
RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    IF NEW.remaining_quantity = 0 AND OLD.remaining_quantity > 0 THEN
        NEW.is_depleted = true;
    END IF;
    RETURN NEW;
END;
$$;

CREATE TRIGGER batches_bu_depletion_trg BEFORE UPDATE ON data.batches
    FOR EACH ROW EXECUTE FUNCTION private.check_batch_depletion();

-- SCHEMA COMMENTS
COMMENT ON SCHEMA data IS 'Data storage layer - tables and indexes';
COMMENT ON SCHEMA private IS 'Internal functions and triggers';
COMMENT ON SCHEMA api IS 'Public API layer - views for application access';

-- ==========================================================================
-- DOWN (rollback)
-- ==========================================================================
-- DROP TABLE IF EXISTS data.uploads CASCADE;
-- DROP TABLE IF EXISTS data.stock_out_details CASCADE;
-- DROP TABLE IF EXISTS data.stock_out CASCADE;
-- DROP TABLE IF EXISTS data.stock_in CASCADE;
-- DROP TABLE IF EXISTS data.batches CASCADE;
-- DROP TABLE IF EXISTS data.suppliers CASCADE;
-- DROP TABLE IF EXISTS data.items CASCADE;
-- DROP TABLE IF EXISTS data.units CASCADE;
-- DROP TABLE IF EXISTS data.categories CASCADE;
-- DROP TABLE IF EXISTS data.refresh_tokens CASCADE;
-- DROP TABLE IF EXISTS data.users CASCADE;
-- DROP SCHEMA IF EXISTS api CASCADE;
-- DROP SCHEMA IF EXISTS private CASCADE;
-- DROP SCHEMA IF EXISTS data CASCADE;
