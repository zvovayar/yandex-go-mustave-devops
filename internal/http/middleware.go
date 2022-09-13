package http

import (
	"log"
	"net"
	"net/http"
	"strings"

	inst "github.com/zvovayar/yandex-go-mustave-devops/internal/storage"
)

func TrustIpFilter(next http.Handler) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {

		// trust all connections
		if inst.TrustedSubnet == "" {
			next.ServeHTTP(w, r)
			return
		}

		rip := net.ParseIP(strings.Split(r.RemoteAddr, ":")[0])
		_, tnet, err := net.ParseCIDR(inst.TrustedSubnet)

		if err != nil {
			log.Fatal(err)
			return
		}

		// allow ip
		if tnet.Contains(rip) {
			inst.Sugar.Infow("Allow TrustIpFilter RemoteAddr=" + r.RemoteAddr +
				" ip=" + rip.String() +
				" TrustSubnet=" + inst.TrustedSubnet +
				" trust net=" + tnet.String())
			next.ServeHTTP(w, r)
			return
		}

		// forbidden
		inst.Sugar.Infow("Forbidden TrustIpFilter RemoteAddr=" + r.RemoteAddr +
			" ip=" + rip.String() +
			" TrustSubnet=" + inst.TrustedSubnet +
			" trust net=" + tnet.String())

		w.WriteHeader(http.StatusForbidden)
		_, err = w.Write([]byte("<h1>access forbidden</h1>"))
		if err != nil {
			log.Fatal(err)
			return
		}

	}

	return http.HandlerFunc(fn)
}
