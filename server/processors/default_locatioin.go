package processors

import (
	"context"
	"main/core"
)

func (p *Processors) DefaultLocation(ctx context.Context, req *core.Request) error {
	if req.Coordinate.Empty() {
		client, err := p.ClientRegistry.Client(ctx, req.Source, req.ClientID, nil)
		if err != nil {
			return err
		}
		req.Coordinate = client.DefaultLocation
	}
	return nil
}
