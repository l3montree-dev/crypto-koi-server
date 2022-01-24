package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"gitlab.com/l3montree/cryptogotchi/clodhopper/graph/generated"
	"gitlab.com/l3montree/cryptogotchi/clodhopper/graph/model"
)

func (r *mutationResolver) HandleNewEvent(ctx context.Context, event model.NewEvent) (*model.Cryptogotchi, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Cryptogotchis(ctx context.Context) ([]*model.Cryptogotchi, error) {
	return r.cryptogotchiResolver.Cryptogotchis(ctx)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
