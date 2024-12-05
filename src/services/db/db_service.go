package dbservice

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

const (
	// Better to set these from .env file next time
	host     = "database"
	port     = 5432
	user     = "postgres"
	password = "password123"
	dbname   = "postgres"
)

type DBService struct {
	context            *context.Context
	currentTransaction *sql.Tx
	database           *sql.DB
}

var (
	ErrTransNotStarted = fmt.Errorf("no transaction started")
	ErrUserNotFound    = fmt.Errorf("could not find user with GUID provided")
)

func New() *DBService {
	psqlconn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	return &DBService{
		context:            &ctx,
		currentTransaction: nil,
		database:           db}
}

func (s *DBService) BeginTransaction() error {
	tx, err := s.database.BeginTx(*s.context, nil)
	if err != nil {
		return err
	}
	s.currentTransaction = tx
	return nil
}

func (s *DBService) CommitTransaction() error {
	if s.currentTransaction == nil {
		return nil
	}
	err := s.currentTransaction.Commit()
	if err != nil {
		return err
	}
	s.currentTransaction = nil
	return nil
}

func (s *DBService) RollbackTransaction() error {
	if s.currentTransaction == nil {
		return nil
	}
	err := s.currentTransaction.Rollback()
	if err != nil {
		return err
	}
	s.currentTransaction = nil
	return nil
}

func (s *DBService) CheckUserPresentByGUID(guid *uuid.UUID) (bool, error) {
	var userPresent = false
	err := s.database.QueryRowContext(*s.context,
		`SELECT true FROM users WHERE guid = $1`, *guid).Scan(&userPresent)
	if err == sql.ErrNoRows {
		err = ErrUserNotFound
	}
	return userPresent, err
}

func (s *DBService) UpdateRefreshTokenHashByGUID(
	guid *uuid.UUID, newHash []byte) (err error) {

	updateStmt := "UPDATE users SET refresh_token_hash = $1 where guid = $2"
	_, err = s.database.ExecContext(*s.context, updateStmt, newHash, guid)
	return
}

func (s *DBService) UpdateRefreshTokenHashByGUIDTx(
	guid *uuid.UUID, newHash []byte) (err error) {

	if s.currentTransaction == nil {
		return ErrTransNotStarted
	}
	updateStmt := "UPDATE users SET refresh_token_hash = $1 where guid = $2"
	_, err = s.currentTransaction.ExecContext(*s.context, updateStmt, newHash, guid)
	return
}

func (s *DBService) GetMailByGUIDTx(guid *uuid.UUID) (string, error) {
	if s.currentTransaction == nil {
		return "", ErrTransNotStarted
	}
	var email = ""
	err := s.currentTransaction.QueryRowContext(*s.context,
		`SELECT email FROM users WHERE guid = $1`, *guid).Scan(&email)
	if err == sql.ErrNoRows {
		err = ErrUserNotFound
	}
	return email, err
}

func (s *DBService) GetRefreshTokenHashByGUIDTx(guid *uuid.UUID) (string, error) {
	if s.currentTransaction == nil {
		return "", ErrTransNotStarted
	}
	var tokenHash = ""
	err := s.currentTransaction.QueryRowContext(*s.context,
		`SELECT refresh_token_hash FROM users WHERE guid = $1`, *guid).Scan(&tokenHash)
	if err == sql.ErrNoRows {
		err = ErrUserNotFound
	}
	return tokenHash, err
}

func (s *DBService) CloseConnection() {
	s.database.Close()
}
