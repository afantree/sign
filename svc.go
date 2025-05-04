package sign

import (
	"fmt"
	"github.com/dromara/carbon/v2"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"math"
)

var timeNow = carbon.Now

func DailyLists(db *gorm.DB, rewards []interface{}, aid int64) (map[string]interface{}, error) {
	result := map[string]interface{}{}
	nowTime := timeNow()
	days := nowTime.DaysInMonth()
	signAt := nowTime.Format("Ymd")
	var signinfos []map[string]interface{}
	if err := db.Model(&SignDaily{}).
		Where("aid = ? and sign_at >= ?", aid, nowTime.Format("Ym01")).
		Select("sign_at").
		Order("sign_at desc").
		Find(&signinfos).Error; err != nil {
		return result, err
	}

	num := len(signinfos)
	if num > 0 {
		if cast.ToString(signinfos[0]["sign_at"]) == signAt {
			result["status"] = 1
		} else {
			result["status"] = 0
		}
	} else {
		result["status"] = 0
	}
	result["sign_num"] = num

	result["list"] = rewards[0:days]
	result["max_num"] = days
	return result, nil
}

func DailySignIn(db *gorm.DB, aid int64) error {
	nowTime := timeNow()
	info := map[string]interface{}{}
	signAt := nowTime.Format("Ymd")
	info["aid"] = aid
	info["sign_at"] = signAt

	var signinfos []map[string]interface{}
	if err := db.Model(&SignDaily{}).
		Where("aid = ? and sign_at >= ?", aid, nowTime.Format("Ym01")).
		Select("sign_at").
		Order("sign_at desc").
		Find(&signinfos).Error; err != nil {
		return err
	}
	num := len(signinfos)
	if num > 0 && cast.ToString(signinfos[0]["sign_at"]) == signAt {
		return fmt.Errorf("sign in already aid(%d)", aid)
	}

	return db.Model(&SignDaily{}).Create(info).Error
}

func CumulativeLists(db *gorm.DB, aid int64) (int, error) {
	nowTime := timeNow()
	result := []map[string]interface{}{}
	if err := db.Model(&SignCumulative{}).
		Where("aid = ? and sign_month = ?", aid, nowTime.Format("Ym")).
		Select("max(num) as maxnum").
		Find(&result).Error; err != nil {
		return -1, err
	}
	if len(result) == 0 {
		return 5, nil
	}
	num := cast.ToFloat64(result[0]["maxnum"])
	if num >= 6 {
		return 31, nil
	}
	rewardDays := math.Min(math.Min((num+1)*5, 30), cast.ToFloat64(nowTime.DaysInMonth()))
	return cast.ToInt(rewardDays), nil
}

func CumulativeSignIn(db *gorm.DB, aid int64) error {
	nowTime := timeNow()
	signMonth := nowTime.Format("Ym")

	// 判断是否领完奖
	result := []map[string]interface{}{}
	if err := db.Model(&SignCumulative{}).
		Where("aid = ? and sign_month = ?", aid, signMonth).
		Select("max(num) as maxnum").
		Find(&result).Error; err != nil {
		return err
	}
	var num float64 = 0
	if len(result) > 0 {
		num = cast.ToFloat64(result[0]["maxnum"])
		if num >= 6 {
			return fmt.Errorf("reward finish aid(%d)", aid)
		}
	}
	rewardDays := math.Min(math.Min((num+1)*5, 30), cast.ToFloat64(nowTime.DaysInMonth()))

	// 判断是否满足领奖条件
	var signNum int64
	if err := db.Model(&SignDaily{}).
		Where("aid = ? and sign_at >= ?", aid, nowTime.Format("Ym01")).
		Count(&signNum).Error; err != nil {
		return err
	}
	if cast.ToInt64(rewardDays) > signNum {
		return fmt.Errorf("不符合领取条件 aid(%d)", aid)
	}

	// 开始记录
	info := map[string]interface{}{}
	info["aid"] = aid
	info["sign_month"] = signMonth
	info["num"] = num + 1
	return db.Model(&SignCumulative{}).Create(info).Error
}
