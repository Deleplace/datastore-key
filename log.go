package datastorekey

import (
	"context"
	"log"
)

func logf(c context.Context, level, msg string, args ...interface{}) {
	log.Printf(level+" "+msg, args...)
}
