package pgstore

import (
	"context"
	"encoding/json"
	"main/core"
	"main/mocks"
	"reflect"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestPGStore_Client(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		ctx      context.Context
		source   string
		clientId string
	}
	type test struct {
		name    string
		args    args
		want    core.Client
		wantErr error
	}
	tests := []test{
		{
			name: "happy",
			args: args{
				ctx:      context.Background(),
				source:   "test source",
				clientId: "client id",
			},
			want: core.Client{
				Id:     "some id",
				Source: "some source",
				UserId: "some user",
				DefaultLocation: core.Coordinate{
					Latitude:  1,
					Longitude: 2,
				},
				Info: map[string]any{
					"k": "v",
				},
			},
		},
		{
			name: "unhappy",
			args: args{
				ctx:      context.Background(),
				source:   "test source",
				clientId: "client id",
			},
			wantErr: errorTestError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newStore(ctrl)
			row := mocks.NewMockRow(ctrl)
			row.EXPECT().
				Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Do(func(dest ...interface{}) {
					*(dest[0].(*string)) = tt.want.Source
					*(dest[1].(*string)) = tt.want.Id
					*(dest[2].(*string)) = tt.want.UserId
					*(dest[3].(*float64)) = tt.want.DefaultLocation.Latitude
					*(dest[4].(*float64)) = tt.want.DefaultLocation.Longitude
					if tt.want.Info != nil {
						infob, _ := json.Marshal(tt.want.Info)
						*(dest[5].(*string)) = string(infob)
					}
				}).
				Return(tt.wantErr)
			s.pool.(*mocks.MockDatabase).EXPECT().
				QueryRow(gomock.Any(), "SELECT source, client_id, user_id, latitude, longitude, info FROM clients_registry WHERE source = $1 AND client_id = $2", tt.args.source, tt.args.clientId).
				Return(row)
			got, err := s.Client(tt.args.ctx, tt.args.source, tt.args.clientId, func(client *core.Client, info string) error {
				var infoMap map[string]any
				err := json.Unmarshal([]byte(info), &infoMap)
				client.Info = infoMap
				return err
			})
			if err != tt.wantErr {
				t.Errorf("PGStore.Client() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PGStore.Client() = %v, want %v", got, tt.want)
			}
		})
	}
}
