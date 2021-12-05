package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/json"
	"pub/graph/generated"
	"pub/graph/model"
	"pub/types"
)

func (r *mutationResolver) SendNotification(ctx context.Context, to string, text string) (*model.SendNotificationResponse, error) {
	m := types.Message{
		Text: text,
	}
	s, err := json.Marshal(m)
	if err != nil {
		return &model.SendNotificationResponse{ReceivedUsers: 0}, err
	}

	// `to` で指定されたチャンネル(=ユーザー)に `text` を送信する
	receivedUsers, err := r.redisClient.Publish(ctx, to, s).Result()
	if err != nil {
		return &model.SendNotificationResponse{ReceivedUsers: 0}, err
	}

	return &model.SendNotificationResponse{ReceivedUsers: int(receivedUsers)}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
