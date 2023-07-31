package zkp_auth

import (
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
	"testing"
)

func TestInsertRowIntoDBAndGetRowByUser(t *testing.T) {
	dbPath := "./test.db"
	defer os.Remove(dbPath) // Remove the test database file after the tests

	// Test data
	user := "test_user"
	y1 := int64(123)
	y2 := int64(456)

	// Test insertRowIntoDB
	err := insertRowIntoDB(dbPath, user, y1, y2)
	if err != nil {
		t.Fatalf("insertRowIntoDB failed: %v", err)
	}

	// Test getRowByUser
	retrievedUser, retrievedY1, retrievedY2, err := getRowByUser(dbPath, user)
	if err != nil {
		t.Fatalf("getRowByUser failed: %v", err)
	}

	// Check if retrieved data matches the original data
	if retrievedUser != user {
		t.Errorf("Expected user: %s, Got: %s", user, retrievedUser)
	}
	if retrievedY1 != y1 {
		t.Errorf("Expected y1: %d, Got: %d", y1, retrievedY1)
	}
	if retrievedY2 != y2 {
		t.Errorf("Expected y2: %d, Got: %d", y2, retrievedY2)
	}

	// Test user not found
	nonExistingUser := "non_existing_user"
	_, _, _, err = getRowByUser(dbPath, nonExistingUser)
	if err == nil {
		t.Errorf("Expected an error for non-existing user")
	}
	if err != nil && err.Error() != "user not found: "+nonExistingUser {
		t.Errorf("Expected 'user not found' error, Got: %v", err)
	}
}

func TestInsertRowIntoDBWithError(t *testing.T) {
	// Test database error
	dbPath := filepath.Join("non_existing_directory", "test.db")
	user := "test_user"
	y1 := int64(123)
	y2 := int64(456)

	err := insertRowIntoDB(dbPath, user, y1, y2)
	if err == nil {
		t.Errorf("Expected an error for non-existing database file")
	}
}

func TestGetRowByUserWithError(t *testing.T) {
	// Test database error
	dbPath := filepath.Join("non_existing_directory", "test.db")
	user := "test_user"

	_, _, _, err := getRowByUser(dbPath, user)
	if err == nil {
		t.Errorf("Expected an error for non-existing database file")
	}
}
