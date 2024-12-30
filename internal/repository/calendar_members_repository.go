package repository

import (
	"schedule_table/internal/interface/request"
	"schedule_table/internal/interface/response"
	"schedule_table/internal/model/dao"
	"schedule_table/util"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MembersRepository interface {
	GetMembers(calendarId string) (*[]response.Member, error)
	GetMemberId(memberId string) (*response.Member, error)
	CreateNewMember(calendarId string, req *request.CreateNewMember) (*response.Member, error)
}

type membersRepository struct {
	db *gorm.DB
}

func (mrp *membersRepository) GetMemberId(memberId string) (*response.Member, error) {
	var member *response.Member

	if err := mrp.db.Model(&dao.Members{}).First(&member, "id = ?", memberId).Error; err != nil {
		return nil, err
	}

	return member, nil
}

func (m *membersRepository) GetMembers(calendarId string) (*[]response.Member, error) {
	var Members *[]response.Member

	if err := m.db.Find(&Members, "calendar_id = ?", calendarId).Error; err != nil {
		return nil, err
	}

	return Members, nil
}

func (m *membersRepository) CreateNewMember(calendarId string, req *request.CreateNewMember) (*response.Member, error) {

	newMember := &dao.Members{
		Id:          uuid.New(),
		ImageURL:    req.File.Filename,
		CalendarId:  util.Must(uuid.Parse(calendarId)),
		Name:        req.Name,
		Nickname:    req.NickName,
		Color:       req.Color,
		Description: req.Description,
		Position:    req.Position,
		Email:       req.Email,
		Telephone:   req.Telephone,
	}

	if err := m.db.Model(&dao.Members{}).Create(&newMember).Error; err != nil {
		return nil, err
	}

	responseMember := &response.Member{}
	if err := m.db.Model(&dao.Members{}).First(&responseMember, newMember.Id).Error; err != nil {
		return nil, err
	}

	return responseMember, nil
}

func NewMemberRepository(db *gorm.DB) MembersRepository {
	return &membersRepository{
		db: db,
	}
}
