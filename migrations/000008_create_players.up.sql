CREATE TABLE players (
    id SERIAL PRIMARY KEY,
    team_id INT NOT NULL REFERENCES teams(id),
    name VARCHAR(255) NOT NULL,
    height_cm NUMERIC(5,2),
    weight_kg NUMERIC(5,2),
    position VARCHAR(20) NOT NULL,
    jersey_number INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    UNIQUE (team_id, jersey_number)
);
