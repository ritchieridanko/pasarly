CREATE TABLE sessions(
    session_id BIGSERIAL PRIMARY KEY,
    auth_id BIGINT NOT NULL,

    parent_id BIGINT,
    token VARCHAR UNIQUE NOT NULL,
    user_agent TEXT NOT NULL,
    ip_address TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,

    FOREIGN KEY (auth_id) REFERENCES auth(auth_id) ON DELETE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES sessions(session_id) ON DELETE CASCADE
);

-- Optimize queries by auth_id
CREATE INDEX idx_sessions_auth_id ON sessions(auth_id);

-- Optimize queries by parent_id if parent_id is not null
CREATE INDEX idx_sessions_parent_id ON sessions(parent_id) WHERE parent_id IS NOT NULL;

-- Optimize queries by token for active (not revoked) records
CREATE INDEX idx_sessions_active ON sessions(token) WHERE revoked_at IS NULL;
