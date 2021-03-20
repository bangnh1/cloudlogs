package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/bangnh1/cloudlogs/pkg/awslogs"
	"github.com/bangnh1/cloudlogs/pkg/common"
	"os"
	"strconv"
	"strings"
)

func main() {

	sess, _ := session.NewSession()
	logsvc := cloudwatchlogs.New(sess)

	listCommand := flag.NewFlagSet("list", flag.ExitOnError)
	getCommand := flag.NewFlagSet("get", flag.ExitOnError)

	logGroup := getCommand.String("group", "", "AWS Cloudwatch logs group name - Required")
	start := getCommand.String("start", "Jan 1, 1970 00:00:00 UTC", "Start Time - The default value is: Jan 1, 1970 00:00:00 UTC")
	end := getCommand.String("end", "", "End Time - Default is current time")
	interval := getCommand.String("interval", "0m", "Batch mode - Write logs to files - The default value is: 0m, Ex: 0 minute, 5 hours, 2 days, 1 weeks")
	filter := getCommand.String("filter", "", "Logs filter pattern")
	limit := getCommand.String("limit", "10000", "Maximum events to return - The default value is: 10000")
	prefix := getCommand.String("prefix", "", "Logs streams prefix")
	streams := getCommand.String("streams", "", "List of streams, seprated  by ;")
	goroutines := getCommand.Int("goroutines", 10, "Number concurrent goroutines tasks - The default value is: 10")
	watchInterval := getCommand.Int("watchInterval", 10, "Watch interval - The default value is: 10s")
	watch := getCommand.Bool("watch", false, "Watch logs in real time - The default value is: false")
	writeToFile := getCommand.Bool("file", false, "Write logs to file - The default value is: false")

	listGroupPrefix := listCommand.String("group", "", "Cloudwatch logs group prefix")
	listStreamPrefix := listCommand.String("prefix", "", "Logs streams prefix")
	listStreamDescending := listCommand.Bool("descending", false, "If the value is true, results are returned in descending order")
	listStreamOrderBy := listCommand.String("orderby", "LogStreamName", "Order by (LastEventTime || LogStreamName) The default value is LogStreamName")
	listLimit := listCommand.String("limit", "50", "Maximum groups/streams to return - The default value is 50")

	if len(os.Args) < 2 {
		fmt.Println("list or count subcommand is required")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "list":
		listCommand.Parse(os.Args[3:])
	case "get":
		getCommand.Parse(os.Args[2:])
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	if getCommand.Parsed() {

		// Required Flags
		if *logGroup == "" {
			getCommand.PrintDefaults()
			os.Exit(1)
		}

		if !*writeToFile {
			setInterval := "0m"
			interval = &setInterval
		}

		if *writeToFile && *start == "Jan 1, 1970 00:00:00 UTC" {
			setStartTime := "1h ago"
			start = &setStartTime
		}

		timePeriod := common.TimeByInterval(start, end, interval)

		var dst awslogs.Destination
		dst.Group = *logGroup
		dst.Limit, _ = strconv.ParseInt(*limit, 10, 64)
		dst.Svc = logsvc
		dst.FilterPattern = *filter
		dst.StreamPrefix = *prefix
		dst.WatchInterval = *watchInterval
		dst.Watch = *watch
		dst.WriteToFile = *writeToFile

		if *streams != "" {
			streamList := strings.Split(*streams, ";")
			var listPStream []*string
			for _, i := range streamList {
				listPStream = append(listPStream, &i)
			}
			dst.Streams = listPStream
		}

		if dst.Watch {
			awslogs.WatchAwsLogs(&dst, timePeriod)
		} else {
			awslogs.FilterAwsLogs(&dst, *goroutines, timePeriod)
		}
	}

	if listCommand.Parsed() {
		var dst awslogs.Destination
		dst.Group = *listGroupPrefix
		dst.Svc = logsvc
		dst.Limit, dst.MaxItems = awslogs.LimitHandler(listLimit)
		dst.Descending = *listStreamDescending
		dst.OrderBy = *listStreamOrderBy
		dst.StreamPrefix = *listStreamPrefix

		if os.Args[2] == "groups" {
			awslogs.ListLogGroups(&dst)
		}

		if os.Args[2] == "streams" {
			awslogs.ListStreams(&dst)
		}
	}

}
