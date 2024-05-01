package postgresdb

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
	"github.com/open-amt-cloud-toolkit/console/pkg/postgres"
)

// IEEE8021xRepo -.
type IEEE8021xRepo struct {
	*postgres.DB
	log logger.Interface
}

// New -.
func NewIEEE8021xRepo(pg *postgres.DB, log logger.Interface) *IEEE8021xRepo {
	return &IEEE8021xRepo{pg, log}
}

// CheckProfileExits -.
func (r *IEEE8021xRepo) CheckProfileExists(ctx context.Context, profileName, tenantID string) (bool, error) {
	sqlQuery, _, err := r.Builder.
		Select("COUNT(*) OVER() AS total_count").
		From("ieee8021xconfigs").
		Where("profile_name and tenant_id = ?", profileName, tenantID).
		ToSql()
	if err != nil {
		return false, fmt.Errorf("IEEE8021xRepo - CheckProfileExists - r.Builder: %w", err)
	}

	var count int

	err = r.Pool.QueryRow(ctx, sqlQuery, tenantID).Scan(&count)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}

		return false, fmt.Errorf("IEEE8021xRepo - CheckProfileExists - r.Pool.QueryRow: %w", err)
	}

	return true, nil
}

// GetCount -.
func (r *IEEE8021xRepo) GetCount(ctx context.Context, tenantID string) (int, error) {
	sqlQuery, _, err := r.Builder.
		Select("COUNT(*) OVER() AS total_count").
		From("ieee8021xconfigs").
		Where("tenant_id = ?", tenantID).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("IEEE8021xRepo - GetCount - r.Builder: %w", err)
	}

	var count int

	err = r.Pool.QueryRow(ctx, sqlQuery, tenantID).Scan(&count)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}

		return 0, fmt.Errorf("IEEE8021xRepo - GetCount - r.Pool.QueryRow: %w", err)
	}

	return count, nil
}

// Get -.
func (r *IEEE8021xRepo) Get(ctx context.Context, top, skip int, tenantID string) ([]entity.IEEE8021xConfig, error) {
	if top == 0 {
		top = 100
	}

	sqlQuery, _, err := r.Builder.
		Select(`
			    profile_name,
        	auth_Protocol,
        	pxe_timeout,
       		wired_interface,
        	tenant_id,
          CAST(xmin as text) as xmin
			`).
		From("ieee8021xconfigs").
		Where("tenant_id = ?", tenantID).
		Limit(uint64(top)).
		Offset(uint64(skip)).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("IEEE8021xRepo - Get - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sqlQuery, tenantID)
	if err != nil {
		return nil, fmt.Errorf("IEEE8021xRepo - Get - r.Pool.Query: %w", err)
	}

	defer rows.Close()

	ieee8021xConfigs := make([]entity.IEEE8021xConfig, 0)

	for rows.Next() {
		p := entity.IEEE8021xConfig{}

		err = rows.Scan(&p.ProfileName, &p.AuthenticationProtocol, &p.PxeTimeout, &p.WiredInterface, &p.TenantID, &p.Version)
		if err != nil {
			return nil, fmt.Errorf("IEEE8021xRepo - Get - rows.Scan: %w", err)
		}

		ieee8021xConfigs = append(ieee8021xConfigs, p)
	}

	return ieee8021xConfigs, nil
}

// GetByName -.
func (r *IEEE8021xRepo) GetByName(ctx context.Context, profileName, tenantID string) (entity.IEEE8021xConfig, error) {
	sqlQuery, _, err := r.Builder.
		Select(`
			    profile_name,
        	auth_Protocol,
        	pxe_timeout,
        	wired_interface,
        	tenant_id,
          CAST(xmin as text) as xmin
			`).
		From("ieee8021xconfigs").
		Where("profile_name = ? and tenant_id = ?", profileName, tenantID).
		ToSql()
	if err != nil {
		return entity.IEEE8021xConfig{}, fmt.Errorf("IEEE8021xRepo - Get - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sqlQuery, profileName, tenantID)
	if err != nil {
		return entity.IEEE8021xConfig{}, fmt.Errorf("IEEE8021xRepo - Get - r.Pool.Query: %w", err)
	}

	defer rows.Close()

	ieee8021xConfigs := make([]entity.IEEE8021xConfig, 0)

	for rows.Next() {
		p := entity.IEEE8021xConfig{}

		err = rows.Scan(&p.ProfileName, &p.AuthenticationProtocol, &p.PxeTimeout, &p.WiredInterface, &p.TenantID, &p.Version)
		if err != nil {
			return p, fmt.Errorf("IEEE8021xRepo - Get - rows.Scan: %w", err)
		}

		ieee8021xConfigs = append(ieee8021xConfigs, p)
	}

	if len(ieee8021xConfigs) == 0 {
		return entity.IEEE8021xConfig{}, nil
	}

	return ieee8021xConfigs[0], nil
}

// Delete -.
func (r *IEEE8021xRepo) Delete(ctx context.Context, profileName, tenantID string) (bool, error) {
	sqlQuery, args, err := r.Builder.
		Delete("ieee8021xconfigs").
		Where("profile_name = ? AND tenant_id = ?", profileName, tenantID).
		ToSql()
	if err != nil {
		return false, fmt.Errorf("IEEE8021xRepo - Delete - r.Builder: %w", err)
	}

	res, err := r.Pool.Exec(ctx, sqlQuery, args...)
	if err != nil {
		return false, fmt.Errorf("IEEE8021xRepo - Delete - r.Pool.Exec: %w", err)
	}

	return res.RowsAffected() > 0, nil
}

// Update -.
func (r *IEEE8021xRepo) Update(ctx context.Context, p *entity.IEEE8021xConfig) (bool, error) {
	sqlQuery, args, err := r.Builder.
		Update("ieee8021xconfigs").
		Set("auth_protocol", p.AuthenticationProtocol).
		Set("servername", p.ServerName).
		Set("domain", p.Domain).
		Set("username", p.Username).
		Set("password", p.Password).
		Set("roaming_identity", p.RoamingIdentity).
		Set("active_in_s0", p.ActiveInS0).
		Set("pxe_timeout", p.PxeTimeout).
		Set("wired_interface", p.WiredInterface).
		Where("profile_name = ? AND tenant_id = ?", p.ProfileName, p.TenantID).
		ToSql()
	if err != nil {
		return false, fmt.Errorf("IEEE8021xRepo - Update - r.Builder: %w", err)
	}

	res, err := r.Pool.Exec(ctx, sqlQuery, args...)
	if err != nil {
		return false, fmt.Errorf("IEEE8021xRepo - Update - r.Pool.Exec: %w", err)
	}

	return res.RowsAffected() > 0, nil
}

// Insert -.
func (r *IEEE8021xRepo) Insert(ctx context.Context, p *entity.IEEE8021xConfig) (string, error) {
	sqlQuery, args, err := r.Builder.
		Insert("ieee8021xconfigs").
		Columns("profile_name", "auth_protocol", "pxe_timeout", "wired_interface", "tenant_id").
		Values(p.ProfileName, p.AuthenticationProtocol, p.PxeTimeout, p.WiredInterface, p.TenantID).
		Suffix("RETURNING xmin::text").
		ToSql()
	if err != nil {
		return "", fmt.Errorf("IEEE8021xRepo - Insert - r.Builder: %w", err)
	}

	var version string

	err = r.Pool.QueryRow(ctx, sqlQuery, args...).Scan(&version)
	if err != nil {
		return "", fmt.Errorf("IEEE8021xRepo - Insert - r.Pool.QueryRow: %w", err)
	}

	return version, nil
}
