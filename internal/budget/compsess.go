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
	if haveLink {
		cancel()
	}
	if ctx.Err() != nil && leftoverResult != nil {
		res.FSPLDB = leftoverResult.FSPLDB
		res.PrDBm = leftoverResult.PrDBm
		res.PrWatts = leftoverResult.PrWatts
		if res.Assessment != nil && leftoverResult.Assessment != nil {
			res.Assessment.SNRDB = leftoverResult.Assessment.SNRDB
			res.Assessment.MarginDB = leftoverResult.Assessment.MarginDB
			res.Assessment.Feasible = leftoverResult.Assessment.Feasible
		}
		return res
	}
	leftoverResult = res
	haveLink = true
	return res
}
