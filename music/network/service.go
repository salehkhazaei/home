package network

import (
	"fmt"
	"net"
)

type Device struct {
	l net.Listener
	c *net.TCPConn

	sampleChannel chan []byte
}

func Bind(address string) (*Device, error) {
	l, err := net.Listen("tcp4", address)
	if err != nil {
		return nil, err
	}

	return &Device{
		l: l,
	}, nil
}

func Connect(address string) (*Device, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}

	return &Device{
		c: conn,
	}, nil
}

func (s *Device) Capture(sampleRate uint32, channels uint32, samplingChannel chan []byte, doneChannel chan interface{}) error {
	sampleSize := sampleRate * channels
	for {
		fmt.Println("Listening for socket")
		c, err := s.l.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("Accepted new socket")

		buf := make([]byte, sampleSize * channels / 100 * 2)
		for {
			n, err := c.Read(buf)
			if err != nil {
				c.Close()
				break
			}

			if n < 0 {
				c.Close()
				break
			}

			if n == 0 {
				continue
			}

			samplingChannel <- buf[:n]

			select {
			case <-doneChannel:
				c.Close()
				return nil
			default:
			}
		}

		select {
		case <-doneChannel:
			return nil
		default:
		}
	}
}

func (s *Device) Play(channels uint32, sampleRate uint32, sampleReadingChannel chan []byte, doneChannel chan interface{}) error {
	// sending to other server to play
	for sample := range sampleReadingChannel {
		n, err := s.c.Write(sample)
		if err != nil {
			return err
		}

		if n != len(sample) {
			fmt.Println("What the hell 1?", n, len(sample))
			s.c.Close()
			return nil
		}

		select {
		case <-doneChannel:
			s.c.Close()
			return nil
		default:
		}
	}
	return nil
}
