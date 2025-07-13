package db

import (
	"context"
	"database/sql"
	"time"
)

// Adversary represents a Daggerheart adversary entity
type Adversary struct {
	ID              int64
	Name            string
	Type            string
	ChallengeRating string
	Size            string
	ArmorClass      int
	HitPoints       int
	Speed           string
	Strength        int
	Dexterity       int
	Constitution    int
	Intelligence    int
	Wisdom          int
	Charisma        int
	Abilities       string
	Actions         string
	Reactions       string
	Description     string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// GetAllAdversaries retrieves all adversaries from the database
func GetAllAdversaries(ctx context.Context, db *sql.DB) ([]*Adversary, error) {
	query := `
		SELECT id, name, type, challenge_rating, size, armor_class, hit_points, 
		       speed, strength, dexterity, constitution, intelligence, wisdom, 
		       charisma, abilities, actions, reactions, description, created_at, updated_at
		FROM adversaries
		ORDER BY name ASC
	`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var adversaries []*Adversary
	for rows.Next() {
		adv := &Adversary{}
		err := rows.Scan(
			&adv.ID, &adv.Name, &adv.Type, &adv.ChallengeRating, &adv.Size,
			&adv.ArmorClass, &adv.HitPoints, &adv.Speed, &adv.Strength,
			&adv.Dexterity, &adv.Constitution, &adv.Intelligence, &adv.Wisdom,
			&adv.Charisma, &adv.Abilities, &adv.Actions, &adv.Reactions,
			&adv.Description, &adv.CreatedAt, &adv.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		adversaries = append(adversaries, adv)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return adversaries, nil
}

// GetAdversaryByID retrieves a single adversary by ID
func GetAdversaryByID(ctx context.Context, db *sql.DB, id int64) (*Adversary, error) {
	query := `
		SELECT id, name, type, challenge_rating, size, armor_class, hit_points, 
		       speed, strength, dexterity, constitution, intelligence, wisdom, 
		       charisma, abilities, actions, reactions, description, created_at, updated_at
		FROM adversaries
		WHERE id = ?
	`

	adv := &Adversary{}
	err := db.QueryRowContext(ctx, query, id).Scan(
		&adv.ID, &adv.Name, &adv.Type, &adv.ChallengeRating, &adv.Size,
		&adv.ArmorClass, &adv.HitPoints, &adv.Speed, &adv.Strength,
		&adv.Dexterity, &adv.Constitution, &adv.Intelligence, &adv.Wisdom,
		&adv.Charisma, &adv.Abilities, &adv.Actions, &adv.Reactions,
		&adv.Description, &adv.CreatedAt, &adv.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return adv, nil
}

// CreateAdversary inserts a new adversary into the database
func CreateAdversary(ctx context.Context, db *sql.DB, adv *Adversary) (int64, error) {
	query := `
		INSERT INTO adversaries (
			name, type, challenge_rating, size, armor_class, hit_points, 
			speed, strength, dexterity, constitution, intelligence, wisdom, 
			charisma, abilities, actions, reactions, description
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := db.ExecContext(
		ctx, query,
		adv.Name, adv.Type, adv.ChallengeRating, adv.Size,
		adv.ArmorClass, adv.HitPoints, adv.Speed, adv.Strength,
		adv.Dexterity, adv.Constitution, adv.Intelligence, adv.Wisdom,
		adv.Charisma, adv.Abilities, adv.Actions, adv.Reactions,
		adv.Description,
	)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// UpdateAdversary updates an existing adversary in the database
func UpdateAdversary(ctx context.Context, db *sql.DB, adv *Adversary) error {
	query := `
		UPDATE adversaries
		SET name = ?, type = ?, challenge_rating = ?, size = ?, 
		    armor_class = ?, hit_points = ?, speed = ?, strength = ?, 
		    dexterity = ?, constitution = ?, intelligence = ?, wisdom = ?, 
		    charisma = ?, abilities = ?, actions = ?, reactions = ?, 
		    description = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := db.ExecContext(
		ctx, query,
		adv.Name, adv.Type, adv.ChallengeRating, adv.Size,
		adv.ArmorClass, adv.HitPoints, adv.Speed, adv.Strength,
		adv.Dexterity, adv.Constitution, adv.Intelligence, adv.Wisdom,
		adv.Charisma, adv.Abilities, adv.Actions, adv.Reactions,
		adv.Description, adv.ID,
	)

	return err
}

// DeleteAdversary removes an adversary from the database
func DeleteAdversary(ctx context.Context, db *sql.DB, id int64) error {
	query := `DELETE FROM adversaries WHERE id = ?`
	_, err := db.ExecContext(ctx, query, id)
	return err
}
