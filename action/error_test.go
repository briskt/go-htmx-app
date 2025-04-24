package action

import (
	"net"
	"net/http"
)

func (s *Suite) Test_getClientIPAddress() {
	request := &http.Request{
		Header: make(http.Header),
	}

	got, err := getClientIPAddress(request)
	s.Error(err)
	s.Nil(got)

	request.RemoteAddr = "192.168.100.1"
	got, err = getClientIPAddress(request)
	s.Error(err)
	s.Nil(got)

	request.RemoteAddr = "192.168.100.1:1234"
	got, err = getClientIPAddress(request)
	s.NoError(err)
	s.Equal(net.ParseIP("192.168.100.1"), got)

	request.Header.Add("CF-Connecting-IP", "192.168.0.1")
	got, err = getClientIPAddress(request)
	s.NoError(err)
	s.Equal(net.ParseIP("192.168.0.1"), got)
}
