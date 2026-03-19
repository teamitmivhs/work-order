package repository

import (
	"database/sql"
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
	UpdateMemberStatus(memberID int, newStatus string) error
	// SetMemberPassword: update password + role member yang sudah ada di DB
	// Dipakai saat register — member sudah ada (nama dikenal), hanya set password
	SetMemberPassword(memberID int, hashedPassword string, role string) error
}

type memberRepository struct{}

func NewMemberRepository() MemberRepository {
	return &memberRepository{}
}

func (r *memberRepository) GetAllMembers() ([]models.Member, error) {
	rows, err := config.DB.Query("SELECT ID, Name, Role, Status, Avatar FROM members")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// FIX: inisialisasi sebagai empty slice bukan nil
	// agar JSON response selalu array [] bukan null
	members := make([]models.Member, 0)

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
	result, err := config.DB.Exec(
		"INSERT INTO members (Name, Password, Role, Status, Avatar) VALUES (?, ?, ?, ?, ?)",
		member.Name, member.Password, member.Role, member.Status, member.Avatar,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	member.ID = int(id)
	return nil
}

func (r *memberRepository) GetMemberByName(name string) (*models.Member, error) {
	row := config.DB.QueryRow(
		"SELECT ID, Name, Password, Role, Status, Avatar FROM members WHERE Name = ?",
		name,
	)

	var m models.Member
	if err := row.Scan(&m.ID, &m.Name, &m.Password, &m.Role, &m.Status, &m.Avatar); err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *memberRepository) GetMemberByID(id int) (*models.Member, error) {
	row := config.DB.QueryRow(
		"SELECT ID, Name, Password, Role, Status, Avatar FROM members WHERE ID = ?",
		id,
	)

	var m models.Member
	if err := row.Scan(&m.ID, &m.Name, &m.Password, &m.Role, &m.Status, &m.Avatar); err != nil {
		return nil, err
	}
	return &m, nil
}

// IsMemberAssigned mengecek apakah member sudah di-assign ke order tertentu
// FIX: nama kolom sudah dibuat eksplisit — order_id dan member_id
// (asumsi nama kolom di DB ikut diperbaiki menjadi order_id dan member_id)
func (r *memberRepository) IsMemberAssigned(orderID int64, memberID int) (bool, error) {
	row := config.DB.QueryRow(
		"SELECT COUNT(*) FROM executors WHERE order_id = ? AND member_id = ?",
		orderID, memberID,
	)

	var count int
	if err := row.Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

// SetMemberPassword menyimpan password hash dan role ke member yang sudah ada
// Dipanggil saat register — tidak membuat member baru, hanya update kolom Password dan Role
func (r *memberRepository) SetMemberPassword(memberID int, hashedPassword string, role string) error {
	_, err := config.DB.Exec(
		"UPDATE members SET Password = ?, Role = ? WHERE ID = ?",
		hashedPassword, role, memberID,
	)
	return err
}

// UpdateMemberStatus memperbarui status member ke database
// dipanggil dari PATCH /api/members/:id/status
func (r *memberRepository) UpdateMemberStatus(memberID int, newStatus string) error {
	validStatuses := map[string]bool{
		"standby":   true,
		"onjob":     true,
		"support":   true,
		"nextshift": true,
		"offduty":   true,
	}
	if !validStatuses[newStatus] {
		return sql.ErrNoRows // pakai sentinel; caller bisa wrap menjadi 400
	}

	_, err := config.DB.Exec(
		"UPDATE members SET Status = ? WHERE ID = ?",
		newStatus, memberID,
	)
	return err
}
