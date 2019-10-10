package model

import (
	"errors"
)

type SshFilterGroupQ struct {
	PaginationQ
	SshFilterGroup
}

//All
func (m SshFilterGroupQ) Search() (pagination PaginationQ, err error) {
	pagination = m.PaginationQ
	pagination.Data = &[]SshFilterGroup{}
	tx := db.Model(m.SshFilterGroup) //.Where("ancestor_path like ?", m.qAncetorPath())
	if m.Remark != "" {
		tx = tx.Where("remark like ?", "%"+m.Remark+"%")
	}
	if m.Name != "" {
		tx = tx.Where("name like ?", "%"+m.Name+"%")
	}
	err = pagination.Search(tx)
	return
}

//SshFilterGroup
type SshFilterGroup struct {
	BaseModel
	Name    string             `gorm:"index" json:"name" form:"name"`
	Remark  string             `gorm:"index" json:"remark" form:"remark"`
	Filters JsonArraySshFilter `gorm:"type:json" json:"filters" form:"filters"`
}

func (m *SshFilterGroup) AfterFind() (err error) {
	return
}

//One
func (m *SshFilterGroup) One() error {
	return db.First(m).Error
}

//Update
func (m *SshFilterGroup) Update() (err error) {
	if m.Id < 1 {
		return errors.New("id必须大于0")
	}
	return db.Model(m).Update(m).Error
}

//Create
func (m *SshFilterGroup) Create() (err error) {
	m.Id = 0
	return db.Create(m).Error
}

//Delete
func (m *SshFilterGroup) Delete() (err error) {
	if m.Id < 1 {
		return errors.New("id必须大于0")
	}
	return crudDelete(m)
}
