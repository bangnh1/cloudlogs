package cmd

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/bangnh1/cloudlogs/pkg/awslogs"
	"github.com/bangnh1/cloudlogs/pkg/common"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
)

var group string
var start string
var end string
var interval string
var filter string
var limit string
var prefix string
var streams string
var goroutines int
var watchInterval int
var watch bool
var file bool

var getCmd = &cobra.Command{
	Use:   "get [flags]",
	Short: "Get logs from logs group",
	Run: func(cmd *cobra.Command, args []string) {
		execGetCmd()
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().StringVarP(&group, "group", "g", "", "AWS Cloudwatch logs group name(string - Required)")
	getCmd.MarkFlagRequired("group")
	getCmd.Flags().StringVarP(&start, "start", "s", "Jan 1, 1970 00:00:00 UTC", "Start Time - The default value is: Jan 1, 1970 00:00:00 UTC (string)")
	getCmd.Flags().StringVarP(&end, "end", "e", "", "End Time - Default is current time (string)")
	getCmd.Flags().StringVarP(&interval, "interval", "i", "0m", "Desired amount of time per file - The default value is: 0m, Ex: 0 minute, 5 hours, 2 days, 1 weeks (string)")
	getCmd.Flags().StringVar(&filter, "filter", "", "Logs filter pattern (string)")
	getCmd.Flags().StringVar(&limit, "limit", "10000", "Maximum events to return - The default value is: 10000 (string)")
	getCmd.Flags().StringVar(&prefix, "prefix", "", "Logs streams prefix (string)")
	getCmd.Flags().StringVar(&streams, "streams", "", "List of streams, seprated  by ; (string)")
	getCmd.Flags().IntVar(&goroutines, "goroutines", 10, "Number concurrent goroutines tasks - The default value is: 10 (int)")
	getCmd.Flags().IntVar(&watchInterval, "watchInterval", 10, "Watch interval - The default value is: 10s (int)")
	getCmd.Flags().BoolVarP(&watch, "watch", "w", false, "Watch logs in real time - The default value is: false (bool)")
	getCmd.Flags().BoolVarP(&file, "file", "f", false, "Write logs to file - The default value is: false (bool)")
}

func execGetCmd() {

	sess, _ := session.NewSession()
	logsvc := cloudwatchlogs.New(sess)

	if !file {
		setInterval := "0m"
		interval = setInterval
	}

	// Workaround for start from head when writing to files
	if file && start == "Jan 1, 1970 00:00:00 UTC" {
		setStartTime := "1h ago"
		start = setStartTime
	}

	timePeriod := common.TimeByInterval(&start, &end, &interval)

	var dst awslogs.Destination
	dst.Group = group
	dst.Limit, _ = strconv.ParseInt(limit, 10, 64)
	dst.Svc = logsvc
	dst.FilterPattern = filter
	dst.StreamPrefix = prefix
	dst.WatchInterval = watchInterval
	dst.Watch = watch
	dst.WriteToFile = file

	if streams != "" {
		streamList := strings.Split(streams, ";")
		var listPStream []*string
		for _, i := range streamList {
			listPStream = append(listPStream, &i)
		}
		dst.Streams = listPStream
	}

	if dst.Watch {
		awslogs.WatchAwsLogs(&dst, timePeriod)
	} else {
		awslogs.FilterAwsLogs(&dst, goroutines, timePeriod)
	}
}
