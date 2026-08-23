-- migrations/001_phase_one.sql

CREATE EXTENSION IF NOT EXISTS pgcrypto;

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

CREATE TYPE network_type AS ENUM (
    'solana',
    'solana-devnet'
);

CREATE TYPE asset_type AS ENUM (
    'USDC',
    'SOL',
    'USDT'
);

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    email VARCHAR(255) NOT NULL UNIQUE,

    password_hash TEXT NOT NULL,

    name VARCHAR(120) NOT NULL,

    status user_status NOT NULL DEFAULT 'active',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

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

CREATE TABLE wallet_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    wallet_id UUID NOT NULL
        REFERENCES wallets(id)
        ON DELETE CASCADE,

    asset asset_type NOT NULL,

    balance NUMERIC(30,9) NOT NULL DEFAULT 0,

    available_balance NUMERIC(30,9) NOT NULL DEFAULT 0,

    reserved_balance NUMERIC(30,9) NOT NULL DEFAULT 0,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(wallet_id, asset),

    CHECK (balance >= 0),

    CHECK (available_balance >= 0),

    CHECK (reserved_balance >= 0)
);

CREATE INDEX idx_wallets_user_id
ON wallets(user_id);

CREATE INDEX idx_wallet_accounts_wallet_id
ON wallet_accounts(wallet_id);

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
