package middleware

import "MarineSnow/framework"

func WrapNext(handler framework.HandlerFunc) framework.HandlerFunc {
	return func(ctx *framework.Context) error {
		handler(ctx)

		// forward to next
		ctx.NextHandler()
		return nil
	}
}
