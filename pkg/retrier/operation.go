package retrier

import "time"

type Try func() error

func Do(number uint8, duration time.Duration, try Try) error {
	var err error

	for range number {
		err = try()

		if err == nil {
			break
		}
	}

	return err
}
