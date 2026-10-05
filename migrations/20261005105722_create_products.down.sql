CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    category_id UUID NOT NULL REFERENCES categories (id),
    created_by UUID NOT NULL REFERENCES users (id),
    name VARCHAR(200) NOT NULL,
    slug VARCHAR(220) NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT '',
    price_minor BIGINT NOT NULL CHECK (
        price_minor BETWEEN 0 AND 1000000000
    ),
    currency VARCHAR(3) NOT NULL DEFAULT 'BDT' CHECK (currency = 'BDT'),
    status VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (
        status IN ('draft', 'active', 'archived')
    ),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
