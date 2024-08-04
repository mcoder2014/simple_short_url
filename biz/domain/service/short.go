package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/mcoder2014/simple_short_url/biz/domain/model"
)

type ShortService struct {
	ShortCache      map[string]*model.ShortURLConfig
	ShortCacheMutex sync.RWMutex
}

var shortConfigMutex sync.Mutex

const (
	shortConfigPath = "./conf/short.json"
)

func NewShortService() (*ShortService, error) {
	s := &ShortService{}
	cfg, err := s.loadConfig()
	if err != nil {
		return nil, err
	}
	s.ShortCache = make(map[string]*model.ShortURLConfig, len(cfg))
	for _, c := range cfg {
		s.ShortCache[c.Short] = c
	}
	return s, nil
}

func (s *ShortService) GetShortUrl(ctx context.Context, short string) (*string, error) {
	s.ShortCacheMutex.RLock()
	defer s.ShortCacheMutex.RUnlock()
	config, ok := s.ShortCache[short]
	if !ok {
		return nil, nil
	}
	if !config.Enable {
		hlog.CtxInfof(ctx, "short=%v is disable", short)
		return nil, nil
	}

	return &config.Long, nil
}

func (s *ShortService) AddConfig(ctx context.Context, short, long string, desp string, creator string) (*model.ShortURLConfig, error) {
	s.ShortCacheMutex.Lock()
	defer s.ShortCacheMutex.Unlock()
	if _, ok := s.ShortCache[short]; ok {
		return nil, fmt.Errorf("short %v already exists", short)
	}

	cfg := &model.ShortURLConfig{
		Short:  short,
		Long:   long,
		Enable: true,
		Desp:   desp,
	}

	s.ShortCache[short] = cfg
	configs, err := s.loadConfig()
	if err != nil {
		return nil, fmt.Errorf("check config failed: %w", err)
	}
	configs = append(configs, cfg)
	return cfg, s.saveConfig(configs)
}

func (s *ShortService) Refresh(ctx context.Context) error {
	configs, err := s.loadConfig()
	if err != nil {
		return err
	}
	s.ShortCacheMutex.Lock()
	defer s.ShortCacheMutex.Unlock()
	s.ShortCache = make(map[string]*model.ShortURLConfig, len(configs))
	for _, cfg := range configs {
		s.ShortCache[cfg.Short] = cfg
	}
	return nil
}

func (s *ShortService) loadConfig() ([]*model.ShortURLConfig, error) {
	var config []*model.ShortURLConfig

	shortConfigMutex.Lock()
	defer shortConfigMutex.Unlock()

	content, err := os.ReadFile(shortConfigPath)
	if err != nil {
		return nil, fmt.Errorf("read short config file: %w", err)
	}
	err = json.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (s *ShortService) saveConfig(config []*model.ShortURLConfig) error {
	shortConfigMutex.Lock()
	defer shortConfigMutex.Unlock()

	content, err := json.Marshal(config)
	if err != nil {
		return err
	}
	return os.WriteFile(shortConfigPath, content, 0644)
}
