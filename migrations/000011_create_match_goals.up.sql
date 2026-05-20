CREATE TABLE match_goals (
    id SERIAL PRIMARY KEY,
    match_result_id INT NOT NULL REFERENCES match_results(id),
    player_id INT NOT NULL REFERENCES players(id),
    minute INT NOT NULL,
    is_own_goal BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
