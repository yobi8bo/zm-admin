package jwtutil

import "testing"

func TestGeneratedTokensContainSessionAndTokenIDs(t *testing.T) {
	const (
		secret    = "test-secret"
		sessionID = "session-1"
	)

	accessToken, err := GenerateAccessToken(7, "tester", sessionID, secret, 60)
	if err != nil {
		t.Fatalf("GenerateAccessToken() error = %v", err)
	}
	refreshToken, err := GenerateRefreshToken(7, "tester", sessionID, secret, 60)
	if err != nil {
		t.Fatalf("GenerateRefreshToken() error = %v", err)
	}

	accessClaims, err := ParseToken(accessToken, secret)
	if err != nil {
		t.Fatalf("ParseToken(access) error = %v", err)
	}
	refreshClaims, err := ParseToken(refreshToken, secret)
	if err != nil {
		t.Fatalf("ParseToken(refresh) error = %v", err)
	}

	if accessClaims.SessionID != sessionID || refreshClaims.SessionID != sessionID {
		t.Fatalf("session IDs = (%q, %q), want %q", accessClaims.SessionID, refreshClaims.SessionID, sessionID)
	}
	if accessClaims.ID == "" || refreshClaims.ID == "" || accessClaims.ID == refreshClaims.ID {
		t.Fatalf("token IDs must be non-empty and unique")
	}
}
