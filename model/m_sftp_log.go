package model

import (
	"errors"
	"time"
)

type SftpLog struct {
	Id        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserId    uint      `gorm:"index" json:"user_id" form:"user_id" comment:"用户ID"`
	MachineId uint      `gorm:"index" json:"machine_id" form:"machine_id" comment:"机器id"`
	SshUser   string    `json:"ssh_user" comment:"sftp登陆ssh账号"`
	ClientIp  string    `json:"client_ip" comment:"客户端浏览器IP"`
	Status    uint      `json:"status" comment:"0-未标记 2-正常 4-警告 8-危险 16-致命"`
	Remark    string    `json:"remark" comment:"备注"`
	Action    string    `json:"action" comment:"操作" form:"action"`
	Path      string    `json:"path" comment:"sftp路径" form:"path"`

	Machine Machine `gorm:"association_autoupdate:false;association_autocreate:false" json:"machine"`
	User    User    `gorm:"association_autoupdate:false;association_autocreate:false" json:"user"`
}

type SftpLogQ struct {
	SftpLog
	PaginationQ
	FromTime string `json:"from_time" form:"from_time"`
	ToTime   string `json:"to_time" form:"to_time"`
}

func (m SftpLogQ) Search(u *User) (list *[]SftpLog, total uint, err error) {
	list = &[]SftpLog{}
	tx := db.Model(m.SftpLog).Preload("User").Preload("Machine")

	if m.Path != "" {
		tx = tx.Where("path like ?", "%"+m.Path+"%")
	}
	if m.Action != "" {
		tx = tx.Where("action = ?", m.Action)
	}
	if m.FromTime != "" && m.ToTime != "" {
		tx = tx.Where("`created_at` BETWEEN ? AND ?", m.FromTime, m.ToTime)
	}
	if u.IsAdmin() {
		if m.UserId > 0 {
			tx = tx.Where("user_id = ?", m.UserId)
		}
		if m.MachineId > 0 {
			tx = tx.Where("machine_id = ?", m.MachineId)
		}
	} else {
		//非管理员 智能看自己的日志
		//不支持搜索搜索
		tx = tx.Where("`user_id` = ?", u.Id)
	}
	total, err = crudAll(&m.PaginationQ, tx, list)
	return
}

func (m *SftpLog) AfterFind() (err error) {
	return
}

//One
func (m *SftpLog) One() error {
	return crudOne(m)
}

//All
func (m SftpLog) All(q *PaginationQ) (list *[]SftpLog, total uint, err error) {
	list = &[]SftpLog{}
	tx := db.Model(m)
	total, err = crudAll(q, tx, list)
	return
}

//Update
func (m *SftpLog) Update() (err error) {
	if m.Id < 1 {
		return errors.New("id必须大于0")
	}
	return db.Model(m).Update(m).Error
}

//Create
func (m *SftpLog) Create() (err error) {
	m.Id = 0
	return db.Create(m).Error
}

//Delete
func (m *SftpLog) Delete(ids []int) (err error) {
	return db.Unscoped().Where("id in (?)", ids).Delete(m).Error
}
