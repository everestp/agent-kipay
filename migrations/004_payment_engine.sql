-- migrations/004_payment_engine.sql

CREATE TYPE payment_status AS ENUM (
    'pending',
    'approved',
    'blocked',
    'submitted',
    'completed',
    'failed'
);

CREATE TABLE api_services (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    name VARCHAR(150) NOT NULL,

    category VARCHAR(100),

    endpoint TEXT NOT NULL,

    price_per_request NUMERIC(30,9) NOT NULL
        CHECK (price_per_request >= 0),

    asset asset_type NOT NULL DEFAULT 'USDC',

    network network_type NOT NULL DEFAULT 'solana-devnet',

    description TEXT,

    provider_reputation NUMERIC(5,2) DEFAULT 0,

    active BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_api_services_endpoint
ON api_services(endpoint);


CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    agent_id UUID NOT NULL
        REFERENCES agents(id)
        ON DELETE RESTRICT,

    session_id UUID
        REFERENCES agent_sessions(id)
        ON DELETE SET NULL,

    service_id UUID
        REFERENCES api_services(id)
        ON DELETE SET NULL,

    amount NUMERIC(30,9) NOT NULL
        CHECK (amount > 0),

    asset asset_type NOT NULL,

    network network_type NOT NULL,

    protocol VARCHAR(32) NOT NULL DEFAULT 'x402',

    status payment_status NOT NULL DEFAULT 'pending',

    policy_decision VARCHAR(64),

    policy_reason TEXT,

    idempotency_key VARCHAR(255) NOT NULL UNIQUE,

    payment_nonce VARCHAR(255),

    tx_hash TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    approved_at TIMESTAMPTZ,

    submitted_at TIMESTAMPTZ,

    completed_at TIMESTAMPTZ,

    failed_at TIMESTAMPTZ
);

CREATE INDEX idx_payments_agent_id
ON payments(agent_id);

CREATE INDEX idx_payments_session_id
ON payments(session_id);

CREATE INDEX idx_payments_service_id
ON payments(service_id);

CREATE INDEX idx_payments_status
ON payments(status);

CREATE INDEX idx_payments_created_at
ON payments(created_at DESC);

CREATE INDEX idx_payments_tx_hash
ON payments(tx_hash);
