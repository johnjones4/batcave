package pgstore

import (
	"context"
	"main/mocks"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestPGStore_ClearScheduledRecurringEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		ctx context.Context
		id  string
	}
	type test struct {
		name    string
		args    args
		wantErr error
	}
	tests := []test{
		{
			name: "happy",
			args: args{
				context.Background(),
				"test",
			},
		},
		{
			name: "unhappy",
			args: args{
				context.Background(),
				"test 1",
			},
			wantErr: errorTestError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newStore(ctrl)
			s.pool.(*mocks.MockDatabase).EXPECT().Exec(gomock.Any(), "DELETE FROM recurring_events WHERE event_id = $1", tt.args.id).Return(tt.wantErr)
			if err := s.ClearScheduledRecurringEvent(tt.args.ctx, tt.args.id); err != tt.wantErr {
				t.Errorf("PGStore.ClearScheduledRecurringEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
