package lib

import (
	// "github.com/hprose/hprose-golang/rpc"
	"context"
	etcd "github.com/coreos/etcd/client"
	"github.com/hprose/hprose-golang/rpc"
	"path"
	"strings"
	"sync"
	"time"
)

var serviceMap *map[string]backend
var lk *sync.RWMutex

func init() {
	lk = &sync.RWMutex{}
}

func clientInit() {
	go get_service_map()
}

func Client(sv string) *rpc.TCPClient {
	lk.RLock()
	defer lk.RUnlock()
	if serviceMap != nil {
		if backend, ok := (*serviceMap)[sv]; ok {
			return backend.Client
		}
	}

	return nil
}

func testEq(a, b []string) bool {

	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func get_service_map() {
	for {
		func() {
			new_service_map := make(map[string]backend)

			kAPI := Cfg.KApi()
			rsp, err := kAPI.Get(context.Background(), "/nodes", &etcd.GetOptions{
				Recursive: true,
			})
			if err != nil {
				return
			}
			for _, srv_n := range rsp.Node.Nodes {
				for _, backend_n := range srv_n.Nodes {
					paths := strings.Split(backend_n.Key, "/")
					srv_name := paths[len(paths)-2]
					endpoint := path.Base(backend_n.Key)

					if _, ok := new_service_map[srv_name]; !ok {
						new_service_map[srv_name] = backend{
							Nodes: []string{},
						}
					}
					nodes := append(new_service_map[srv_name].Nodes, "tcp4://"+endpoint)
					backend := new_service_map[srv_name]
					backend.Nodes = nodes

					if serviceMap == nil {
						backend.Client = rpc.NewTCPClient(backend.Nodes...)
					} else if c, ok := (*serviceMap)[srv_name]; ok {
						backend.Client = c.Client
						if testEq(c.Nodes, backend.Nodes) == false {
							backend.Client.SetURIList(backend.Nodes)
						}
					} else {
						backend.Client = rpc.NewTCPClient(backend.Nodes...)
					}
					backend.Client = rpc.NewTCPClient(backend.Nodes...)
					new_service_map[srv_name] = backend
				}
			}

			lk.Lock()
			defer lk.Unlock()

			serviceMap = &new_service_map
		}()
		time.Sleep(time.Second)
	}
}

type backend struct {
	Client *rpc.TCPClient
	Nodes  []string
}
