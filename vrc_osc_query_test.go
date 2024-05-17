package vrc_osc_query_test

import (
	"net"
	"testing"

	vrc_osc_query "github.com/StarLight-Oliver/vrc_osc_query"
	"github.com/stretchr/testify/assert"
)

func TestIsPortInUse(t *testing.T) {
	l, _ := net.Listen("tcp", "localhost:0")
	defer l.Close()

	port := l.Addr().(*net.TCPAddr).Port
	assert.True(t, vrc_osc_query.IsPortInUse(port, "tcp"))
	assert.False(t, vrc_osc_query.IsPortInUse(port+1, "tcp"))
}

func TestGetFreeTCPPort(t *testing.T) {
	port, err := vrc_osc_query.GetFreeTCPPort()
	assert.Nil(t, err)
	assert.False(t, vrc_osc_query.IsPortInUse(port, "tcp"))
}

func TestGetFreeUDPPort(t *testing.T) {
	// This test is currently failing, due to how we test for port availability
	port, err := vrc_osc_query.GetFreeUDPPort()
	assert.Nil(t, err)
	assert.False(t, vrc_osc_query.IsPortInUse(port, "udp"))
}

func TestNewVRCOSCService(t *testing.T) {
	service, err := vrc_osc_query.NewVRCOSCService("TestService", 0)
	assert.Nil(t, err)
	assert.Equal(t, "TestService", service.OscTree.HostInfo.Name)
	assert.Equal(t, "127.0.0.1", service.OscTree.HostInfo.Ip)
	assert.Equal(t, "UDP", service.OscTree.HostInfo.Transport)
	assert.True(t, service.OscTree.HostInfo.Extensions["ACCESS"])
}

func TestAddHandler(t *testing.T) {
	service, _ := vrc_osc_query.NewVRCOSCService("TestService", 0)
	service.AddHandler("/test", 1, "test desc", func(msg *vrc_osc_query.Message) {})

	assert.Equal(t, "test desc", service.OscTree.Contents["test"].Description)
	assert.Equal(t, "/test", service.OscTree.Contents["test"].FullPath)
	assert.Equal(t, "i", service.OscTree.Contents["test"].Type)
}
