package dal

import (
	"golang.org/x/net/context"
	"time"
)

func defaultBatchStats(ctx context.Context, b *BatchStats) {
}

func (b *BatchStats) valid() error {
	return nil
}

func ListBatchStatsByBatchIdRange(ctx context.Context, batchId int) ([]*BatchStats, []*BatchStats) {
	batch := FindBatch(ctx, batchId)
	if batch == nil {
		return nil, nil
	}
	dayStats := []*BatchStats{}
	days := listDays(batch.Sid, 8)
	for i, _ := range days {
		if i == len(days)-1 {
			break
		}
		start := days[i]
		end := days[i+1]
		batchStats := statsBatchByStartAndEnd(ctx, batch, start, end)
		dayStats = append(dayStats, batchStats)
	}
	weekStats := []*BatchStats{}
	weeks := listWeeks(batch.Sid, 8)
	for i, _ := range weeks {
		if i == len(weeks)-1 {
			break
		}
		start := weeks[i]
		end := weeks[i+1]
		batchStats := statsBatchByStartAndEnd(ctx, batch, start, end)
		weekStats = append(weekStats, batchStats)
	}
	return dayStats, weekStats
}

func statsBatchByStartAndEnd(ctx context.Context, batch *Batch, start time.Time, end time.Time) *BatchStats {
	batchStats := FindBatchStatsByBatchIdStartEnd(ctx, batch.Id, start, end)
	if batchStats != nil {
		return batchStats
	}
	batchStats = NewBatchStats(ctx)
	batchStats.BatchId = batch.Id
	batchStats.Start = start
	batchStats.End = end
	// stats
	batchStats.Success = CountTestdataByBatchSuccessStartEnd(ctx, batch.Sid, start, end)
	batchStats.RightFirstTime = CountTestdataByBatchRightFirstTimeStartEnd(ctx, batch.Sid, start, end, batch.Sid, start, end)
	batchStats.Failed = CountTestdataByBatchFailedStartEnd(ctx, batch.Sid, start, end)
	batchStats.Rejected = CountTestdataByBatchRejectedStartEnd(ctx, batch.Sid, start, end)
	batchStats.Cnt = batchStats.Success + batchStats.Rejected
	if batchStats.Cnt > 0 {
		// pct %
		batchStats.SuccessPct = batchStats.Success * 10000 / batchStats.Cnt
		batchStats.RightFirstTimePct = batchStats.RightFirstTime * 10000 / batchStats.Cnt
		batchStats.FailedPct = batchStats.Failed * 10000 / batchStats.Cnt
		batchStats.RejectedPct = batchStats.Rejected * 10000 / batchStats.Cnt
	}
	batchStats.Save()
	return batchStats
}

func getMinNight(batchSid string) time.Time {
	testdatas := ListTestdataByBatch(nil, batchSid, 0, 1)
	t := time.Now()
	if len(testdatas) > 0 {
		t = testdatas[0].Created
	}
	_, offset := t.Zone()
	minNight := t.Truncate(24 * time.Hour).Add(-time.Duration(offset) * time.Second)
	return minNight
}

func listDays(batchSid string, num int) []time.Time {
	days := make([]time.Time, num)
	minNight := getMinNight(batchSid)
	for i, _ := range days {
		days[num-i-1] = minNight
		minNight = minNight.Add(-24 * time.Hour)
	}
	return days
}

func listWeeks(batchSid string, num int) []time.Time {
	minNight := getMinNight(batchSid)
	weekday := int(minNight.Weekday()) - 1
	weekMinNight := minNight.Add(-time.Duration(weekday) * 24 * time.Hour)
	weeks := make([]time.Time, num)
	for i, _ := range weeks {
		weeks[num-i-1] = weekMinNight
		weekMinNight = weekMinNight.Add(-7 * 24 * time.Hour)
	}
	return weeks
}
