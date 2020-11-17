package consul

import (
	"fmt"
	reg "github.com/clearcodecn/carim/registry"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestRegistry_Init(t *testing.T) {
	var opts []Option
	opts = append(opts, Address("http://127.0.0.1:8500"))
	opts = append(opts, TTL(time.Second*30))

	var services = []*reg.Service{
		{
			Id:       "1",
			Name:     "service",
			Endpoint: "192.168.2.48:1111",
			Version:  "v1",
		},
		{
			Id:       "2",
			Name:     "service",
			Endpoint: "192.168.2.48:1112",
			Version:  "v1",
		},
		{
			Id:       "3",
			Name:     "service",
			Endpoint: "192.168.2.48:1113",
			Version:  "v1",
		},
		{
			Id:       "4",
			Name:     "service",
			Endpoint: "192.168.2.48:1114",
			Version:  "v1",
		},
	}

	ctx := RegistryOption(opts...)
	for _, s := range services {
		srv := s
		go func() {
			_, port, _ := net.SplitHostPort(srv.Endpoint)
			mux := http.NewServeMux()
			mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
				writer.Write([]byte(fmt.Sprintf("from service %s", srv.String())))
			})
			http.ListenAndServe(fmt.Sprintf(":%s", port), mux)
		}()

		t.Run(s.Id, func(t *testing.T) {
			r := new(registry)
			if err := r.Init(ctx); err != nil {
				t.Fatal(err)
			}
			if err := r.Register(srv); err != nil {
				t.Fatal(err)
			}
			resp, meta, err := r.client.DiscoveryChain().Get("service", nil, nil)
			if err != nil {
				t.Fatal(err)
				return
			}
			_ = meta
			for k,v := range resp.Chain.Targets {
				fmt.Println(k)
				fmt.Println(v)
			}
		})
	}

	time.Sleep(100 * time.Hour)
}
