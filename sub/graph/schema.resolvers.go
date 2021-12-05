package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/json"
	"sub/graph/generated"
	"sub/graph/model"
	"time"

	"github.com/google/uuid"
)

func (r *subscriptionResolver) OnNotificationReceived(ctx context.Context, userName string) (<-chan *model.Notification, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	ch := make(chan *model.Notification, 1)

	u, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	connID := u.String()

	if _, ok := r.observers[userName]; !ok {
		r.observers[userName] = make(map[string]chan<- *model.Notification)
	}
	r.observers[userName][connID] = ch

	pubsub := r.redisClient.Subscribe(ctx, userName)

	go func() {
		pubsubCh := pubsub.Channel()

		for sub := range pubsubCh {
			n := &model.Notification{}
			if err := json.Unmarshal([]byte(sub.Payload), n); err != nil {
				continue
			}
			n.Timestamp = int(time.Now().Unix())

			// メッセージを配信
			r.mutex.Lock()
			for _, observer := range r.observers[userName] {
				observer <- n
			}
			r.mutex.Unlock()
		}
	}()

	go func() {
		<-ctx.Done()
		r.mutex.Lock()
		delete(r.observers[userName], connID)
		r.mutex.Unlock()

		pubsub.Close()
	}()

	return ch, nil
}

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type subscriptionResolver struct{ *Resolver }
