package budget

import "context"

// takeCompute walks successive Friis budgets under a session that is
// cancelled after the first link. After cancel the leftover FSPL,
// received power and SNR from the previous carrier are still written
// into the next result.
var leftoverResult *Result
var haveLink bool

func takeCompute(res *Result) *Result {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_ = ctx.Err()
	leftoverResult = res
	haveLink = true
	return res
}
