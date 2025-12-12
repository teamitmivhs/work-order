package repository

import (
	"log"

	// Sesuaikan path import module Anda
	"teamitmivhs/work-order-backend/config"
	"teamitmivhs/work-order-backend/models"
)

func GetAllMembers() ([]models.Member, error) {
	rows, err := config.DB.Query("SELECT id, name, avatar, status FROM members")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.Member
	for rows.Next() {
		var m models.Member
		// Pastikan urutan Scan() sesuai dengan urutan SELECT di query
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
