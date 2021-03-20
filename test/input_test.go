package unittest

import (
	"github.com/bangnh1/cloudlogs/pkg/awslogs"
	"github.com/bangnh1/cloudlogs/pkg/common"
	"testing"
)

func TestTimeStringToInt64(t *testing.T) {

	var input = []string{
		"2021-03-01T00:00:00.000Z",
		"03/01/2021 12:00:00 AM",
		"03/01/2021",
		"Mon Mar 01 00:00:00 UTC 2021",
	}

	for _, i := range input {
		time, _ := common.TimeStringToInt64(&i)
		if time != 1614556800000 {
			t.Errorf("Output expect 1614556800000 instead of %v", time)
		}
	}
}

func TestTimeShift(t *testing.T) {
	before := common.TimeShift(false, 1614621600000, "1", "h")
	if before != 1614618000000 {
		t.Errorf("Output expect 1614618000000 instead of %v", before)
	}
	after := common.TimeShift(true, 1616594400000, "1", "d")
	if after != 1616680800000 {
		t.Errorf("Output expect 1616680800000 instead of %v", after)
	}
}

func TestParseDateTime(t *testing.T) {
	start := "2021-03-01T00:00:00.000Z"
	end := "2021-03-02T00:00:00.000Z"
	startTime, endTime := common.ParseDateTime(&start, &end)

	if startTime != 1614556800000 {
		t.Errorf("Output expect 1614556800000 instead of %v", startTime)
	}

	if endTime != 1614643200000 {
		t.Errorf("Output expect 1614643200000 instead of %v", endTime)
	}

}

func TestTimeByInterval(t *testing.T) {
	start := "2021-03-01T00:00:00.000Z"
	end := "2021-03-02T00:00:00.000Z"
	interval := "12h"
	output := common.TimeByInterval(&start, &end, &interval)

	for _, i := range output {
		if i.Start != 1614556800000 && i.Start != 1614600000000 {
			t.Errorf("Output expect 1614556800000 || 1614600000000 instead of %v", i.Start)
		}
		if i.End != 1614600000000 && i.End != 1614643200000 {
			t.Errorf("Output expect 1614600000000 || 1614643200000 instead of %v", i.End)
		}
	}
}

func TestLimitHandler(t *testing.T) {

	input := [5]string{"-10", "0", "20", "50", "70"}
	limitOutput := [5]int64{0, 0, 20, 50, 50}
	maxOutput := [5]int64{0, 0, 20, 50, 70}

	for i := 1; i < 5; i++ {
		limit, max := awslogs.LimitHandler(&input[i])
		if limit != limitOutput[i] {
			t.Errorf("With Input %v -> Output expect %v instead of %v", limit, limitOutput[i], limit)
		}
		if max != maxOutput[i] {
			t.Errorf("With Input %v -> Output expect %v instead of %v", max, maxOutput[i], max)
		}
	}
}
