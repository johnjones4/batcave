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

func TestPGStore_ClientsForUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		userId string
	}
	tests := []struct {
		name     string
		args     args
		want     []core.Client
		wantErr  error
		queryErr error
		scanErr  error
	}{
		{
			name: "happy",
			args: args{
				userId: "some user",
			},
			want: []core.Client{
				{
					Id:     "some id 1",
					Source: "some source 1",
					UserId: "some user",
					DefaultLocation: core.Coordinate{
						Latitude:  1,
						Longitude: 2,
					},
					Info: map[string]any{
						"k": "v",
					},
				},
				{
					Id:     "some id 2",
					Source: "some source 2",
					UserId: "some user",
					DefaultLocation: core.Coordinate{
						Latitude:  3,
						Longitude: 4,
					},
					Info: map[string]any{
						"k": "vv",
					},
				},
			},
		},
		{
			name: "query error",
			args: args{
				userId: "some user",
			},
			want:     nil,
			queryErr: errorTestError,
			wantErr:  errorTestError,
		},
		{
			name: "scan error",
			args: args{
				userId: "some user",
			},
			want:    nil,
			scanErr: errorTestError,
			wantErr: errorTestError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newStore(ctrl)
			rows := mocks.NewMockRows(ctrl)
			i := 0
			if tt.queryErr == nil {
				scans := len(tt.want)
				if tt.scanErr != nil {
					scans = 1
				}
				rows.EXPECT().
					Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Do(func(dest ...interface{}) {
						if i >= len(tt.want) {
							return
						}
						*(dest[0].(*string)) = tt.want[i].Source
						*(dest[1].(*string)) = tt.want[i].Id
						*(dest[2].(*string)) = tt.want[i].UserId
						*(dest[3].(*float64)) = tt.want[i].DefaultLocation.Latitude
						*(dest[4].(*float64)) = tt.want[i].DefaultLocation.Longitude
						if tt.want[i].Info != nil {
							infob, _ := json.Marshal(tt.want[i].Info)
							*(dest[5].(*string)) = string(infob)
						}
						i++
					}).
					Return(tt.scanErr).Times(scans)
				rows.EXPECT().Close()
				rows.EXPECT().Next().Times(len(tt.want) + 1).DoAndReturn(func() bool {
					if tt.scanErr == nil {
						return i < len(tt.want)
					} else {
						return true
					}
				})
			}
			s.pool.(*mocks.MockDatabase).EXPECT().
				Query(gomock.Any(), "SELECT source, client_id, user_id, latitude, longitude, info FROM clients_registry WHERE user_id = $1", tt.args.userId).
				Return(rows, tt.queryErr)
			got, err := s.ClientsForUser(context.Background(), tt.args.userId, func(client *core.Client, info string) error {
				var infoMap map[string]any
				err := json.Unmarshal([]byte(info), &infoMap)
				client.Info = infoMap
				return err
			})
			if err != tt.wantErr {
				t.Errorf("PGStore.ClientsForUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PGStore.ClientsForUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
