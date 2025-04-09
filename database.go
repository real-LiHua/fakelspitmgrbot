package main

import (
	"crypto/rand"
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

const (
	FlagNotMember = 0x0
	FlagIsMember  = 0x1
	FlagBanned    = 0x2
)

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
			telegram_id INTEGER PRIMARY KEY,
			github_id INTEGER DEFAULT 0,
			challenge_code TEXT DEFAULT '',
			github_username TEXT DEFAULT '',
			flag INTEGER DEFAULT 0,
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
	_, err := db.Conn.Exec("INSERT INTO users (telegram_id, github_id, challenge_code, github_username, flag) VALUES (?, ?, ?, ?, ?)",
		user.TelegramID,
		user.GithubID,
		user.ChallengeCode,
		user.GithubUsername,
		user.Flag,
	)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) UpdateUser(user *User) error {
	_, err := db.Conn.Exec("UPDATE users SET github_id = ?, github_username = ?, flag = ? WHERE telegram_id = ?",
		user.GithubID,
		user.GithubUsername,
		user.Flag,
		user.TelegramID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) UpdateChallengeCode(telegramID int64) string {
	challengeCode := rand.Text()
	_, err := db.Conn.Exec("UPDATE users SET challenge_code = ? WHERE telegram_id = ?", challengeCode, telegramID)
	if err != nil {
		db.AddUser(&User{
			TelegramID:    telegramID,
			ChallengeCode: challengeCode,
		})
	}
	return challengeCode
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
	_, err := db.Conn.Exec("UPDATE users SET flag = ? WHERE telegram_id = ?", FlagBanned, user.TelegramID)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) UnbanUser(user User) error {
	if user.GithubID == 0 {
		return nil
	}
	_, err := db.Conn.Exec("DELETE FROM users WHERE telegram_id = ?", user.TelegramID)
	if err != nil {
		return err
	}
	return nil
}
