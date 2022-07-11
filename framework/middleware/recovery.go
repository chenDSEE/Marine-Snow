package middleware

import (
	"MarineSnow/framework"
	"fmt"
)

func Recovery() framework.HandlerFunc {
	return func(ctx *framework.Context) error {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("Recovery: ", err)
			}
		}()

		// call next Handler
		ctx.NextHandler()
		return nil
	}
}
