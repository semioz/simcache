# simcache

Semantic cache for your LLM apps in Go!

# Usage of simcache via LangchainGo and Upstash

```
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/semioz/simcache"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/upstash/vector-go"
)

func main() {
	index := vector.NewIndex("UPSTASH_URL", "UPSTASH_TOKEN")
	simCache := simcache.NewSimCache(simcache.UpstashConfig{
		Index:        index,
		MinProximity: 0.9,
	})

	ctx := context.Background()
	llm, err := openai.New()
	if err != nil {
		log.Fatal(err)
	}
	prompt := "What are the some vector databases that I can use?"
	response, err := llms.GenerateFromSinglePrompt(ctx, llm, prompt)

	simCache.Set(prompt, response)

	result, err := simCache.Get("List some vector databases I can use when building a project")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
}
```
