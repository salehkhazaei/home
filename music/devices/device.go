//+build !linux

package devices

import (
	"fmt"
	"github.com/salehkhazaei/malgo"
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

	deviceConfig := malgo.DeviceConfig{
		DeviceType: d.DeviceType,
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

	// This is the function that's used for sending more data to the device for playback.
	onSamples := func(pOutputSample, pInputSamples []byte, framecount uint32) {
		select {
		case bytes := <-sampleReadingChannel:
			copy(pOutputSample, bytes)
		default:
			//fmt.Println("missed one")
		}
	}

	deviceCallbacks := malgo.DeviceCallbacks{
		Data: onSamples,
	}

	fmt.Println("Init device", d.Name)
	device, err := malgo.InitDevice(ctx.Context, deviceConfig, deviceCallbacks)
	if err != nil {
		return err
	}
	defer device.Uninit()

	fmt.Println("Starting device", d.Name)
	err = device.Start()
	if err != nil {
		fmt.Println("Error on starting device", err)
		return err
	}

	fmt.Println("Playing on device", d.Name)
	<-doneChannel
	fmt.Println("DONE with device", d.Name)
	return nil
}
