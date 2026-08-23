CREATE TYPE transaction_status AS ENUM (
    'pending',
    'completed',
    'failed',
    'blocked'
);

CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    payment_id UUID UNIQUE
        REFERENCES payments(id)
        ON DELETE SET NULL,

    settlement_id UUID UNIQUE
        REFERENCES settlements(id)
        ON DELETE SET NULL,

    agent_id UUID NOT NULL
        REFERENCES agents(id)
        ON DELETE RESTRICT,

    service_id UUID
        REFERENCES api_services(id)
        ON DELETE SET NULL,

    amount NUMERIC(30,9) NOT NULL,

    asset asset_type NOT NULL,

    network network_type NOT NULL,

    protocol VARCHAR(32) NOT NULL DEFAULT 'x402',

    status transaction_status NOT NULL DEFAULT 'pending',

    tx_hash TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    completed_at TIMESTAMPTZ
);

CREATE INDEX idx_transactions_agent
ON transactions(agent_id);

CREATE INDEX idx_transactions_status
ON transactions(status);

CREATE INDEX idx_transactions_created
ON transactions(created_at DESC);

CREATE INDEX idx_transactions_tx_hash
ON transactions(tx_hash);
CREATE TABLE security_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    agent_id UUID
        REFERENCES agents(id)
        ON DELETE SET NULL,

    type VARCHAR(64) NOT NULL,

    title VARCHAR(255) NOT NULL,

    description TEXT,

    amount NUMERIC(30,9),

    result VARCHAR(32) NOT NULL,

    ip_address INET,

    metadata JSONB,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_security_events_agent
ON security_events(agent_id);

CREATE INDEX idx_security_events_type
ON security_events(type);

CREATE INDEX idx_security_events_created
ON security_events(created_at DESC);
