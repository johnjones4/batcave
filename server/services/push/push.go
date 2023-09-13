package push

import (
	"context"
	"encoding/json"
	"errors"
	"main/core"
	"time"

	"github.com/sirupsen/logrus"
)

type PushClient struct {
	Source string `json:"source"`
	Id     string `json:"id"`
}

type PushConfiguration struct {
	ClienIdClusters [][]PushClient `json:"clienIdClusters"`
}

type Push struct {
	PushConfiguration PushConfiguration
	ClientProviders   []core.ClientProvider
	Scheduler         core.Scheduler
	Log               logrus.FieldLogger
	PushLogger        core.PushLogger
}

type pushEventInfo struct {
	Source   string           `json:"source"`
	ClientId string           `json:"clientID"`
	Message  core.PushMessage `json:"message"`
}

var (
	ErrorClientDoesNotSupportPush = errors.New("client id does not support push")
)

const (
	eventType = "push"
)

func (a *Push) SendLater(ctx context.Context, when time.Time, source string, clientId string, message core.PushMessage) error {
	return a.Scheduler.ScheduleEvent(ctx, &core.ScheduledEvent{
		EventType: eventType,
		Scheduled: when,
		Info: pushEventInfo{
			Source:   source,
			ClientId: clientId,
			Message:  message,
		},
	})
}

func (a *Push) Start(ctx context.Context) error {
	ticker := time.NewTicker(time.Second * 10) //TODO look ahead
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			events, err := a.Scheduler.GetReadyEvents(ctx, eventType, func(event *core.ScheduledEvent, info string) error {
				var receiver pushEventInfo
				err := json.Unmarshal([]byte(info), &receiver)
				if err != nil {
					return err
				}
				event.Info = receiver
				return nil
			})
			if err != nil {
				a.Log.Errorf("Error getting scheduled events: %e", err)
				continue
			}

			for _, event := range events {
				info, ok := event.Info.(pushEventInfo)
				if !ok {
					continue
				}
				err = a.Send(ctx, info.Source, info.ClientId, info.Message)
				if err != nil {
					a.Log.Errorf("Error sending push message: %e", err)
					continue
				}
				err = a.Scheduler.ClearScheduledEvent(ctx, event.ID)
				if err != nil {
					a.Log.Errorf("Error clearing push message: %e", err)
					continue
				}
			}
		}
	}
}

func (a *Push) Send(ctx context.Context, source string, clientId string, message core.PushMessage) error {
	err := a.sendToClent(ctx, clientId, message)
	if err != nil && err != ErrorClientDoesNotSupportPush {
		return err
	}
	if err == nil {
		return nil
	}

	clusterIndex := -1
	skipIndex := -1
clusterLoop:
	for i, cluster := range a.PushConfiguration.ClienIdClusters {
		for j, clusterClientId := range cluster {
			if clusterClientId.Source == source && clusterClientId.Id == clientId {
				clusterIndex = i
				skipIndex = j
				break clusterLoop
			}
		}
	}

	if clusterIndex < 0 {
		return errors.New("client id does not have a cluster")
	}

	for i, id := range a.PushConfiguration.ClienIdClusters[clusterIndex] {
		if i == skipIndex {
			continue
		}
		err = a.sendToClent(ctx, id.Id, message)
		if err != nil && err != ErrorClientDoesNotSupportPush {
			return err
		}
	}

	return ErrorClientDoesNotSupportPush
}

func (a *Push) sendToClent(ctx context.Context, clientId string, message core.PushMessage) error {
	for _, provider := range a.ClientProviders {
		ok, err := provider.SendToClient(ctx, clientId, message)
		if ok {
			return nil
		}
		if err != nil {
			return err
		}
		err = a.PushLogger.LogPush(ctx, clientId, &message)
		if err != nil {
			return err
		}
	}
	return ErrorClientDoesNotSupportPush
}
