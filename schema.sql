CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- =========================================================
-- ENUMS
-- =========================================================

CREATE TYPE user_status AS ENUM (
    'active',
    'suspended',
    'deleted'
);

CREATE TYPE wallet_status AS ENUM (
    'active',
    'frozen',
    'closed'
);

CREATE TYPE agent_status AS ENUM (
    'active',
    'paused',
    'error',
    'idle'
);

CREATE TYPE policy_status AS ENUM (
    'active',
    'review',
    'draft',
    'expired'
);

CREATE TYPE session_status AS ENUM (
    'active',
    'expired',
    'revoked'
);

CREATE TYPE payment_status AS ENUM (
    'pending',
    'completed',
    'failed',
    'blocked'
);

CREATE TYPE payment_type AS ENUM (
    'incoming',
    'outgoing',
    'agent',
    'manual',
    'x402',
    'failed'
);

CREATE TYPE protocol_type AS ENUM (
    'x402'
);

CREATE TYPE network_type AS ENUM (
    'solana',
    'solana-devnet'
);

CREATE TYPE asset_type AS ENUM (
    'USDC',
    'SOL',
    'USDT'
);

CREATE TYPE settlement_status AS ENUM (
    'pending',
    'processing',
    'confirmed',
    'failed'
);

CREATE TYPE ledger_entry_type AS ENUM (
    'debit',
    'credit'
);

CREATE TYPE api_service_status AS ENUM (
    'active',
    'disabled'
);

CREATE TYPE api_key_status AS ENUM (
    'active',
    'revoked'
);

-- =========================================================
-- USERS
-- =========================================================

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    email VARCHAR(255) NOT NULL UNIQUE,

    password_hash TEXT NOT NULL,

    name VARCHAR(120) NOT NULL,

    status user_status NOT NULL DEFAULT 'active',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email
ON users(email);

-- =========================================================
-- WALLETS
-- =========================================================

CREATE TABLE wallets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,

    name VARCHAR(120) NOT NULL DEFAULT 'Main Wallet',

    address VARCHAR(128) NOT NULL UNIQUE,

    network network_type NOT NULL,

    status wallet_status NOT NULL DEFAULT 'active',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_wallets_user_id
ON wallets(user_id);

-- =========================================================
-- WALLET ACCOUNTS
-- One wallet can contain multiple assets
-- =========================================================

CREATE TABLE wallet_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    wallet_id UUID NOT NULL
        REFERENCES wallets(id)
        ON DELETE CASCADE,

    asset asset_type NOT NULL,

    balance NUMERIC(30,9) NOT NULL DEFAULT 0
        CHECK (balance >= 0),

    available_balance NUMERIC(30,9) NOT NULL DEFAULT 0
        CHECK (available_balance >= 0),

    reserved_balance NUMERIC(30,9) NOT NULL DEFAULT 0
        CHECK (reserved_balance >= 0),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(wallet_id, asset)
);

CREATE INDEX idx_wallet_accounts_wallet_id
ON wallet_accounts(wallet_id);

-- =========================================================
-- AGENTS
-- =========================================================

CREATE TABLE agents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,

    wallet_id UUID NOT NULL
        REFERENCES wallets(id),

    name VARCHAR(120) NOT NULL,

    description TEXT,

    status agent_status NOT NULL DEFAULT 'idle',

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

-- =========================================================
-- POLICIES
-- =========================================================

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

    status policy_status NOT NULL DEFAULT 'draft',

    expires_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_policies_agent_id
ON policies(agent_id);

CREATE UNIQUE INDEX idx_one_active_policy_per_agent
ON policies(agent_id)
WHERE status = 'active';

-- =========================================================
-- POLICY ALLOWED ASSETS
-- =========================================================

CREATE TABLE policy_allowed_assets (
    policy_id UUID NOT NULL
        REFERENCES policies(id)
        ON DELETE CASCADE,

    asset asset_type NOT NULL,

    PRIMARY KEY(policy_id, asset)
);

-- =========================================================
-- POLICY ALLOWED NETWORKS
-- =========================================================

CREATE TABLE policy_allowed_networks (
    policy_id UUID NOT NULL
        REFERENCES policies(id)
        ON DELETE CASCADE,

    network network_type NOT NULL,

    PRIMARY KEY(policy_id, network)
);

-- =========================================================
-- POLICY ALLOWED MERCHANTS
-- =========================================================

CREATE TABLE policy_allowed_merchants (
    policy_id UUID NOT NULL
        REFERENCES policies(id)
        ON DELETE CASCADE,

    merchant_identifier VARCHAR(255) NOT NULL,

    PRIMARY KEY(policy_id, merchant_identifier)
);

-- =========================================================
-- POLICY BLOCKED SERVICES
-- =========================================================

CREATE TABLE policy_blocked_services (
    policy_id UUID NOT NULL
        REFERENCES policies(id)
        ON DELETE CASCADE,

    service_identifier VARCHAR(255) NOT NULL,

    PRIMARY KEY(policy_id, service_identifier)
);

-- =========================================================
-- SESSIONS
-- Temporary scoped authority for agents
-- =========================================================

CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    agent_id UUID NOT NULL
        REFERENCES agents(id)
        ON DELETE CASCADE,

    status session_status NOT NULL DEFAULT 'active',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    expires_at TIMESTAMPTZ NOT NULL,

    spending_limit NUMERIC(30,9) NOT NULL
        CHECK (spending_limit >= 0),

    spent_amount NUMERIC(30,9) NOT NULL DEFAULT 0
        CHECK (spent_amount >= 0),

    key_hash TEXT NOT NULL,

    revoked_at TIMESTAMPTZ
);

CREATE INDEX idx_sessions_agent_id
ON sessions(agent_id);

CREATE INDEX idx_sessions_status
ON sessions(status);

-- =========================================================
-- SESSION ALLOWED ASSETS
-- =========================================================

CREATE TABLE session_allowed_assets (
    session_id UUID NOT NULL
        REFERENCES sessions(id)
        ON DELETE CASCADE,

    asset asset_type NOT NULL,

    PRIMARY KEY(session_id, asset)
);

-- =========================================================
-- SESSION ALLOWED NETWORKS
-- =========================================================

CREATE TABLE session_allowed_networks (
    session_id UUID NOT NULL
        REFERENCES sessions(id)
        ON DELETE CASCADE,

    network network_type NOT NULL,

    PRIMARY KEY(session_id, network)
);

-- =========================================================
-- SESSION ALLOWED SERVICES
-- =========================================================

CREATE TABLE session_allowed_services (
    session_id UUID NOT NULL
        REFERENCES sessions(id)
        ON DELETE CASCADE,

    service_identifier VARCHAR(255) NOT NULL,

    PRIMARY KEY(session_id, service_identifier)
);

-- =========================================================
-- API SERVICES
-- =========================================================

CREATE TABLE api_services (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    provider_user_id UUID
        REFERENCES users(id)
        ON DELETE SET NULL,

    name VARCHAR(160) NOT NULL,

    category VARCHAR(100) NOT NULL,

    description TEXT,

    endpoint TEXT NOT NULL,

    price_per_request NUMERIC(30,9) NOT NULL
        CHECK (price_per_request >= 0),

    asset asset_type NOT NULL DEFAULT 'USDC',

    network network_type NOT NULL DEFAULT 'solana-devnet',

    protocol protocol_type NOT NULL DEFAULT 'x402',

    latency_ms INTEGER NOT NULL DEFAULT 0
        CHECK (latency_ms >= 0),

    success_rate NUMERIC(5,2) NOT NULL DEFAULT 100
        CHECK (
            success_rate >= 0
            AND success_rate <= 100
        ),

    provider_reputation NUMERIC(5,2) NOT NULL DEFAULT 0
        CHECK (
            provider_reputation >= 0
            AND provider_reputation <= 100
        ),

    requests_30d BIGINT NOT NULL DEFAULT 0,

    revenue_30d NUMERIC(30,9) NOT NULL DEFAULT 0,

    status api_service_status NOT NULL DEFAULT 'active',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_api_services_provider_user_id
ON api_services(provider_user_id);

CREATE INDEX idx_api_services_status
ON api_services(status);

-- =========================================================
-- API KEYS
-- =========================================================

CREATE TABLE api_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,

    name VARCHAR(120) NOT NULL,

    publishable_key VARCHAR(255) NOT NULL UNIQUE,

    secret_key_hash TEXT NOT NULL,

    status api_key_status NOT NULL DEFAULT 'active',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    revoked_at TIMESTAMPTZ
);

CREATE INDEX idx_api_keys_user_id
ON api_keys(user_id);

-- =========================================================
-- PAYMENTS
-- x402 payment intent/execution
-- =========================================================

CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    agent_id UUID
        REFERENCES agents(id)
        ON DELETE SET NULL,

    session_id UUID
        REFERENCES sessions(id)
        ON DELETE SET NULL,

    api_service_id UUID
        REFERENCES api_services(id)
        ON DELETE SET NULL,

    wallet_id UUID
        REFERENCES wallets(id)
        ON DELETE SET NULL,

    amount NUMERIC(30,9) NOT NULL
        CHECK (amount > 0),

    asset asset_type NOT NULL,

    network network_type NOT NULL,

    protocol protocol_type NOT NULL DEFAULT 'x402',

    status payment_status NOT NULL DEFAULT 'pending',

    type payment_type NOT NULL DEFAULT 'x402',

    idempotency_key VARCHAR(255) NOT NULL UNIQUE,

    policy_decision VARCHAR(32),

    verification_status VARCHAR(32),

    settlement_status settlement_status,

    latency_ms INTEGER,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    completed_at TIMESTAMPTZ
);

CREATE INDEX idx_payments_agent_id
ON payments(agent_id);

CREATE INDEX idx_payments_session_id
ON payments(session_id);

CREATE INDEX idx_payments_api_service_id
ON payments(api_service_id);

CREATE INDEX idx_payments_status
ON payments(status);

CREATE INDEX idx_payments_created_at
ON payments(created_at DESC);

-- =========================================================
-- TRANSACTIONS
-- Financial transaction record
-- =========================================================

CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    payment_id UUID UNIQUE
        REFERENCES payments(id)
        ON DELETE SET NULL,

    agent_id UUID
        REFERENCES agents(id)
        ON DELETE SET NULL,

    wallet_id UUID
        REFERENCES wallets(id)
        ON DELETE SET NULL,

    service_id UUID
        REFERENCES api_services(id)
        ON DELETE SET NULL,

    amount NUMERIC(30,9) NOT NULL
        CHECK (amount > 0),

    asset asset_type NOT NULL,

    network network_type NOT NULL,

    status payment_status NOT NULL DEFAULT 'pending',

    type payment_type NOT NULL,

    protocol protocol_type NOT NULL DEFAULT 'x402',

    policy_decision TEXT,

    verification_status TEXT,

    settlement_status settlement_status,

    tx_hash VARCHAR(255),

    block_number BIGINT,

    sender_address VARCHAR(128),

    receiver_address VARCHAR(128),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    settled_at TIMESTAMPTZ
);

CREATE INDEX idx_transactions_payment_id
ON transactions(payment_id);

CREATE INDEX idx_transactions_agent_id
ON transactions(agent_id);

CREATE INDEX idx_transactions_wallet_id
ON transactions(wallet_id);

CREATE INDEX idx_transactions_service_id
ON transactions(service_id);

CREATE INDEX idx_transactions_created_at
ON transactions(created_at DESC);

CREATE INDEX idx_transactions_status
ON transactions(status);

CREATE INDEX idx_transactions_tx_hash
ON transactions(tx_hash);

-- =========================================================
-- LEDGER ACCOUNTS
-- =========================================================

CREATE TABLE ledger_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    wallet_id UUID
        REFERENCES wallets(id)
        ON DELETE CASCADE,

    agent_id UUID
        REFERENCES agents(id)
        ON DELETE CASCADE,

    name VARCHAR(160) NOT NULL,

    currency asset_type NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CHECK (
        wallet_id IS NOT NULL
        OR agent_id IS NOT NULL
    )
);

CREATE INDEX idx_ledger_accounts_wallet_id
ON ledger_accounts(wallet_id);

CREATE INDEX idx_ledger_accounts_agent_id
ON ledger_accounts(agent_id);

-- =========================================================
-- LEDGER TRANSACTIONS
-- One financial event
-- =========================================================

CREATE TABLE ledger_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    payment_id UUID
        REFERENCES payments(id)
        ON DELETE SET NULL,

    transaction_id UUID
        REFERENCES transactions(id)
        ON DELETE SET NULL,

    reference VARCHAR(255),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_ledger_transactions_payment_id
ON ledger_transactions(payment_id);

CREATE INDEX idx_ledger_transactions_transaction_id
ON ledger_transactions(transaction_id);

-- =========================================================
-- LEDGER ENTRIES
-- DOUBLE ENTRY ACCOUNTING
-- =========================================================

CREATE TABLE ledger_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    ledger_transaction_id UUID NOT NULL
        REFERENCES ledger_transactions(id)
        ON DELETE CASCADE,

    ledger_account_id UUID NOT NULL
        REFERENCES ledger_accounts(id),

    entry_type ledger_entry_type NOT NULL,

    amount NUMERIC(30,9) NOT NULL
        CHECK (amount > 0),

    asset asset_type NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_ledger_entries_account_id
ON ledger_entries(ledger_account_id);

CREATE INDEX idx_ledger_entries_transaction_id
ON ledger_entries(ledger_transaction_id);

-- =========================================================
-- SETTLEMENTS
-- Actual Solana settlement
-- =========================================================

CREATE TABLE settlements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    payment_id UUID UNIQUE
        REFERENCES payments(id)
        ON DELETE SET NULL,

    transaction_id UUID UNIQUE
        REFERENCES transactions(id)
        ON DELETE SET NULL,

    network network_type NOT NULL,

    asset asset_type NOT NULL,

    amount NUMERIC(30,9) NOT NULL
        CHECK (amount > 0),

    sender_address VARCHAR(128),

    receiver_address VARCHAR(128),

    tx_hash VARCHAR(255) UNIQUE,

    block_number BIGINT,

    status settlement_status NOT NULL DEFAULT 'pending',

    error_message TEXT,

    attempts INTEGER NOT NULL DEFAULT 0,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    confirmed_at TIMESTAMPTZ
);

CREATE INDEX idx_settlements_status
ON settlements(status);

CREATE INDEX idx_settlements_tx_hash
ON settlements(tx_hash);

-- =========================================================
-- APPROVALS
-- =========================================================

CREATE TABLE approvals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    agent_id UUID NOT NULL
        REFERENCES agents(id)
        ON DELETE CASCADE,

    payment_id UUID
        REFERENCES payments(id)
        ON DELETE CASCADE,

    requested_amount NUMERIC(30,9) NOT NULL
        CHECK (requested_amount > 0),

    asset asset_type NOT NULL,

    reason TEXT,

    status VARCHAR(32) NOT NULL DEFAULT 'pending',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    resolved_at TIMESTAMPTZ
);

CREATE INDEX idx_approvals_agent_id
ON approvals(agent_id);

CREATE INDEX idx_approvals_payment_id
ON approvals(payment_id);

CREATE INDEX idx_approvals_status
ON approvals(status);

-- =========================================================
-- SECURITY EVENTS
-- =========================================================

CREATE TABLE security_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    agent_id UUID
        REFERENCES agents(id)
        ON DELETE SET NULL,

    payment_id UUID
        REFERENCES payments(id)
        ON DELETE SET NULL,

    type VARCHAR(64) NOT NULL,

    title VARCHAR(255) NOT NULL,

    description TEXT,

    amount NUMERIC(30,9),

    result VARCHAR(32) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_security_events_agent_id
ON security_events(agent_id);

CREATE INDEX idx_security_events_payment_id
ON security_events(payment_id);

CREATE INDEX idx_security_events_created_at
ON security_events(created_at DESC);

-- =========================================================
-- ACTIVITY EVENTS
-- Dashboard activity feed
-- =========================================================

CREATE TABLE activity_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID
        REFERENCES users(id)
        ON DELETE CASCADE,

    agent_id UUID
        REFERENCES agents(id)
        ON DELETE SET NULL,

    category VARCHAR(64) NOT NULL,

    type VARCHAR(64) NOT NULL,

    title VARCHAR(255) NOT NULL,

    description TEXT,

    amount NUMERIC(30,9),

    severity VARCHAR(32) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_activity_events_user_id
ON activity_events(user_id);

CREATE INDEX idx_activity_events_agent_id
ON activity_events(agent_id);

CREATE INDEX idx_activity_events_created_at
ON activity_events(created_at DESC);

-- =========================================================
-- WEBHOOKS
-- =========================================================

CREATE TABLE webhooks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID
        REFERENCES users(id)
        ON DELETE CASCADE,

    url TEXT NOT NULL,

    secret_hash TEXT NOT NULL,

    events JSONB NOT NULL DEFAULT '[]'::jsonb,

    active BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_webhooks_user_id
ON webhooks(user_id);

CREATE INDEX idx_webhooks_active
ON webhooks(active);

-- =========================================================
-- WEBHOOK DELIVERIES
-- =========================================================

CREATE TABLE webhook_deliveries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    webhook_id UUID NOT NULL
        REFERENCES webhooks(id)
        ON DELETE CASCADE,

    event_type VARCHAR(100) NOT NULL,

    payload JSONB NOT NULL,

    status VARCHAR(32) NOT NULL DEFAULT 'pending',

    attempts INTEGER NOT NULL DEFAULT 0,

    response_code INTEGER,

    response_body TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    delivered_at TIMESTAMPTZ
);

CREATE INDEX idx_webhook_deliveries_webhook_id
ON webhook_deliveries(webhook_id);

CREATE INDEX idx_webhook_deliveries_status
ON webhook_deliveries(status);

CREATE INDEX idx_webhook_deliveries_created_at
ON webhook_deliveries(created_at DESC);

-- =========================================================
-- UPDATED_AT TRIGGER
-- =========================================================

CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER wallets_updated_at
BEFORE UPDATE ON wallets
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER wallet_accounts_updated_at
BEFORE UPDATE ON wallet_accounts
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER agents_updated_at
BEFORE UPDATE ON agents
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER policies_updated_at
BEFORE UPDATE ON policies
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER api_services_updated_at
BEFORE UPDATE ON api_services
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER transactions_updated_at
BEFORE UPDATE ON transactions
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER webhooks_updated_at
BEFORE UPDATE ON webhooks
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();


ALTER TABLE sessions
ADD COLUMN name TEXT NOT NULL;
