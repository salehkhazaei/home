package devices

import "errors"

type ChannelReader struct {
	Channel chan []byte
}

func (c ChannelReader) Read(buf []byte) (int, error) {
	var sample []byte
	select {
		case sample = <- c.Channel:
	default:
		return 0, nil
	}

	if len(sample) > len(buf) {
		return -1, errors.New("larger sample")
	}

	for i := 0; i < len(buf); i ++ {
		buf[i] = 0
	}

	copy(buf, sample)
	return len(sample), nil
}

