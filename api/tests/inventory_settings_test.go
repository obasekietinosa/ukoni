package tests

import (
	"context"
	"testing"
	"ukoni/internal/models"
)

func TestInventorySettingsModel(t *testing.T) {
	// Ensure clean state
	clearDB()

	ctx := context.Background()
	db := testDB

	// 1. Create a User
	user := &models.User{
		Email:        "test_settings@example.com",
		Name:         "Settings User",
		PasswordHash: "hash",
	}
	userModel := &models.UserModel{DB: db}
	if err := userModel.Insert(user); err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	// 2. Create an Inventory
	inventory := &models.Inventory{
		Name:        "Settings Inventory",
		OwnerUserID: user.ID,
	}
	inventoryModel := &models.InventoryModel{DB: db}
	if err := inventoryModel.Create(ctx, db, inventory); err != nil {
		t.Fatalf("failed to create inventory: %v", err)
	}

	// 3. Initialize Settings Model
	settingsModel := &models.InventorySettingsModel{DB: db}

	// 4. Test Get (should be nil initially)
	s, err := settingsModel.Get(ctx, inventory.ID)
	if err != nil {
		t.Fatalf("failed to get settings: %v", err)
	}
	if s != nil {
		t.Errorf("expected nil settings, got %v", s)
	}

	// 5. Test Upsert (Insert)
	provider := "openai"
	key := "sk-test-key"
	newSettings := &models.InventorySettings{
		InventoryID: inventory.ID,
		LLMProvider: &provider,
		LLMAPIKey:   &key,
	}
	err = settingsModel.Upsert(ctx, db, newSettings)
	if err != nil {
		t.Fatalf("failed to upsert settings: %v", err)
	}

	// 6. Verify Insert
	s, err = settingsModel.Get(ctx, inventory.ID)
	if err != nil {
		t.Fatalf("failed to get settings: %v", err)
	}
	if s == nil {
		t.Fatal("expected settings, got nil")
	}
	if s.LLMProvider == nil || *s.LLMProvider != provider {
		t.Errorf("expected provider %s, got %v", provider, s.LLMProvider)
	}
	if s.LLMAPIKey == nil || *s.LLMAPIKey != key {
		t.Errorf("expected key %s, got %v", key, s.LLMAPIKey)
	}

	// 7. Test Upsert (Update)
	newProvider := "gemini"
	newKey := "new-key"
	updateSettings := &models.InventorySettings{
		InventoryID: inventory.ID,
		LLMProvider: &newProvider,
		LLMAPIKey:   &newKey,
	}
	err = settingsModel.Upsert(ctx, db, updateSettings)
	if err != nil {
		t.Fatalf("failed to update settings: %v", err)
	}

	// 8. Verify Update
	s, err = settingsModel.Get(ctx, inventory.ID)
	if err != nil {
		t.Fatalf("failed to get settings: %v", err)
	}
	if s == nil {
		t.Fatal("expected settings, got nil")
	}
	if s.LLMProvider == nil || *s.LLMProvider != newProvider {
		t.Errorf("expected provider %s, got %v", newProvider, s.LLMProvider)
	}
	if s.LLMAPIKey == nil || *s.LLMAPIKey != newKey {
		t.Errorf("expected key %s, got %v", newKey, s.LLMAPIKey)
	}
}
