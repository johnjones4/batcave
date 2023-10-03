package push

import (
	"context"
	"encoding/json"
	"errors"
	"main/core"
	"time"

	"github.com/google/uuid"
	"github.com/gorhill/cronexpr"
	"github.com/sirupsen/logrus"
)

type Push struct {
	ClientSenders     []core.ClientSender
	ClientRegistry    core.ClientRegistry
	Scheduler         core.Scheduler
	Log               logrus.FieldLogger
	PushLogger        core.PushLogger
	PushIntentFactory core.PushIntentFactory
}

type pushEventInfo struct {
	Message core.PushMessage `json:"message"`
}

var (
	ErrorClientDoesNotSupportPush = errors.New("client id does not support push")
)

const (
	eventType              = "push"
	singleCheckInterval    = time.Second * 30
	recurringCheckInterval = time.Minute
)

func (a *Push) SendRecurring(ctx context.Context, source string, clientId string, schedule string, intent string, info map[string]any) error {
	return a.Scheduler.ScheduleRecurringEvent(ctx, &core.ScheduledRecurringEvent{
		ScheduledEventCore: core.ScheduledEventCore{
			Source:   source,
			ClientId: clientId,
		},
		Info:      info,
		Intent:    intent,
		Scheduled: schedule,
		LastRun:   time.Now(),
	})
}

func (a *Push) SendLater(ctx context.Context, when time.Time, source string, clientId string, message core.PushMessage) error {
	wait := time.Until(when)
	if wait <= singleCheckInterval {
		go a.sendScheduledAsync(context.Background(), when, "", source, clientId, message)
		return nil
	} else {
		return a.Scheduler.ScheduleEvent(ctx, &core.ScheduledEvent{
			ScheduledEventCore: core.ScheduledEventCore{
				Source:   source,
				ClientId: clientId,
			},
			Info: pushEventInfo{
				Message: message,
			},
			EventType: eventType,
			Scheduled: when,
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

func (a *Push) doSingleEvents(ctx context.Context) error {
	events, err := a.Scheduler.ReadyEvents(ctx, time.Now().Add(singleCheckInterval), eventType, func(event *core.ScheduledEvent, info string) error {
		var receiver pushEventInfo
		err := json.Unmarshal([]byte(info), &receiver)
		if err != nil {
			return err
		}
		event.Info = receiver
		return nil
	})
	if err != nil {
		return err
	}

	for _, event := range events {
		info, ok := event.Info.(pushEventInfo)
		if !ok {
			continue
		}
		go a.sendScheduledAsync(context.Background(), event.Scheduled, event.ID, event.Source, event.ClientId, info.Message)
	}

	return nil
}

func (a *Push) doRecurringEvents(ctx context.Context) error {
	events, err := a.Scheduler.RecurringEvents(ctx)
	if err != nil {
		return err
	}

	limit := time.Now().Add(recurringCheckInterval)

	for _, event := range events {
		parsed, err := cronexpr.Parse(event.Scheduled)
		if err != nil {
			a.Log.Error(err)
			continue
		}
		nextTime := parsed.Next(event.LastRun)
		a.Log.Debug(event, nextTime)
		if nextTime.Before(limit) {
			intent := a.PushIntentFactory.PushIntent(event.Intent)
			if intent != nil {
				push, err := intent.ActOnAsyncIntent(ctx, event.Source, event.ClientId, &core.IntentMetadata{
					IntentParseReceiver: event.Info,
				})
				if err != nil {
					a.Log.Error(err)
					continue
				}
				go a.sendScheduledAsync(context.Background(), nextTime, "", event.Source, event.ClientId, push)
				err = a.Scheduler.UpdateRecurringEventTimestamp(ctx, event.ID, time.Now())
				if err != nil {
					a.Log.Error(err)
					continue
				}
			}
		}
	}

	return nil
}

func (a *Push) Start(ctx context.Context) error {
	singleTicker := time.NewTicker(singleCheckInterval)
	recurringTicker := time.NewTicker(recurringCheckInterval)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-singleTicker.C:
			err := a.doSingleEvents(ctx)
			if err != nil {
				a.Log.Errorf("Error processing scheduled events: %e", err)
				continue
			}
		case <-recurringTicker.C:
			err := a.doRecurringEvents(ctx)
			if err != nil {
				a.Log.Errorf("Error processing recurring events: %e", err)
				continue
			}
		}
	}
}

func (a *Push) Send(ctx context.Context, source string, clientId string, message core.PushMessage) error {
	user, err := a.ClientRegistry.UserForClient(ctx, source, clientId)
	if err != nil {
		return err
	}

	clients, err := a.ClientRegistry.ClientsForUser(ctx, user.Id, nil)
	if err != nil {
		return err
	}

	sends := 0
	for _, client := range clients {
		message.EventId = uuid.NewString()
		//TODO make better use of source
		err = a.sendToClent(ctx, client.Id, message)
		if err != nil && err != ErrorClientDoesNotSupportPush {
			return err
		}
		if err == nil || err == ErrorClientDoesNotSupportPush {
			sends++
		}
	}

	if sends == 0 {
		return ErrorClientDoesNotSupportPush
	}

	return nil
}

func (a *Push) sendToClent(ctx context.Context, clientId string, message core.PushMessage) error {
	sends := 0
	for _, provider := range a.ClientSenders {
		ok, err := provider.SendToClient(ctx, clientId, message)
		if ok {
			err = a.PushLogger.LogPush(ctx, clientId, &message)
			if err != nil {
				return err
			}
			sends++
		}
		if err != nil {
			return err
		}
	}
	if sends == 0 {
		return ErrorClientDoesNotSupportPush
	}
	return nil
}
