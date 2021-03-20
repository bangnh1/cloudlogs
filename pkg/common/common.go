package common

import (
	// "fmt"
	// "log"
	// "os"
	"regexp"
	"strconv"
	"time"

	"github.com/araddon/dateparse"
	"github.com/aws/aws-sdk-go/aws"
)

type TimeDuration struct {
	Start int64
	End   int64
}

func TimeStringToInt64(inputTime *string) (int64, error) {

	convertedTime, err := dateparse.ParseAny(*inputTime)

	if err != nil {
		panic(err)
	}
	unixTime := aws.TimeUnixMilli(convertedTime)

	return unixTime, err
}

func TimeInt64ToString(inputTime int64) string {
	t := time.Unix(inputTime/1000, 0)
	strDate := t.Format(time.RFC3339Nano)

	return strDate
}

func ParseDateTime(start *string, end *string) (int64, int64) {

	var startTime, endTime int64

	ago_regexp := `(\d+)\s?(m|minute|minutes|h|hour|hours|d|day|days|w|weeks|weeks)(?: ago)?`
	match, _ := regexp.MatchString(ago_regexp, *start)

	if match {
		re := regexp.MustCompile(ago_regexp)
		startRegexp := re.FindStringSubmatch(*start)
		endTime := time.Now().Unix() * 1000
		startTime = TimeShift(false, endTime, startRegexp[1], startRegexp[2])
		return startTime, endTime
	} else {
		if *end == "" {
			endTime := time.Now().Unix() * 1000
			startTime, _ = TimeStringToInt64(start)
			return startTime, endTime
		} else {
			endTime, _ = TimeStringToInt64(end)
			startTime, _ = TimeStringToInt64(start)
			return startTime, endTime
		}
	}
}

func TimeShift(moreOrLess bool, start int64, tshift string, unit string) int64 {

	unitList := map[string]int64{"m": 60, "h": 3600, "d": 86400, "w": 604800}
	var nextTime int64

	for key, value := range unitList {
		if key == unit {
			timeShift, _ := strconv.ParseInt(tshift, 10, 64)
			if moreOrLess {
				nextTime = start + value*timeShift*1000
				return nextTime
			} else {
				nextTime = start - value*timeShift*1000
				return nextTime
			}
		}
	}
	return nextTime
}

func TimeByInterval(start *string, end *string, interval *string) []TimeDuration {

	var timeList []TimeDuration
	intervalRegexp := `(\d+)\s?(m|minute|minutes|h|hour|hours|d|day|days|w|weeks|weeks)?`
	re := regexp.MustCompile(intervalRegexp)
	duration := re.FindStringSubmatch(*interval)
	startTime, endTime := ParseDateTime(start, end)

	for {
		nextTime := TimeShift(true, startTime, duration[1], duration[2])

		if nextTime < startTime {
			break
		}

		if nextTime >= endTime || nextTime == startTime {
			nextDuration := TimeDuration{
				Start: startTime,
				End:   endTime,
			}

			timeList = append(timeList, nextDuration)
			break
		}

		nextDuration := TimeDuration{
			Start: startTime,
			End:   nextTime,
		}

		timeList = append(timeList, nextDuration)
		startTime = nextTime
	}

	return timeList
}
