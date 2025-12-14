package wifi_test

import (
	"fmt"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
)

type MockWiFiHandle struct {
	mock.Mock
}

func (m *MockWiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	args := m.Called()

	var ifaces []*wifi.Interface

	if v := args.Get(0); v != nil {
		cast, ok := v.([]*wifi.Interface)
		if !ok {
			return nil, fmt.Errorf("type assertion failed")
		}

		ifaces = cast
	}

	if err := args.Error(1); err != nil {
		return nil, fmt.Errorf("mock error: %w", err)
	}

	return ifaces, nil
}
