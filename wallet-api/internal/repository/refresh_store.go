package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Hash      string
	UserAgent string
	IP        string
	IssuedAt  time.Time
	UsedAt    sql.NullTime
	Revoked   bool
}

type Store struct {
	db *sql.DB
}

func NewRefreshStore(db *sql.DB) *Store {
	return &Store{db}
}

func (s *Store) Save(rt *RefreshToken) error {
	_, err := s.db.Exec(
		`INSERT INTO refresh_tokens(id,user_id,token_hash,user_agent,ip,issued_at,revoked)
         VALUES($1,$2,$3,$4,$5,$6,$7)`,
		rt.ID, rt.UserID, rt.Hash, rt.UserAgent, rt.IP, rt.IssuedAt, false,
	)
	return err
}

func (s *Store) FindValid(userID uuid.UUID) (*RefreshToken, error) {
	row := s.db.QueryRow(
		`SELECT id,user_id,token_hash,user_agent,ip,issued_at,used_at,revoked
         FROM refresh_tokens WHERE user_id=$1 AND revoked=false AND used_at IS NULL`,
		userID,
	)
	var rt RefreshToken
	err := row.Scan(&rt.ID, &rt.UserID, &rt.Hash, &rt.UserAgent, &rt.IP, &rt.IssuedAt, &rt.UsedAt, &rt.Revoked)
	return &rt, err
}

func (s *Store) MarkUsed(id uuid.UUID) error {
	_, err := s.db.Exec(
		`UPDATE refresh_tokens SET used_at=$1, revoked=true WHERE id=$2`,
		time.Now(), id,
	)
	return err
}

func (s *Store) RevokeAll(userID uuid.UUID) error {
	_, err := s.db.Exec(
		`UPDATE refresh_tokens SET revoked=true WHERE user_id=$1`,
		userID,
	)
	return err
}
