package restutils

import (
	"context"
	"fmt"

	"github.com/Fierro-Labs/Fierro/src/models"
)

func ContextWithToken(jwtToken string) context.Context {
	ctx := context.WithValue(context.Background(), models.TOKEN_KEY, jwtToken)
	return ctx
}

func GetTokenFromContext(ctx context.Context) string {
	return fmt.Sprintf("%s", ctx.Value(models.TOKEN_KEY))
}
