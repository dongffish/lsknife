package lsdatetime

import (
	"time"
)

// ExecDateDiff calc date difference from bengin to end
func ExecDateDiff(begin time.Time, end time.Time) time.Duration {
	return end.Sub(begin)
}

// ExecTime 计算二个日期之间时间差，返回time.Duration。
func ExecTime(begin time.Time, end time.Time) time.Duration {
	return end.Sub(begin)
}

// ExecTimeSeconds 计算二个日期之间时间差，返回:秒
func ExecTimeSeconds(begin string, end string) (int64, error) {
	const baseFormat = "2006-01-02 15:04:05"
	var err error
	var beginTime, endTime time.Time
	beginTime, err = time.Parse(baseFormat, begin)
	if err != nil {
		return 0, err
	}
	endTime, err = time.Parse(baseFormat, end)
	if err != nil {
		return 0, err
	}
	dt := ExecTime(beginTime, endTime)
	return int64(dt.Seconds()), err
}

// SetTimeLocation 设置时区。如果名字为空，则默认Asia/Shanghai
func SetTimeLocation(name string) (*time.Location, error) {
	if name == "" {
		name = "Asia/Shanghai"
	}
	nyc, err := time.LoadLocation(name) // 设置时区
	if err != nil {
		nyc = time.FixedZone("CST", 8*3600) //替换上海时区
	}
	time.Local = nyc
	return nyc, err
}

// FirstDateOfMonth 获取输入时间点的月初日期
func FirstDateOfMonth(enterTime time.Time) string {
	firstDate := enterTime.Format("2006-01") + "-01"
	return firstDate
}

// LastDateOfMonth 获取输入时间点的月底日期
func LastDateOfMonth(enterTime time.Time) string {
	firstDate := enterTime.Format("2006-01") + "-01"
	middle, _ := time.ParseInLocation("2006-01-02", firstDate, time.Local)
	lastDate := middle.AddDate(0, 1, -1).Format("2006-01-02")
	return lastDate
}
