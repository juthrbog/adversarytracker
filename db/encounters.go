package db

import (
	"context"
	"database/sql"
	"time"
)

// Encounter represents a combat encounter with adversaries
type Encounter struct {
	ID          int64
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Adversaries []*EncounterAdversary
}

// EncounterAdversary represents an adversary in an encounter with a count
type EncounterAdversary struct {
	ID          int64
	EncounterID int64
	AdversaryID int64
	Count       int
	Adversary   *Adversary
}

// GetAllEncounters retrieves all encounters from the database
func GetAllEncounters(ctx context.Context, db *sql.DB) ([]*Encounter, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM encounters
		ORDER BY name ASC
	`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var encounters []*Encounter
	for rows.Next() {
		enc := &Encounter{}
		err := rows.Scan(
			&enc.ID, &enc.Name, &enc.Description, &enc.CreatedAt, &enc.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		encounters = append(encounters, enc)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Load adversaries for each encounter
	for _, encounter := range encounters {
		adversaries, err := GetEncounterAdversaries(ctx, db, encounter.ID)
		if err != nil {
			return nil, err
		}
		encounter.Adversaries = adversaries
	}

	return encounters, nil
}

// GetEncounterByID retrieves a single encounter by ID
func GetEncounterByID(ctx context.Context, db *sql.DB, id int64) (*Encounter, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM encounters
		WHERE id = ?
	`

	enc := &Encounter{}
	err := db.QueryRowContext(ctx, query, id).Scan(
		&enc.ID, &enc.Name, &enc.Description, &enc.CreatedAt, &enc.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	// Load adversaries for the encounter
	adversaries, err := GetEncounterAdversaries(ctx, db, enc.ID)
	if err != nil {
		return nil, err
	}
	enc.Adversaries = adversaries

	return enc, nil
}

// GetEncounterAdversaries retrieves all adversaries for an encounter
func GetEncounterAdversaries(ctx context.Context, db *sql.DB, encounterID int64) ([]*EncounterAdversary, error) {
	query := `
		SELECT ea.id, ea.encounter_id, ea.adversary_id, ea.count,
		       a.id, a.name, a.type, a.challenge_rating, a.size, a.armor_class, a.hit_points, 
		       a.speed, a.strength, a.dexterity, a.constitution, a.intelligence, a.wisdom, 
		       a.charisma, a.abilities, a.actions, a.reactions, a.description, a.created_at, a.updated_at
		FROM encounter_adversaries ea
		JOIN adversaries a ON ea.adversary_id = a.id
		WHERE ea.encounter_id = ?
		ORDER BY a.name ASC
	`

	rows, err := db.QueryContext(ctx, query, encounterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var adversaries []*EncounterAdversary
	for rows.Next() {
		ea := &EncounterAdversary{
			Adversary: &Adversary{},
		}
		err := rows.Scan(
			&ea.ID, &ea.EncounterID, &ea.AdversaryID, &ea.Count,
			&ea.Adversary.ID, &ea.Adversary.Name, &ea.Adversary.Type, &ea.Adversary.ChallengeRating, 
			&ea.Adversary.Size, &ea.Adversary.ArmorClass, &ea.Adversary.HitPoints, &ea.Adversary.Speed, 
			&ea.Adversary.Strength, &ea.Adversary.Dexterity, &ea.Adversary.Constitution, 
			&ea.Adversary.Intelligence, &ea.Adversary.Wisdom, &ea.Adversary.Charisma, 
			&ea.Adversary.Abilities, &ea.Adversary.Actions, &ea.Adversary.Reactions, 
			&ea.Adversary.Description, &ea.Adversary.CreatedAt, &ea.Adversary.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		adversaries = append(adversaries, ea)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return adversaries, nil
}

// CreateEncounter inserts a new encounter into the database
func CreateEncounter(ctx context.Context, db *sql.DB, enc *Encounter) (int64, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// Insert encounter
	query := `
		INSERT INTO encounters (name, description)
		VALUES (?, ?)
	`

	result, err := tx.ExecContext(ctx, query, enc.Name, enc.Description)
	if err != nil {
		return 0, err
	}

	// Get the new encounter ID
	encounterID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Insert encounter adversaries if any
	if len(enc.Adversaries) > 0 {
		for _, ea := range enc.Adversaries {
			ea.EncounterID = encounterID
			_, err := AddAdversaryToEncounter(ctx, tx, ea)
			if err != nil {
				return 0, err
			}
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return encounterID, nil
}

// UpdateEncounter updates an existing encounter in the database
func UpdateEncounter(ctx context.Context, db *sql.DB, enc *Encounter) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update encounter
	query := `
		UPDATE encounters
		SET name = ?, description = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err = tx.ExecContext(ctx, query, enc.Name, enc.Description, enc.ID)
	if err != nil {
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// DeleteEncounter removes an encounter from the database
func DeleteEncounter(ctx context.Context, db *sql.DB, id int64) error {
	query := `DELETE FROM encounters WHERE id = ?`
	_, err := db.ExecContext(ctx, query, id)
	return err
}

// AddAdversaryToEncounter adds an adversary to an encounter
func AddAdversaryToEncounter(ctx context.Context, tx *sql.Tx, ea *EncounterAdversary) (int64, error) {
	// Check if the adversary is already in the encounter
	query := `
		SELECT id, count FROM encounter_adversaries
		WHERE encounter_id = ? AND adversary_id = ?
	`

	var existingID int64
	var existingCount int
	err := tx.QueryRowContext(ctx, query, ea.EncounterID, ea.AdversaryID).Scan(&existingID, &existingCount)
	
	if err == sql.ErrNoRows {
		// Insert new adversary to encounter
		query = `
			INSERT INTO encounter_adversaries (encounter_id, adversary_id, count)
			VALUES (?, ?, ?)
		`
		result, err := tx.ExecContext(ctx, query, ea.EncounterID, ea.AdversaryID, ea.Count)
		if err != nil {
			return 0, err
		}
		return result.LastInsertId()
	} else if err != nil {
		return 0, err
	}

	// Update existing adversary count
	query = `
		UPDATE encounter_adversaries
		SET count = ?
		WHERE id = ?
	`
	_, err = tx.ExecContext(ctx, query, ea.Count, existingID)
	if err != nil {
		return 0, err
	}

	return existingID, nil
}

// RemoveAdversaryFromEncounter removes an adversary from an encounter
func RemoveAdversaryFromEncounter(ctx context.Context, db *sql.DB, encounterAdversaryID int64) error {
	query := `DELETE FROM encounter_adversaries WHERE id = ?`
	_, err := db.ExecContext(ctx, query, encounterAdversaryID)
	return err
}
