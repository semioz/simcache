package simcache

import (
	"testing"

	"github.com/upstash/vector-go"
)

func TestNewSimCache(t *testing.T) {
	config := UpstashOptions{
		MinProximity: 0.5,
		Index:        vector.NewIndex("URL", "TOKEN"),
	}
	cache := NewSimCache(config)
	if cache.minProximity != 0.5 {
		t.Errorf("Expected minProximity to be 0.5, but got %f", cache.minProximity)
	}
	if cache.index == nil {
		t.Errorf("Expected index to be not nil, but got nil")
	}
}

func TestSet(t *testing.T) {
	config := UpstashOptions{
		MinProximity: 0.5,
		Index:        vector.NewIndex("URL", "TOKEN"),
	}
	cache := NewSimCache(config)
	err := cache.Set("Current CEO of Apple", "Tim Cook")
	if err != nil {
		t.Errorf("Expected Set('key1', 'value1') to succeed, but got error %v", err)
	}
	err = cache.Set("Current CEO of Microsoft", "Satya Nadella")
	if err != nil {
		t.Errorf("Expected Set('key2', 'value2') to succeed, but got error %v", err)
	}
	err = cache.Set("Current CEO of Google", "Sundar Pichai")
	if err != nil {
		t.Errorf("Expected Set('key3', 'value3') to succeed, but got error %v", err)
	}
}

func TestGet(t *testing.T) {
	config := UpstashOptions{
		MinProximity: 0.5,
		Index:        vector.NewIndex("URL", "TOKEN"),
	}
	cache := NewSimCache(config)
	err := cache.index.UpsertData(vector.UpsertData{
		Id:   "tr",
		Data: "Republic of Turkey",
		Metadata: map[string]interface{}{
			"value": "Mustafa Kemal Atatürk",
		},
	})
	if err != nil {
		return
	}
	err = cache.index.UpsertData(vector.UpsertData{
		Id:   "us",
		Data: "United States of America",
		Metadata: map[string]interface{}{
			"value": "George Washington",
		},
	})
	if err != nil {
		return
	}
	err = cache.index.UpsertData(vector.UpsertData{
		Id:   "fr",
		Data: "French Republic",
		Metadata: map[string]interface{}{
			"value": "Napoleon Bonaparte",
		},
	})
	if err != nil {
		return
	}
	if val, err := cache.Get("tr"); err != nil || val != "Mustafa Kemal Atatürk" {
		t.Errorf("Expected Get('tr') to return Ataturk, but got error %v", err)
	}
	if val, err := cache.Get("us"); err != nil || val != "George Washington" {
		t.Errorf("Expected Get('us') to return Washington, but got error %v", err)
	}
	if val, err := cache.Get("fr"); err != nil || val != "Napoleon Bonaparte" {
		t.Errorf("Expected Get('fr') to return Napoleon, but got error %v", err)
	}
	// Test multiple keys
	if val, err := cache.Get([]string{"tr", "us", "fr"}); err != nil || val.([]interface{})[0] != "Mustafa Kemal Atatürk" || val.([]interface{})[1] != "George Washington" || val.([]interface{})[2] != "Napoleon Bonaparte" {
		t.Errorf("Expected Get(['tr', 'us', 'fr']) to return Ataturk, Washington, Napoleon, but got error %v", err)
	}
}

func TestDelete(t *testing.T) {
	config := UpstashOptions{
		MinProximity: 0.5,
		Index:        vector.NewIndex("URL", "TOKEN"),
	}
	cache := NewSimCache(config)
	err := cache.index.UpsertData(vector.UpsertData{
		Id:   "key1",
		Data: "data1",
		Metadata: map[string]interface{}{
			"value": "value1",
		},
	})
	if err != nil {
		return
	}
	err = cache.Delete("key1")
	if err != nil {
		t.Errorf("Expected Delete('key1') to succeed, but got error %v", err)
	}
	err = cache.Delete("nonexistent")
	if err == nil {
		t.Errorf("Expected Delete('nonexistent') to return an error")
	}
}

func TestBulkDelete(t *testing.T) {
	config := UpstashOptions{
		MinProximity: 0.5,
		Index:        vector.NewIndex("URL", "TOKEN"),
	}
	cache := NewSimCache(config)
	err := cache.index.UpsertData(vector.UpsertData{
		Id:   "key1",
		Data: "data1",
		Metadata: map[string]interface{}{
			"value": "value1",
		},
	})
	if err != nil {
		return
	}
	err = cache.index.UpsertData(vector.UpsertData{
		Id:   "key2",
		Data: "data2",
		Metadata: map[string]interface{}{
			"value": "value2",
		},
	})
	if err != nil {
		return
	}
	err = cache.BulkDelete([]string{"key1", "key2"})
	if err != nil {
		t.Errorf("Expected BulkDelete(['key1', 'key2']) to succeed, but got error %v", err)
	}
	err = cache.BulkDelete([]string{"nonexistent", "nonexistent2"})
	if err == nil {
		t.Errorf("Expected BulkDelete(['nonexistent', 'nonexistent2']) to return an error")
	}
}

func TestFlush(t *testing.T) {
	config := UpstashOptions{
		MinProximity: 0.5,
		Index:        vector.NewIndex("URL", "TOKEN"),
	}
	cache := NewSimCache(config)
	err := cache.index.UpsertData(vector.UpsertData{
		Id:   "key1",
		Data: "data1",
		Metadata: map[string]interface{}{
			"value": "value1",
		},
	})
	if err != nil {
		return
	}
	err = cache.Flush()
	if err != nil {
		t.Errorf("Expected Flush() to succeed, but got error %v", err)
	}
}
