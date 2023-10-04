package pgstore

import (
	"context"
	"main/core"
	"main/mocks"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestPGStore_LogPush(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		clientId string
		push     *core.PushMessage
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "happy",
			args: args{
				clientId: "some id",
				push: &core.PushMessage{
					OutboundMessage: core.OutboundMessage{
						Message: core.Message{
							Text: "message text",
						},
						Media: core.Media{
							Type: "media type",
							URL:  "media url",
						},
					},
				},
			},
		},
		{
			name: "happy",
			args: args{
				clientId: "some id",
				push: &core.PushMessage{
					OutboundMessage: core.OutboundMessage{
						Message: core.Message{
							Text: "message text",
						},
						Media: core.Media{
							Type: "media type",
							URL:  "media url",
						},
					},
				},
			},
			wantErr: errorTestError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newStore(ctrl)
			s.pool.(*mocks.MockDatabase).
				EXPECT().
				Exec(gomock.Any(), "INSERT INTO pushes (event_id, timestamp, client_id, message_text, media_url, media_type) VALUES ($1,$2,$3,$4,$5,$6)", tt.args.push.EventId, gomock.Any(), tt.args.clientId, tt.args.push.Message.Text, tt.args.push.Media.URL, tt.args.push.Media.Type).
				Return(tt.wantErr)
			if err := s.LogPush(context.Background(), tt.args.clientId, tt.args.push); err != tt.wantErr {
				t.Errorf("PGStore.LogPush() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
