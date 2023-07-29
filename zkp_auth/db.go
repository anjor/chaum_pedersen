package zkp_auth

import (
	"database/sql"
	"fmt"
)

func insertRowIntoDB(dbPath string, user string, y1 int64, y2 int64) error {
	// Open the database file
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	// Create the table if it doesn't exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS user_table (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			u TEXT,
			y1 INTEGER,
			y2 INTEGER 
		);`)
	if err != nil {
		return err
	}

	// Insert the data into the table
	_, err = db.Exec("INSERT INTO user_table (u, y1, y2) VALUES (?, ?, ?)",
		user, y1, y2)
	if err != nil {
		return err
	}

	fmt.Printf("Row inserted successfully for user: %s\n", user)
	return nil
}

func getRowByUser(dbPath, user string) (string, int64, int64, error) {
	// Open the database file
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return "", 0, 0, err
	}
	defer db.Close()

	// Query the row with the specified user
	row := db.QueryRow("SELECT u, y1, y2 FROM user_table WHERE u = ?", user)

	// Extract the data from the row
	var retrievedUser string
	var y1, y2 int64
	err = row.Scan(&retrievedUser, &y1, &y2)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", 0, 0, fmt.Errorf("user not found: %s", user)
		}
		return "", 0, 0, err
	}

	return retrievedUser, y1, y2, nil
}
