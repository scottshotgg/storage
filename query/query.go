package query

import (
	"context"
	"errors"
	"time"

	"github.com/pizzahutdigital/storage/storage"
)

type Query struct {
	ctx        context.Context
	cancelFunc context.CancelFunc

	queryFuncs []func()

	results *[]storage.Item
	err     error
	future  <-chan struct{}
	timeout time.Duration
	async   bool

	stage int

	// TODO: do something with this; need to check these
	// started bool
	// finished bool

	// TODO: should prob put a mutex here that is checked
}

var (
	ErrEmptyCancelFunc = errors.New("Attempted to set cancel func to nil value")
	ErrEmptyOutput     = errors.New("Attempted to set output to nil value")
	ErrEmptyContext    = errors.New("Attempted to use nil context")
	ErrEmptyQueryFunc  = errors.New("Attempted to set query func to nil value")

	ErrNilResults = errors.New("Nil results")
)

func New(qf func()) *Query {
	ctx, cancelFunc := context.WithCancel(context.Background())

	return &Query{
		ctx:        ctx,
		cancelFunc: cancelFunc,
		queryFuncs: []func(){qf},
		results:    &[]storage.Item{},
		future:     make(chan struct{}),
	}
}

func (q *Query) Start() error {
	if q.err != nil {
		return q.err
	}

	// start an err chan?
	// start the streaming chans?

	go func() {
		for i, qf := range q.queryFuncs {
			qf()

			// If the last query triggered an error then stop the process
			if q.err != nil {
				q.Cancel()
				q.stage = i
				break
			}
		}
	}()

	// There was no error in starting the query itself
	return nil
}

func (q *Query) Run() error {
	err := q.Start()
	if err != nil {
		return err
	}

	futureChan := q.Future()
	queryDoneChan := q.Done()

	doneChan := make(chan struct{})

	go func() {
		defer close(doneChan)

		for {
			select {
			case <-futureChan:
				return

			case <-queryDoneChan:
				return
			}
		}
	}()

	<-doneChan

	if q.err != nil {
		return q.err
	}

	if q.results != nil {
		return ErrNilResults
	}

	return nil
}

func (q *Query) AddQueryFunc(qf func()) *Query {
	if qf == nil {
		q.err = ErrEmptyQueryFunc
		return q
	}

	q.queryFuncs = append(q.queryFuncs, qf)
	return q
}

func (q *Query) WithQueryFunc(qf func()) *Query {
	if qf == nil {
		q.err = ErrEmptyQueryFunc
		return q
	}

	q.queryFuncs = []func(){qf}
	return q
}

func (q *Query) WithContext(ctx context.Context) *Query {
	if ctx == nil {
		q.err = ErrEmptyContext
		return q
	}

	q.ctx = ctx
	return q
}

func (q *Query) WithCancelFunc(cf context.CancelFunc) *Query {
	if cf == nil {
		q.err = ErrEmptyCancelFunc
		return q
	}

	q.cancelFunc = cf
	return q
}

func (q *Query) WithTimeout(timeout time.Duration) *Query {
	q.timeout = timeout
	return q
}

func (q *Query) Async(async bool) *Query {
	q.async = async
	return q
}

func (q *Query) WithOutput(output *[]storage.Item) *Query {
	if output == nil {
		q.err = ErrEmptyOutput
		return q
	}

	q.results = output
	return q
}

// func (q *Query) IsStream(stream bool) {
// 	q.stream = stream
// }

// func (q *Query) Wait() {}
// func (q *Query) WaitForResults() {}

func (q *Query) Done() <-chan struct{} {
	return q.ctx.Done()
}

func (q *Query) Results() []storage.Item {
	return *q.results
}

func (q *Query) Err() error {
	return q.err
}

func (q *Query) Future() <-chan struct{} {
	return q.future
}

func (q *Query) Cancel() {
	if q.cancelFunc != nil {
		q.cancelFunc()
	}
}

func (q *Query) IsAsync() {
	return q.async
}
