package cmd

import (
	"fmt"
	"github.com/go-programming-tour-book/tour/internal/timer"
	"github.com/spf13/cobra"
	"log"
	"strconv"
	"strings"
	"time"
)

var calculateTime string
var duration string

var timeCmd = &cobra.Command{
	Use: "time",
	Short: "时间格式处理",
	Long: "时间格式处理long",
	Run: func(cmd *cobra.Command, args []string) {
		var curTime time.Time
		var layout = "2006-01-02 15:04:05"
		if calculateTime == "" {
			curTime = timer.GetNowTime()
		} else {
			var err error
			space := strings.Count(calculateTime, " ")
			if space == 0 {
				layout = "2006-01-02"
			}
			if space == 1 {
				layout = "2006-01-02 15:04"
			}
			curTime, err = time.Parse(layout, calculateTime)
			if err != nil {
				t, _ := strconv.Atoi(calculateTime)
				curTime = time.Unix(int64(t), 0)
			}
		}
		t, err := timer.GetCalculateTime(curTime, duration)
		if err != nil {
			log.Fatalf("timer.GetCalculateTime err: %v", err)
		}
		log.Printf("输出结果: %s, %d", t.Format(layout), t.Unix())
	},
}

var nowCmd = &cobra.Command{
	Use: "now",
	Short: "获取当前时间",
	Long: "获取当前时间long",
	Run: func(cmd *cobra.Command, args []string) {
		nowTime := timer.GetNowTime()
		fmt.Printf("输出结果： %s %d", nowTime.Format("2006-01-02 15:04:05"), nowTime.Unix())
	},
}

func init() {
	timeCmd.Flags().StringVarP(&calculateTime, "basetime", "b", "", "需要计算的时间，有效单位为时间戳或已格式化后的时间")
	timeCmd.Flags().StringVarP(&duration, "duration", "d", "", `持续时间，有效时间单位为"ns", "us" (or "µ s"), "ms", "s", "m", "h"`)
}
