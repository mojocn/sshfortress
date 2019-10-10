package model

import (
	"errors"
	"github.com/jinzhu/gorm"
	"time"
)

//PaginationQ gin handler query binding struct
type PaginationQ struct {
	Size  uint        `form:"size" json:"size"`
	Page  uint        `form:"page" json:"page"`
	Total uint        `json:"total" form:"-"`
	Data  interface{} `json:"data" form:"-" comment:"data 必须是[]interface{} 指针"`
	Ok    bool        `json:"ok" form:"-"`
}

func (p *PaginationQ) Search(tx *gorm.DB) (err error) {
	p.Ok = true
	if p.Size == 9999 {
		return tx.Find(p.Data).Error
	}
	if p.Size < 1 {
		p.Size = 10
	}
	if p.Page < 1 {
		p.Page = 1
	}

	var total uint
	err = tx.Count(&total).Error
	if err != nil {
		return err
	}
	offset := p.Size * (p.Page - 1)
	err = tx.Limit(p.Size).Offset(offset).Find(p.Data).Error
	if err != nil {
		return err
	}
	p.Total = total
	return
}

func crudAll(p *PaginationQ, queryTx *gorm.DB, list interface{}) (uint, error) {
	if p.Size == 9999 {
		return 0, queryTx.Find(list).Error
	}
	if p.Size < 1 {
		p.Size = 10
	}
	if p.Page < 1 {
		p.Page = 1
	}

	var total uint
	err := queryTx.Count(&total).Error
	if err != nil {
		return 0, err
	}
	offset := p.Size * (p.Page - 1)
	err = queryTx.Limit(p.Size).Offset(offset).Find(list).Error
	if err != nil {
		return 0, err
	}
	return total, err
}

func crudOne(m interface{}) (err error) {
	if db.First(m).RecordNotFound() {
		return errors.New("resource is not found")
	}
	return nil
}

func crudDelete(m interface{}) (err error) {
	//WARNING When delete a record, you need to ensure it’s primary field has value, and GORM will use the primary key to delete the record, if primary field’s blank, GORM will delete all records for the model
	//primary key must be not zero value
	db := db.Unscoped().Delete(m)
	if err = db.Error; err != nil {
		return
	}
	if db.RowsAffected != 1 {
		return errors.New("resource is not found to destroy")
	}
	return nil
}
func mysqlTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}
