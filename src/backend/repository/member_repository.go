package repository

import (
	"log"

	"teamitmivhs/work-order-backend/config"
	"teamitmivhs/work-order-backend/models"
)

type MemberRepository interface {
	GetAllMembers() ([]models.Member, error)
	CreateMember(member *models.Member) error
	GetMemberByName(name string) (*models.Member, error)
	GetMemberByID(id int) (*models.Member, error)
	IsMemberAssigned(orderID int64, memberID int) (bool, error)
}

type memberRepository struct {
	db interface{}
}

func NewMemberRepository() MemberRepository {
	return &memberRepository{
		db: config.DB,
	}
}

func (r *memberRepository) GetAllMembers() ([]models.Member, error) {
	rows, err := config.DB.Query("SELECT ID, Name, Role, Status, Avatar FROM members")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.Member
	for rows.Next() {
		var m models.Member
		if err := rows.Scan(&m.ID, &m.Name, &m.Role, &m.Status, &m.Avatar); err != nil {
			log.Printf("Error scanning member row: %v", err)
			return nil, err
		}
		members = append(members, m)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return members, nil
}

func (r *memberRepository) CreateMember(member *models.Member) error {
	result, err := config.DB.Exec("INSERT INTO members (Name, Password, Role, Status, Avatar) VALUES (?, ?, ?, ?, ?)", member.Name, member.Password, member.Role, member.Status, member.Avatar)
	if err != nil {
		return err
	}

	// Get inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	member.ID = int(id)
	return nil
}

func (r *memberRepository) GetMemberByName(name string) (*models.Member, error) {
	row := config.DB.QueryRow("SELECT ID, Name, Password, Role, Status, Avatar FROM members WHERE Name = ?", name)

	var m models.Member
	if err := row.Scan(&m.ID, &m.Name, &m.Password, &m.Role, &m.Status, &m.Avatar); err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *memberRepository) GetMemberByID(id int) (*models.Member, error) {
	row := config.DB.QueryRow("SELECT ID, Name, Password, Role, Status, Avatar FROM members WHERE ID = ?", id)

	var m models.Member
	if err := row.Scan(&m.ID, &m.Name, &m.Password, &m.Role, &m.Status, &m.Avatar); err != nil {
		return nil, err
	}
	return &m, nil
}

// IsMemberAssigned checks if member is assigned to an order
func (r *memberRepository) IsMemberAssigned(orderID int64, memberID int) (bool, error) {
	row := config.DB.QueryRow("SELECT COUNT(*) FROM executors WHERE ID = ? AND Executors = ?", orderID, memberID)

	var count int
	if err := row.Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}
