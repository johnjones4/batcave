package push

import (
	"context"
	"encoding/json"
	"errors"
	"main/core"
	"time"

	"github.com/sirupsen/logrus"
)

type Push struct {
	ClientSenders  []core.ClientSender
	ClientRegistry core.ClientRegistry
	Scheduler      core.Scheduler
	Log            logrus.FieldLogger
	PushLogger     core.PushLogger
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
	eventType     = "push"
	checkInterval = time.Second * 30
)

func (a *Push) SendLater(ctx context.Context, when time.Time, source string, clientId string, message core.PushMessage) error {
	wait := time.Until(when)
	if wait <= checkInterval {
		go a.sendScheduledAsync(context.Background(), when, message.EventId, source, clientId, message)
		return nil
	} else {
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
}

func (a *Push) sendScheduledAsync(ctx context.Context, when time.Time, eventId string, source string, clientId string, message core.PushMessage) {
	wait := time.Until(when)
	if wait > 0 {
		time.Sleep(wait)
	}
	err := a.Send(ctx, source, clientId, message)
	if err != nil {
		a.Log.Errorf("Error sending push message: %e", err)
		return
	}
	if eventId != "" {
		err = a.Scheduler.ClearScheduledEvent(ctx, eventId)
		if err != nil {
			a.Log.Errorf("Error clearing push message: %e", err)
			return
		}
	}
}

func (a *Push) Start(ctx context.Context) error {
	ticker := time.NewTicker(checkInterval)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			events, err := a.Scheduler.ReadyEvents(ctx, time.Now().Add(checkInterval), eventType, func(event *core.ScheduledEvent, info string) error {
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
				go a.sendScheduledAsync(context.Background(), event.Scheduled, event.ID, info.Source, info.ClientId, info.Message)
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

	user, err := a.ClientRegistry.UserForClient(ctx, source, clientId)
	if err != nil {
		return err
	}

	clients, err := a.ClientRegistry.ClientsForUser(ctx, user.Id, nil)
	if err != nil {
		return err
	}

	skipIndex := -1
	for i, client := range clients {
		if client.Source == source && client.Id == clientId {
			skipIndex = i
			break
		}
	}
	if skipIndex < 0 {
		return errors.New("client id does not have a cluster")
	}

	for i, client := range clients {
		if i == skipIndex {
			continue
		}
		err = a.sendToClent(ctx, client.Id, message)
		if err != nil && err != ErrorClientDoesNotSupportPush {
			return err
		}
	}

	return ErrorClientDoesNotSupportPush
}

func (a *Push) sendToClent(ctx context.Context, clientId string, message core.PushMessage) error {
	for _, provider := range a.ClientSenders {
		ok, err := provider.SendToClient(ctx, clientId, message)
		if ok {
			err = a.PushLogger.LogPush(ctx, clientId, &message)
			if err != nil {
				return err
			}
			return nil
		}
		if err != nil {
			return err
		}

	}
	return ErrorClientDoesNotSupportPush
}
