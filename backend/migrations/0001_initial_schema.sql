-- Migration 0001: Initial schema for StackGenome catalog (staging)
-- Applied with: wrangler d1 migrations apply DB --local

CREATE TABLE IF NOT EXISTS resources (
    id             TEXT PRIMARY KEY,         -- e.g. "tool:opentelemetry"
    type           TEXT NOT NULL,            -- "tool" | "library" | "service" | "platform"
    canonical_name TEXT NOT NULL,
    summary        TEXT NOT NULL,
    canonical_url  TEXT,
    ecosystem      TEXT NOT NULL DEFAULT '[]', -- JSON array of language/infra strings
    infra_targets  TEXT NOT NULL DEFAULT '[]', -- JSON array of infra strings
    license        TEXT,
    status         TEXT NOT NULL DEFAULT 'active', -- "active" | "deprecated" | "archived"
    verified_at    TEXT,
    updated_at     TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS technologies (
    id      TEXT PRIMARY KEY,
    kind    TEXT NOT NULL,  -- "language" | "framework" | "infra"
    name    TEXT NOT NULL,
    aliases TEXT NOT NULL DEFAULT '[]'  -- JSON array
);

CREATE TABLE IF NOT EXISTS resource_technologies (
    resource_id   TEXT NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
    technology_id TEXT NOT NULL REFERENCES technologies(id) ON DELETE CASCADE,
    relation      TEXT NOT NULL DEFAULT 'targets',
    PRIMARY KEY (resource_id, technology_id)
);

CREATE TABLE IF NOT EXISTS ranking_snapshots (
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    algorithm_version TEXT NOT NULL,
    weights_json      TEXT NOT NULL,
    catalog_version   TEXT NOT NULL,
    created_at        TEXT NOT NULL
);

-- Full-text search virtual table for resources
CREATE VIRTUAL TABLE IF NOT EXISTS resources_fts USING fts5(
    id UNINDEXED,
    canonical_name,
    summary,
    content='resources',
    content_rowid='rowid'
);

-- Triggers to keep FTS in sync
CREATE TRIGGER IF NOT EXISTS resources_ai AFTER INSERT ON resources BEGIN
    INSERT INTO resources_fts(rowid, id, canonical_name, summary)
    VALUES (new.rowid, new.id, new.canonical_name, new.summary);
END;

CREATE TRIGGER IF NOT EXISTS resources_ad AFTER DELETE ON resources BEGIN
    INSERT INTO resources_fts(resources_fts, rowid, id, canonical_name, summary)
    VALUES ('delete', old.rowid, old.id, old.canonical_name, old.summary);
END;

CREATE TRIGGER IF NOT EXISTS resources_au AFTER UPDATE ON resources BEGIN
    INSERT INTO resources_fts(resources_fts, rowid, id, canonical_name, summary)
    VALUES ('delete', old.rowid, old.id, old.canonical_name, old.summary);
    INSERT INTO resources_fts(rowid, id, canonical_name, summary)
    VALUES (new.rowid, new.id, new.canonical_name, new.summary);
END;

CREATE INDEX IF NOT EXISTS idx_resources_status ON resources(status);
CREATE INDEX IF NOT EXISTS idx_resources_type   ON resources(type);
CREATE INDEX IF NOT EXISTS idx_resource_tech    ON resource_technologies(resource_id);
CREATE INDEX IF NOT EXISTS idx_tech_kind        ON technologies(kind);

-- Seed the initial ranking snapshot (algorithm v1)
INSERT OR IGNORE INTO ranking_snapshots (algorithm_version, weights_json, catalog_version, created_at)
VALUES (
    'v1',
    '{"language_match":0.40,"infra_match":0.30,"framework_match":0.20,"novelty":0.10}',
    'staging-v1',
    datetime('now')
);
