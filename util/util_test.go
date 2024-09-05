package util

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestGenerateToken(t *testing.T) {
	userID := 186
	token, err := GenerateToken(userID)
	if err != nil {
		t.Errorf("generate token failed, err:%v", err)
		return
	}
	t.Logf("token:%s", token)
	claims, err := ParseToken(token)
	if err != nil {
		t.Errorf("parse token failed, err:%v", err)
		return
	}
	assert.Equal(t, claims.UserID, userID)
}
