CREATE TABLE IF NOT EXISTS
    talvi.twofactor (
        id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid (),
        parent_account_hash TEXT NOT NULL REFERENCES talvi.accounts (email_provider_hash) ON DELETE CASCADE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        secret TEXT NOT NULL,
        enabled BOOLEAN DEFAULT FALSE
    );