package main

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type User struct {
	GithubID       int64  `json:"github_id"`
	TelegramID     int64  `json:"telegram_id"`
	ChallengeCode  string `json:"challenge_code"`
	GithubUsername string `json:"github_username"`
	Flag           int    `json:"flag"`
}

type Database struct {
	Conn *sql.DB
}

func NewDatabase() (*Database, error) {
	db, err := sql.Open("sqlite", "database.db")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			github_id INTEGER PRIMARY KEY,
			telegram_id INTEGER,
			challenge_code TEXT,
			github_username TEXT,
			flag INTEGER
		);
		`)
	if err != nil {
		return nil, err
	}
	return &Database{Conn: db}, nil
}

func (db *Database) Close() error {
	return db.Conn.Close()
}

func (db *Database) AddUser(user *User) error {
	_, err := db.Conn.Exec("INSERT INTO users (github_id, telegram_id, challenge_code, github_username, flag) VALUES (?, ?, ?, ?, ?)",
		user.GithubID,
		user.TelegramID,
		user.ChallengeCode,
		user.GithubUsername,
		user.Flag,
	)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) GetUserByTelegramID(telegramID int64) User {
	user := User{}
	err := db.Conn.QueryRow("SELECT * FROM users WHERE telegram_id = ?", telegramID).Scan(
		&user.GithubID,
		&user.TelegramID,
		&user.ChallengeCode,
		&user.GithubUsername,
		&user.Flag,
	)
	if err != nil {
		return User{}
	}
	return user
}

func (db *Database) GetUserByGithubID(githubID int64) User {
	user := User{}
	err := db.Conn.QueryRow("SELECT * FROM users WHERE github_id = ?", githubID).Scan(
		&user.GithubID,
		&user.TelegramID,
		&user.ChallengeCode,
		&user.GithubUsername,
		&user.Flag,
	)
	if err != nil {
		return User{}
	}
	return user
}

func (db *Database) BanUser(user User) error {
	if user.GithubID == 0 {
		return nil
	}
	_, err := db.Conn.Exec("UPDATE users SET flag = 1 WHERE github_id = ?", user.GithubID)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) UnbanUser(user User) error {
	if user.GithubID == 0 {
		return nil
	}
	_, err := db.Conn.Exec("UPDATE users SET flag = 0 WHERE github_id = ?", user.GithubID)
	if err != nil {
		return err
	}
	return nil
}
