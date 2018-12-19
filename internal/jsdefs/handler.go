package jsdefs

import "github.com/kataras/iris/context"

func Handler(ctx context.Context) {
	ctx.JSON(GetDefs())
}
