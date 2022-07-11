package middleware

import (
	"MarineSnow/framework"
	"fmt"
	"time"
)

func Cost() framework.HandlerFunc {
	return func(ctx *framework.Context) error {
		startTime := time.Now()

		ctx.NextHandler()

		endTime := time.Now()
		fmt.Printf("Time cost:[%v us]\n", endTime.Sub(startTime).Microseconds())
		return ctx.NextHandler()
	}
}
