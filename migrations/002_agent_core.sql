-- migrations/002_agent_core.sql

CREATE TYPE policy_decision AS ENUM (
    'allowed',
    'blocked',
    'approval_required'
);

CREATE TABLE agents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,

    wallet_id UUID NOT NULL
        REFERENCES wallets(id)
        ON DELETE RESTRICT,

    name VARCHAR(120) NOT NULL,

    description TEXT,

    status VARCHAR(32) NOT NULL DEFAULT 'idle'
        CHECK (status IN ('active', 'paused', 'error', 'idle')),

    asset asset_type NOT NULL DEFAULT 'USDC',

    network network_type NOT NULL DEFAULT 'solana-devnet',

    color VARCHAR(32),

    auto_payments BOOLEAN NOT NULL DEFAULT TRUE,

    last_active_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_agents_user_id
ON agents(user_id);

CREATE INDEX idx_agents_wallet_id
ON agents(wallet_id);

CREATE TABLE policies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    agent_id UUID NOT NULL
        REFERENCES agents(id)
        ON DELETE CASCADE,

    daily_limit NUMERIC(30,9) NOT NULL DEFAULT 0
        CHECK (daily_limit >= 0),

    per_transaction_limit NUMERIC(30,9) NOT NULL DEFAULT 0
        CHECK (per_transaction_limit >= 0),

    weekly_limit NUMERIC(30,9) NOT NULL DEFAULT 0
        CHECK (weekly_limit >= 0),

    require_approval_above NUMERIC(30,9) NOT NULL DEFAULT 0
        CHECK (require_approval_above >= 0),

    expiration_days INTEGER NOT NULL DEFAULT 30
        CHECK (expiration_days > 0),

    require_approval_new_merchants BOOLEAN NOT NULL DEFAULT TRUE,

    require_approval_new_assets BOOLEAN NOT NULL DEFAULT TRUE,

    block_unknown_apis BOOLEAN NOT NULL DEFAULT TRUE,

    auto_payments BOOLEAN NOT NULL DEFAULT TRUE,

    status VARCHAR(32) NOT NULL DEFAULT 'draft'
        CHECK (status IN ('active', 'review', 'draft', 'expired')),

    expires_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_policies_agent_id
ON policies(agent_id);

CREATE UNIQUE INDEX idx_active_policy_per_agent
ON policies(agent_id)
WHERE status = 'active';

CREATE TABLE policy_allowed_assets (
    policy_id UUID NOT NULL
        REFERENCES policies(id)
        ON DELETE CASCADE,

    asset asset_type NOT NULL,

    PRIMARY KEY(policy_id, asset)
);

CREATE TABLE policy_allowed_networks (
    policy_id UUID NOT NULL
        REFERENCES policies(id)
        ON DELETE CASCADE,

    network network_type NOT NULL,

    PRIMARY KEY(policy_id, network)
);

CREATE TABLE policy_allowed_merchants (
    policy_id UUID NOT NULL
        REFERENCES policies(id)
        ON DELETE CASCADE,

    merchant_identifier VARCHAR(255) NOT NULL,

    PRIMARY KEY(policy_id, merchant_identifier)
);

CREATE TABLE policy_blocked_services (
    policy_id UUID NOT NULL
        REFERENCES policies(id)
        ON DELETE CASCADE,

    service_identifier VARCHAR(255) NOT NULL,

    PRIMARY KEY(policy_id, service_identifier)
);

CREATE INDEX idx_policy_allowed_merchants_policy
ON policy_allowed_merchants(policy_id);

CREATE INDEX idx_policy_blocked_services_policy
ON policy_blocked_services(policy_id);

CREATE TRIGGER agents_updated_at
BEFORE UPDATE ON agents
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER policies_updated_at
BEFORE UPDATE ON policies
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();
