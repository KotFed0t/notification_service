package utils

import (
	"context"
	"github.com/google/uuid"
)

func CreateCtxWithRqID(c context.Context) context.Context {
	rqId := uuid.New().String()
	ctx := context.WithValue(c, "rqId", rqId)
	return ctx
}

func GetRequestIdFromCtx(ctx context.Context) string {
	rqId, ok := ctx.Value("rqId").(string)
	if !ok {
		return ""
	}
	return rqId
}
