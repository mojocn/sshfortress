package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"regexp"
)

type CmdFilter struct {
	Command string `json:"command" comment:"需要过滤的ssh命令 默认正则匹配" form:"command"`
	Msg     string `json:"msg" comment:"告警消息" form:"msg"`
	Enable  uint   `gorm:"default:2" json:"enable" form:"enable" comment:"2-开启 4-关闭"`
}

func (m CmdFilter) IsMatchFilterRule(command string) bool {
	if m.Enable == 4 {
		return false
	}
	re := regexp.MustCompile(fmt.Sprintf(`\b%s\b`, m.Command))
	return re.MatchString(command)
}

type JsonArraySshFilter []CmdFilter

//MatchTest 检测命令是否符合
func (fr JsonArraySshFilter) MatchTest(command string) (rule CmdFilter, isMatch bool) {
	for _, sf := range fr {
		if isMatch = sf.IsMatchFilterRule(command); isMatch {
			rule = sf
			return
		}
	}
	return
}

func (o JsonArraySshFilter) Value() (driver.Value, error) {
	b, err := json.Marshal(o)
	return string(b), err
}

func (o *JsonArraySshFilter) Scan(input interface{}) (err error) {

	switch v := input.(type) {
	case []byte:
		return json.Unmarshal(v, o)
	case string:
		return json.Unmarshal([]byte(v), o)
	default:
		err = fmt.Errorf("unexpected type %T in JsonArraySshFilter", v)
	}
	return err
}
