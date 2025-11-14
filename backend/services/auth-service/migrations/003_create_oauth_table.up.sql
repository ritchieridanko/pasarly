CREATE TABLE oauth(
    oauth_id BIGSERIAL PRIMARY KEY,
    auth_id BIGINT NOT NULL,
    
    provider VARCHAR NOT NULL,
    provider_uid VARCHAR NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    FOREIGN KEY (auth_id) REFERENCES auth(auth_id)
);

-- Enforce uniqueness of auth_id for active (not soft-deleted) records
CREATE UNIQUE INDEX idx_oauth_unique_auth_id ON oauth(auth_id) WHERE deleted_at IS NULL;

-- Enforce uniqueness of provider_uid per provider for active (not soft-deleted) records
CREATE UNIQUE INDEX idx_oauth_unique_provider_provider_uid ON oauth(provider, provider_uid) WHERE deleted_at IS NULL;
