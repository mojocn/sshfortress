package handler

import (
	"fmt"
)

type ldapObj struct {
	End      string `json:"end"`
	Sid      string `json:"sid"`
	Ref      string `json:"ref"`
	Mail     string `json:"mail" comment:"eg dejavuzhou@qq.com"`
	MemberOf string `json:"memberOf" comment:"eg CN=g-cmp,OU=MailGroup,OU=BJHT,DC=corp,DC=mojotv,DC=net"`
	User     string `json:"user" comment:"eg zhouxing1"`
	Display  string `json:"display" comment:"eg Eric Zhou"`
}

func (l *ldapObj) String() string {
	return fmt.Sprintf("%s?sid=%s&ref=%s", l.End, l.Sid, l.Ref)
}
