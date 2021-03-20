package awslogs

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/bangnh1/cloudlogs/pkg/common"
)

type Destination struct {
	_ struct{} `type:"structure"`

	Group         string
	NextToken     string
	Limit         int64
	MaxItems      int64
	FilterPattern string
	Streams       []*string
	StreamPrefix  string
	WatchInterval int
	Watch         bool
	Descending    bool
	OrderBy       string
	WriteToFile   bool
	Svc           *cloudwatchlogs.CloudWatchLogs
}

func (dst *Destination) NewListGroup() *cloudwatchlogs.DescribeLogGroupsInput {

	if dst.Group == "" {
		return &cloudwatchlogs.DescribeLogGroupsInput{
			Limit: &dst.Limit,
		}
	} else {
		return &cloudwatchlogs.DescribeLogGroupsInput{
			Limit:              &dst.Limit,
			LogGroupNamePrefix: &dst.Group,
		}
	}
}

func (dst *Destination) ListLogGroups() []*string {

	var nextToken string = ""
	var logGroups []*string
	describeLogGroupsInput := dst.NewListGroup()

	for {
		if nextToken != "" {
			describeLogGroupsInput.NextToken = aws.String(nextToken)
		}

		describeLogGroupsOutput, err := dst.Svc.DescribeLogGroups(describeLogGroupsInput)

		for _, value := range describeLogGroupsOutput.LogGroups {
			if int64(len(logGroups)) >= dst.MaxItems {
				break
			}
			logGroups = append(logGroups, value.LogGroupName)
		}

		if err != nil {
			panic(err)
		}

		if describeLogGroupsOutput.NextToken == nil || *describeLogGroupsOutput.NextToken == "" {
			break
		}

		nextToken = *describeLogGroupsOutput.NextToken

	}

	return logGroups
}

func (dst *Destination) NewStreams() *cloudwatchlogs.DescribeLogStreamsInput {
	input := &cloudwatchlogs.DescribeLogStreamsInput{
		Limit:        &dst.Limit,
		LogGroupName: &dst.Group,
		Descending:   &dst.Descending,
		OrderBy:      &dst.OrderBy,
	}

	if dst.StreamPrefix != "" {
		input.LogStreamNamePrefix = &dst.StreamPrefix
	}

	return input
}

func (dst *Destination) ListStreams() []*string {

	var nextToken string = ""
	var logStreams []*string
	describeLogStreamsInput := dst.NewStreams()

	for {
		if nextToken != "" {
			describeLogStreamsInput.NextToken = aws.String(nextToken)
		}

		describeLogStreamsOutput, err := dst.Svc.DescribeLogStreams(describeLogStreamsInput)

		for _, value := range describeLogStreamsOutput.LogStreams {
			if int64(len(logStreams)) >= dst.MaxItems {
				break
			}
			logStreams = append(logStreams, value.LogStreamName)
		}

		if err != nil {
			panic(err)
		}

		if describeLogStreamsOutput.NextToken == nil || *describeLogStreamsOutput.NextToken == "" {
			break
		}

		nextToken = *describeLogStreamsOutput.NextToken

	}

	return logStreams

}

func (dst *Destination) NewFilter(startTime, endTime int64) *cloudwatchlogs.FilterLogEventsInput {
	input := &cloudwatchlogs.FilterLogEventsInput{
		LogGroupName: &dst.Group,
		Limit:        &dst.Limit,
		StartTime:    &startTime,
		EndTime:      &endTime,
	}

	if dst.StreamPrefix != "" {
		input.LogStreamNamePrefix = &dst.StreamPrefix
	} else {
		input.LogStreamNames = dst.Streams
	}

	if dst.FilterPattern != "" {
		input.FilterPattern = &dst.FilterPattern
	}
	return input
}

func (dst *Destination) FilterAwsLogs(startTime, endTime int64) ([]*cloudwatchlogs.FilterLogEventsOutput, error) {

	filterLogEventsInput := dst.NewFilter(startTime, endTime)
	var filterLogsEventsOutput *cloudwatchlogs.FilterLogEventsOutput
	reps := make([]*cloudwatchlogs.FilterLogEventsOutput, 0)
	var err error

	nextToken := ""
	for {

		if nextToken != "" {
			filterLogEventsInput.NextToken = aws.String(nextToken)
		}

		filterLogsEventsOutput, err = dst.Svc.FilterLogEvents(filterLogEventsInput)

		reps = append(reps, filterLogsEventsOutput)
		if err != nil {
			panic(err)
		}

		if filterLogsEventsOutput.NextToken == nil || *filterLogsEventsOutput.NextToken == "" {
			break
		}
		nextToken = *filterLogsEventsOutput.NextToken
	}

	return reps, err
}

func ListLogGroups(dst *Destination) {
	listGroups := dst.ListLogGroups()

	for _, group := range listGroups {
		fmt.Println(*group)
	}
}

func ListStreams(dst *Destination) {
	listStreams := dst.ListStreams()

	for _, stream := range listStreams {
		fmt.Println(*stream)
	}
}

func FilterAwsLogs(dst *Destination, goroutines int, timePeriod []common.TimeDuration) {
	var wg sync.WaitGroup
	maxGoroutines := goroutines
	guard := make(chan struct{}, maxGoroutines)

	for _, value := range timePeriod {

		strDate := common.TimeInt64ToString(value.Start)
		filename := strDate + ".txt"
		start := value.Start
		end := value.End
		wg.Add(1)
		guard <- struct{}{}
		go func() {
			filterLogsEventsOutput, err := dst.FilterAwsLogs(start, end)
			if err != nil {
				panic(err)
			}
			if dst.WriteToFile {
				WriteToFile(filename, filterLogsEventsOutput)
			} else {
				WriteToCli(filterLogsEventsOutput)
			}
			<-guard
			defer wg.Done()

		}()
	}

	wg.Wait()
}

func WatchAwsLogs(dst *Destination, timePeriod []common.TimeDuration) {

	startTime := timePeriod[0].Start
	endTime := timePeriod[0].End
	for {
		filterLogsEventsOutput, err := dst.FilterAwsLogs(startTime, endTime)
		if err != nil {
			panic(err)
		}

		for _, value := range filterLogsEventsOutput {
			for _, event := range value.Events {
				message := formatMessage(event)
				fmt.Println(message)
			}
		}

		fmt.Println("\r- Press Ctrl+C to exit Terminal")
		time.Sleep(time.Duration(dst.WatchInterval) * time.Second)
		startTime = endTime
		endTime = time.Now().Unix() * 1000
		SetupCloseHandler()
	}
}

func SetupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		os.Exit(0)
	}()
}
