package simcache

import (
	"errors"

	"github.com/upstash/vector-go"
)

type CacheConfig struct {
	MinProximity float32 `json:"minProximity"`
	Index        *vector.Index
}

type SimCache struct {
	minProximity float32
	index        *vector.Index
}

func NewSimCache(config CacheConfig) *SimCache {
	if config.MinProximity == 0 {
		config.MinProximity = 0.8
	}
	return &SimCache{
		minProximity: config.MinProximity,
		index:        config.Index,
	}
}

func (cache *SimCache) Get(keyOrKeys interface{}) (interface{}, error) {
	switch key := keyOrKeys.(type) {
	case string:
		return cache.queryKey(key)
	case []string:
		res := make([]interface{}, len(key))
		for i, k := range key {
			value, err := cache.queryKey(k)
			if err != nil {
				return "", err
			}
			res[i] = value
		}
		return res, nil
	}
	return "", errors.New("invalid types or lengths")
}

func (cache *SimCache) queryKey(key string) (interface{}, error) {
	res, err := cache.index.QueryData(vector.QueryData{
		Data:            key,
		TopK:            2,
		IncludeVectors:  true,
		IncludeMetadata: true,
	})
	if err != nil {
		return "", err
	}
	if len(res) > 0 && res[0].Score > cache.minProximity {
		return res[0].Metadata["value"], nil
	}
	return "", nil
}

func (cache *SimCache) Set(keyOrKeys interface{}, valueOrValues interface{}) error {
	switch key := keyOrKeys.(type) {
	case string:
		if value, ok := valueOrValues.(string); ok {
			err := cache.index.UpsertData(vector.UpsertData{
				Id:   key,
				Data: key,
				Metadata: map[string]interface{}{
					"value": value,
				},
			})
			if err != nil {
				return err
			}
			return nil
		}
	case []string:
		if values, ok := valueOrValues.([]string); ok {
			for i, key := range key {
				err := cache.index.UpsertData(vector.UpsertData{
					Id:   key,
					Data: key,
					Metadata: map[string]interface{}{
						"value": values[i],
					},
				})
				if err != nil {
					return err
				}
			}
			return nil
		}
	}
	return errors.New("invalid types or lengths")
}

func (cache *SimCache) Delete(key string) error {
	_, err := cache.index.Delete(key)
	if err != nil {
		return err
	}
	return nil
}

func (cache *SimCache) BulkDelete(keys []string) error {
	_, err := cache.index.DeleteMany(keys)
	if err != nil {
		return err
	}
	return nil
}

func (cache *SimCache) Flush() error {
	err := cache.index.Reset()
	if err != nil {
		return err
	}
	return nil
}
