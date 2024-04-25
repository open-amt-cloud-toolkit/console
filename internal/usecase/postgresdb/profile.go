package postgresdb

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	"github.com/open-amt-cloud-toolkit/console/pkg/postgres"
)

// ProfileRepo -.

type ProfileRepo struct {
	*postgres.DB
}

// New -.

func NewProfileRepo(pg *postgres.DB) *ProfileRepo {
	return &ProfileRepo{pg}
}

// GetCount -.

func (r *ProfileRepo) GetCount(ctx context.Context, tenantID string) (int, error) {
	sqlQuery, _, err := r.Builder.
		Select("COUNT(*) OVER() AS total_count").
		From("profiles").
		Where("tenant_id = ?", tenantID).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("ProfileRepo - GetCount - r.Builder: %w", err)
	}

	var count int

	err = r.Pool.QueryRow(ctx, sqlQuery, tenantID).Scan(&count)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}

		return 0, fmt.Errorf("ProfileRepo - GetCount - r.Pool.QueryRow: %w", err)
	}

	return count, nil
}

// Get -.

func (r *ProfileRepo) Get(ctx context.Context, top, skip int, tenantID string) ([]entity.Profile, error) {
	if top == 0 {
		top = 100
	}

	sqlQuery, _, err := r.Builder.
		Select(`profile_name,
            activation,
            amt_password,
            generate_random_password,
            cira_config_name,
            mebx_password,
            generate_random_mebx_password,
            tags,
            dhcp_enabled,
            tenant_id,
            tls_mode,
            user_consent,
            ider_enabled,
            kvm_enabled,
            sol_enabled,
            tls_signing_authority,
            ip_sync_enabled,
            local_wifi_sync_enabled,
            ieee8021x_profile_name,
            CAST(xmin as text) as xmin
				`).
		From("profiles").
		Where("tenant_id = ?", tenantID).
		OrderBy("profile_name").
		Limit(uint64(top)).
		Offset(uint64(skip)).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("ProfileRepo - Get - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sqlQuery, tenantID)
	if err != nil {
		return nil, fmt.Errorf("ProfileRepo - Get - r.Pool.Query: %w", err)
	}

	defer rows.Close()

	profiles := make([]entity.Profile, 0)

	for rows.Next() {
		p := entity.Profile{}

		err = rows.Scan(&p.ProfileName, &p.Activation, &p.AMTPassword, &p.GenerateRandomPassword,
			&p.CIRAConfigName, &p.MEBXPassword,
			&p.GenerateRandomMEBxPassword, &p.Tags, &p.DhcpEnabled, &p.TenantID, &p.TLSMode,
			&p.UserConsent, &p.IDEREnabled, &p.KVMEnabled, &p.SOLEnabled, &p.TLSSigningAuthority,
			&p.IPSyncEnabled, &p.LocalWifiSyncEnabled, &p.Ieee8021xProfileName, &p.Version)
		if err != nil {
			return nil, fmt.Errorf("ProfileRepo - Get - rows.Scan: %w", err)
		}

		profiles = append(profiles, p)
	}

	return profiles, nil
}

// GetByName -.

func (r *ProfileRepo) GetByName(ctx context.Context, profileName, tenantID string) (entity.Profile, error) {
	sqlQuery, _, err := r.Builder.
		Select(`profile_name,
            activation,
            amt_password,
            generate_random_password,
            cira_config_name,
            mebx_password,
            generate_random_mebx_password,
            tags,
            dhcp_enabled,
            tenant_id,
            tls_mode,
            user_consent,
            ider_enabled,
            kvm_enabled,
            sol_enabled,
            tls_signing_authority,
            ip_sync_enabled,
            local_wifi_sync_enabled,
            ieee8021x_profile_name,
            CAST(xmin as text) as xmin
				`).
		From("profiles").
		Where("profile_name = ? and tenant_id = ?", profileName, tenantID).
		ToSql()
	if err != nil {
		return entity.Profile{}, fmt.Errorf("ProfileRepo - GetByName - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sqlQuery, profileName, tenantID)
	if err != nil {
		return entity.Profile{}, fmt.Errorf("ProfileRepo - GetByName - r.Pool.Query: %w", err)
	}

	defer rows.Close()

	profiles := make([]entity.Profile, 0)

	for rows.Next() {
		p := entity.Profile{}

		err = rows.Scan(&p.ProfileName, &p.Activation, &p.AMTPassword, &p.GenerateRandomPassword,
			&p.CIRAConfigName, &p.MEBXPassword,
			&p.GenerateRandomMEBxPassword, &p.Tags, &p.DhcpEnabled, &p.TenantID, &p.TLSMode,
			&p.UserConsent, &p.IDEREnabled, &p.KVMEnabled, &p.SOLEnabled, &p.TLSSigningAuthority,
			&p.IPSyncEnabled, &p.LocalWifiSyncEnabled, &p.Ieee8021xProfileName, &p.Version)
		if err != nil {
			return p, fmt.Errorf("ProfileRepo - GetByName - rows.Scan: %w", err)
		}

		profiles = append(profiles, p)
	}

	if len(profiles) == 0 {
		return entity.Profile{}, fmt.Errorf("ProfileRepo - GetByName - Not Found: %w", err)
	}

	return profiles[0], nil
}

// Delete -.

func (r *ProfileRepo) Delete(ctx context.Context, profileName, tenantID string) (bool, error) {
	sqlQuery, args, err := r.Builder.
		Delete("profiles").
		Where("profile_name = ? AND tenant_id = ?", profileName, tenantID).
		ToSql()
	if err != nil {
		return false, fmt.Errorf("ProfileRepo - Delete - r.Builder: %w", err)
	}

	res, err := r.Pool.Exec(ctx, sqlQuery, args...)
	if err != nil {
		return false, fmt.Errorf("ProfileRepo - Delete - r.Pool.Exec: %w", err)
	}

	return res.RowsAffected() > 0, nil
}

// Update -.

func (r *ProfileRepo) Update(ctx context.Context, p *entity.Profile) (bool, error) {
	sqlQuery, args, err := r.Builder.
		Update("profiles").
		Set("activation", p.Activation).
		Set("amt_password", p.AMTPassword).
		Set("generate_random_password", p.GenerateRandomPassword).
		Set("cira_config_name", p.CIRAConfigName).
		Set("mebx_password", p.MEBXPassword).
		Set("generate_random_mebx_password", p.GenerateRandomMEBxPassword).
		Set("tags", p.Tags).
		Set("dhcp_enabled", p.DhcpEnabled).
		Set("tls_mode", p.TLSMode).
		Set("user_consent", p.UserConsent).
		Set("ider_enabled", p.IDEREnabled).
		Set("kvm_enabled", p.KVMEnabled).
		Set("sol_enabled", p.SOLEnabled).
		Set("tls_signing_authority", p.TLSSigningAuthority).
		Set("ieee8021x_profile_name", p.Ieee8021xProfileName).
		Set("ip_sync_enabled", p.IPSyncEnabled).
		Set("local_wifi_sync_enabled", p.LocalWifiSyncEnabled).
		Where("name = ? AND tenant_id = ?", p.ProfileName, p.TenantID).
		Suffix("AND xmin::text = ?", p.Version).
		ToSql()
	if err != nil {
		return false, fmt.Errorf("ProfileRepo - Update - r.Builder: %w", err)
	}

	res, err := r.Pool.Exec(ctx, sqlQuery, args...)
	if err != nil {
		return false, fmt.Errorf("ProfileRepo - Update - r.Pool.Exec: %w", err)
	}

	return res.RowsAffected() > 0, nil
}

// Insert -.

func (r *ProfileRepo) Insert(ctx context.Context, p *entity.Profile) (string, error) {
	ciraConfigName := p.CIRAConfigName

	ieee8021xProfileName := p.Ieee8021xProfileName

	if *p.CIRAConfigName == "" {
		ciraConfigName = nil
	}

	if *p.Ieee8021xProfileName == "" {
		ieee8021xProfileName = nil
	}

	sqlQuery, args, err := r.Builder.
		Insert("profiles").
		Columns("profile_name", "activation", "amt_password", "generate_random_password", "cira_config_name", "mebx_password", "generate_random_mebx_password", "tags", "dhcp_enabled", "tls_mode", "user_consent", "ider_enabled", "kvm_enabled", "sol_enabled", "tls_signing_authority", "ieee8021x_profile_name", "ip_sync_enabled", "local_wifi_sync_enabled", "tenant_id").
		Values(p.ProfileName, p.Activation, p.AMTPassword, p.GenerateRandomPassword, ciraConfigName, p.MEBXPassword, p.GenerateRandomMEBxPassword, p.Tags, p.DhcpEnabled, p.TLSMode, p.UserConsent, p.IDEREnabled, p.KVMEnabled, p.SOLEnabled, p.TLSSigningAuthority, ieee8021xProfileName, p.IPSyncEnabled, p.LocalWifiSyncEnabled, p.TenantID).
		Suffix("RETURNING xmin::text").
		ToSql()
	if err != nil {
		return "", fmt.Errorf("ProfileRepo - Insert - r.Builder: %w", err)
	}

	var version string

	err = r.Pool.QueryRow(ctx, sqlQuery, args...).Scan(&version)
	if err != nil {
		return "", fmt.Errorf("ProfileRepo - Insert - r.Pool.QueryRow: %w", err)
	}

	return version, nil
}
