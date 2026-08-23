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
CREATE TABLE IF NOT EXISTS api_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,

    name VARCHAR(100) NOT NULL,

    publishable_key VARCHAR(100) NOT NULL UNIQUE,

    secret_key_hash VARCHAR(128) NOT NULL,

    status VARCHAR(20) NOT NULL DEFAULT 'active'
        CHECK (status IN ('active', 'revoked')),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_api_keys_user_id
ON api_keys(user_id);

CREATE INDEX IF NOT EXISTS idx_api_keys_publishable_key
ON api_keys(publishable_key);
CREATE TABLE IF NOT EXISTS settlements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,

    payment_id UUID NOT NULL
        REFERENCES payments(id)
        ON DELETE RESTRICT,

    tx_hash VARCHAR(128) NOT NULL,

    network VARCHAR(30) NOT NULL
        CHECK (network IN ('solana', 'solana-devnet')),

    status VARCHAR(30) NOT NULL DEFAULT 'pending'
        CHECK (
            status IN (
                'pending',
                'confirmed',
                'failed'
            )
        ),

    amount NUMERIC(30, 9) NOT NULL DEFAULT 0,

    asset VARCHAR(20) NOT NULL,

    receiver VARCHAR(128) NOT NULL,

    block_number BIGINT NOT NULL DEFAULT 0,

    message TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    confirmed_at TIMESTAMPTZ,

    UNIQUE(network, tx_hash)
);

CREATE INDEX idx_settlements_user_id
ON settlements(user_id);

CREATE INDEX idx_settlements_payment_id
ON settlements(payment_id);

CREATE INDEX idx_settlements_status
ON settlements(status);

CREATE INDEX idx_settlements_tx_hash
ON settlements(tx_hash);
CREATE TABLE IF NOT EXISTS webhook_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    event_id VARCHAR(255) NOT NULL UNIQUE,

    event_type VARCHAR(100) NOT NULL,

    network VARCHAR(50) NOT NULL,

    signature TEXT,

    payload JSONB NOT NULL DEFAULT '{}'::jsonb,

    status VARCHAR(30) NOT NULL DEFAULT 'received',

    error_message TEXT,

    received_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    processed_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_webhook_events_event_type
ON webhook_events(event_type);

CREATE INDEX IF NOT EXISTS idx_webhook_events_network
ON webhook_events(network);

CREATE INDEX IF NOT EXISTS idx_webhook_events_status
ON webhook_events(status);

CREATE INDEX IF NOT EXISTS idx_webhook_events_created_at
ON webhook_events(created_at);
CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    payment_id UUID NOT NULL,

    agent_id UUID NOT NULL,

    session_id UUID,

    api_service_id UUID,

    type VARCHAR(30) NOT NULL,

    protocol VARCHAR(30) NOT NULL DEFAULT 'x402',

    network VARCHAR(50) NOT NULL,

    asset VARCHAR(20) NOT NULL,

    amount NUMERIC(30, 10) NOT NULL CHECK (amount > 0),

    sender_address TEXT NOT NULL,

    receiver_address TEXT NOT NULL,

    status VARCHAR(30) NOT NULL DEFAULT 'pending',

    tx_hash TEXT,

    block_number BIGINT,

    verification_status VARCHAR(30)
        NOT NULL DEFAULT 'pending',

    settlement_status VARCHAR(30)
        NOT NULL DEFAULT 'pending',

    created_at TIMESTAMPTZ
        NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMPTZ
        NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_transactions_payment_id
ON transactions(payment_id);

CREATE INDEX IF NOT EXISTS idx_transactions_agent_id
ON transactions(agent_id);

CREATE INDEX IF NOT EXISTS idx_transactions_session_id
ON transactions(session_id);

CREATE INDEX IF NOT EXISTS idx_transactions_tx_hash
ON transactions(tx_hash);

CREATE INDEX IF NOT EXISTS idx_transactions_status
ON transactions(status);

CREATE INDEX IF NOT EXISTS idx_transactions_created_at
ON transactions(created_at);
