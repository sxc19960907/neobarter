package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreditLevel(t *testing.T) {
	cases := []struct {
		score int
		want  string
	}{
		{0, CreditLevelNormal},
		{100, CreditLevelNormal},
		{119, CreditLevelNormal},
		{120, CreditLevelSilver}, // 边界
		{199, CreditLevelSilver},
		{200, CreditLevelGold}, // 边界
		{499, CreditLevelGold},
		{500, CreditLevelDiamond}, // 边界
		{1000, CreditLevelDiamond},
	}
	for _, c := range cases {
		assert.Equal(t, c.want, CreditLevel(c.score), "score=%d", c.score)
	}
}

func TestUserMarshalJSON_IncludesCreditLevel(t *testing.T) {
	u := User{ID: 1, Phone: "13800000000", CreditScore: 250}
	data, err := json.Marshal(u)
	require.NoError(t, err)

	var m map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &m))

	assert.Equal(t, "gold", m["credit_level"], "250分应为金牌")
	assert.Equal(t, float64(250), m["credit_score"], "原字段保留")
	assert.Equal(t, "13800000000", m["phone"])
	// 敏感字段不应出现
	_, hasRealName := m["RealName"]
	assert.False(t, hasRealName)
}
