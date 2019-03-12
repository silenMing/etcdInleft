package lib

import (
	"context"
	"strings"
	"net"
	"log"
	"time"
	etcd "github.com/coreos/etcd/client"
)

type EtcdConfig struct {
	Client        etcd.Client
	Name          string
	MyAddrs       []string
	config        map[string]string
	global_config map[string]string
}

var Cfg *EtcdConfig

func ConnectEtcd(name, etcd_addr, myaddr string) {
	addr_splits := strings.Split(myaddr, ":")
	my_addrs := []string{}
	if addr_splits[0] == "0.0.0.0" {
		ifaces, err := net.Interfaces()
		if err != nil {
			log.Panic(err)
		}
		for _, i := range ifaces {
			addrs, err := i.Addrs()
			if err != nil {
				log.Panic(err)
			}
			for _, addr := range addrs {
				func() {
					var ip string
					switch v := addr.(type) {
					case *net.IPNet:
						ip = v.IP.String()
					case *net.IPAddr:
						ip = v.IP.String()
					}
					if ip != "127.0.0.1" && strings.Contains(ip, ".") && (len(ip) < 8 || ip[0:8] != "169.254.") {
						my_addrs = append(my_addrs, ip+":"+addr_splits[1])
					}
				}()
			}
		}
	} else {
		my_addrs = []string{myaddr}
	}

	Cfg = &EtcdConfig{
		Name:          name,
		MyAddrs:       my_addrs,
		config:        make(map[string]string),
		global_config: make(map[string]string),
	}
	Cfg.connect(etcd_addr)
	//client_init()
}

func (e *EtcdConfig) OnConfigSet(key string, handler func()) {
}

func (e *EtcdConfig) OnGlobalConfigSet(key string, handler func()) {
}

func (e *EtcdConfig) Get(key string) string {
	return e.config[key]
}

func (e *EtcdConfig) GetGlobal(key string) string {
	return e.global_config[key]
}

func (e *EtcdConfig) KApi() etcd.KeysAPI {
	return etcd.NewKeysAPI(e.Client)
}

func (e *EtcdConfig) connect(etcd_addr string) (err error) {
	cfg := etcd.Config{
		Endpoints: []string{"http://" + etcd_addr},
		Transport: etcd.DefaultTransport,
	}

	e.Client, err = etcd.New(cfg)
	if err != nil {
		panic(err)
	}

	err = e.load_env()
	go e.start_heartbeat()
	return
}

func (e *EtcdConfig) start_heartbeat() {
	for {
		time.Sleep(5 * time.Second)
		kAPI := etcd.NewKeysAPI(e.Client)
		for _, addr := range e.MyAddrs {
			_, err := kAPI.Set(context.Background(), "/etcdInLeft/nodes/"+e.Name+"/"+addr, "ok", &etcd.SetOptions{
				TTL: time.Second * 10,
			})

			if err != nil {
				log.Println("etcd-reg", err)
			}
		}
	}
}

func (e *EtcdConfig) load_env() (err error) {
	kAPI := etcd.NewKeysAPI(e.Client)
	if err != nil {
		panic(err)
	}
	e.config, _ = e.LoadConfig(kAPI, "/etcdInLeft/config/"+e.Name)
	e.global_config, _ = e.LoadConfig(kAPI, "/etcdInLeft/config/global")
	return nil
}

func (e *EtcdConfig) LoadConfig(kAPI etcd.KeysAPI, path string) (cfg map[string]string, err error) {
	cfg = make(map[string]string)
	rsp, err := kAPI.Get(context.Background(), path, &etcd.GetOptions{
		Recursive: true,
	})
	if err != nil {
		return
	}
	for _, n := range rsp.Node.Nodes {
		paths := strings.Split(n.Key, "/")
		cfg[paths[len(paths)-1]] = n.Value
	}
	return
}
