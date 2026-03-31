CREATE TABLE IF NOT EXISTS applicants (
    id TEXT PRIMARY KEY,

    team_name TEXT NOT NULL,
    age_group TEXT NOT NULL,
    institution TEXT,
    city_or_region TEXT NOT NULL,

    members JSONB NOT NULL DEFAULT '[]',
    responsible_person JSONB NOT NULL,

    how_heard_about TEXT,
    comments TEXT,

    confirm_truthful BOOLEAN NOT NULL DEFAULT false,
    confirm_rules BOOLEAN NOT NULL DEFAULT false,
    confirm_media BOOLEAN NOT NULL DEFAULT false,

    applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    status TEXT NOT NULL DEFAULT 'pending'
);