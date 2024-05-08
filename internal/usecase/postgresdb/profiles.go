package postgresdb

import (
	"context"
	"database/sql"
	"errors"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	"github.com/open-amt-cloud-toolkit/console/pkg/consoleerrors"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
	"github.com/open-amt-cloud-toolkit/console/pkg/postgres"
)

// ProfileRepo -.

type ProfileRepo struct {
	*postgres.DB
	log logger.Interface
}

var (
	ErrProfileDatabase  = consoleerrors.DatabaseError{Console: consoleerrors.CreateConsoleError("ProfileRepo")}
	ErrProfileNotUnique = consoleerrors.DatabaseError{Console: consoleerrors.CreateConsoleError("ProfileRepo")}
)

// New -.

func NewProfileRepo(pg *postgres.DB, log logger.Interface) *ProfileRepo {
	return &ProfileRepo{pg, log}
}

// GetCount -.

func (r *ProfileRepo) GetCount(_ context.Context, tenantID string) (int, error) {
	sqlQuery, _, err := r.Builder.
		Select("COUNT(*) OVER() AS total_count").
		From("profiles").
		Where("tenant_id = ?", tenantID).
		ToSql()
	if err != nil {
		return 0, ErrProfileDatabase.Wrap("GetCount", "r.Builder", err)
	}

	var count int

	err = r.Pool.QueryRow(sqlQuery, tenantID).Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}

		return 0, ErrProfileDatabase.Wrap("GetCount", "r.Pool.QueryRow", err)
	}

	return count, nil
}

// Get -.

func (r *ProfileRepo) Get(_ context.Context, top, skip int, tenantID string) ([]entity.Profile, error) {
	if top == 0 {
		top = 100
	}

	sqlQuery, _, err := r.Builder.
		Select(
			"profile_name",
			"activation",
			"amt_password",
			"generate_random_password",
			"cira_config_name",
			"mebx_password",
			"generate_random_mebx_password",
			"tags",
			"dhcp_enabled",
			"tenant_id",
			"tls_mode",
			"user_consent",
			"ider_enabled",
			"kvm_enabled",
			"sol_enabled",
			"tls_signing_authority",
			"ip_sync_enabled",
			"local_wifi_sync_enabled",
			"ieee8021x_profile_name",
			"CAST(xmin as text) as xmin").
		From("profiles").
		Where("tenant_id = ?", tenantID).
		OrderBy("profile_name").
		Limit(uint64(top)).
		Offset(uint64(skip)).
		ToSql()
	if err != nil {
		return nil, ErrProfileDatabase.Wrap("Get", "r.Builder", err)
	}

	rows, err := r.Pool.Query(sqlQuery, tenantID)
	if err != nil {
		return nil, ErrProfileDatabase.Wrap("Get", "r.Pool.Query", err)
	}

	defer rows.Close()

	if rows.Err() != nil {
		return nil, ErrDeviceDatabase.Wrap("Get", "rows.Err", rows.Err())
	}

	profiles := make([]entity.Profile, 0)

	for rows.Next() {
		p := entity.Profile{}

		err = rows.Scan(&p.ProfileName, &p.Activation, &p.AMTPassword, &p.GenerateRandomPassword,
			&p.CIRAConfigName, &p.MEBXPassword,
			&p.GenerateRandomMEBxPassword, &p.Tags, &p.DHCPEnabled, &p.TenantID, &p.TLSMode,
			&p.UserConsent, &p.IDEREnabled, &p.KVMEnabled, &p.SOLEnabled, &p.TLSSigningAuthority,
			&p.IPSyncEnabled, &p.LocalWiFiSyncEnabled, &p.IEEE8021xProfileName, &p.Version)
		if err != nil {
			return nil, ErrProfileDatabase.Wrap("Get", "rows.Scan", err)
		}

		profiles = append(profiles, p)
	}

	return profiles, nil
}

// GetByName -.

func (r *ProfileRepo) GetByName(_ context.Context, profileName, tenantID string) (*entity.Profile, error) {
	sqlQuery, _, err := r.Builder.
		Select("profile_name",
			"activation",
			"amt_password",
			"generate_random_password",
			"cira_config_name",
			"mebx_password",
			"generate_random_mebx_password",
			"tags",
			"dhcp_enabled",
			"tenant_id",
			"tls_mode",
			"user_consent",
			"ider_enabled",
			"kvm_enabled",
			"sol_enabled",
			"tls_signing_authority",
			"ip_sync_enabled",
			"local_wifi_sync_enabled",
			"ieee8021x_profile_name",
			"CAST(xmin as text) as xmin").
		From("profiles").
		Where("profile_name = ? and tenant_id = ?", profileName, tenantID).
		ToSql()
	if err != nil {
		return nil, ErrProfileDatabase.Wrap("GetByName", "r.Builder", err)
	}

	rows, err := r.Pool.Query(sqlQuery, profileName, tenantID)
	if err != nil {
		return nil, ErrProfileDatabase.Wrap("GetByName", "r.Pool.Query", err)
	}

	defer rows.Close()

	if rows.Err() != nil {
		return nil, ErrDeviceDatabase.Wrap("Get", "rows.Err", rows.Err())
	}

	profiles := make([]*entity.Profile, 0)

	for rows.Next() {
		p := &entity.Profile{}

		err = rows.Scan(&p.ProfileName, &p.Activation, &p.AMTPassword, &p.GenerateRandomPassword,
			&p.CIRAConfigName, &p.MEBXPassword,
			&p.GenerateRandomMEBxPassword, &p.Tags, &p.DHCPEnabled, &p.TenantID, &p.TLSMode,
			&p.UserConsent, &p.IDEREnabled, &p.KVMEnabled, &p.SOLEnabled, &p.TLSSigningAuthority,
			&p.IPSyncEnabled, &p.LocalWiFiSyncEnabled, &p.IEEE8021xProfileName, &p.Version)
		if err != nil {
			return p, ErrProfileDatabase.Wrap("GetByName", "rows.Scan", err)
		}

		profiles = append(profiles, p)
	}

	if len(profiles) == 0 {
		return nil, ErrProfileDatabase.Wrap("GetByName", "Not Found", err)
	}

	return profiles[0], nil
}

// Delete -.

func (r *ProfileRepo) Delete(_ context.Context, profileName, tenantID string) (bool, error) {
	sqlQuery, args, err := r.Builder.
		Delete("profiles").
		Where("profile_name = ? AND tenant_id = ?", profileName, tenantID).
		ToSql()
	if err != nil {
		return false, ErrProfileDatabase.Wrap("Delete", "r.Builder", err)
	}

	res, err := r.Pool.Exec(sqlQuery, args...)
	if err != nil {
		return false, ErrProfileDatabase.Wrap("Delete", "r.Pool.Exec", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, ErrProfileDatabase.Wrap("Delete", "res.RowsAffected", err)
	}

	return rowsAffected > 0, nil
}

// Update -.

func (r *ProfileRepo) Update(_ context.Context, p *entity.Profile) (bool, error) {
	sqlQuery, args, err := r.Builder.
		Update("profiles").
		Set("activation", p.Activation).
		Set("amt_password", p.AMTPassword).
		Set("generate_random_password", p.GenerateRandomPassword).
		Set("cira_config_name", p.CIRAConfigName).
		Set("mebx_password", p.MEBXPassword).
		Set("generate_random_mebx_password", p.GenerateRandomMEBxPassword).
		Set("tags", p.Tags).
		Set("dhcp_enabled", p.DHCPEnabled).
		Set("tls_mode", p.TLSMode).
		Set("user_consent", p.UserConsent).
		Set("ider_enabled", p.IDEREnabled).
		Set("kvm_enabled", p.KVMEnabled).
		Set("sol_enabled", p.SOLEnabled).
		Set("tls_signing_authority", p.TLSSigningAuthority).
		Set("ieee8021x_profile_name", p.IEEE8021xProfileName).
		Set("ip_sync_enabled", p.IPSyncEnabled).
		Set("local_wifi_sync_enabled", p.LocalWiFiSyncEnabled).
		Where("profile_name = ? AND tenant_id = ?", p.ProfileName, p.TenantID).
		Suffix("AND xmin::text = ?", p.Version).
		ToSql()
	if err != nil {
		return false, ErrProfileDatabase.Wrap("Update", "r.Builder", err)
	}

	res, err := r.Pool.Exec(sqlQuery, args...)
	if err != nil {
		return false, ErrProfileDatabase.Wrap("Update", "r.Pool.Exec", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, ErrProfileDatabase.Wrap("Update", "res.RowsAffected", err)
	}

	return rowsAffected > 0, nil
}

// Insert -.

func (r *ProfileRepo) Insert(_ context.Context, p *entity.Profile) (string, error) {
	ciraConfigName := p.CIRAConfigName

	ieee8021xProfileName := p.IEEE8021xProfileName

	if ciraConfigName != nil {
		if *p.CIRAConfigName == "" {
			ciraConfigName = nil
		}
	}

	if ieee8021xProfileName != nil {
		if *p.IEEE8021xProfileName == "" {
			ieee8021xProfileName = nil
		}
	}

	sqlQuery, args, err := r.Builder.
		Insert("profiles").
		Columns("profile_name", "activation", "amt_password", "generate_random_password", "cira_config_name", "mebx_password", "generate_random_mebx_password", "tags", "dhcp_enabled", "tls_mode", "user_consent", "ider_enabled", "kvm_enabled", "sol_enabled", "tls_signing_authority", "ieee8021x_profile_name", "ip_sync_enabled", "local_wifi_sync_enabled", "tenant_id").
		Values(p.ProfileName, p.Activation, p.AMTPassword, p.GenerateRandomPassword, ciraConfigName, p.MEBXPassword, p.GenerateRandomMEBxPassword, p.Tags, p.DHCPEnabled, p.TLSMode, p.UserConsent, p.IDEREnabled, p.KVMEnabled, p.SOLEnabled, p.TLSSigningAuthority, ieee8021xProfileName, p.IPSyncEnabled, p.LocalWiFiSyncEnabled, p.TenantID).
		Suffix("RETURNING xmin::text").
		ToSql()
	if err != nil {
		return "", ErrProfileDatabase.Wrap("Insert", "r.Builder", err)
	}

	var version string

	err = r.Pool.QueryRow(sqlQuery, args...).Scan(&version)
	if err != nil {
		if postgres.CheckNotUnique(err) {
			return "", ErrProfileNotUnique
		}

		return "", ErrProfileDatabase.Wrap("Insert", "r.Pool.QueryRow", err)
	}

	return version, nil
}
