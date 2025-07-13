-- Daggerheart Adversary Tracker Schema

-- Adversaries table
CREATE TABLE IF NOT EXISTS adversaries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    challenge_rating TEXT NOT NULL,
    size TEXT NOT NULL,
    armor_class INTEGER NOT NULL,
    hit_points INTEGER NOT NULL,
    speed TEXT NOT NULL,
    strength INTEGER NOT NULL,
    dexterity INTEGER NOT NULL,
    constitution INTEGER NOT NULL,
    intelligence INTEGER NOT NULL,
    wisdom INTEGER NOT NULL,
    charisma INTEGER NOT NULL,
    abilities TEXT,
    actions TEXT,
    reactions TEXT,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Encounters table
CREATE TABLE IF NOT EXISTS encounters (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Encounter_adversaries junction table
CREATE TABLE IF NOT EXISTS encounter_adversaries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    encounter_id INTEGER NOT NULL,
    adversary_id INTEGER NOT NULL,
    count INTEGER NOT NULL DEFAULT 1,
    FOREIGN KEY (encounter_id) REFERENCES encounters(id) ON DELETE CASCADE,
    FOREIGN KEY (adversary_id) REFERENCES adversaries(id) ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_adversaries_name ON adversaries(name);
CREATE INDEX IF NOT EXISTS idx_adversaries_type ON adversaries(type);
CREATE INDEX IF NOT EXISTS idx_adversaries_cr ON adversaries(challenge_rating);
CREATE INDEX IF NOT EXISTS idx_encounters_name ON encounters(name);
CREATE INDEX IF NOT EXISTS idx_encounter_adversaries_encounter_id ON encounter_adversaries(encounter_id);
CREATE INDEX IF NOT EXISTS idx_encounter_adversaries_adversary_id ON encounter_adversaries(adversary_id);
