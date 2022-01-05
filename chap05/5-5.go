package chap05

import "context"

// Dummy API with two endpoints

type APIConnection struct {}

func Open() *APIConnection {
	return &APIConnection{}
}

func (a *APIConnection) ReadFile(ctx context.Context) error {
	// Do something
	return nil
}

func (a *APIConnection) ResolveAddress(ctx context.Context) error {
	return nil
}


