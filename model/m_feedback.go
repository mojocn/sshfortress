package model

import (
	"errors"
)

type FeedbackQ struct {
	PaginationQ
	Feedback
}

//Feedback 用户反馈 也可以当作工单系统
type Feedback struct {
	BaseModel
	PageUrl  string `json:"page_url" form:"page_url"`
	Email    string `json:"email" form:"email"`
	UserName string `json:"user_name" form:"user_name"`
	Content  string `gorm:"type:text" json:"content" form:"content"`
	Status   uint   `json:"status" Comment:"处理状态:0-unhandle 2-ignore 4-useful 8-done 16-todo"`
	UserId   uint   `gorm:"index;default:'0'" json:"user_id" comment:"添加者"`
	AtUserId uint   `gorm:"index;default:'0'" json:"at_user_id" comment:"指派给用户"`
	User     User   `gorm:"association_autoupdate:false;association_autocreate:false" json:"user"`
}

func (m *Feedback) AfterFind() (err error) {
	return
}

//One
func (m *Feedback) One() error {
	return crudOne(m)
}

//All
func (m Feedback) All(q *PaginationQ) (list *[]Feedback, total uint, err error) {
	tx := db.Model(m)
	list = &[]Feedback{}
	if m.PageUrl != "" {
		tx = tx.Where("page_url like ?", "%"+m.PageUrl+"%")
	}
	if m.Email != "" {
		tx = tx.Where("email like ?", "%"+m.Email+"%")
	}
	if m.UserName != "" {
		tx = tx.Where("user_name like ?", "%"+m.UserName+"%")
	}
	if m.Content != "" {
		tx = tx.Where("content like ?", "%"+m.Content+"%")
	}
	total, err = crudAll(q, tx, list)
	return
}

//Update
func (m *Feedback) Update() (err error) {
	if m.Id < 1 {
		return errors.New("id必须大于0")
	}
	return db.Model(m).Update(m).Error
}

//Create
func (m *Feedback) Create() (err error) {
	m.Id = 0
	return db.Create(m).Error
}

func (m *Feedback) Delete(ids []int) (err error) {
	return db.Unscoped().Where("id in (?)", ids).Delete(m).Error
}
