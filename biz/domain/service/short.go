package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/mcoder2014/simple_short_url/biz/config"
	"github.com/mcoder2014/simple_short_url/biz/domain/model"
	"github.com/samber/lo"
)

type ShortService struct {
	ShortCache      map[string]*model.ShortURLConfig
	ShortCacheMutex sync.RWMutex
}

var shortConfigMutex sync.Mutex

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
		Short:      short,
		Long:       long,
		Enable:     true,
		Desp:       desp,
		CreateTime: time.Now().Unix(),
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
	var configs []*model.ShortURLConfig

	shortConfigMutex.Lock()
	defer shortConfigMutex.Unlock()

	content, err := os.ReadFile(config.GetConfig().ShortURLFile)
	if err != nil {
		return nil, fmt.Errorf("read short configs file: %w", err)
	}
	err = json.Unmarshal(content, &configs)
	if err != nil {
		return nil, err
	}
	return configs, nil
}

func (s *ShortService) saveConfig(configs []*model.ShortURLConfig) error {
	shortConfigMutex.Lock()
	defer shortConfigMutex.Unlock()

	content, err := json.Marshal(configs)
	if err != nil {
		return err
	}
	return os.WriteFile(config.GetConfig().ShortURLFile, content, 0644)
}

func (s *ShortService) ListConfig(ctx context.Context, offset, limit int) ([]*model.ShortURLConfig, bool, error) {
	if offset > len(s.ShortCache) {
		return nil, false, nil
	}

	configs, err := s.loadConfig()
	if err != nil {
		return nil, false, err
	}

	end := lo.Min([]int{offset + limit, len(configs)})
	return configs[offset:end], end < len(configs), nil
}

func (s *ShortService) DeleteConfig(ctx context.Context, short string) (err error) {
	s.ShortCacheMutex.Lock()
	defer s.ShortCacheMutex.Unlock()

	if _, ok := s.ShortCache[short]; !ok {
		hlog.CtxInfof(ctx, "short=%v is not found", short)
		return nil
	}

	defer func() {
		if err == nil {
			delete(s.ShortCache, short)
		}
	}()

	configs, err := s.loadConfig()
	if err != nil {
		return fmt.Errorf("check config failed: %w", err)
	}
	var res []*model.ShortURLConfig
	for _, cfg := range configs {
		if cfg.Short != short {
			res = append(res, cfg)
			continue
		}
		hlog.CtxInfof(ctx, "delete config=%v, long=%v", short, cfg.Long)
	}
	return s.saveConfig(res)
}
