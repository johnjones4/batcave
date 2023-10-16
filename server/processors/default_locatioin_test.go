package processors

import (
	"context"
	"errors"
	"main/core"
	"main/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var errorTestError = errors.New("test error")

func TestProcessors_DefaultLocation(t *testing.T) {
	ctrl := gomock.NewController(t)
	registry := mocks.NewMockClientRegistry(ctrl)
	tests := []struct {
		name       string
		req        core.Request
		wantCoords core.Coordinate
		wantErr    error
	}{
		{
			name: "happy 1",
			req: core.Request{
				Source:   "source 1",
				ClientID: "client 1",
			},
			wantCoords: core.Coordinate{
				Latitude:  1,
				Longitude: 2,
			},
		},
		{
			name: "happy 2",
			req: core.Request{
				Source:   "source 1",
				ClientID: "client 1",
				Coordinate: core.Coordinate{
					Latitude:  3,
					Longitude: 4,
				},
			},
			wantCoords: core.Coordinate{
				Latitude:  3,
				Longitude: 4,
			},
		},
		{
			name: "unhappy",
			req: core.Request{
				Source:   "source 1",
				ClientID: "client 1",
			},
			wantErr: errorTestError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Processors{
				ClientRegistry: registry,
			}
			ctx := context.TODO()
			if tt.req.Coordinate.Empty() {
				registry.EXPECT().Client(ctx, tt.req.Source, tt.req.ClientID, nil).Return(core.Client{
					DefaultLocation: tt.wantCoords,
				}, tt.wantErr)
			}
			if err := p.DefaultLocation(ctx, &tt.req); err != tt.wantErr {
				t.Errorf("Processors.DefaultLocation() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.wantCoords.Latitude, tt.req.Coordinate.Latitude)
			assert.Equal(t, tt.wantCoords.Longitude, tt.req.Coordinate.Longitude)
		})
	}
}
