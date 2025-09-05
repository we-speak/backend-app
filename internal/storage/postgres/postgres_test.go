package postgres_test

import (
	"backend-app/internal/storage/models"
	"backend-app/internal/storage/postgres"
	"fmt"
	"testing"

	driver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB() (*postgres.Storage, error) {
	dsn := "host=localhost user=postgres password=postgres dbname=authdb port=5432 sslmode=disable TimeZone=Europe/Moscow"
	db, err := gorm.Open(driver.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.User{})

	return &postgres.Storage{DB: db}, nil
}

func TestCreateUser(t *testing.T) {
	storage, err := setupTestDB()
	if err != nil {
		t.Fatalf("Fa`iled to set up test DB: %v", err)
	}

	user := &models.User{
		Username: "tes3tusedr33333346576",
		Password: "password1233",
		Email:    "testu3ser3333@example.com",
		Country:  "Testland",
	}

	err = storage.CreateUser(user)
	if err != nil {
		t.Errorf("Failed to create user: %v", err)
	}

	if user.ID == 0 {
		t.Errorf("Expected user ID to be set, got 0")
	}
}

func TestGetUserByID(t *testing.T) {
	storage, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test DB: %v", err)
	}

	user := &models.User{
		Username: "testuser56565365",
		Password: "password123",
		Email:    "testus33er@example.com",
		Country:  "Testland",
	}

	storage.CreateUser(user)

	retrievedUser, err := storage.GetUserByID(user.ID)
	if err != nil {
		t.Errorf("Failed to get user by ID: %v", err)
	}

	if retrievedUser.Username != user.Username {
		t.Errorf("Expected username %s, got %s", user.Username, retrievedUser.Username)
	}
}

func TestUpdateUser(t *testing.T) {
	storage, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test DB: %v", err)
	}

	user := &models.User{
		Username: "testuse3r66767676",
		Password: "password123",
		Email:    "tes3tuse4r@example.com",
		Country:  "Testland",
	}

	storage.CreateUser(user)

	user.Username = "updateduser"
	err = storage.UpdateUser(user)
	if err != nil {
		t.Errorf("Failed to update user: %v", err)
	}

	updatedUser, _ := storage.GetUserByID(user.ID)
	if updatedUser.Username != "updateduser" {
		t.Errorf("Expected username to be updated to 'updateduser', got %s", updatedUser.Username)
	}
}

func TestDeleteUser(t *testing.T) {
	storage, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test DB: %v", err)
	}

	user := &models.User{
		Username: "test43user54353453",
		Password: "password123",
		Email:    "testu34ser@example.com",
		Country:  "Testland",
	}

	storage.CreateUser(user)

	err = storage.DeleteUser(user.ID)
	if err != nil {
		t.Errorf("Failed to delete user: %v", err)
	}

	_, err = storage.GetUserByID(user.ID)
	if err == nil {
		t.Errorf("Expected error when retrieving deleted user, got nil")
	}
}

func TestGetAllUsers(t *testing.T) {
	storage, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test DB: %v", err)
	}

	users := []models.User{
		{Username: "user1", Password: "pass1", Email: "user1@example.com", Country: "Country1"},
		{Username: "user2", Password: "pass2", Email: "user2@example.com", Country: "Country2"},
	}

	for _, user := range users {
		storage.CreateUser(&user)
	}

	retrievedUsers, err := storage.GetAllUsers()
	if err != nil {
		t.Errorf("Failed to get all users: %v", err)
	}

	if len(retrievedUsers) != len(users) {
		fmt.Printf("Expected %d users, got %d", len(users), len(retrievedUsers))
	}
}
