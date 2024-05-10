package sqldb

import (
	"context"
	"database/sql"
	"errors"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	"github.com/open-amt-cloud-toolkit/console/pkg/consoleerrors"
	"github.com/open-amt-cloud-toolkit/console/pkg/db"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
)

// IEEE8021xRepo -.
type IEEE8021xRepo struct {
	*db.SQL
	log logger.Interface
}

var (
	ErrIEEE8021xDatabase  = DatabaseError{Console: consoleerrors.CreateConsoleError("IEEE8021xRepo")}
	ErrIEEE8021xNotUnique = DatabaseError{Console: consoleerrors.CreateConsoleError("IEEE8021xRepo")}
)

// New -.
func NewIEEE8021xRepo(database *db.SQL, log logger.Interface) *IEEE8021xRepo {
	return &IEEE8021xRepo{database, log}
}

// CheckProfileExits -.
func (r *IEEE8021xRepo) CheckProfileExists(_ context.Context, profileName, tenantID string) (bool, error) {
	sqlQuery, _, err := r.Builder.
		Select("COUNT(*) OVER() AS total_count").
		From("ieee8021xconfigs").
		Where("profile_name and tenant_id = ?", profileName, tenantID).
		ToSql()
	if err != nil {
		return false, ErrIEEE8021xDatabase.Wrap("CheckProfileExists", "r.Builder: ", err)
	}

	var count int

	err = r.Pool.QueryRow(sqlQuery, tenantID).Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, ErrIEEE8021xDatabase.Wrap("CheckProfileExists", "r.Pool.QueryRow", err)
	}

	return true, nil
}

// GetCount -.
func (r *IEEE8021xRepo) GetCount(_ context.Context, tenantID string) (int, error) {
	sqlQuery, _, err := r.Builder.
		Select("COUNT(*) OVER() AS total_count").
		From("ieee8021xconfigs").
		Where("tenant_id = ?", tenantID).
		ToSql()
	if err != nil {
		return 0, ErrIEEE8021xDatabase.Wrap("GetCount", "r.Builder: ", err)
	}

	var count int

	err = r.Pool.QueryRow(sqlQuery, tenantID).Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}

		return 0, ErrIEEE8021xDatabase.Wrap("GetCount", "r.Pool.QueryRow", err)
	}

	return count, nil
}

// Get -.
func (r *IEEE8021xRepo) Get(_ context.Context, top, skip int, tenantID string) ([]entity.IEEE8021xConfig, error) {
	if top == 0 {
		top = 100
	}

	sqlQuery, _, err := r.Builder.
		Select("profile_name",
			"auth_Protocol",
			"pxe_timeout",
			"wired_interface",
			"tenant_id",
		).
		From("ieee8021xconfigs").
		Where("tenant_id = ?", tenantID).
		Limit(uint64(top)).
		Offset(uint64(skip)).
		ToSql()
	if err != nil {
		return nil, ErrIEEE8021xDatabase.Wrap("Get", "r.Builder: ", err)
	}

	rows, err := r.Pool.Query(sqlQuery, tenantID)
	if err != nil {
		return nil, ErrIEEE8021xDatabase.Wrap("Get", "r.Pool.Query", err)
	}

	defer rows.Close()

	if rows.Err() != nil {
		return nil, ErrDeviceDatabase.Wrap("Get", "rows.Err", rows.Err())
	}

	ieee8021xConfigs := make([]entity.IEEE8021xConfig, 0)

	for rows.Next() {
		p := entity.IEEE8021xConfig{}

		err = rows.Scan(&p.ProfileName, &p.AuthenticationProtocol, &p.PXETimeout, &p.WiredInterface, &p.TenantID)
		if err != nil {
			return nil, ErrIEEE8021xDatabase.Wrap("Get", "rows.Scan: ", err)
		}

		ieee8021xConfigs = append(ieee8021xConfigs, p)
	}

	return ieee8021xConfigs, nil
}

// GetByName -.
func (r *IEEE8021xRepo) GetByName(_ context.Context, profileName, tenantID string) (*entity.IEEE8021xConfig, error) {
	sqlQuery, _, err := r.Builder.
		Select("profile_name",
			"auth_Protocol",
			"pxe_timeout",
			"wired_interface",
			"tenant_id",
		).
		From("ieee8021xconfigs").
		Where("profile_name = ? and tenant_id = ?", profileName, tenantID).
		ToSql()
	if err != nil {
		return nil, ErrIEEE8021xDatabase.Wrap("Get", "r.Builder: ", err)
	}

	rows, err := r.Pool.Query(sqlQuery, profileName, tenantID)
	if err != nil {
		return nil, ErrIEEE8021xDatabase.Wrap("Get", "r.Pool.Query", err)
	}

	defer rows.Close()

	if rows.Err() != nil {
		return nil, ErrDeviceDatabase.Wrap("Get", "rows.Err", rows.Err())
	}

	ieee8021xConfigs := make([]*entity.IEEE8021xConfig, 0)

	for rows.Next() {
		p := &entity.IEEE8021xConfig{}

		err = rows.Scan(&p.ProfileName, &p.AuthenticationProtocol, &p.PXETimeout, &p.WiredInterface, &p.TenantID)
		if err != nil {
			return p, ErrIEEE8021xDatabase.Wrap("Get", "rows.Scan: ", err)
		}

		ieee8021xConfigs = append(ieee8021xConfigs, p)
	}

	if len(ieee8021xConfigs) == 0 {
		return nil, nil
	}

	return ieee8021xConfigs[0], nil
}

// Delete -.
func (r *IEEE8021xRepo) Delete(_ context.Context, profileName, tenantID string) (bool, error) {
	sqlQuery, args, err := r.Builder.
		Delete("ieee8021xconfigs").
		Where("profile_name = ? AND tenant_id = ?", profileName, tenantID).
		ToSql()
	if err != nil {
		return false, ErrIEEE8021xDatabase.Wrap("Delete", "r.Builder: ", err)
	}

	res, err := r.Pool.Exec(sqlQuery, args...)
	if err != nil {
		return false, ErrIEEE8021xDatabase.Wrap("Delete", "r.Pool.Exec", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, ErrIEEE8021xDatabase.Wrap("Delete", "res.RowsAffected", err)
	}

	return rowsAffected > 0, nil
}

// Update -.
func (r *IEEE8021xRepo) Update(_ context.Context, p *entity.IEEE8021xConfig) (bool, error) {
	sqlQuery, args, err := r.Builder.
		Update("ieee8021xconfigs").
		Set("auth_protocol", p.AuthenticationProtocol).
		Set("servername", p.ServerName).
		Set("domain", p.Domain).
		Set("username", p.Username).
		Set("password", p.Password).
		Set("roaming_identity", p.RoamingIdentity).
		Set("active_in_s0", p.ActiveInS0).
		Set("pxe_timeout", p.PXETimeout).
		Set("wired_interface", p.WiredInterface).
		Where("profile_name = ? AND tenant_id = ?", p.ProfileName, p.TenantID).
		ToSql()
	if err != nil {
		return false, ErrIEEE8021xDatabase.Wrap("Update", "r.Builder: ", err)
	}

	res, err := r.Pool.Exec(sqlQuery, args...)
	if err != nil {
		return false, ErrIEEE8021xDatabase.Wrap("Update", "r.Pool.Exec", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, ErrIEEE8021xDatabase.Wrap("Update", "res.RowsAffected", err)
	}

	return rowsAffected > 0, nil
}

// Insert -.
func (r *IEEE8021xRepo) Insert(_ context.Context, p *entity.IEEE8021xConfig) (string, error) {
	sqlQuery, args, err := r.Builder.
		Insert("ieee8021xconfigs").
		Columns("profile_name", "auth_protocol", "pxe_timeout", "wired_interface", "tenant_id").
		Values(p.ProfileName, p.AuthenticationProtocol, p.PXETimeout, p.WiredInterface, p.TenantID).
		ToSql()
	if err != nil {
		return "", ErrIEEE8021xDatabase.Wrap("Insert", "r.Builder: ", err)
	}

	version := ""

	if r.IsEmbedded {
		_, err = r.Pool.Exec(sqlQuery, args...)
	} else {
		err = r.Pool.QueryRow(sqlQuery, args...).Scan(&version)
	}

	if err != nil {
		if db.CheckNotUnique(err) {
			return "", ErrIEEE8021xNotUnique
		}

		return "", ErrIEEE8021xDatabase.Wrap("Insert", "r.Pool.QueryRow", err)
	}

	return version, nil
}
