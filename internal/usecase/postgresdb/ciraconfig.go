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

// CIRARepo -.
type CIRARepo struct {
	*postgres.DB
	log logger.Interface
}

// New -.
func NewCIRARepo(pg *postgres.DB, log logger.Interface) *CIRARepo {
	return &CIRARepo{pg, log}
}

var (
	ErrCIRARepo          = consoleerrors.CreateConsoleError("CIRARepo")
	ErrCIRARepoDatabase  = consoleerrors.DatabaseError{Console: consoleerrors.CreateConsoleError("CIRARepo")}
	ErrCIRARepoNotFound  = consoleerrors.NotFoundError{Console: consoleerrors.CreateConsoleError("CIRARepo")}
	ErrCIRARepoNotUnique = consoleerrors.NotUniqueError{Console: consoleerrors.CreateConsoleError("CIRARepo")}
)

// GetCount -.
func (r *CIRARepo) GetCount(_ context.Context, tenantID string) (int, error) {
	sqlQuery, _, err := r.Builder.
		Select("COUNT(*) OVER() AS total_count").
		From("ciraconfigs").
		Where("tenant_id = ?", tenantID).
		ToSql()
	if err != nil {
		return 0, ErrCIRARepoDatabase.Wrap("GetCount", "r.Builder", err)
	}

	var count int

	err = r.Pool.QueryRow(sqlQuery, tenantID).Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}

		return 0, ErrCIRARepoDatabase.Wrap("GetCount", "r.Pool.QueryRow", err)
	}

	return count, nil
}

// Get -.
func (r *CIRARepo) Get(_ context.Context, top, skip int, tenantID string) ([]entity.CIRAConfig, error) {
	if top == 0 {
		top = 100
	}

	sqlQuery, _, err := r.Builder.
		Select("cira_config_name",
			"mps_server_address",
			"mps_port",
			"user_name",
			"password",
			"common_name",
			"server_address_format",
			"auth_method",
			"mps_root_certificate",
			"proxydetails",
			"tenant_id",
			"CAST(xmin as text) as xmin").
		From("ciraconfigs").
		Where("tenant_id = ?", tenantID).
		OrderBy("cira_config_name").
		Limit(uint64(top)).
		Offset(uint64(skip)).
		ToSql()
	if err != nil {
		return nil, ErrCIRARepoDatabase.Wrap("Get", "r.Builder", err)
	}

	rows, err := r.Pool.Query(sqlQuery, tenantID)
	if err != nil {
		return nil, ErrCIRARepoDatabase.Wrap("Get", "r.Pool.Query", err)
	}

	defer rows.Close()

	if rows.Err() != nil {
		return nil, ErrDeviceDatabase.Wrap("Get", "rows.Err", rows.Err())
	}

	configs := make([]entity.CIRAConfig, 0)

	for rows.Next() {
		p := entity.CIRAConfig{}

		err = rows.Scan(&p.ConfigName, &p.MPSAddress, &p.MPSPort, &p.Username, &p.Password, &p.CommonName, &p.ServerAddressFormat, &p.AuthMethod, &p.MPSRootCertificate, &p.ProxyDetails, &p.TenantID, &p.Version)
		if err != nil {
			return nil, ErrCIRARepoDatabase.Wrap("Get", "rows.Scan", err)
		}

		configs = append(configs, p)
	}

	return configs, nil
}

// GetByName -.
func (r *CIRARepo) GetByName(_ context.Context, configName, tenantID string) (*entity.CIRAConfig, error) {
	sqlQuery, _, err := r.Builder.
		Select("cira_config_name",
			"mps_server_address",
			"mps_port",
			"user_name",
			"password",
			"common_name",
			"server_address_format",
			"auth_method",
			"mps_root_certificate",
			"proxydetails",
			"tenant_id",
			"CAST(xmin as text) as xmin").
		From("ciraconfigs").
		Where("cira_config_name = ? and tenant_id = ?", configName, tenantID).
		ToSql()
	if err != nil {
		return nil, ErrCIRARepoDatabase.Wrap("GetByName", "r.Builder", err)
	}

	rows, err := r.Pool.Query(sqlQuery, configName, tenantID)
	if err != nil {
		return nil, ErrCIRARepoDatabase.Wrap("GetByName", "r.Pool.Query", err)
	}

	defer rows.Close()

	if rows.Err() != nil {
		return nil, ErrDeviceDatabase.Wrap("Get", "rows.Err", rows.Err())
	}

	configs := make([]*entity.CIRAConfig, 0)

	for rows.Next() {
		p := &entity.CIRAConfig{}

		err = rows.Scan(&p.ConfigName, &p.MPSAddress, &p.MPSPort, &p.Username, &p.Password, &p.CommonName, &p.ServerAddressFormat, &p.AuthMethod, &p.MPSRootCertificate, &p.ProxyDetails, &p.TenantID, &p.Version)
		if err != nil {
			return p, ErrCIRARepoDatabase.Wrap("GetByName", "rows.Scan", err)
		}

		configs = append(configs, p)
	}

	if len(configs) == 0 {
		return nil, nil
	}

	return configs[0], nil
}

// Delete -.
func (r *CIRARepo) Delete(_ context.Context, configName, tenantID string) (bool, error) {
	sqlQuery, args, err := r.Builder.
		Delete("ciraconfigs").
		Where("cira_config_name = ? AND tenant_id = ?", configName, tenantID).
		ToSql()
	if err != nil {
		return false, ErrCIRARepoDatabase.Wrap("Delete", "r.Builder", err)
	}

	res, err := r.Pool.Exec(sqlQuery, args...)
	if err != nil {
		return false, ErrCIRARepoDatabase.Wrap("Delete", "r.Pool.Exec", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, ErrCIRARepoDatabase.Wrap("Delete", "res.RowsAffected", err)
	}

	return rowsAffected > 0, nil
}

// Update -.
func (r *CIRARepo) Update(_ context.Context, p *entity.CIRAConfig) (bool, error) {
	sqlQuery, args, err := r.Builder.
		Update("ciraconfigs").
		Set("mps_server_address", p.MPSAddress).
		Set("mps_port", p.MPSPort).
		Set("user_name", p.Username).
		Set("password", p.Password).
		Set("common_name", p.CommonName).
		Set("server_address_format", p.ServerAddressFormat).
		Set("auth_method", p.AuthMethod).
		Set("mps_root_certificate", p.MPSRootCertificate).
		Set("proxydetails", p.ProxyDetails).
		Where("cira_config_name = ? AND tenant_id = ?", p.ConfigName, p.TenantID).
		ToSql()
	if err != nil {
		return false, ErrCIRARepoDatabase.Wrap("Update", "r.Builder", err)
	}

	res, err := r.Pool.Exec(sqlQuery, args...)
	if err != nil {
		return false, ErrCIRARepoDatabase.Wrap("Update", "r.Pool.Exec", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, ErrCIRARepoDatabase.Wrap("Delete", "res.RowsAffected", err)
	}

	return rowsAffected > 0, nil
}

// Insert -.
func (r *CIRARepo) Insert(_ context.Context, p *entity.CIRAConfig) (string, error) {
	sqlQuery, args, err := r.Builder.
		Insert("ciraconfigs").
		Columns("cira_config_name", "mps_server_address", "mps_port", "user_name", "password", "common_name", "server_address_format", "auth_method", "mps_root_certificate", "proxydetails", "tenant_id").
		Values(p.ConfigName, p.MPSAddress, p.MPSPort, p.Username, p.Password, p.CommonName, p.ServerAddressFormat, p.AuthMethod, p.MPSRootCertificate, p.ProxyDetails, p.TenantID).
		Suffix("RETURNING xmin::text").
		ToSql()
	if err != nil {
		return "", ErrCIRARepoDatabase.Wrap("Insert", "r.Builder", err)
	}

	var version string

	err = r.Pool.QueryRow(sqlQuery, args...).Scan(&version)
	if err != nil {
		if postgres.CheckNotUnique(err) {
			return "", ErrCIRARepoNotUnique.Wrap("Insert", "r.Pool.QueryRow", err)
		}

		return "", ErrCIRARepoDatabase.Wrap("Insert", "r.Pool.QueryRow", err)
	}

	return version, nil
}
