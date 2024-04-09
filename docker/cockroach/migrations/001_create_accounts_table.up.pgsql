CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE ROLE AS ENUM('admin', 'user');

CREATE TABLE IF NOT EXISTS
    talvi.accounts (
        id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid (),
        created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        name TEXT NOT NULL,
        role ROLE NOT NULL,
        email TEXT NOT NULL,
        provider TEXT NOT NULL,
        email_provider_hash TEXT NOT NULL UNIQUE
    );