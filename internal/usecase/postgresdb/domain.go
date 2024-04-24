package postgresdb

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	"github.com/open-amt-cloud-toolkit/console/pkg/postgres"
)

// DomainRepo -.
type DomainRepo struct {
	*postgres.DB
}

// New -.
func NewDomainRepo(pg *postgres.DB) *DomainRepo {
	return &DomainRepo{pg}
}

// GetCount -.
func (r *DomainRepo) GetCount(ctx context.Context, tenantID string) (int, error) {
	sql, _, err := r.Builder.
		Select("COUNT(*) OVER() AS total_count").
		From("domains").
		Where("tenant_id = ?", tenantID).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("DomainRepo - GetCount - r.Builder: %w", err)
	}

	var count int

	err = r.Pool.QueryRow(ctx, sql, tenantID).Scan(&count)
	if err != nil {
		if err.Error() == NoRowsInResultSet {
			return 0, nil
		}

		return 0, fmt.Errorf("DomainRepo - GetCount - r.Pool.QueryRow: %w", err)
	}

	return count, nil
}

// Get -.
func (r *DomainRepo) Get(ctx context.Context, top, skip int, tenantID string) ([]entity.Domain, error) {
	if top == 0 {
		top = 100
	}

	sql, _, err := r.Builder.
		Select(`name as "profileName",
				domain_suffix as "domainSuffix",
				provisioning_cert as "provisioningCert",
				provisioning_cert_storage_format as "provisioningCertStorageFormat",
				provisioning_cert_key as "provisioningCertPassword",
				tenant_id as "tenantId",
				xmin as "version"`).
		From("domains").
		Where("tenant_id = ?", tenantID).
		OrderBy("name").
		Limit(uint64(top)).
		Offset(uint64(skip)).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("DomainRepo - Get - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, tenantID)
	if err != nil {
		return nil, fmt.Errorf("DomainRepo - Get - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	domains := make([]entity.Domain, 0, top)

	for rows.Next() {
		d := entity.Domain{}

		err = rows.Scan(&d.ProfileName, &d.DomainSuffix, &d.ProvisioningCert, &d.ProvisioningCertStorageFormat, &d.ProvisioningCertPassword, &d.TenantID, &d.Version)
		if err != nil {
			return nil, fmt.Errorf("DomainRepo - Get - rows.Scan: %w", err)
		}

		domains = append(domains, d)
	}

	return domains, nil
}

// GetDomainByDomainSuffix -.
//
//nolint:dupl // sql statements are going to be similar
func (r *DomainRepo) GetDomainByDomainSuffix(ctx context.Context, domainSuffix, tenantID string) (*entity.Domain, error) {
	sql, _, err := r.Builder.
		Select(`name as "profileName",
				domain_suffix as "domainSuffix",
				provisioning_cert as "provisioningCert",
				provisioning_cert_storage_format as "provisioningCertStorageFormat",
				provisioning_cert_key as "provisioningCertPassword",
				tenant_id as "tenantId",
				xmin as "version"`).
		From("domains").
		Where("domain_suffix = ? AND tenant_id = ?", domainSuffix, tenantID).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("DomainRepo - GetDomainByDomainSuffix - r.Builder: %w", err)
	}

	row := r.Pool.QueryRow(ctx, sql)
	d := entity.Domain{}

	err = row.Scan(&d.ProfileName, &d.DomainSuffix, &d.ProvisioningCert, &d.ProvisioningCertStorageFormat, &d.ProvisioningCertPassword, &d.TenantID, &d.Version)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("DomainRepo - GetDomainByDomainSuffix - row.Scan: %w", err)
	}

	return &d, nil
}

// GetByName -.
//
//nolint:dupl // sql statements are going to be similar
func (r *DomainRepo) GetByName(ctx context.Context, domainName, tenantID string) (*entity.Domain, error) {
	sql, _, err := r.Builder.
		Select(`name as "profileName",
				domain_suffix as "domainSuffix",
				provisioning_cert as "provisioningCert",
				provisioning_cert_storage_format as "provisioningCertStorageFormat",
				provisioning_cert_key as "provisioningCertPassword",
				tenant_id as "tenantId",
				xmin as "version"`).
		From("domains").
		Where("name = ? AND tenant_id = ?", domainName, tenantID).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("DomainRepo - GetByName - r.Builder: %w", err)
	}

	row := r.Pool.QueryRow(ctx, sql)
	d := entity.Domain{}

	err = row.Scan(&d.ProfileName, &d.DomainSuffix, &d.ProvisioningCert, &d.ProvisioningCertStorageFormat, &d.ProvisioningCertPassword, &d.TenantID, &d.Version)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("DomainRepo - GetByName - row.Scan: %w", err)
	}

	return &d, nil
}

// Delete -.
func (r *DomainRepo) Delete(ctx context.Context, domainName, tenantID string) (bool, error) {
	sql, _, err := r.Builder.
		Delete("domains").
		Where("name = ? AND tenant_id = ?", domainName, tenantID).
		ToSql()
	if err != nil {
		return false, fmt.Errorf("DomainRepo - Delete - r.Builder: %w", err)
	}

	res, err := r.Pool.Exec(ctx, sql)
	if err != nil {
		return false, fmt.Errorf("DomainRepo - Delete - r.Pool.Exec: %w", err)
	}

	return res.RowsAffected() > 0, nil
}

// Update -.
func (r *DomainRepo) Update(ctx context.Context, d *entity.Domain) (bool, error) {
	sql, args, err := r.Builder.
		Update("domains").
		Set("name", d.ProfileName).
		Set("domain_suffix", d.DomainSuffix).
		Set("provisioning_cert", d.ProvisioningCert).
		Set("provisioning_cert_storage_format", d.ProvisioningCertStorageFormat).
		Set("provisioning_cert_key", d.ProvisioningCertPassword).
		Where("name = ? AND tenant_id = ?", d.ProfileName, d.TenantID).
		Suffix("AND xmin::text = ?", d.Version).
		ToSql()
	if err != nil {
		return false, fmt.Errorf("DomainRepo - Update - r.Builder: %w", err)
	}

	res, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return false, fmt.Errorf("DomainRepo - Update - r.Pool.Exec: %w", err)
	}

	return res.RowsAffected() > 0, nil
}

// Insert -.
func (r *DomainRepo) Insert(ctx context.Context, d *entity.Domain) (string, error) {
	sql, args, err := r.Builder.
		Insert("domains").
		Columns("name", "domain_suffix", "provisioning_cert", "provisioning_cert_storage_format", "provisioning_cert_key", "tenant_id").
		Values(d.ProfileName, d.DomainSuffix, d.ProvisioningCert, d.ProvisioningCertStorageFormat, d.ProvisioningCertPassword, d.TenantID).
		Suffix("RETURNING xmin::text").
		ToSql()
	if err != nil {
		return "", fmt.Errorf("DomainRepo - Insert - r.Builder: %w", err)
	}

	var version string

	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&version)
	if err != nil {
		return "", fmt.Errorf("DomainRepo - Insert - r.Pool.QueryRow: %w", err)
	}

	return version, nil
}
