CREATE TABLE matches (
    id SERIAL PRIMARY KEY,
    home_team_id INT NOT NULL REFERENCES teams(id),
    away_team_id INT NOT NULL REFERENCES teams(id),
    match_date DATE NOT NULL,
    match_time TIME NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'scheduled',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT different_teams CHECK (home_team_id != away_team_id)
);
