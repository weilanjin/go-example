package concurrency_pattern

import "time"

type Item struct {
	Title, Channel, GUID string
}

type Fetcher interface {
	Fetch() (item []Item, next time.Time, err error)
}

type Subscription interface {
	Updates() <-chan Item
	Close() error
}

func Subscribe(fetcher Fetcher) Subscription {
	s := &sub{
		fetcher: fetcher,
		updates: make(chan Item),
		closing: make(chan chan error),
	}
	go s.loop()
	return s
}

type sub struct {
	fetcher Fetcher
	updates chan Item
	closing chan chan error
}

func (s *sub) Updates() <-chan Item {
	return s.updates
}

func (s *sub) Close() error {
	errc := make(chan error)
	s.closing <- errc
	return <-errc
}

func (s *sub) loop() {
	var err error
	for {
		select {
		case errc := <-s.closing:
			errc <- err
			close(s.updates)
			return
		}
	}
}

func (s *sub) loopFetchOnly() {
	var (
		pending []Item
		next    time.Time
		err     error
	)
	for {
		var fetchDelay time.Duration
		if now := time.Now(); next.After(now) {
			fetchDelay = next.Sub(now)
		}
		startFetch := time.After(fetchDelay)

		select {
		case <-startFetch:
			var fetched []Item
			fetched, next, err = s.fetcher.Fetch()
			if err != nil {
				next = time.Now().Add(10 * time.Second)
				break
			}
			pending = append(pending, fetched...)
		}
	}
}
