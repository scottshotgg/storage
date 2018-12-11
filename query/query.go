package ouery

import (
	"context"
	"errors"
	"time"

	"github.com/pizzahutdigital/storage/storage"
)

type Operation struct {
	ctx        context.Context
	cancelFunc context.CancelFunc

	oueryFuncs []func()

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

type OpFunc func()

var (
	ErrEmptyCancelFunc    = errors.New("Attempted to set cancel func to nil value")
	ErrEmptyOutput        = errors.New("Attempted to set output to nil value")
	ErrEmptyContext       = errors.New("Attempted to use nil context")
	ErrEmptyOperationFunc = errors.New("Attempted to set ouery func to nil value")

	ErrNilResults = errors.New("Nil results")
)

func New(of OpFunc) *Operation {
	ctx, cancelFunc := context.WithCancel(context.Background())

	return &Operation{
		ctx:        ctx,
		cancelFunc: cancelFunc,
		oueryFuncs: []func(){of},
		results:    &[]storage.Item{},
		// future:     make(chan struct{}),
	}
}

func (o *Operation) Start() error {
	if o.err != nil {
		return o.err
	}

	// start an err chan?
	// start the streaming chans?

	go func() {
		for i := range o.oueryFuncs {
			o.oueryFuncs[i]()

			// If the last ouery triggered an error then stop the process
			if o.err != nil {
				o.Cancel()
				o.stage = i
				break
			}
		}
	}()

	// There was no error in starting the ouery itself
	return nil
}

// TODO: this is not fleshed out yet
func (o *Operation) Run() error {
	if o.err != nil {
		return o.err
	}

	var err = o.Start()
	if err != nil {
		return err
	}

	var (
		// TODO: might need to close these
		futureChan    = o.Future()
		oueryDoneChan = o.Done()
		doneChan      = make(chan struct{})
	)

	go func() {
		defer close(doneChan)

		for {
			select {
			case <-futureChan:
				return

			case <-oueryDoneChan:
				return
			}
		}
	}()

	<-doneChan

	if o.err != nil {
		return o.err
	}

	if o.results != nil {
		return ErrNilResults
	}

	return nil
}

func (o *Operation) AddOperation(of func()) *Operation {
	if of == nil {
		o.err = ErrEmptyOperationFunc
		return o
	}

	o.oueryFuncs = append(o.oueryFuncs, of)
	return o
}

func (o *Operation) WithOperation(of func()) *Operation {
	if of == nil {
		o.err = ErrEmptyOperationFunc
		return o
	}

	o.oueryFuncs = []func(){of}
	return o
}

func (o *Operation) WithContext(ctx context.Context) *Operation {
	if ctx == nil {
		o.err = ErrEmptyContext
		return o
	}

	o.ctx = ctx
	return o
}

func (o *Operation) WithCancelFunc(cf context.CancelFunc) *Operation {
	if cf == nil {
		o.err = ErrEmptyCancelFunc
		return o
	}

	o.cancelFunc = cf
	return o
}

func (o *Operation) WithTimeout(timeout time.Duration) *Operation {
	o.timeout = timeout
	return o
}

func (o *Operation) Async(async bool) *Operation {
	o.async = async
	return o
}

func (o *Operation) WithOutput(output *[]storage.Item) *Operation {
	if output == nil {
		o.err = ErrEmptyOutput
		return o
	}

	o.results = output
	return o
}

// func (o *Operation) IsStream(stream bool) {
// 	o.stream = stream
// }

// func (o *Operation) Wait() {}
// func (o *Operation) WaitForResults() {}

func (o *Operation) Done() <-chan struct{} {
	return o.ctx.Done()
}

func (o *Operation) Results() []storage.Item {
	return *o.results
}

func (o *Operation) Err() error {
	return o.err
}

func (o *Operation) Future() <-chan struct{} {
	return o.future
}

func (o *Operation) Cancel() {
	if o.cancelFunc != nil {
		o.cancelFunc()
	}
}

func (o *Operation) IsAsync() {
	return o.async
}
