package service

import (
	"hal9000/types"
	"os"
	"strings"
	"sync"
)

type alertQueueConcrete struct {
	mutex sync.Mutex
	queue []types.ResponseMessage
}

func InitAlertQueue(runtime *types.Runtime) (types.AlertQueue, error) {
	q := alertQueueConcrete{
		mutex: sync.Mutex{},
		queue: make([]types.ResponseMessage, 0),
	}

	alertedUserNames := strings.Split(os.Getenv("ALERTED_USER_NAMES"), ",")
	users := make([]*types.Person, len(alertedUserNames))
	for i, alertedUserName := range alertedUserNames {
		person, err := (*(*runtime).People()).GetPersonByName(alertedUserName)
		if err != nil {
			return nil, err
		}
		users[i] = person
	}
	go (func() {
		for {
			q.mutex.Lock()
			for _, m := range q.queue {
				for _, user := range users {
					err := (*(*runtime).People()).SendMessageToPerson(runtime, user, m)
					if err != nil {
						(*(*runtime).Logger()).LogError(err)
					}
				}
			}
			q.queue = make([]types.ResponseMessage, 0)
			q.mutex.Unlock()
		}
	})()

	return &q, nil
}

func (q *alertQueueConcrete) Enqueue(m types.ResponseMessage) {
	q.mutex.Lock()
	q.queue = append(q.queue, m)
	q.mutex.Unlock()
}
