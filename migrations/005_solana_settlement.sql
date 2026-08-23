CREATE TYPE settlement_status AS ENUM (
    'pending',
    'processing',
    'confirmed',
    'failed'
);

CREATE TABLE settlements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    payment_id UUID NOT NULL UNIQUE
        REFERENCES payments(id)
        ON DELETE RESTRICT,

    tx_hash TEXT,

    network network_type NOT NULL,

    amount NUMERIC(30,9) NOT NULL,

    asset asset_type NOT NULL,

    sender_address TEXT NOT NULL,

    receiver_address TEXT NOT NULL,

    block_number BIGINT,

    status settlement_status NOT NULL DEFAULT 'pending',

    confirmations INTEGER NOT NULL DEFAULT 0,

    error_message TEXT,

    submitted_at TIMESTAMPTZ,

    confirmed_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_settlements_payment_id
ON settlements(payment_id);

CREATE INDEX idx_settlements_tx_hash
ON settlements(tx_hash);

CREATE INDEX idx_settlements_status
ON settlements(status);


CREATE TABLE ledger_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    agent_id UUID
        REFERENCES agents(id)
        ON DELETE RESTRICT,

    wallet_id UUID
        REFERENCES wallets(id)
        ON DELETE RESTRICT,

    asset asset_type NOT NULL,

    balance NUMERIC(30,9) NOT NULL DEFAULT 0,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(agent_id, asset)
);


CREATE TABLE ledger_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    payment_id UUID
        REFERENCES payments(id)
        ON DELETE SET NULL,

    settlement_id UUID
        REFERENCES settlements(id)
        ON DELETE SET NULL,

    account_id UUID NOT NULL
        REFERENCES ledger_accounts(id)
        ON DELETE RESTRICT,

    entry_type VARCHAR(32) NOT NULL,

    debit NUMERIC(30,9) NOT NULL DEFAULT 0,

    credit NUMERIC(30,9) NOT NULL DEFAULT 0,

    balance_after NUMERIC(30,9) NOT NULL,

    reference TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CHECK (
        debit >= 0
        AND credit >= 0
    ),

    CHECK (
        NOT (debit > 0 AND credit > 0)
    )
);

CREATE INDEX idx_ledger_entries_account
ON ledger_entries(account_id);

CREATE INDEX idx_ledger_entries_payment
ON ledger_entries(payment_id);

CREATE INDEX idx_ledger_entries_created
ON ledger_entries(created_at DESC);
