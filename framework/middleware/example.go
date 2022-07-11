package middleware

import (
	"MarineSnow/framework"
	"fmt"
)

func Example1() framework.HandlerFunc {
	return func(ctx *framework.Context) error {
		fmt.Println("Before middleware Example1()")
		if err := ctx.NextHandler(); err != nil { // call next Handler
			return err
		}
		fmt.Println("After middleware Example1()")

		return nil
	}
}

func Example2() framework.HandlerFunc {
	return func(ctx *framework.Context) error {
		fmt.Println("Before middleware Example2()")
		if err := ctx.NextHandler(); err != nil { // call next Handler
			return err
		}
		fmt.Println("After middleware Example2()")

		return nil
	}
}

func Example3() framework.HandlerFunc {
	return func(ctx *framework.Context) error {
		fmt.Println("Before middleware Example3()")
		if err := ctx.NextHandler(); err != nil { // call next Handler
			return err
		}
		fmt.Println("After middleware Example3()")

		return nil
	}
}
