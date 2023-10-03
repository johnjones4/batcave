package pgstore

import (
	"context"
	"encoding/json"
	"main/mocks"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestPGStore_ClosestMatchingIntent(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		embedding []float32
	}
	tests := []struct {
		name     string
		args     args
		distance float32
		intent   string
		want     string
		wantErr  error
	}{
		{
			args: args{
				[]float32{1, 2, 3, 4},
			},
			distance: 0.1,
			intent:   "intent",
			want:     "intent",
		},
		{
			args: args{
				[]float32{1, 2, 3, 4},
			},
			distance: 0.5,
			intent:   "intent",
			want:     "",
		},
		{
			args: args{
				[]float32{1, 2, 3, 4},
			},
			wantErr: errorTestError,
			want:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newStore(ctrl)

			if tt.wantErr == nil {
				s.log.(*mocks.MockFieldLogger).EXPECT().Debugf(gomock.Any(), gomock.Any(), gomock.Any())
			}

			row := mocks.NewMockRow(ctrl)
			row.EXPECT().
				Scan(gomock.Any(), gomock.Any()).
				DoAndReturn(func(dest ...interface{}) error {
					*(dest[0].(*string)) = tt.intent
					*(dest[1].(*float32)) = tt.distance
					return tt.wantErr
				})

			embeddingjson, _ := json.Marshal(tt.args.embedding)
			s.pool.(*mocks.MockDatabase).
				EXPECT().
				QueryRow(gomock.Any(), "SELECT intent_label, embedding <=> $1 as distance FROM intents ORDER BY distance LIMIT 1", string(embeddingjson)).
				Return(row)
			got, err := s.ClosestMatchingIntent(context.Background(), tt.args.embedding)
			if err != tt.wantErr {
				t.Errorf("PGStore.ClosestMatchingIntent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGStore.ClosestMatchingIntent() = %v, want %v", got, tt.want)
			}
		})
	}
}
