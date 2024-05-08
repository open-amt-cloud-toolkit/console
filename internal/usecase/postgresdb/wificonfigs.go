package postgresdb

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	"github.com/open-amt-cloud-toolkit/console/pkg/consoleerrors"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
	"github.com/open-amt-cloud-toolkit/console/pkg/postgres"
)

// WirelessRepo -.
type WirelessRepo struct {
	*postgres.DB
	logger.Interface
}

var (
	ErrWiFiDatabase  = consoleerrors.DatabaseError{Console: consoleerrors.CreateConsoleError("WirelessRepo")}
	ErrWiFiNotUnique = consoleerrors.DatabaseError{Console: consoleerrors.CreateConsoleError("WirelessRepo")}
)

// New -.
func NewWirelessRepo(pg *postgres.DB, log logger.Interface) *WirelessRepo {
	return &WirelessRepo{pg, log}
}

// CheckProfileExits -.
func (r *WirelessRepo) CheckProfileExists(_ context.Context, profileName, tenantID string) (bool, error) {
	sqlQuery, _, err := r.Builder.
		Select("COUNT(*) OVER() AS total_count").
		From("wirelessconfigs").
		Where("wireless_profile_name and tenant_id = ?", profileName, tenantID).
		ToSql()
	if err != nil {
		return false, ErrWiFiDatabase.Wrap("CheckProfileExists", "r.Builder", err)
	}

	var count int

	err = r.Pool.QueryRow(sqlQuery, tenantID).Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, ErrWiFiDatabase.Wrap("CheckProfileExists", "r.Pool.QueryRow", err)
	}

	return true, nil
}

// GetCount -.
func (r *WirelessRepo) GetCount(_ context.Context, tenantID string) (int, error) {
	sqlQuery, _, err := r.Builder.
		Select("COUNT(*) OVER() AS total_count").
		From("wirelessconfigs").
		Where("tenant_id = ?", tenantID).
		ToSql()
	if err != nil {
		return 0, ErrWiFiDatabase.Wrap("GetCount", "r.Builder", err)
	}

	var count int

	err = r.Pool.QueryRow(sqlQuery, tenantID).Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}

		return 0, ErrWiFiDatabase.Wrap("GetCount", "r.Pool.QueryRow", err)
	}

	return count, nil
}

// Get -.
func (r *WirelessRepo) Get(_ context.Context, top, skip int, tenantID string) ([]entity.WirelessConfig, error) {
	if top == 0 {
		top = 100
	}

	sqlQuery, _, err := r.Builder.
		Select(
			"wireless_profile_name",
			"authentication_method",
			"encryption_method",
			"ssid",
			"psk_value",
			"psk_passphrase",
			"link_policy",
			"tenant_id",
			"ieee8021x_profile_name",
			"CAST(xmin as text) as xmin").
		From("wirelessconfigs").
		Where("tenant_id = ?", tenantID).
		OrderBy("wireless_profile_name").
		Limit(uint64(top)).
		Offset(uint64(skip)).
		ToSql()
	if err != nil {
		return nil, ErrWiFiDatabase.Wrap("Get", "r.Builder", err)
	}

	rows, err := r.Pool.Query(sqlQuery, tenantID)
	if err != nil {
		return nil, ErrWiFiDatabase.Wrap("Get", "r.Pool.Query", err)
	}

	defer rows.Close()

	if rows.Err() != nil {
		return nil, ErrDeviceDatabase.Wrap("Get", "rows.Err", rows.Err())
	}

	wirelessConfigs := make([]entity.WirelessConfig, 0)

	for rows.Next() {
		p := entity.WirelessConfig{}

		err = rows.Scan(&p.ProfileName, &p.AuthenticationMethod, &p.EncryptionMethod, &p.SSID, &p.PSKValue, &p.PSKPassphrase, &p.LinkPolicy, &p.TenantID, &p.IEEE8021xProfileName, &p.Version)
		if err != nil {
			return nil, ErrWiFiDatabase.Wrap("Get", "rows.Scan", err)
		}

		wirelessConfigs = append(wirelessConfigs, p)
	}

	return wirelessConfigs, nil
}

// GetByName -.
func (r *WirelessRepo) GetByName(_ context.Context, profileName, tenantID string) (*entity.WirelessConfig, error) {
	sqlQuery, _, err := r.Builder.
		Select(
			"wireless_profile_name",
			"authentication_method",
			"encryption_method",
			"ssid",
			"psk_value",
			// "psk_passphrase",
			"link_policy",
			"tenant_id",
			"ieee8021x_profile_name",
			"CAST(xmin as text) as xmin").
		From("wirelessconfigs").
		Where("wireless_profile_name = ? and tenant_id = ?", profileName, tenantID).
		ToSql()
	if err != nil {
		return nil, ErrWiFiDatabase.Wrap("GetByName", "r.Builder", err)
	}

	rows, err := r.Pool.Query(sqlQuery, profileName, tenantID)
	if err != nil {
		return nil, ErrWiFiDatabase.Wrap("GetByName", "r.Pool.Query", err)
	}

	defer rows.Close()

	if rows.Err() != nil {
		return nil, ErrDeviceDatabase.Wrap("Get", "rows.Err", rows.Err())
	}

	wirelessConfigs := make([]*entity.WirelessConfig, 0)

	for rows.Next() {
		p := &entity.WirelessConfig{}

		err = rows.Scan(&p.ProfileName, &p.AuthenticationMethod, &p.EncryptionMethod, &p.SSID, &p.PSKValue, &p.LinkPolicy, &p.TenantID, &p.IEEE8021xProfileName, &p.Version)
		if err != nil {
			return p, ErrWiFiDatabase.Wrap("GetByName", "rows.Scan", err)
		}

		wirelessConfigs = append(wirelessConfigs, p)
	}

	if len(wirelessConfigs) == 0 {
		return nil, nil
	}

	return wirelessConfigs[0], nil
}

// Delete -.
func (r *WirelessRepo) Delete(_ context.Context, profileName, tenantID string) (bool, error) {
	sqlQuery, args, err := r.Builder.
		Delete("wirelessconfigs").
		Where("wireless_profile_name = ? AND tenant_id = ?", profileName, tenantID).
		ToSql()
	if err != nil {
		return false, ErrWiFiDatabase.Wrap("Delete", "r.Builder", err)
	}

	res, err := r.Pool.Exec(sqlQuery, args...)
	if err != nil {
		return false, ErrWiFiDatabase.Wrap("Delete", "r.Pool.Exec", err)
	}

	result, err := res.RowsAffected()
	if err != nil {
		return false, ErrDomainDatabase.Wrap("Delete", "res.RowsAffected", err)
	}

	return result > 0, nil
}

// Update -.
func (r *WirelessRepo) Update(_ context.Context, p *entity.WirelessConfig) (bool, error) {
	sqlQuery, args, err := r.Builder.
		Update("wirelessconfigs").
		Set("authentication_method", p.AuthenticationMethod).
		Set("encryption_method", p.EncryptionMethod).
		Set("ssid", p.SSID).
		Set("psk_value", p.PSKValue).
		Set("psk_passphrase", p.PSKPassphrase).
		Set("link_policy", p.LinkPolicy).
		Set("ieee8021x_profile_name", p.IEEE8021xProfileName).
		Where("wireless_profile_name = ? AND tenant_id = ?", p.ProfileName, p.TenantID).
		Suffix("AND xmin::text = ?", p.Version).
		ToSql()
	if err != nil {
		return false, ErrWiFiDatabase.Wrap("Update", "r.Builder", err)
	}

	res, err := r.Pool.Exec(sqlQuery, args...)
	if err != nil {
		return false, ErrWiFiDatabase.Wrap("Update", "r.Pool.Exec", err)
	}

	result, err := res.RowsAffected()
	if err != nil {
		return false, ErrDomainDatabase.Wrap("Update", "res.RowsAffected", err)
	}

	return result > 0, nil
}

// Insert -.
func (r *WirelessRepo) Insert(_ context.Context, p *entity.WirelessConfig) (string, error) {
	date := time.Now().Format("2006-01-02 15:04:05")

	ieeeProfileName := p.IEEE8021xProfileName

	if p.IEEE8021xProfileName != nil {
		if *p.IEEE8021xProfileName == "" {
			ieeeProfileName = nil
		}
	}

	sqlQuery, args, err := r.Builder.
		Insert("wirelessconfigs").
		Columns("wireless_profile_name", "authentication_method", "encryption_method", "ssid", "psk_value", "psk_passphrase", "link_policy", "creation_date", "tenant_id", "ieee8021x_profile_name").
		Values(p.ProfileName, p.AuthenticationMethod, p.EncryptionMethod, p.SSID, p.PSKValue, p.PSKPassphrase, p.LinkPolicy, date, p.TenantID, ieeeProfileName).
		Suffix("RETURNING xmin::text").
		ToSql()
	if err != nil {
		return "", ErrWiFiDatabase.Wrap("Insert", "r.Builder", err)
	}

	var version string

	err = r.Pool.QueryRow(sqlQuery, args...).Scan(&version)
	if err != nil {
		if postgres.CheckNotUnique(err) {
			return "", ErrWiFiNotUnique
		}

		return "", ErrWiFiDatabase.Wrap("Insert", "r.Pool.QueryRow", err)
	}

	return version, nil
}
