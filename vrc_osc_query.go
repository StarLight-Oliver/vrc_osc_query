package vrc_osc_query

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/StarLight-Oliver/go-osc/osc"
	"github.com/StarLight-Oliver/vrc_osc_query/osc_query"
	"github.com/grandcat/zeroconf"
)

type VRCOSCService struct {
	OscPort  int
	OscTree  *osc_query.OscNode
	HttpPort int
	oscD     *osc.StandardDispatcher
}

const (
	OscTypeInt    = 1
	OscTypeFloat  = 2
	OscTypeBool   = 3
	OscTypeString = 4
)

func IsPortInUse(port int, netType string) bool {

	if netType == "" {
		netType = "tcp"
	}

	l, err := net.Listen(netType, fmt.Sprintf("localhost:%d", port))
	if err != nil {
		fmt.Println(err)
		return true
	}
	defer l.Close()
	return false
}

func GetFreeTCPPort() (port int, err error) {
	var a *net.TCPAddr
	if a, err = net.ResolveTCPAddr("tcp", "localhost:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer l.Close()
			return l.Addr().(*net.TCPAddr).Port, nil
		}
	}
	return
}

func GetFreeUDPPort() (port int, err error) {
	var a *net.UDPAddr
	if a, err = net.ResolveUDPAddr("udp", "localhost:0"); err == nil {
		var l *net.UDPConn
		if l, err = net.ListenUDP("udp", a); err == nil {
			defer l.Close()
			return l.LocalAddr().(*net.UDPAddr).Port, nil
		}
	}
	return
}

func NewVRCOSCService(name string, oscport int) (*VRCOSCService, error) {

	if oscport == 0 {
		newOscport, err := GetFreeUDPPort()
		if err != nil {
			return nil, err
		}
		oscport = newOscport
	}

	OscTree := osc_query.NewOscNodeTree(name, "127.0.0.1", oscport)

	httpPort := oscport

	// see if the tcp port is already in use
	if IsPortInUse(oscport, "tcp") {
		// find a free port
		tHttpPort, err := GetFreeTCPPort()
		if err != nil {
			return nil, err
		}
		httpPort = tHttpPort
	}

	return &VRCOSCService{
		OscPort:  oscport,
		OscTree:  OscTree,
		HttpPort: httpPort,
		oscD:     osc.NewStandardDispatcher(),
	}, nil
}

func (osc_service *VRCOSCService) ListenAndServe() error {
	OscTree := osc_service.OscTree

	r := http.NewServeMux()

	hostName, hostErr := os.Hostname()
	if hostErr != nil {
		return hostErr
	}

	server, err := zeroconf.RegisterProxy(OscTree.HostInfo.Name, "_oscjson._tcp", "local.", osc_service.HttpPort, hostName, []string{"127.0.0.1"}, []string{"txtver=1"}, nil)
	server2, err2 := zeroconf.RegisterProxy(OscTree.HostInfo.Name, "_osc._udp", "local.", osc_service.OscPort, hostName, []string{"127.0.0.1"}, []string{"txtvers=1"}, nil)
	if err != nil {
		return err
	}

	if err2 != nil {
		return err2
	}

	defer server.Shutdown()
	defer server2.Shutdown()

	r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if _, ok := r.URL.Query()["HOST_INFO"]; ok {
			host_info := OscTree.HostInfo
			host_info_json, err := json.Marshal(host_info)
			if err != nil {
				fmt.Println(err)
			}
			w.Write(host_info_json)
			return
		}

		json_string, err := json.Marshal(OscTree)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
			return
		}

		w.Write(json_string)
	})

	sig := make(chan os.Signal, 1)
	errSignal := make(chan error, 1)
	go func() {
		err := http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", osc_service.HttpPort), r)
		if err != nil {
			errSignal <- err
		}

	}()

	go func() {
		server := &osc.Server{
			Addr:       fmt.Sprintf("127.0.0.1:%d", osc_service.OscPort),
			Dispatcher: osc_service.oscD,
		}

		err := server.ListenAndServe()
		if err != nil {
			errSignal <- err
		}

	}()

	defer server.Shutdown()
	defer server2.Shutdown()

	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	select {
	case <-sig:
		break
	case err := <-errSignal:
		return err
	}

	return nil
}

type Message osc.Message

func (osc_service *VRCOSCService) AddHandler(fullpath string, value int, desc string, method func(*Message)) {
	osc_service.OscTree.AddChild(fullpath, value, desc)
	osc_service.oscD.AddMsgHandler(fullpath, func(msg *osc.Message) {

		message := Message{
			Arguments: msg.Arguments,
			Address:   msg.Address,
		}

		method(&message)
	})
}
