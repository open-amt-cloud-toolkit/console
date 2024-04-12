package postgresdb

import (
	"context"
	"fmt"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	"github.com/open-amt-cloud-toolkit/console/pkg/postgres"
)

// CIRARepo -.
type CIRARepo struct {
	*postgres.DB
}

// New -.
func NewCIRARepo(pg *postgres.DB) *CIRARepo {
	return &CIRARepo{pg}
}

// GetCount -.
func (r *CIRARepo) GetCount(ctx context.Context, tenantID string) (int, error) {
	sql, _, err := r.Builder.
		Select("COUNT(*) OVER() AS total_count").
		From("ciraconfigs").
		Where("tenant_id = ?", tenantID).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("CIRARepo - GetCount - r.Builder: %w", err)
	}

	var count int

	err = r.Pool.QueryRow(ctx, sql, tenantID).Scan(&count)
	if err != nil {
		if err.Error() == "no rows in result set" {
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

	sql, _, err := r.Builder.
		Select(`
			cira_config_name as "configName", 
			mps_server_address as "mpsServerAddress", 
			mps_port as "mpsPort", 
			user_name as "username", 
			password as "password", 
			common_name as "commonName",
			server_address_format as "serverAddressFormat", 
			auth_method as "authMethod", 
			mps_root_certificate as "mpsRootCertificate", 
			proxydetails as "proxyDetails", 
			tenant_id  as "tenantId",
			xmin as "version"
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

	rows, err := r.Pool.Query(ctx, sql, tenantID)
	if err != nil {
		return nil, fmt.Errorf("CIRARepo - Get - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	configs := make([]entity.CIRAConfig, 0, top)

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
	sql, _, err := r.Builder.
		Select(`
			cira_config_name as "configName", 
			mps_server_address as "mpsServerAddress", 
			mps_port as "mpsPort", 
			user_name as "username", 
			password as "password", 
			common_name as "commonName", 
			server_address_format as "serverAddressFormat", 
			auth_method as "authMethod", 
			mps_root_certificate as "mpsRootCertificate", 
			proxydetails as "proxyDetails", 
			tenant_id as "tenantId",
			xmin as "version"
			`).
		From("ciraconfigs").
		Where("cira_config_name = ? and tenant_id = ?", configName, tenantID).
		ToSql()
	if err != nil {
		return entity.CIRAConfig{}, fmt.Errorf("CIRARepo - Get - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, configName, tenantID)
	if err != nil {
		return entity.CIRAConfig{}, fmt.Errorf("CIRARepo - Get - r.Pool.Query: %w", err)
	}

	defer rows.Close()

	configs := make([]entity.CIRAConfig, 0, 1)
	for rows.Next() {
		p := entity.CIRAConfig{}

		err = rows.Scan(&p.ConfigName, &p.MPSServerAddress, &p.MpsPort, &p.Username, &p.Password, &p.CommonName, &p.ServerAddressFormat, &p.AuthMethod, &p.MpsRootCertificate, &p.ProxyDetails, &p.TenantID, &p.Version)
		if err != nil {
			return p, fmt.Errorf("CIRARepo - Get - rows.Scan: %w", err)
		}

		configs = append(configs, p)
	}

	if len(configs) == 0 {
		return entity.CIRAConfig{}, nil
	}

	return configs[0], nil
}

// Delete -.
func (r *CIRARepo) Delete(ctx context.Context, configName, tenantID string) (bool, error) {
	sql, _, err := r.Builder.
		Delete("ciraconfigs").
		Where("cira_config_name = ? AND tenant_id = ?", configName, tenantID).
		ToSql()
	if err != nil {
		return false, fmt.Errorf("CIRARepo - Delete - r.Builder: %w", err)
	}

	res, err := r.Pool.Exec(ctx, sql)
	if err != nil {
		return false, fmt.Errorf("CIRARepo - Delete - r.Pool.Exec: %w", err)
	}

	return res.RowsAffected() > 0, nil
}

// Update -.
func (r *CIRARepo) Update(ctx context.Context, p *entity.CIRAConfig) (bool, error) {
	sql, args, err := r.Builder.
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
		Suffix("AND xmin::text = ?", p.Version).
		ToSql()
	if err != nil {
		return false, fmt.Errorf("CIRARepo - Update - r.Builder: %w", err)
	}

	res, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return false, fmt.Errorf("CIRARepo - Update - r.Pool.Exec: %w", err)
	}

	return res.RowsAffected() > 0, nil
}

// Insert -.
func (r *CIRARepo) Insert(ctx context.Context, p *entity.CIRAConfig) (string, error) {
	sql, args, err := r.Builder.
		Insert("ciraconfigs").
		Columns("cira_config_name", "mps_server_address", "mps_port", "user_name", "password", "common_name", "server_address_format", "auth_method", "mps_root_certificate", "proxydetails", "tenant_id").
		Values(p.ConfigName, p.MPSServerAddress, p.MpsPort, p.Username, p.Password, p.CommonName, p.ServerAddressFormat, p.AuthMethod, p.MpsRootCertificate, p.ProxyDetails, p.TenantID).
		Suffix("RETURNING xmin::text").
		ToSql()
	if err != nil {
		return "", fmt.Errorf("CIRARepo - Insert - r.Builder: %w", err)
	}

	var version string

	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&version)
	if err != nil {
		return "", fmt.Errorf("CIRARepo - Insert - r.Pool.QueryRow: %w", err)
	}

	return version, nil
}
