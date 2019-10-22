package model

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"sshfortress/util"
	"strings"
	"time"
)

const (
	UserRoleAdmin = 2
	UserRoleUser  = 4
)

type UserQ struct {
	User
	PaginationQ
}

//User
type User struct {
	BaseModel
	Name             string         `gorm:"type:varchar(60);unique_index" json:"name" comment:"ssh账号前缀sshfortress_" form:"name"`
	Email            string         `gorm:"type:varchar(100);unique_index" json:"email" form:"email"`
	RealName         string         `json:"real_name" form:"real_name"`
	Mobile           string         `gorm:"type:varchar(20)" json:"mobile" form:"mobile"`
	Password         string         `json:"-"`
	SshPassword      string         `json:"ssh_password" comment:"初始ssh密码为随机12为字母数字"`
	Role             uint           `json:"role" comment:"2-管理员 4-用户"`
	ParentId         uint           `json:"parent_id" comment:"父级ID"`
	AncestorPath     string         `gorm:"index;default:'0'" json:"ancestor_path" comment:"祖先用户路径"`
	ExpiredAt        *time.Time     `json:"expired_at" comment:"过期时间"`
	MemberOf         string         `gorm:"default:''" json:"member_of"`
	InputPassword    string         `gorm:"-" json:"input_password,omitempty"`
	InputSshPassword string         `gorm:"-" json:"input_ssh_password,omitempty"`
	Avatar           string         `gorm:"default:'//p1.ssl.qhimg.com/t01ff98c4a29f7a7db5.png'" json:"avatar"`
	GithubToken      string         `json:"github_token"`
	SshFilterGroupId uint           `gorm:"index" json:"ssh_filter_group_id" form:"ssh_filter_group_id"`
	SshFilterGroup   SshFilterGroup `json:"ssh_filter_group"`
}

func (m User) SshUserName() string {
	prefix := viper.GetString("app.ssh_user_prefix")
	if prefix == "" {
		prefix = "sshfortress_"
	}
	return fmt.Sprintf("%s%s", prefix, strings.TrimSpace(m.Name))
}
func (m User) IsAdmin() bool {
	return m.Role == UserRoleAdmin
}
func (m *User) AfterFind() (err error) {
	return
}

//One
func (m *User) One() error {
	return crudOne(m)
}
func (m User) qAncetorPath() string {
	return fmt.Sprintf("%s/%d%%", m.AncestorPath, m.Id)
}

//All
func (m User) All(q *PaginationQ) (list *[]User, total uint, err error) {
	list = &[]User{}
	tx := db.Model(m).Preload("SshFilterGroup")
	if m.Name != "" {
		tx = tx.Where("name like ?", "%"+m.Name+"%")
	}
	if m.Email != "" {
		tx = tx.Where("email like ?", "%"+m.Email+"%")
	}
	if m.Mobile != "" {
		tx = tx.Where("mobile like ?", "%"+m.Mobile+"%")
	}
	if m.RealName != "" {
		tx = tx.Where("real_name like ?", "%"+m.RealName+"%")
	}
	total, err = crudAll(q, tx, list)
	return
}

//Update
func (m *User) Update() (err error) {
	if m.Id < 1 {
		return errors.New("id必须大于0")
	}
	m.makePassword()
	return db.Model(m).Update(m).Error
}

//Create
func (m *User) Create() (err error) {
	m.Id = 0
	m.makePassword()

	return db.Create(m).Error
}

//Delete
func (m *User) Delete() (err error) {
	if m.Id < 2 {
		return errors.New("id must be larger than 1")
	}
	return crudDelete(m)
}

//Login
func (m *User) Login(ip string) (*jwtObj, error) {
	if m.InputPassword == "" {
		return nil, errors.New("password is required")
	}

	err := db.Where("email = ?", m.Email).Where("expired_at > ?", mysqlTime(time.Now())).First(&m).Error
	if err != nil {
		return nil, err
	}
	//password is set to bcrypt check
	if err := bcrypt.CompareHashAndPassword([]byte(m.Password), []byte(m.InputPassword)); err != nil {
		return nil, err
	}
	m.Password = ""
	data, err := jwtGenerateToken(m)
	return data, err
}

func (m *User) makePassword() {
	if m.InputPassword != "" {
		if bytes, err := bcrypt.GenerateFromPassword([]byte(m.InputPassword), bcrypt.DefaultCost); err != nil {
			logrus.WithError(err).Error("crypt password failed")
		} else {
			m.Password = string(bytes)
		}
	}
	if m.InputSshPassword != "" {
		m.SshPassword = m.InputSshPassword
	} else {
		m.SshPassword = util.RandomDigitAndLetters(12)
	}
}

func (m User) LoginLdap(email, userName, realName, memberOf string) (*jwtObj, error) {
	ex := time.Now().Add(time.Hour * 24 * 3650)
	pw := "123456"
	u := &User{
		Role:             4,
		MemberOf:         memberOf,
		Name:             userName,
		RealName:         realName,
		Email:            email,
		AncestorPath:     "0/1",
		ParentId:         1,
		ExpiredAt:        &ex,
		InputSshPassword: pw,
		InputPassword:    pw,
	}
	u.makePassword()
	err := db.Model(m).FirstOrCreate(u, User{Email: email, Name: userName}).Error
	if err != nil {
		return nil, err
	}
	return jwtGenerateToken(u)
}

func (m User) LoginGithub(email, userName, realName, memberOf, avatar, token string) (*jwtObj, error) {
	ex := time.Now().Add(time.Hour * 24 * 3650)
	pw := "123456"
	u := &User{
		Role:             4,
		MemberOf:         memberOf,
		Name:             userName,
		RealName:         realName,
		Email:            email,
		AncestorPath:     "0/1",
		ParentId:         1,
		ExpiredAt:        &ex,
		InputSshPassword: pw,
		InputPassword:    pw,
		Avatar:           avatar,
		GithubToken:      token,
	}
	u.makePassword()
	err := db.Model(m).FirstOrCreate(u, User{Email: email, Name: userName}).Error
	if err != nil {
		return nil, err
	}
	return jwtGenerateToken(u)
}

func (m User) MustSshFilterGroup() (g *SshFilterGroup) {
	if m.SshFilterGroupId < 1 {
		logrus.Error("SshFilterGroupId is zero")
		return
	}

	g = &SshFilterGroup{}
	err := db.First(g, m.SshFilterGroupId).Error
	if err != nil {
		logrus.WithError(err).Error("SshFilterGroup db first failed")
	}
	return
}

func FsshUserAuth(user, password string) bool {
	return true
}
