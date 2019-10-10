package model

import (
	"time"
)

type SigninLogQ struct {
	SigninLog
	PaginationQ
	FromTime string `json:"from_time" form:"from_time"`
	ToTime   string `json:"to_time" form:"to_time"`
}

func (m SigninLogQ) Search(q *PaginationQ) (list *[]SigninLog, total uint, err error) {
	tx := db.Model(m.SigninLog).Preload("User")
	if m.LoginType != "" {
		tx = tx.Where("login_type = ?", m.LoginType)
	}
	if m.ClientIp != "" {
		tx = tx.Where("client_ip like ?", "%"+m.ClientIp+"%")
	}
	if m.UserName != "" {
		tx = tx.Where("user_name like ?", "%"+m.UserName+"%")
	}
	if m.Email != "" {
		tx = tx.Where("email like ?", "%"+m.Email+"%")
	}
	if m.FromTime != "" && m.ToTime != "" {
		tx = tx.Where("`created_at` BETWEEN ? AND ?", m.FromTime, m.ToTime)
	}
	list = &[]SigninLog{}
	total, err = crudAll(q, tx, list)
	return
}

type SigninLog struct {
	Id        uint      `gorm:"primary_key" json:"id"`
	UserId    uint      `gorm:"index" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	ClientIp  string    `json:"client_ip" form:"client_ip"`
	UserName  string    `json:"user_name" form:"user_name"`
	Email     string    `json:"email" form:"email"`
	LoginType string    `json:"login_type" comment:"密码 or LDAP or oss" form:"login_type"`
	UserAgent string    `json:"user_agent"`
	User      User      `gorm:"association_autoupdate:false;association_autocreate:false" json:"user"`
}

func (m *SigninLog) AfterFind() (err error) {
	return
}

//All
func (m SigninLog) All(q *PaginationQ) (list *[]SigninLog, total uint, err error) {
	tx := db.Model(m).Preload("User")
	if m.LoginType != "" {
		tx = tx.Where("login_type = ?", m.LoginType)
	}
	if m.ClientIp != "" {
		tx = tx.Where("client_ip like ?", "%"+m.ClientIp+"%")
	}
	if m.UserName != "" {
		tx = tx.Where("user_name like ?", "%"+m.UserName+"%")
	}
	if m.Email != "" {
		tx = tx.Where("email like ?", "%"+m.Email+"%")
	}
	list = &[]SigninLog{}
	total, err = crudAll(q, tx, list)
	return
}

//Create
func (m *SigninLog) Create() (err error) {
	m.Id = 0
	return db.Create(m).Error
}

//Delete
func (m *SigninLog) Delete(ids []int) (err error) {
	return db.Unscoped().Where("id in (?)", ids).Delete(m).Error
}
