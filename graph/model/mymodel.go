package model

type MyError struct {
	Message string
}

func (m *MyError) Error() string {
	return m.Message
}

func (m *MyError) ReturnError() error {
	return m
}
