package user

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
)

func TestProfileResponseContract(t *testing.T) {
	w := httptest.NewRecorder()
	writeError(w, fmt.Errorf("wrapped: %w", ErrUserNotFound))
	if w.Code != 404 {
		t.Fatalf("expected 404, got %d", w.Code)
	}
	got, err := json.Marshal(profileResponse(&User{ID: "id", Name: "Alice"}))
	if err != nil {
		t.Fatal(err)
	}
	want := `{"id":"id","email":"","name":"Alice","role":"","status":"","photo":null,"createdAt":"0001-01-01T00:00:00Z","updatedAt":"0001-01-01T00:00:00Z"}`
	if string(got) != want {
		t.Fatalf("profile JSON changed: %s", got)
	}
}
