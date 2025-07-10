package retrier

import "time"

func Connect[T any](retry uint8, sleep uint, connector func() (T, error)) (T, error) {
	var (
		out T
		err error
	)

	for range retry {
		out, err = connector()

		if err == nil {
			return out, nil
		}

		time.Sleep(time.Duration(sleep) * time.Second)
	}

	return out, err
}
