CREATE TABLE users(
    auth_id BIGINT PRIMARY KEY,
    user_id UUID UNIQUE NOT NULL,

    name VARCHAR NOT NULL,
    bio TEXT,
    sex VARCHAR, -- "male" and "female"
    birthdate DATE,
    phone VARCHAR,
    profile_picture VARCHAR,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
