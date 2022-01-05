package chap04

import (
	"context"
	"fmt"
	"log"
	"time"
)

func printGreeting(done <-chan interface{}) error {
	greeting, err := genGreeting(done)
	if err != nil {
		return err
	}
	log.Printf("%s world!\n", greeting)
	return nil
}

func printFarewell(done <-chan interface{}) error {
	farewell, err := genFarewell(done)
	if err != nil {
		return err
	}
	log.Printf("%s world!\n", farewell)
	return nil
}


func genGreeting(done <-chan interface{}) (string, error) {
	switch locale, err := locale(done); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "hello", nil
	}
	return "", fmt.Errorf("unsupported locale")
}

func genFarewell(done <-chan interface{}) (string, error) {
	switch locale, err := locale(done); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "goodbye", nil
	}
	return "", fmt.Errorf("unsupported locale")
}

func locale(done <-chan interface{}) (string, error) {
	select {
	case <-done:
		return "", fmt.Errorf("canceld")
	case <-time.After(1*time.Minute):
	}
	return "EN/US", nil
}

// 以下はctx使った場合の例


func printGreetingWithCtx(ctx context.Context) error {
	greeting, err := genGreetingWithCtx(ctx)
	if err != nil {
		return err
	}
	log.Printf("%s world!\n", greeting)
	return nil
}

func printFarewellWithCtx(ctx context.Context) error {
	farewell, err := genFarewellWithCtx(ctx)
	if err != nil {
		return err
	}
	log.Printf("%s world!\n", farewell)
	return nil
}


func genGreetingWithCtx(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	switch locale, err := localeWithCtx(ctx); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "hello", nil
	}
	return "", fmt.Errorf("unsupported locale")
}

func genFarewellWithCtx(ctx context.Context) (string, error) {
	switch locale, err := localeWithCtx(ctx); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "goodbye", nil
	}
	return "", fmt.Errorf("unsupported locale")
}

func localeWithCtx(ctx context.Context) (string, error) {
	if deadline, ok := ctx.Deadline(); ok {
		if deadline.Sub(time.Now().Add(1*time.Minute)) <= 0 {
			return "", context.DeadlineExceeded
		}
	}

	select {
	case <-ctx.Done():
		return "", fmt.Errorf("canceld")
	case <-time.After(1*time.Minute):
	}
	return "EN/US", nil
}

// Context内でデータの受け渡し

// パッケージ内で独自の型を定義するのが安全に使うにはいいらしい
type ctxKey int

const (
	ctxUserID ctxKey = iota
	ctxAuthToken
)

func UserID(c context.Context) string {
	return c.Value(ctxUserID).(string)
}

func AuthToken(c context.Context) string {
	return c.Value(ctxAuthToken).(string)
}


func ProcessRequest(userID, authToken string) {
	ctx := context.WithValue(context.Background(), ctxUserID, userID)
	ctx = context.WithValue(ctx, ctxAuthToken, authToken)
	HandleResponse(ctx)
}

func HandleResponse(ctx context.Context) {
	log.Printf(
		"handling response for %v (%v)",
		UserID(ctx),
		AuthToken(ctx),
	)
}