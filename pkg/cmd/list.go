package cmd

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/bangnh1/cloudlogs/pkg/awslogs"
	"github.com/spf13/cobra"
)

var groupPrefix string
var StreamPrefix string
var descending bool
var orderBy string
var streamLimit string

var listCmd = &cobra.Command{
	Use:   "list [commands]",
	Short: "List logs groups or logs streams",
}

var listGroupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "List logs groups in Cloudwatch Logs",
	Run: func(cmd *cobra.Command, args []string) {
		execListGroupsCmd()
	},
}

var listStreamsCmd = &cobra.Command{
	Use:   "streams [flags]",
	Short: "List logs streams in a group",
	Run: func(cmd *cobra.Command, args []string) {
		execListStreamsCmd()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.AddCommand(listGroupsCmd)
	listCmd.AddCommand(listStreamsCmd)

	listStreamsCmd.Flags().StringVarP(&groupPrefix, "group", "g", "", "Cloudwatch logs group prefix")
	listStreamsCmd.Flags().StringVar(&StreamPrefix, "prefix", "", "Logs streams prefix")
	listStreamsCmd.Flags().BoolVar(&descending, "descending", false, "If the value is true, results are returned in descending order")
	listStreamsCmd.Flags().StringVar(&orderBy, "orderby", "LogStreamName", "If the value is true, results are returned in descending order")
	listStreamsCmd.Flags().StringVar(&streamLimit, "limit", "50", "If the value is true, results are returned in descending order")
}

func execListGroupsCmd() {

	sess, _ := session.NewSession()
	logsvc := cloudwatchlogs.New(sess)

	var dst awslogs.Destination
	dst.Group = groupPrefix
	dst.Svc = logsvc
	dst.Limit, dst.MaxItems = awslogs.LimitHandler(&streamLimit)

	awslogs.ListLogGroups(&dst)
}

func execListStreamsCmd() {

	sess, _ := session.NewSession()
	logsvc := cloudwatchlogs.New(sess)

	var dst awslogs.Destination
	dst.Group = groupPrefix
	dst.Svc = logsvc
	dst.Limit, dst.MaxItems = awslogs.LimitHandler(&streamLimit)
	dst.Descending = descending
	dst.OrderBy = orderBy
	dst.StreamPrefix = StreamPrefix

	awslogs.ListStreams(&dst)
}
