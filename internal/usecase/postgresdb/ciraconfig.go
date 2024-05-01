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

// CIRARepo -.
type CIRARepo struct {
	*postgres.DB
	log logger.Interface
}

// New -.
func NewCIRARepo(pg *postgres.DB, log logger.Interface) *CIRARepo {
	return &CIRARepo{pg, log}
}

// GetCount -.
func (r *CIRARepo) GetCount(ctx context.Context, tenantID string) (int, error) {
	sqlQuery, _, err := r.Builder.
		Select("COUNT(*) OVER() AS total_count").
		From("ciraconfigs").
		Where("tenant_id = ?", tenantID).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("CIRARepo - GetCount - r.Builder: %w", err)
	}

	var count int

	err = r.Pool.QueryRow(ctx, sqlQuery, tenantID).Scan(&count)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}

		return 0, fmt.Errorf("CIRARepo - GetCount - r.Pool.QueryRow: %w", err)
	}

	return count, nil
}

// Get -.
func (r *CIRARepo) Get(ctx context.Context, top, skip int, tenantID string) ([]entity.CIRAConfig, error) {
	if top == 0 {
		top = 100
	}

	sqlQuery, _, err := r.Builder.
		Select(`
			cira_config_name,
			mps_server_address,
			mps_port,
			user_name,
			password,
			common_name,
			server_address_format,
			auth_method,
			mps_root_certificate,
			proxydetails,
			tenant_id,
      		CAST(xmin as text) as xmin
		`).
		From("ciraconfigs").
		Where("tenant_id = ?", tenantID).
		OrderBy("cira_config_name").
		Limit(uint64(top)).
		Offset(uint64(skip)).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("CIRARepo - Get - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sqlQuery, tenantID)
	if err != nil {
		return nil, fmt.Errorf("CIRARepo - Get - r.Pool.Query: %w", err)
	}

	defer rows.Close()

	configs := make([]entity.CIRAConfig, 0)

	for rows.Next() {
		p := entity.CIRAConfig{}

		err = rows.Scan(&p.ConfigName, &p.MPSServerAddress, &p.MpsPort, &p.Username, &p.Password, &p.CommonName, &p.ServerAddressFormat, &p.AuthMethod, &p.MpsRootCertificate, &p.ProxyDetails, &p.TenantID, &p.Version)
		if err != nil {
			return nil, fmt.Errorf("CIRARepo - Get - rows.Scan: %w", err)
		}

		configs = append(configs, p)
	}

	return configs, nil
}

// GetByName -.
func (r *CIRARepo) GetByName(ctx context.Context, configName, tenantID string) (entity.CIRAConfig, error) {
	sqlQuery, _, err := r.Builder.
		Select(`
			cira_config_name,
			mps_server_address,
			mps_port,
			user_name,
			password,
			common_name,
			server_address_format,
			auth_method,
			mps_root_certificate,
			proxydetails,
			tenant_id,
      		CAST(xmin as text) as xmin
		`).
		From("ciraconfigs").
		Where("cira_config_name = ? and tenant_id = ?", configName, tenantID).
		ToSql()
	if err != nil {
		return entity.CIRAConfig{}, fmt.Errorf("CIRARepo - GetByName - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sqlQuery, configName, tenantID)
	if err != nil {
		return entity.CIRAConfig{}, fmt.Errorf("CIRARepo - GetByName - r.Pool.Query: %w", err)
	}

	defer rows.Close()

	configs := make([]entity.CIRAConfig, 0)

	for rows.Next() {
		p := entity.CIRAConfig{}

		err = rows.Scan(&p.ConfigName, &p.MPSServerAddress, &p.MpsPort, &p.Username, &p.Password, &p.CommonName, &p.ServerAddressFormat, &p.AuthMethod, &p.MpsRootCertificate, &p.ProxyDetails, &p.TenantID, &p.Version)
		if err != nil {
			return p, fmt.Errorf("CIRARepo - GetByName - rows.Scan: %w", err)
		}

		configs = append(configs, p)
	}

	if len(configs) == 0 {
		return entity.CIRAConfig{}, fmt.Errorf("CIRARepo - GetByName - NotFound: %w", err)
	}

	return configs[0], nil
}

// Delete -.
func (r *CIRARepo) Delete(ctx context.Context, configName, tenantID string) (bool, error) {
	sqlQuery, args, err := r.Builder.
		Delete("ciraconfigs").
		Where("cira_config_name = ? AND tenant_id = ?", configName, tenantID).
		ToSql()
	if err != nil {
		return false, fmt.Errorf("CIRARepo - Delete - r.Builder: %w", err)
	}

	res, err := r.Pool.Exec(ctx, sqlQuery, args...)
	if err != nil {
		return false, fmt.Errorf("CIRARepo - Delete - r.Pool.Exec: %w", err)
	}

	return res.RowsAffected() > 0, nil
}

// Update -.
func (r *CIRARepo) Update(ctx context.Context, p *entity.CIRAConfig) (bool, error) {
	sqlQuery, args, err := r.Builder.
		Update("ciraconfigs").
		Set("mps_server_address", p.MPSServerAddress).
		Set("mps_port", p.MpsPort).
		Set("user_name", p.Username).
		Set("password", p.Password).
		Set("common_name", p.CommonName).
		Set("server_address_format", p.ServerAddressFormat).
		Set("auth_method", p.AuthMethod).
		Set("mps_root_certificate", p.MpsRootCertificate).
		Set("proxydetails", p.ProxyDetails).
		Where("cira_config_name = ? AND tenant_id = ?", p.ConfigName, p.TenantID).
		ToSql()
	if err != nil {
		return false, fmt.Errorf("CIRARepo - Update - r.Builder: %w", err)
	}

	res, err := r.Pool.Exec(ctx, sqlQuery, args...)
	if err != nil {
		return false, fmt.Errorf("CIRARepo - Update - r.Pool.Exec: %w", err)
	}

	return res.RowsAffected() > 0, nil
}

// Insert -.
func (r *CIRARepo) Insert(ctx context.Context, p *entity.CIRAConfig) (string, error) {
	sqlQuery, args, err := r.Builder.
		Insert("ciraconfigs").
		Columns("cira_config_name", "mps_server_address", "mps_port", "user_name", "password", "common_name", "server_address_format", "auth_method", "mps_root_certificate", "proxydetails", "tenant_id").
		Values(p.ConfigName, p.MPSServerAddress, p.MpsPort, p.Username, p.Password, p.CommonName, p.ServerAddressFormat, p.AuthMethod, p.MpsRootCertificate, p.ProxyDetails, p.TenantID).
		Suffix("RETURNING xmin::text").
		ToSql()
	if err != nil {
		return "", fmt.Errorf("CIRARepo - Insert - r.Builder: %w", err)
	}

	var version string

	err = r.Pool.QueryRow(ctx, sqlQuery, args...).Scan(&version)
	if err != nil {
		return "", fmt.Errorf("CIRARepo - Insert - r.Pool.QueryRow: %w", err)
	}

	return version, nil
}
