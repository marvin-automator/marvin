package jsdefs

import "github.com/kataras/iris/v12/context"

func Handler(ctx context.Context) {
	ctx.JSON(GetDefs())
}
