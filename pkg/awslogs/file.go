package awslogs

import (
	"bufio"
	"fmt"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"os"
	"strconv"
)

/**
 * @todo FORMAT MESSAGE
 */
func formatMessage(event *cloudwatchlogs.FilteredLogEvent) string {
	// timeStamp := *event.Timestamp
	message := *event.Message

	// timeStampString := common.TimeInt64ToString(timeStamp)
	formatedMessage := message + "\n"

	return formatedMessage
}

func WriteToFile(filename string, data []*cloudwatchlogs.FilterLogEventsOutput) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		panic(err)
	}

	datawriter := bufio.NewWriter(file)
	for _, value := range data {

		for _, event := range value.Events {
			message := formatMessage(event)
			_, _ = datawriter.WriteString(message)
		}
	}
	datawriter.Flush()
	file.Close()
}

func WriteToCli(data []*cloudwatchlogs.FilterLogEventsOutput) {
	for _, value := range data {

		for _, event := range value.Events {
			message := formatMessage(event)
			fmt.Println(message)
		}
	}
}

func LimitHandler(num *string) (int64, int64) {
	limit, _ := strconv.ParseInt(*num, 10, 64)
	count := limit / 50

	switch count {
	case 0:
		return limit, limit
	default:
		return 50, limit
	}
}
