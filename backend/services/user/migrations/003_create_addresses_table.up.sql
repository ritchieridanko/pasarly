CREATE TABLE addresses(
    address_id BIGSERIAL PRIMARY KEY,
    auth_id BIGINT NOT NULL,

    recipient VARCHAR NOT NULL,
    phone VARCHAR NOT NULL,
    label VARCHAR NOT NULL, -- "Home", "Work", "School", etc.
    notes TEXT,
    is_primary BOOLEAN DEFAULT FALSE,

    country VARCHAR NOT NULL,
    subdivision_1 TEXT,
    subdivision_2 TEXT,
    subdivision_3 TEXT,
    subdivision_4 TEXT,
    street TEXT NOT NULL,
    postcode VARCHAR NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    location GEOGRAPHY(Point, 4326),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Optimize queries of user's addresses by auth_id
CREATE INDEX idx_addresses_auth_id ON addresses(auth_id);
