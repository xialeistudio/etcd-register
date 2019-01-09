package discovery

import (
	"go.etcd.io/etcd/client"
	"time"
	"context"
	"log"
)

type Options struct {
	Lifetime time.Duration // 微服务生存时长
	Interval time.Duration // 微服务注册间隔
}

type Service struct {
	client  client.KeysAPI
	ticker  *time.Ticker
	key     string
	val     string
	options *Options
}

func NewService(endpoints []string, key, val string, options *Options) (*Service, error) {
	cfg := client.Config{
		Endpoints:               endpoints,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := client.New(cfg)
	if err != nil {
		return nil, err
	}
	if options == nil {
		options = &Options{
			Lifetime: time.Minute,
			Interval: time.Second * 30,
		}
	}
	s := &Service{
		client:  client.NewKeysAPI(c),
		key:     key,
		val:     val,
		options: options,
	}
	return s, nil
}

// 启动服务
func (s *Service) Start() error {
	err := s.register()
	if err != nil {
		return err
	}

	log.Printf("[discovery] Service Start successfully key:[%s], val:[%s]", s.key, s.val)

	go func() {
		s.ticker = time.NewTicker(s.options.Interval)
		for range s.ticker.C {
			s.register()
		}
	}()
	return nil
}

// 注册服务
func (s *Service) register() error {
	_, err := s.client.Set(context.TODO(), s.key, s.val, &client.SetOptions{
		TTL: s.options.Lifetime,
	})
	if err != nil {
		log.Printf("[discovery] Service register key:[%s], val:[%s] error:%s", s.key, s.val, err.Error())
	}
	return err
}
