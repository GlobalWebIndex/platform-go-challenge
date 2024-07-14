package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"gwi-platform/database"
	"gwi-platform/models"
	"gwi-platform/utils"
	"net/http"
	"time"
)

func GetUserStatus(ctx context.Context, db database.DB, userID int) (string, error) {
	var fetchedUserID int
	var status string

	utils.InfoLogger.Println("GetUserStatus: Getting status for userID", userID)

	err := db.QueryRowContext(ctx, sqlGetUser, userID).Scan(&fetchedUserID, &status)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("user with ID %d not found", userID)
		}
		return "", fmt.Errorf("failed to fetch user data: %w", err)
	}

	utils.InfoLogger.Println("GetUserStatus: Returning status for userID", userID)
	return status, nil
}

func AddUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate user data
	if user.Username == "" || user.Email == "" {
		http.Error(w, "Username and email are required", http.StatusBadRequest)
		return
	}

	utils.InfoLogger.Println("AddUserHandler: Adding user", user.Username)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Get database connection
	db := database.GetDB()

	// Start a transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to begin transaction: %v", err)
		http.Error(w, "Failed to begin transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Insert the user into the database
	result, err := tx.ExecContext(ctx, sqlAddUser, user.Username, user.Email)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to add user: %v", err)
		http.Error(w, "Failed to insert user", http.StatusInternalServerError)
		return
	}

	// Get the ID of the inserted user
	userID, err := result.LastInsertId()
	if err != nil {
		utils.ErrorLogger.Printf("Failed to get user ID: %v", err)
		http.Error(w, "Failed to get inserted user ID", http.StatusInternalServerError)
		return
	}

	// Set the ID in the user struct
	user.UserID = models.ID(userID)
	// by default all new users are active
	user.Status = "A"

	// Call the stored procedure to create user shards
	_, err = tx.ExecContext(ctx, sqlAddUserShards, userID, 3) // Creating 3 shards, adjust as needed
	if err != nil {
		utils.ErrorLogger.Printf("Failed to add user shards: %v", err)
		http.Error(w, "Failed to create user shards", http.StatusInternalServerError)
		return
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		utils.ErrorLogger.Printf("Failed to commit transaction: %v", err)
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	// Return the inserted user as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
