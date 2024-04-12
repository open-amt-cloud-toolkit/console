package postgresdb

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	"github.com/open-amt-cloud-toolkit/console/pkg/postgres"
)

// DeviceRepo -.
type DeviceRepo struct {
	*postgres.DB
}

// New -.
func NewDeviceRepo(pg *postgres.DB) *DeviceRepo {
	return &DeviceRepo{pg}
}

// GetCount -.
func (r *DeviceRepo) GetCount(ctx context.Context, tenantID string) (int, error) {
	sql, _, err := r.Builder.
		Select("COUNT(*) OVER() AS total_count").
		From("devices").
		Where("tenantid = ?", tenantID).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("DeviceRepo - GetCount - r.Builder: %w", err)
	}

	var count int

	err = r.Pool.QueryRow(ctx, sql, tenantID).Scan(&count)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return 0, nil
		}

		return 0, fmt.Errorf("DeviceRepo - GetCount - r.Pool.QueryRow: %w", err)
	}

	return count, nil
}

// Get -.
func (r *DeviceRepo) Get(ctx context.Context, top, skip int, tenantID string) ([]entity.Device, error) {
	if top == 0 {
		top = 100
	}

	sql, _, err := r.Builder.
		Select(`guid as "guid",
				hostname as "hostname",
				tags as "tags",
				mpsinstance as "mpsInstance",
				connectionstatus as "connectionStatus",
				mpsusername as "mpsusername",
				tenantid as "tenantId",
				friendlyname as "friendlyName",
				dnssuffix as "dnsSuffix",
				deviceinfo as "deviceInfo",
				username as "username",
				password as "password",
				usetls as "useTLS",
				allowselfsigned as "allowSelfSigned"
				`).
		From("devices").
		Where("tenantid = ?", tenantID).
		OrderBy("guid").
		Limit(uint64(top)).
		Offset(uint64(skip)).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("DeviceRepo - Get - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, tenantID)
	if err != nil {
		return nil, fmt.Errorf("DeviceRepo - Get - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	domains := make([]entity.Device, 0, top)

	for rows.Next() {
		d := entity.Device{}

		err = rows.Scan(&d.GUID, &d.Hostname, &d.Tags, &d.MpsInstance, &d.ConnectionStatus, &d.Mpsusername, &d.TenantID, &d.FriendlyName, &d.DNSSuffix, &d.DeviceInfo, &d.Username, &d.Password, &d.UseTLS, &d.AllowSelfSigned)
		if err != nil {
			return nil, fmt.Errorf("DeviceRepo - Get - rows.Scan: %w", err)
		}

		domains = append(domains, d)
	}

	return domains, nil
}

// GetByID -.
func (r *DeviceRepo) GetByID(ctx context.Context, guid, tenantID string) (entity.Device, error) {
	sql, _, err := r.Builder.
		Select(`guid as "guid",
				hostname as "hostname",
				tags as "tags",
				mpsinstance as "mpsInstance",
				connectionstatus as "connectionStatus",
				mpsusername as "mpsusername",
				tenantid as "tenantId",
				friendlyname as "friendlyName",
				dnssuffix as "dnsSuffix",
				deviceinfo as "deviceInfo",
				username as "username",
				password as "password",
				usetls as "useTLS",
				allowselfsigned as "allowSelfSigned"
				`).
		From("devices").
		Where("guid = ? and tenantid = ?").
		ToSql()
	if err != nil {
		return entity.Device{}, fmt.Errorf("DeviceRepo - Get - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, guid, tenantID)
	if err != nil {
		return entity.Device{}, fmt.Errorf("DeviceRepo - Get - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	devices := make([]entity.Device, 0, 1)

	for rows.Next() {
		d := entity.Device{}

		err = rows.Scan(&d.GUID, &d.Hostname, &d.Tags, &d.MpsInstance, &d.ConnectionStatus, &d.Mpsusername, &d.TenantID, &d.FriendlyName, &d.DNSSuffix, &d.DeviceInfo, &d.Username, &d.Password, &d.UseTLS, &d.AllowSelfSigned)
		if err != nil {
			return d, fmt.Errorf("DeviceRepo - Get - rows.Scan: %w", err)
		}

		devices = append(devices, d)
	}

	if len(devices) == 0 {
		return entity.Device{}, nil
	}

	return devices[0], nil
}

// Delete -.
func (r *DeviceRepo) Delete(ctx context.Context, name, tenantID string) (bool, error) {
	sql, _, err := r.Builder.
		Delete("devices").
		Where("name = ? AND tenantid = ?", name, tenantID).
		ToSql()
	if err != nil {
		return false, fmt.Errorf("DeviceRepo - Delete - r.Builder: %w", err)
	}

	res, err := r.Pool.Exec(ctx, sql)
	if err != nil {
		return false, fmt.Errorf("DeviceRepo - Delete - r.Pool.Exec: %w", err)
	}

	return res.RowsAffected() > 0, nil
}

// Update -.
func (r *DeviceRepo) Update(ctx context.Context, d *entity.Device) (bool, error) {
	sql, args, err := r.Builder.
		Update("devices").
		Set("guid", d.GUID).
		Set("hostname", d.Hostname).
		Set("tags", d.Tags).
		Set("mpsinstance", d.MpsInstance).
		Set("connectionstatus", d.ConnectionStatus).
		Set("mpsusername", d.Mpsusername).
		Set("tenantid", d.TenantID).
		Set("friendlyname", d.FriendlyName).
		Set("dnssuffix", d.DNSSuffix).
		Set("deviceinfo", d.DeviceInfo).
		Set("username", d.Username).
		Set("password", d.Password).
		Set("useTLS", d.UseTLS).
		Set("allowSelfSigned", d.AllowSelfSigned).
		Where("guid = ? AND tenantid = ?", d.GUID, d.TenantID).
		ToSql()
	if err != nil {
		return false, fmt.Errorf("DeviceRepo - Update - r.Builder: %w", err)
	}

	res, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return false, fmt.Errorf("DeviceRepo - Update - r.Pool.Exec: %w", err)
	}

	return res.RowsAffected() > 0, nil
}

// Insert -.
func (r *DeviceRepo) Insert(ctx context.Context, d *entity.Device) (string, error) {
	d.GUID = uuid.New().String()
	sql, args, err := r.Builder.
		Insert("devices").
		Columns("guid", "hostname", "tags", "mpsinstance", "connectionstatus", "mpsusername", "tenantid", "friendlyname", "dnssuffix", "deviceinfo", "username", "password", "usetls", "allowselfsigned").
		Values(d.GUID, d.Hostname, d.Tags, d.MpsInstance, d.ConnectionStatus, d.Mpsusername, d.TenantID, d.FriendlyName, d.DNSSuffix, d.DeviceInfo, d.Username, d.Password, d.UseTLS, d.AllowSelfSigned).
		Suffix("RETURNING xmin::text").
		ToSql()
	if err != nil {
		return "", fmt.Errorf("DeviceRepo - Insert - r.Builder: %w", err)
	}

	var version string

	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&version)
	if err != nil {
		return "", fmt.Errorf("DeviceRepo - Insert - r.Pool.QueryRow: %w", err)
	}

	return version, nil
}
