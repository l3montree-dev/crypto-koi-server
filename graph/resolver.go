package graph

import resolver "gitlab.com/l3montree/cryptogotchi/clodhopper/internal/resolvers"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	cryptogotchiResolver *resolver.CryptogotchiResolver
}

func NewResolver(cryptogotchiResolver *resolver.CryptogotchiResolver) Resolver {
	return Resolver{
		cryptogotchiResolver: cryptogotchiResolver,
	}
}
