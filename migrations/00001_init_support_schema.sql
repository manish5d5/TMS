-- +goose Up
-- +goose StatementBegin

-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- =====================
-- Users table
-- =====================
CREATE TABLE users (
    id              BIGSERIAL PRIMARY KEY,
    uuid            UUID NOT NULL DEFAULT gen_random_uuid(),

    name            VARCHAR(100) NOT NULL,
    email           VARCHAR(150) UNIQUE NOT NULL,
    password_hash   TEXT NOT NULL,

    created_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uq_users_uuid UNIQUE (uuid)
);

CREATE INDEX idx_users_email ON users(email);

-- =====================
-- Ticket enums
-- =====================
CREATE TYPE ticket_status AS ENUM ('open', 'closed');
CREATE TYPE ticket_priority AS ENUM ('low', 'medium', 'high');

-- =====================
-- Tickets table
-- =====================
CREATE TABLE tickets (
    id              BIGSERIAL PRIMARY KEY,
    uuid            UUID NOT NULL DEFAULT gen_random_uuid(),

    title           VARCHAR(200) NOT NULL,
    description     TEXT NOT NULL,
    status          ticket_status DEFAULT 'open',
    priority        ticket_priority DEFAULT 'medium',

    created_by_id   BIGINT NOT NULL,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uq_tickets_uuid UNIQUE (uuid),

    CONSTRAINT fk_ticket_user
        FOREIGN KEY (created_by_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_tickets_created_by ON tickets(created_by_id);
CREATE INDEX idx_tickets_status ON tickets(status);
CREATE INDEX idx_tickets_priority ON tickets(priority);

-- =====================
-- Comments table
-- =====================
CREATE TABLE comments (
    id              BIGSERIAL PRIMARY KEY,
    uuid            UUID NOT NULL DEFAULT gen_random_uuid(),

    ticket_id       BIGINT NOT NULL,
    created_by_id   BIGINT NOT NULL,

    comment_text    TEXT NOT NULL,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uq_comments_uuid UNIQUE (uuid),

    CONSTRAINT fk_comment_ticket
        FOREIGN KEY (ticket_id)
        REFERENCES tickets(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_comment_user
        FOREIGN KEY (created_by_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_comments_ticket_id ON comments(ticket_id);
CREATE INDEX idx_comments_created_by ON comments(created_by_id);

-- +goose StatementEnd



-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS tickets;

DROP TYPE IF EXISTS ticket_status;
DROP TYPE IF EXISTS ticket_priority;

DROP TABLE IF EXISTS users;

-- +goose StatementEnd
