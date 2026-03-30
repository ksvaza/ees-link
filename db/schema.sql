CREATE TABLE IF NOT EXISTS applicants (
    ID          TEXT        PRIMARY KEY,
    team_name   TEXT        NOT NULL,
    school      TEXT        NOT NULL,
    members     INTEGER     NOT NULL,
    supervisor  TEXT        NOT NULL,
    applied_at  DATE        NOT NULL DEFAULT CURRENT_DATE,  
    status      TEXT        NOT NULL DEFAULT 'pending'
);