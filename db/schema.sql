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

CREATE TABLE IF NOT EXISTS konti (
    key TEXT PRIMARY KEY,
    fullname TEXT NOT NULL,
    date_of_birth TEXT NOT NULL,
    role TEXT NOT NULL,
    id TEXT,
    password TEXT NOT NULL,
    username TEXT NOT NULL,
    email TEXT NOT NULL,
    phone_number TEXT NOT NULL,
    salt TEXT NOT NULL,
    educational_institution TEXT NOT NULL DEFAULT '',
    class_or_year TEXT NOT NULL DEFAULT '',
    pending_team_id TEXT NOT NULL DEFAULT '',
    registered BOOLEAN NOT NULL DEFAULT false
);

CREATE TABLE IF NOT EXISTS kontu_pieteikumi (
    fullname TEXT NOT NULL,
    date_of_birth TEXT NOT NULL,
    role TEXT NOT NULL,
    password TEXT NOT NULL,
    username TEXT NOT NULL,
    email TEXT NOT NULL,
    phone_number TEXT NOT NULL,
    team_name TEXT NOT NULL,
    key TEXT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS admini (
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    salt TEXT NOT NULL,
    superadmin BOOLEAN NOT NULL DEFAULT false,
    key TEXT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS komandas (
    key TEXT PRIMARY KEY,
    car_id TEXT NOT NULL,
    team_name TEXT NOT NULL,
    team_members TEXT[] NOT NULL DEFAULT '{}',
    age_group TEXT NOT NULL,
    institution TEXT,
    city_or_region TEXT NOT NULL,
    responsible_person JSONB NOT NULL DEFAULT '[]',
    avatar TEXT NOT NULL,
    car_data JSONB NOT NULL DEFAULT '{}'
);

CREATE TABLE IF NOT EXISTS mqtt_logs (
    topic TEXT NOT NULL,
    payload JSONB NOT NULL,
    received_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);