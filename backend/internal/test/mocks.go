package test

import "fmt"

type SuccessClient struct{ Calls []string }

func (c *SuccessClient) SendMessage(_, _, text string) error {
	c.Calls = append(c.Calls, text)
	return nil
}

type FailClient struct{}

func (c *FailClient) SendMessage(_, _, _ string) error {
	return fmt.Errorf("telegram unavailable")
}
