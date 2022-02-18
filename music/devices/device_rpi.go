//+build linux

package devices

import (
	"fmt"
	"github.com/hajimehoshi/oto/v2"
	"github.com/salehkhazaei/malgo"
	"runtime"
	"time"
	"unsafe"
)

type IDevice interface {
	Capture(sampleRate uint32, channels uint32, samplingChannel chan []byte, doneChannel chan interface{}) error
	Play(channels uint32, sampleRate uint32, sampleReadingChannel chan []byte, doneChannel chan interface{}) error
}

type Device struct {
	IDPtr         unsafe.Pointer
	RealID        malgo.DeviceID
	ID            string
	Name          string
	DeviceType    malgo.DeviceType
	ChannelsMin   uint32
	ChannelsMax   uint32
	SampleRateMin uint32
	SampleRateMax uint32
	IsDefault     bool
}

func (d *Device) String() string {
	deviceTypeStr := ""
	switch d.DeviceType {
	case malgo.Playback:
		deviceTypeStr = "Playback"
	case malgo.Capture:
		deviceTypeStr = "Capture"
	case malgo.Duplex:
		deviceTypeStr = "Duplex"
	case malgo.Loopback:
		deviceTypeStr = "Loopback"

	}
	return fmt.Sprintf("ID: %s, Name: %s, DeviceType: %s, Channels: %d-%d, SampleRate: %d-%d, IsDefault: %v",
		d.ID, d.Name, deviceTypeStr, d.ChannelsMin, d.ChannelsMax, d.SampleRateMin, d.SampleRateMax, d.IsDefault)
}

func (d *Device) Capture(sampleRate uint32, channels uint32, samplingChannel chan []byte, doneChannel chan interface{}) error {
	fmt.Println("Capturing device", d.Name)

	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, func(message string) {
		fmt.Printf("LOG <%v>\n", message)
	})
	if err != nil {
		return err
	}
	defer func() {
		_ = ctx.Uninit()
		ctx.Free()
	}()

	deviceType := d.DeviceType
	if deviceType == malgo.Playback {
		deviceType = malgo.Loopback
	}

	deviceConfig := malgo.DeviceConfig{
		DeviceType: deviceType,
		SampleRate: sampleRate,
		Capture: malgo.SubConfig{
			DeviceID: d.IDPtr,
			Format:   malgo.FormatS16,
			Channels: channels,
		},
		Playback: malgo.SubConfig{
			DeviceID: d.IDPtr,
			Format:   malgo.FormatS16,
			Channels: channels,
		},
	}

	onRecvFrames := func(pSample2, pSample []byte, frameCount uint32) {
		samplingChannel <- pSample
	}

	captureCallbacks := malgo.DeviceCallbacks{
		Data: onRecvFrames,
	}
	device, err := malgo.InitDevice(ctx.Context, deviceConfig, captureCallbacks)
	if err != nil {
		return err
	}

	err = device.Start()
	if err != nil {
		return err
	}

	<-doneChannel
	device.Uninit()

	return nil
}

func (d *Device) Play(channels uint32, sampleRate uint32, sampleReadingChannel chan []byte, doneChannel chan interface{}) error {
	fmt.Println("Playing device", d.Name)

	ctx, ready, err := oto.NewContext(int(sampleRate), int(channels), 2)
	if err != nil {
		return err
	}
	<-ready
	fmt.Println("Device is ready", d.Name)

	ch := ChannelReader{Channel: sampleReadingChannel}
	p := ctx.NewPlayer(ch)
	p.Play()
	runtime.KeepAlive(p)

	go func() {
		for {
			fmt.Println("===", p.IsPlaying(), p.UnplayedBufferSize(), p.Err())
			time.Sleep(2 * time.Second)
		}
	}()

	fmt.Println("Playing on", d.Name)
	<-doneChannel
	fmt.Println("DONE with device", d.Name)
	return nil
}
