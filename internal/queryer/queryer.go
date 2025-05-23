package queryer

import (
	"context"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/JosephNaberhaus/hacker-news-font-page/internal/hackernews"
)

type Queryer struct {
	WaitMillisecondsMin, WaitMillisecondsMax int64

	lastQuery time.Time
}

func (q *Queryer) GetTitles(ctx context.Context, year, month, day int) ([30]string, error) {
	// See if we should wait a bit before our next call.
	durSinceLastCall := time.Since(q.lastQuery)
	wait := q.generateWait()
	if durSinceLastCall.Milliseconds() < wait {
		remainingWait := time.Duration(wait-durSinceLastCall.Milliseconds()) * time.Millisecond
		fmt.Printf("Waiting %d...\n", remainingWait.Milliseconds())
		time.Sleep(remainingWait)

		select {
		case <-ctx.Done():
			return [30]string{}, ctx.Err()
		case <-time.NewTimer(remainingWait).C:
			// Do the request.
		}
	}
	q.lastQuery = time.Now()

	fmt.Printf("Querying for %d %d %d\n", year, month, day)

	page, err := hackernews.LoadPage(ctx, year, month, day)
	if err != nil {
		return [30]string{}, err
	}

	return page.ParseTitles()
}

func (q *Queryer) generateWait() int64 {
	return rand.N[int64](q.WaitMillisecondsMax-q.WaitMillisecondsMin) + q.WaitMillisecondsMin
}
