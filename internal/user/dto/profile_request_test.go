package dto

import (
	"encoding/json"
	"github.com/saurav11sarkar/go/internal/utils"
	"testing"
)

func TestPartialProfileValidation(t *testing.T) {
	for _, tc := range []struct {
		body  string
		valid bool
	}{
		{`{"name":"Saurav"}`, true},
		{`{"role":"user","status":"active"}`, true},
		{`{"name":"Saurav","role":"admin","status":"inactive"}`, true},
		{`{"name":""}`, false},
		{`{"role":"owner"}`, false},
		{`{"status":"invalid"}`, false},
	} {
		var req UpdateProfileRequestDTO
		if err := json.Unmarshal([]byte(tc.body), &req); err != nil {
			t.Fatal(err)
		}
		if err := utils.ValidateStruct(req); (err == nil) != tc.valid {
			t.Fatalf("%s: %v", tc.body, err)
		}
		if tc.body == `{"name":"Saurav"}` && (req.Role != nil || req.Status != nil) {
			t.Fatal("omitted fields must remain nil")
		}
	}
}
