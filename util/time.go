/**
 * ImindLab
 *
 * Create by SongLi on 2023/12/28
 * Copyright © 2023 imind.tech All rights reserved.
 */

package util

import (
	"fmt"
	"time"
)

const (
	DefaultFmt = "2006-01-02T15:04:05-07:00"

	DateTimeFmt = "2006-01-02 15:04:05"

	DateFmt = "2006-01-02"

	TimeFmt = "15:04:05"

	DateTimeFmtCn = "2006年01月02日 15时04分05秒"

	DateFmtCn = "2006年01月02日"

	TimeFmtCn = "15时04分05秒"
)

func GetNowWithMillisecond() int64 {
	return time.Now().UnixNano() / 1e6
}

func GetMonthDays(year, month int) int {
	switch month {
	case 1, 3, 5, 7, 8, 10, 12:
		return 31
	case 4, 6, 9, 11:
		return 30
	case 2:
		if IsLeapYear(year) {
			return 29
		}
		return 28
	default:
		panic(fmt.Sprintf("Illegal month:%d", month))
	}
}

func IsLeapYear(year int) bool {
	if year%100 == 0 {
		return year%400 == 0
	}
	return year%4 == 0
}
