// Authored and revised by YOC team, 2016-2018
// License placeholder #1

package http

import (
	"fmt"
	"net/http"

	"github.com/Yocoin15/Yocoin_Sources/log"
)

/*
http roundtripper to register for bzz url scheme
see https://github.com/Yocoin15/Yocoin_Sources/issues/2040
Usage:

import (
 "github.com/Yocoin15/Yocoin_Sources/common/httpclient"
 "github.com/Yocoin15/Yocoin_Sources/swarm/api/http"
)
client := httpclient.New()
// for (private) swarm proxy running locally
client.RegisterScheme("bzz", &http.RoundTripper{Port: port})
client.RegisterScheme("bzz-immutable", &http.RoundTripper{Port: port})
client.RegisterScheme("bzz-raw", &http.RoundTripper{Port: port})

The port you give the Roundtripper is the port the swarm proxy is listening on.
If Host is left empty, localhost is assumed.

Using a public gateway, the above few lines gives you the leanest
bzz-scheme aware read-only http client. You really only ever need this
if you need go-native swarm access to bzz addresses.
*/

type RoundTripper struct {
	Host string
	Port string
}

func (self *RoundTripper) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	host := self.Host
	if len(host) == 0 {
		host = "localhost"
	}
	url := fmt.Sprintf("http://%s:%s/%s:/%s/%s", host, self.Port, req.Proto, req.URL.Host, req.URL.Path)
	log.Info(fmt.Sprintf("roundtripper: proxying request '%s' to '%s'", req.RequestURI, url))
	reqProxy, err := http.NewRequest(req.Method, url, req.Body)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(reqProxy)
}
