-- migrations/003_session_security.sql

CREATE TYPE session_status AS ENUM (
    'active',
    'expired',
    'revoked'
);

CREATE TABLE agent_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    agent_id UUID NOT NULL
        REFERENCES agents(id)
        ON DELETE CASCADE,

    name VARCHAR(120) NOT NULL,

    key_hash TEXT NOT NULL UNIQUE,

    status session_status NOT NULL DEFAULT 'active',

    spending_limit NUMERIC(30,9) NOT NULL DEFAULT 0
        CHECK (spending_limit >= 0),

    spent NUMERIC(30,9) NOT NULL DEFAULT 0
        CHECK (spent >= 0),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    expires_at TIMESTAMPTZ NOT NULL,

    revoked_at TIMESTAMPTZ
);

CREATE INDEX idx_agent_sessions_agent_id
ON agent_sessions(agent_id);

CREATE INDEX idx_agent_sessions_key_hash
ON agent_sessions(key_hash);

CREATE INDEX idx_agent_sessions_status
ON agent_sessions(status);


CREATE TABLE session_allowed_assets (
    session_id UUID NOT NULL
        REFERENCES agent_sessions(id)
        ON DELETE CASCADE,

    asset asset_type NOT NULL,

    PRIMARY KEY(session_id, asset)
);


CREATE TABLE session_allowed_networks (
    session_id UUID NOT NULL
        REFERENCES agent_sessions(id)
        ON DELETE CASCADE,

    network network_type NOT NULL,

    PRIMARY KEY(session_id, network)
);


CREATE TABLE session_allowed_services (
    session_id UUID NOT NULL
        REFERENCES agent_sessions(id)
        ON DELETE CASCADE,

    service_identifier VARCHAR(255) NOT NULL,

    PRIMARY KEY(session_id, service_identifier)
);


CREATE TABLE security_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    agent_id UUID
        REFERENCES agents(id)
        ON DELETE SET NULL,

    session_id UUID
        REFERENCES agent_sessions(id)
        ON DELETE SET NULL,

    type VARCHAR(64) NOT NULL,

    title VARCHAR(255) NOT NULL,

    description TEXT NOT NULL,

    amount NUMERIC(30,9),

    result VARCHAR(32) NOT NULL
        CHECK (
            result IN (
                'blocked',
                'approved',
                'warning',
                'info'
            )
        ),

    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_security_events_agent_id
ON security_events(agent_id);

CREATE INDEX idx_security_events_session_id
ON security_events(session_id);

CREATE INDEX idx_security_events_created_at
ON security_events(created_at DESC);
