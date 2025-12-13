package repository

import (
	"log"

	"teamitmivhs/work-order-backend/config"
	"teamitmivhs/work-order-backend/models"
)

func GetAllMembers() ([]models.Member, error) {
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

func CreateMember(member *models.Member) error {
	_, err := config.DB.Exec("INSERT INTO members (Name, Password, Role, Status, Avatar) VALUES (?, ?, ?, ?, ?)", member.Name, member.Password, member.Role, member.Status, member.Avatar)
	return err
}

func GetMemberByName(name string) (*models.Member, error) {
	row := config.DB.QueryRow("SELECT ID, Name, Password, Role, Status, Avatar FROM members WHERE Name = ?", name)

	var m models.Member
	if err := row.Scan(&m.ID, &m.Name, &m.Password, &m.Role, &m.Status, &m.Avatar); err != nil {
		return nil, err
	}
	return &m, nil
}