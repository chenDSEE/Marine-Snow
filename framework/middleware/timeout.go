package middleware

import (
	"MarineSnow/framework"
	"context"
	"fmt"
	"time"
)

// FIXME: how to handle the data had been write into framework.Context.ResponseWriter before timeout ?
// FIXME: can not call next handler after ctx.NextHandler() in the main call flow
func Timeout(duration time.Duration) framework.HandlerFunc {
	return func(ctx *framework.Context) error {
		panicChan := make(chan interface{}, 1) // 1 buffer channel not to block goroutine
		finishChan := make(chan struct{}, 1)

		timeCtx, cancel := context.WithTimeout(ctx.BaseContext(), duration)
		defer cancel()

		/* do work in another goroutine within time limit */
		// You can only write to framework.Context.ResponseWriter after you are sure that you can do it all
		go func() {
			defer func() {
				if err := recover(); err != nil {
					panicChan <- err
				}
			}()

			// call real handler
			ctx.NextHandler()
			finishChan <- struct{}{}
		}()

		/* monitor the new goroutine */
		select {
		case <-panicChan:
			fmt.Println("handler goroutine meet a panic and exit")
		case <-finishChan:
			fmt.Println("call finish")
		case <-timeCtx.Done():
			fmt.Println("handler timeout")
		}

		return nil
	}
}
