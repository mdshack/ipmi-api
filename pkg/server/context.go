package server

import (
	"context"

	"github.com/mdshack/ipmi-api/pkg/ipmi"
)

// contextKey is a custom type for context keys
type contextKey string

// ipmiClientKey is the context key for the IPMI client
const ipmiClientKey contextKey = "ipmiClient"

// contextWithIPMIClient adds an IPMI client to the context
func contextWithIPMIClient(ctx context.Context, client *ipmi.Client) context.Context {
	return context.WithValue(ctx, ipmiClientKey, client)
}

// ipmiClientFromContext retrieves an IPMI client from the context
func ipmiClientFromContext(ctx context.Context) (*ipmi.Client, bool) {
	client, ok := ctx.Value(ipmiClientKey).(*ipmi.Client)
	return client, ok
}
