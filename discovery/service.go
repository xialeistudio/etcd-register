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
func (s *Service) Start(ctxs ...context.Context) error {
	err := s.register(ctxs...)
	if err != nil {
		return err
	}

	log.Printf("[discovery] Service Start successfully key:[%s], val:[%s]", s.key, s.val)

	go func() {
		s.ticker = time.NewTicker(s.options.Interval)
		for range s.ticker.C {
			s.register(ctxs...)
		}
	}()
	return nil
}

func (s *Service) Shutdown(ctxs ...context.Context) error {
	return s.unregister(ctxs...)
}

// 注册服务
func (s *Service) register(ctxs ...context.Context) error {
	var ctx context.Context
	if len(ctxs) == 0 {
		ctx = context.TODO()
	} else {
		ctx = ctxs[0]
	}
	_, err := s.client.Set(ctx, s.key, s.val, &client.SetOptions{
		TTL: s.options.Lifetime,
	})
	if err != nil {
		log.Printf("[discovery] Service register key:[%s], val:[%s] error:%s", s.key, s.val, err.Error())
	}
	return err
}

func (s *Service) unregister(ctxs ...context.Context) error {
	var ctx context.Context
	if len(ctxs) == 0 {
		ctx = context.TODO()
	} else {
		ctx = ctxs[0]
	}
	_, err := s.client.Delete(ctx, s.key, &client.DeleteOptions{})
	if err != nil {
		log.Printf("[discovery] Service unregister key:[%s], val:[%s] error:%s", s.key, s.val, err.Error())
	}
	return err
}
