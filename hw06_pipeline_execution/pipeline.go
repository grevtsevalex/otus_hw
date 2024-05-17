package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in
	res := make([]interface{}, 0)

	for _, stage := range stages {
		out = stage(out)
	}

OUTER:
	for {
		select {
		case <-done:
			empty := make(Bi)
			close(empty)
			return empty
		case r, ok := <-out:
			if !ok {
				break OUTER
			}
			res = append(res, r)
		}
	}

	outCh := make(Bi, len(res))
	defer close(outCh)

	for _, v := range res {
		outCh <- v
	}

	return outCh
}
