package devices

import (
	"github.com/salehkhazaei/malgo"
	"sync"
)

var once sync.Once
var instance *Service

type Service struct {

}

func Instance () *Service {
	once.Do(func() {
		instance = &Service{}
	})

	return instance
}

func (s *Service) FindDeviceById(devId string) (*Device, error) {
	devs, err := s.ListPlaybackDevices()
	if err != nil {
		return nil, err
	}

	for _, dev := range devs {
		if dev.ID == devId {
			return dev, nil
		}
	}

	devs, err = s.ListCaptureDevices()
	if err != nil {
		return nil, err
	}

	for _, dev := range devs {
		if dev.ID == devId {
			return dev, nil
		}
	}

	return nil, nil
}

func (s *Service) ListPlaybackDevices() ([]*Device, error) {
	return s.ListDevices(malgo.Playback)
}
func (s *Service) ListCaptureDevices() ([]*Device, error) {
	return s.ListDevices(malgo.Capture)
}
func (s *Service) ListDuplexDevices() ([]*Device, error) {
	return s.ListDevices(malgo.Duplex)
}
func (s *Service) ListLoopbackDevices() ([]*Device, error) {
	return s.ListDevices(malgo.Loopback)
}

func (s *Service) ListDevices(deviceType malgo.DeviceType) ([]*Device, error) {
	context, err := malgo.InitContext(nil, malgo.ContextConfig{}, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = context.Uninit()
		context.Free()
	}()

	infos, err := context.Devices(deviceType)
	if err != nil {
		return nil, err
	}

	var devices []*Device
	for _, info := range infos {
		full, err := context.DeviceInfo(deviceType, info.ID, malgo.Shared)
		if err != nil {
			return nil, err
		}

		isDefault := false
		if info.IsDefault > 0 {
			isDefault = true
		}

		devices = append(devices, &Device{
			IDPtr: full.ID.Pointer(),
			ID: BuildID(RemoveAfterNil(full.ID.String())),
			RealID: info.ID,
			Name: RemoveAfterNil(info.Name()),
			DeviceType: deviceType,
			ChannelsMin: full.MinChannels,
			ChannelsMax: full.MaxChannels,
			SampleRateMin: full.MinSampleRate,
			SampleRateMax: full.MaxSampleRate,
			IsDefault: isDefault,
		})
	}

	return devices, nil
}