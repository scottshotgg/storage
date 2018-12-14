package query

import (
	"context"
	"errors"
	"time"

	"github.com/scottshotgg/storage/storage"
)

type Operation struct {
	ctx        context.Context
	cancelFunc context.CancelFunc

	queryFuncs []OpFunc

	results *[]storage.Item
	err     error
	future  <-chan struct{}
	timeout time.Duration
	async   bool

	// Status
	stage   int
	running bool

	// TODO: should prob put a mutex here that is checked
}

// TODO: Switch these to be private later
type Value struct {
	Name   string
	TypeOf string
	Value  interface{}
}

type Values map[string]Value

type OpFunc func(o *Operation, v Values) (Values, error)

var (
	ErrEmptyCancelFunc    = errors.New("Attempted to set cancel func to nil value")
	ErrEmptyOutput        = errors.New("Attempted to set output to nil value")
	ErrEmptyContext       = errors.New("Attempted to use nil context")
	ErrEmptyOperationFunc = errors.New("Attempted to set query func to nil value")

	ErrNilResults = errors.New("Nil results")
)

func New() *Operation {
	var ctx, cancelFunc = context.WithCancel(context.Background())

	return &Operation{
		ctx:        ctx,
		cancelFunc: cancelFunc,
		queryFuncs: []OpFunc{},
		results:    &[]storage.Item{},
		// future:     make(chan struct{}),
	}
}

func Raw() *Operation {
	return &Operation{}
}

func (o *Operation) Start() error {
	// If the operations have already produced an error then don't start
	if o.err != nil {
		return o.err
	}

	// If there aren't any query functions to be run then return an error
	if o.queryFuncs == nil {
		return errors.New("No query functions set")
	}

	// start an err chan?
	// start the streaming chans?

	// Set the started flag
	o.running = true

	go func() {
		defer func() {
			o.running = false
			o.Cancel()
		}()

		for i, op := range o.queryFuncs {
			if op == nil {
				// log
				break
			}

			// This is the stage it errored on
			o.stage = i

			// Run the query operation
			_, o.err = op(o, nil)
			if o.err != nil {
				break
			}

			// select {
			// // If the context has been canceled then we are done
			// case <-o.Done():
			// 	return

			// default:
			// }

			// If the last query triggered an error then stop the cascade of operations
			if o.err != nil {
				break
			}
		}
	}()

	// There was no error in starting the query itself
	return nil
}

func (o *Operation) Run() error {
	// If the operations have already produced an error then don't start
	if o.err != nil {
		return o.err
	}

	// Start the operations
	var err = o.Start()
	if err != nil {
		return err
	}

	// Wait for the operations to end
	o.Wait()

	return o.err
}

func (o *Operation) AddOperation(of OpFunc) *Operation {
	if of == nil {
		o.err = ErrEmptyOperationFunc
		return o
	}

	o.queryFuncs = append(o.queryFuncs, of)
	return o
}

func (o *Operation) WithOperations(ofs []OpFunc) *Operation {
	if ofs == nil {
		o.err = ErrEmptyOperationFunc
		return o
	}

	o.queryFuncs = ofs
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

func (o *Operation) IsAsync() bool {
	return o.async
}

func (o *Operation) IsRunning() bool {
	return o.running
}

func (o *Operation) Stage() int {
	return o.stage
}

func (o *Operation) Wait() {
	for {
		select {
		case <-o.ctx.Done():
			return

		default:
			// If it's not running anymore or an error was produced then return out
			switch {
			case !o.running:
				return

			case o.err != nil:
				return
			}
		}
	}
}

func (o *Operation) WaitWithTimeout(timeout time.Duration) error {
	var (
		done  = o.ctx.Done()
		timer = time.After(timeout)
	)

	for {
		select {
		case <-done:
			return o.err

		case <-timer:
			return errors.New("Timeout hit")

		default:
			// If it's not running anymore or an error was produced then return out
			if !o.running || o.err != nil {
				return o.err
			}
		}
	}
}

// func (o *Operation) IsStream(stream bool) {
// 	o.stream = stream
// }

// func (o *Operation) WaitForResults() {}
