package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"teamitmivhs/work-order-backend/config"
	"teamitmivhs/work-order-backend/models"
	"teamitmivhs/work-order-backend/utils"
)

var ErrWorkOrderManagedStatus = errors.New("on job status is managed by work orders")

func manualMemberStatusAllowed(status string) bool {
	switch status {
	case "standby", "nextshift", "offduty":
		return true
	default:
		return false
	}
}

type MemberRepository interface {
	GetAllMembers() ([]models.Member, error)
	GetAdminMembers(filter string) ([]models.Member, error)
	HasActiveAdmin() (bool, error)
	CreateMember(member *models.Member) error
	GetMemberByName(name string) (*models.Member, error)
	GetMemberByID(id int) (*models.Member, error)
	IsMemberAssigned(orderID int64, memberID int) (bool, error)
	UpdateMemberStatus(memberID int, newStatus string) error
	UpdateMemberName(memberID int, name string) error
	SetMemberPassword(memberID int, hashedPassword string, role string, division string, batchYear string) error
	UpdateMemberPassword(memberID int, hashedPassword string) error
	UpdateMemberAvatar(memberID int, avatarFilename string) error
	ApproveMember(memberID int, role string, division string, batchYear string, canHandle bool, approvedBy int) error
	RejectMember(memberID int) error
	DisableMember(memberID int) error
	MarkMemberAsAlumni(memberID int, graduationYear int) error
	GraduateBatch(batchYear string, graduationYear int) (int64, error)
	ChangeRole(memberID int, role string) error
}

type memberRepository struct{}

func NewMemberRepository() MemberRepository {
	return &memberRepository{}
}

func (r *memberRepository) GetAllMembers() ([]models.Member, error) {
	rows, err := config.DB.Query(`
		SELECT ID, Name, COALESCE(Role, ''), COALESCE(Division, ''), Status, Avatar, AccountStatus, MembershipStatus,
		       COALESCE(BatchYear, ''), GraduationYear, CanHandleWorkOrder,
		       DATE_FORMAT(RegisteredAt, '%Y-%m-%d %H:%i:%s'),
		       DATE_FORMAT(ApprovedAt, '%Y-%m-%d %H:%i:%s'),
		       ApprovedBy
		FROM members
		WHERE AccountStatus = 'active'
		  AND MembershipStatus = 'active'
		  AND CanHandleWorkOrder = 1
		  AND COALESCE(Role, '') <> 'Guest'
		ORDER BY Name ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members := make([]models.Member, 0)
	for rows.Next() {
		var m models.Member
		if err := scanMember(rows, &m); err != nil {
			log.Printf("Error scanning member row: %v", err)
			return nil, err
		}
		members = append(members, m)
	}
	return members, rows.Err()
}

func (r *memberRepository) GetAdminMembers(filter string) ([]models.Member, error) {
	query := `
		SELECT ID, Name, COALESCE(Role, ''), COALESCE(Division, ''), Status, Avatar, AccountStatus, MembershipStatus,
		       COALESCE(BatchYear, ''), GraduationYear, CanHandleWorkOrder,
		       DATE_FORMAT(RegisteredAt, '%Y-%m-%d %H:%i:%s'),
		       DATE_FORMAT(ApprovedAt, '%Y-%m-%d %H:%i:%s'),
		       ApprovedBy
		FROM members
		WHERE COALESCE(Role, '') <> 'Guest'
	`
	args := []any{}
	switch filter {
	case "pending":
		query += " AND AccountStatus = 'pending'"
	case "active":
		query += " AND AccountStatus = 'active' AND MembershipStatus = 'active'"
	case "alumni":
		query += " AND MembershipStatus = 'alumni'"
	case "disabled":
		query += " AND AccountStatus = 'disabled'"
	case "rejected":
		query += " AND AccountStatus = 'rejected'"
	}
	query += " ORDER BY RegisteredAt DESC, Name ASC"

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members := make([]models.Member, 0)
	for rows.Next() {
		var m models.Member
		if err := scanMember(rows, &m); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, rows.Err()
}

func (r *memberRepository) HasActiveAdmin() (bool, error) {
	var count int
	err := config.DB.QueryRow(`
		SELECT COUNT(*)
		FROM members
		WHERE Role IN ('Admin', 'Guru')
		  AND AccountStatus = 'active'
		  AND MembershipStatus = 'active'
	`).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *memberRepository) CreateMember(member *models.Member) error {
	result, err := config.DB.Exec(
		`INSERT INTO members
			(Name, Password, Role, Division, Status, Avatar, AccountStatus, MembershipStatus, BatchYear, CanHandleWorkOrder)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		member.Name, member.Password, member.Role, nullableString(member.Division), member.Status, member.Avatar,
		member.AccountStatus, member.MembershipStatus, nullableString(member.BatchYear), member.CanHandleWorkOrder,
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
		`SELECT ID, Name, Password, COALESCE(Role, ''), COALESCE(Division, ''), Status, Avatar, AccountStatus, MembershipStatus,
		        COALESCE(BatchYear, ''), GraduationYear, CanHandleWorkOrder,
		        DATE_FORMAT(RegisteredAt, '%Y-%m-%d %H:%i:%s'),
		        DATE_FORMAT(ApprovedAt, '%Y-%m-%d %H:%i:%s'),
		        ApprovedBy
		 FROM members WHERE Name = ?`,
		name,
	)
	var m models.Member
	if err := scanMemberWithPassword(row, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *memberRepository) GetMemberByID(id int) (*models.Member, error) {
	row := config.DB.QueryRow(
		`SELECT ID, Name, Password, COALESCE(Role, ''), COALESCE(Division, ''), Status, Avatar, AccountStatus, MembershipStatus,
		        COALESCE(BatchYear, ''), GraduationYear, CanHandleWorkOrder,
		        DATE_FORMAT(RegisteredAt, '%Y-%m-%d %H:%i:%s'),
		        DATE_FORMAT(ApprovedAt, '%Y-%m-%d %H:%i:%s'),
		        ApprovedBy
		 FROM members WHERE ID = ?`,
		id,
	)
	var m models.Member
	if err := scanMemberWithPassword(row, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

// IsMemberAssigned mengecek apakah member sudah di-assign ke order
// FIX: pakai kolom order_id dan member_id (bukan ID dan Executors yang sudah di-rename)
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

// UpdateMemberStatus memperbarui status member
func (r *memberRepository) UpdateMemberStatus(memberID int, newStatus string) error {
	if newStatus == "onjob" {
		return ErrWorkOrderManagedStatus
	}
	if !manualMemberStatusAllowed(newStatus) {
		return fmt.Errorf("invalid status: %s", newStatus)
	}
	tx, err := config.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var currentStatus string
	if err := tx.QueryRow(
		"SELECT Status FROM members WHERE ID = ? FOR UPDATE",
		memberID,
	).Scan(&currentStatus); err != nil {
		return err
	}
	if currentStatus == "onjob" {
		return ErrWorkOrderManagedStatus
	}

	if _, err := tx.Exec(
		"UPDATE members SET Status = ? WHERE ID = ?",
		newStatus, memberID,
	); err != nil {
		return err
	}
	if err := enqueueEvent(tx, "member.status_updated", "member", int64(memberID)); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *memberRepository) UpdateMemberName(memberID int, name string) error {
	tx, err := config.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec("UPDATE members SET Name = ? WHERE ID = ?", name, memberID); err != nil {
		return err
	}
	if err := enqueueEvent(tx, "member.profile_updated", "member", int64(memberID)); err != nil {
		return err
	}
	return tx.Commit()
}

// SetMemberPassword update password + role member yang sudah ada
// Dipakai saat register — member sudah ada (nama dikenal), hanya set password
func (r *memberRepository) SetMemberPassword(memberID int, hashedPassword string, role string, division string, batchYear string) error {
	_, err := config.DB.Exec(
		"UPDATE members SET Password = ?, Role = ?, Division = ?, BatchYear = ? WHERE ID = ?",
		hashedPassword, role, nullableString(division), nullableString(batchYear), memberID,
	)
	return err
}

func (r *memberRepository) UpdateMemberPassword(memberID int, hashedPassword string) error {
	_, err := config.DB.Exec(
		"UPDATE members SET Password = ? WHERE ID = ?",
		hashedPassword, memberID,
	)
	return err
}

func (r *memberRepository) ApproveMember(memberID int, role string, division string, batchYear string, canHandle bool, approvedBy int) error {
	_, err := config.DB.Exec(`
		UPDATE members
		SET Role = ?,
		    Division = ?,
		    BatchYear = ?,
		    CanHandleWorkOrder = ?,
		    AccountStatus = 'active',
		    MembershipStatus = 'active',
		    GraduationYear = NULL,
		    ApprovedAt = NOW(),
		    ApprovedBy = ?
		WHERE ID = ?
	`, role, nullableString(division), nullableString(batchYear), canHandle, approvedBy, memberID)
	return err
}

func (r *memberRepository) RejectMember(memberID int) error {
	_, err := config.DB.Exec(`
		UPDATE members
		SET AccountStatus = 'rejected',
		    CanHandleWorkOrder = 0,
		    Status = 'offduty'
		WHERE ID = ?
	`, memberID)
	return err
}

func (r *memberRepository) DisableMember(memberID int) error {
	_, err := config.DB.Exec(`
		UPDATE members
		SET AccountStatus = 'disabled',
		    CanHandleWorkOrder = 0,
		    Status = 'offduty'
		WHERE ID = ?
	`, memberID)
	return err
}

func (r *memberRepository) MarkMemberAsAlumni(memberID int, graduationYear int) error {
	_, err := config.DB.Exec(`
		UPDATE members
		SET MembershipStatus = 'alumni',
		    AccountStatus = 'disabled',
		    CanHandleWorkOrder = 0,
		    GraduationYear = ?,
		    Status = 'offduty'
		WHERE ID = ?
	`, graduationYear, memberID)
	return err
}

func (r *memberRepository) ChangeRole(memberID int, role string) error {
	_, err := config.DB.Exec(`
		UPDATE members
		SET Role = ?
		WHERE ID = ?
		  AND AccountStatus = 'active'
		  AND MembershipStatus = 'active'
	`, role, memberID)
	return err
}

func (r *memberRepository) GraduateBatch(batchYear string, graduationYear int) (int64, error) {
	result, err := config.DB.Exec(`
		UPDATE members
		SET MembershipStatus = 'alumni',
		    AccountStatus = 'disabled',
		    CanHandleWorkOrder = 0,
		    GraduationYear = ?,
		    Status = 'offduty'
		WHERE BatchYear = ?
		  AND MembershipStatus = 'active'
		  AND COALESCE(Role, '') <> 'Guest'
	`, graduationYear, batchYear)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

type memberScanner interface {
	Scan(dest ...any) error
}

func scanMember(rows memberScanner, m *models.Member) error {
	var graduationYear sql.NullInt64
	var registeredAt sql.NullString
	var approvedAt sql.NullString
	var approvedBy sql.NullInt64

	if err := rows.Scan(
		&m.ID,
		&m.Name,
		&m.Role,
		&m.Division,
		&m.Status,
		&m.Avatar,
		&m.AccountStatus,
		&m.MembershipStatus,
		&m.BatchYear,
		&graduationYear,
		&m.CanHandleWorkOrder,
		&registeredAt,
		&approvedAt,
		&approvedBy,
	); err != nil {
		return err
	}

	applyNullableMemberFields(m, graduationYear, registeredAt, approvedAt, approvedBy)
	return nil
}

func scanMemberWithPassword(rows memberScanner, m *models.Member) error {
	var graduationYear sql.NullInt64
	var registeredAt sql.NullString
	var approvedAt sql.NullString
	var approvedBy sql.NullInt64

	if err := rows.Scan(
		&m.ID,
		&m.Name,
		&m.Password,
		&m.Role,
		&m.Division,
		&m.Status,
		&m.Avatar,
		&m.AccountStatus,
		&m.MembershipStatus,
		&m.BatchYear,
		&graduationYear,
		&m.CanHandleWorkOrder,
		&registeredAt,
		&approvedAt,
		&approvedBy,
	); err != nil {
		return err
	}

	applyNullableMemberFields(m, graduationYear, registeredAt, approvedAt, approvedBy)
	return nil
}

func applyNullableMemberFields(m *models.Member, graduationYear sql.NullInt64, registeredAt sql.NullString, approvedAt sql.NullString, approvedBy sql.NullInt64) {
	m.Name = utils.SanitizeText(m.Name)
	m.Role = utils.SanitizeText(m.Role)
	m.Division = utils.SanitizeText(m.Division)
	m.Status = utils.SanitizeText(m.Status)
	m.Avatar = utils.SanitizeText(m.Avatar)
	m.AccountStatus = utils.SanitizeText(m.AccountStatus)
	m.MembershipStatus = utils.SanitizeText(m.MembershipStatus)
	m.BatchYear = utils.SanitizeText(m.BatchYear)
	if graduationYear.Valid {
		v := int(graduationYear.Int64)
		m.GraduationYear = &v
	}
	if registeredAt.Valid {
		m.RegisteredAt = registeredAt.String
	}
	if approvedAt.Valid {
		m.ApprovedAt = approvedAt.String
	}
	if approvedBy.Valid {
		v := int(approvedBy.Int64)
		m.ApprovedBy = &v
	}
}

func nullableString(value string) any {
	if value == "" {
		return nil
	}
	return value
}

// UpdateMemberAvatar update nama file avatar member
func (r *memberRepository) UpdateMemberAvatar(memberID int, avatarFilename string) error {
	_, err := config.DB.Exec(
		"UPDATE members SET Avatar = ? WHERE ID = ?",
		avatarFilename, memberID,
	)
	if err != nil {
		return fmt.Errorf("failed to update avatar: %w", err)
	}
	return nil
}
