package ch04

import "log"

type Monitor struct {
	*log.Logger
}

func (m Monitor) Write(b []byte) (int, error) {
	err := m.Output(2, string(b))
	if err != nil {
		log.Println(err)
	}

	return len(b), nil
}
