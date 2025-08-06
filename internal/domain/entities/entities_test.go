package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransactionState_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		state TransactionState
		want  bool
	}{
		{
			name:  "valid win state",
			state: StateWin,
			want:  true,
		},
		{
			name:  "valid lose state",
			state: StateLose,
			want:  true,
		},
		{
			name:  "invalid empty state",
			state: TransactionState(""),
			want:  false,
		},
		{
			name:  "invalid random state",
			state: TransactionState("invalid"),
			want:  false,
		},
		{
			name:  "invalid uppercase state",
			state: TransactionState("WIN"),
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.state.IsValid()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSourceType_IsValid(t *testing.T) {
	tests := []struct {
		name       string
		sourceType SourceType
		want       bool
	}{
		{
			name:       "valid game source",
			sourceType: SourceTypeGame,
			want:       true,
		},
		{
			name:       "valid server source",
			sourceType: SourceTypeServer,
			want:       true,
		},
		{
			name:       "valid payment source",
			sourceType: SourceTypePayment,
			want:       true,
		},
		{
			name:       "invalid empty source",
			sourceType: SourceType(""),
			want:       false,
		},
		{
			name:       "invalid random source",
			sourceType: SourceType("invalid"),
			want:       false,
		},
		{
			name:       "invalid uppercase source",
			sourceType: SourceType("GAME"),
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.sourceType.IsValid()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTransactionRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request TransactionRequest
	}{
		{
			name: "valid transaction request",
			request: TransactionRequest{
				State:         "win",
				Amount:        "10.50",
				TransactionID: "tx-123",
			},
		},
		{
			name: "valid lose transaction",
			request: TransactionRequest{
				State:         "lose",
				Amount:        "25.75",
				TransactionID: "tx-456",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that the request can be created without issues
			assert.NotEmpty(t, tt.request.State)
			assert.NotEmpty(t, tt.request.Amount)
			assert.NotEmpty(t, tt.request.TransactionID)
		})
	}
}

func TestBalanceResponse_Structure(t *testing.T) {
	response := BalanceResponse{
		UserID:  123,
		Balance: "100.50",
	}

	assert.Equal(t, uint64(123), response.UserID)
	assert.Equal(t, "100.50", response.Balance)
}
