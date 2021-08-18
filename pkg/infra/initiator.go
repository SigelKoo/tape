package infra

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/pkg/errors"
	"golang.org/x/time/rate"
)

const value = "awb3jKtk5bAFsMUCtMpwgo92QDFoaDq32squM9kYMhzvuAfL7y1lScsmNSIZ2ozzbfsrZjLJp8vuCbSVJYkQZvAAnWt7NRtZu4qkIFb5TY9YOd89Qca63VSDEAL5YPrApGUcn30aZbg0aNPhSDiJoAyctRC9KEyjkDOhdyUHox1lk08UCjIoST1wBi58ge63ey8Th2XK5XiUzO3uSvGOQVbrQZIaPkkb6EXA0b1WMNQYonZTcRBUHu5emkncyLnqRP9zgc21clOGKz7vUqHZ1iSHrIzFx5Tp8gtwjAYS5DBC9sFvpHWF0YYhdMX5erbMnL0Ny4YLpYwFMVky5jC26UQi0N0EjaQzLsTBEPh9zpwxg0ISfbVvVxwGQDfi1mtEvGBzONXweVsX0FwQzDnAF0wUhGuIUddzivXnYN5IC4nyimuMkDtJt8CL6mujvT58uqismuBeb1SUAbj4QuduDBJEH0SoDSETDf0HM0vKzeAMDMsdkNfEK19IUW9KmU1jVxWMWs8HlAvKSW5FNRVQPsqe47csntD8FPdFq6dckkaqO0lNv6jjZsS3Um2ZcayijmyNB2IRuAn7Eosv275HHlRGELMNAw5yA0XZzVvcQu50oLBnyWu0KzkRRVzWQlZ58d3REijCF4mLMNXVlbkcEMSYUzGKeWGLpd4EqVJOBhO1gfxqO982OZ70Q8NYKCYYH6d3EWquewsEhrOWokOVKALm60Vm7b8J2kCrsK8nxPIo7KtoJdA1Tb9yCB75TzwTWbwKbVQJ2nZvyqUo7KSYnnOqRt6SENW8op9dy32z8wousJrarj8SnNK7PV3Bl8j9hiQ6HqiKGUVIKArgZO7xPWR6O9EAh1UC8JZEFdLrhJgI9UhL4fPzZcBRIRuwj4uaQ9Pt1KjPTrLTznG3mWexXTvu3yYGtWYdC9WOHKsCg1z5A8lD5dpgoUds9c3lMpNoCzGNjirQRo0fBUlRZ8XfR3iAJgMTgmmJJz916SDqJ5mWpbWNlBjXPX5y9lQtdJA4"

func StartCreateProposal(num int, burst int, r float64, config Config, crypto *Crypto, raw chan *Elements, errorCh chan error) {
	limit := rate.Inf
	ctx := context.Background()
	if r > 0 {
		limit = rate.Limit(r)
	}
	limiter := rate.NewLimiter(limit, burst)

	for i := 0; i < num; i++ {
		k := rand.Intn(10000000000)
		key := fmt.Sprintf("%032d%032d", k, num)
		config.Args[1] = key
		config.Args[2] = value
		prop, err := CreateProposal(
			crypto,
			config.Channel,
			config.Chaincode,
			config.Version,
			config.Args...,
		)
		if err != nil {
			errorCh <- errors.Wrapf(err, "error creating proposal")
			return
		}

		if err = limiter.Wait(ctx); err != nil {
			errorCh <- errors.Wrapf(err, "error creating proposal")
			return
		}

		raw <- &Elements{Proposal: prop}
	}
}
