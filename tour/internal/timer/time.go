package timer

import "time"

func GetNowTime() time.Time {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	return time.Now().In(loc)
}

func GetCalculateTime(currTime time.Time, d string) (time.Time, error) {
	dutration, err := time.ParseDuration(d)
	if err != nil {
		return currTime, err
	}
	return currTime.Add(dutration), nil
}
