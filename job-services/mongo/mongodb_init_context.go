package mongo

import (
	"context"
	"time"
)

func InitContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(60000)*time.Millisecond)
}
