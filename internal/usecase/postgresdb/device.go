package postgresdb

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
	"github.com/open-amt-cloud-toolkit/console/pkg/postgres"
)

// DeviceRepo -.
type DeviceRepo struct {
	*postgres.DB
	log logger.Interface
}

// New -.
func NewDeviceRepo(pg *postgres.DB, log logger.Interface) *DeviceRepo {
	return &DeviceRepo{pg, log}
}

// GetCount -.
func (r *DeviceRepo) GetCount(ctx context.Context, tenantID string) (int, error) {
	sqlQuery, _, err := r.Builder.
		Select("COUNT(*) OVER() AS total_count").
		From("devices").
		Where("tenantid = ?", tenantID).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("DeviceRepo - GetCount - r.Builder: %w", err)
	}

	var count int

	err = r.Pool.QueryRow(ctx, sqlQuery, tenantID).Scan(&count)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
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

	sqlQuery, _, err := r.Builder.
		Select(`guid, 
				hostname, 
				tags, 
				mpsinstance, 
				connectionstatus, 
				mpsusername, 
				tenantid, 
				friendlyname, 
				dnssuffix, 
				deviceinfo, 
				username, 
				password, 
				usetls, 
				allowselfsigned 
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

	rows, err := r.Pool.Query(ctx, sqlQuery, tenantID)
	if err != nil {
		return nil, fmt.Errorf("DeviceRepo - Get - r.Pool.Query: %w", err)
	}

	defer rows.Close()

	domains := make([]entity.Device, 0)

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
	sqlQuery, _, err := r.Builder.
		Select(`guid,
				hostname,
				tags,
				mpsinstance,
				connectionstatus,
				mpsusername,
				tenantid,
				friendlyname,
				dnssuffix,
				deviceinfo,
				username,
				password,
				usetls,
				allowselfsigned
		`).
		From("devices").
		Where("guid = ? and tenantid = ?").
		ToSql()
	if err != nil {
		return entity.Device{}, fmt.Errorf("DeviceRepo - Get - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sqlQuery, guid, tenantID)
	if err != nil {
		return entity.Device{}, fmt.Errorf("DeviceRepo - Get - r.Pool.Query: %w", err)
	}

	defer rows.Close()

	devices := make([]entity.Device, 0)

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

func (r *DeviceRepo) GetDistinctTags(ctx context.Context, tenantID string) ([]string, error) {
	sqlQuery, _, err := r.Builder.
		Select("DISTINCT unnest(tags) as tag").
		From("devices").
		Where("tenantid = ?", tenantID).
		ToSql()
	if err != nil {
		return []string{}, fmt.Errorf("DeviceRepo - GetDistinctTags - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sqlQuery, tenantID)
	if err != nil {
		return []string{}, fmt.Errorf("DeviceRepo - GetDistinctTags - r.Pool.Query: %w", err)
	}

	defer rows.Close()

	tags := make([]string, 0)

	for rows.Next() {
		var tag string

		err = rows.Scan(&tag)
		if err != nil {
			return []string{tag}, fmt.Errorf("DeviceRepo - GetDistinctTags - rows.Scan: %w", err)
		}

		tags = append(tags, tag)
	}

	if len(tags) == 0 {
		return []string{}, nil
	}

	return tags, nil
}

func (r *DeviceRepo) GetByTags(ctx context.Context, tags []string, method string, limit, offset int, tenantID string) ([]entity.Device, error) {
	builder := r.Builder.
		Select(`guid,
            hostname,
            tags,
            mpsinstance,
            connectionstatus,
            mpsusername,
            tenantid,
            friendlyname,
            dnssuffix,
            deviceinfo`).
		From("devices")
	if method == "AND" {
		builder = builder.Where("tags @> ? and tenantId = ?", tags, tenantID)
	} else {
		builder = builder.Where("tags && ? and tenantId = ?", tags, tenantID)
	}

	sqlQuery, _, err := builder.OrderBy("guid").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		ToSql()
	if err != nil {
		return []entity.Device{}, fmt.Errorf("DeviceRepo - GetByTags - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sqlQuery, tags, tenantID)
	if err != nil {
		return []entity.Device{}, fmt.Errorf("DeviceRepo - GetByTags - r.Pool.Query: %w", err)
	}

	defer rows.Close()

	devices := make([]entity.Device, 0)

	for rows.Next() {
		d := entity.Device{}

		err = rows.Scan(&d.GUID, &d.Hostname, &d.Tags, &d.MpsInstance, &d.ConnectionStatus, &d.Mpsusername, &d.TenantID, &d.FriendlyName, &d.DNSSuffix, &d.DeviceInfo)
		if err != nil {
			return []entity.Device{d}, fmt.Errorf("DeviceRepo - GetByTags - rows.Scan: %w", err)
		}

		devices = append(devices, d)
	}

	if len(devices) == 0 {
		return []entity.Device{}, nil
	}

	return devices, nil
}

// Delete -.
func (r *DeviceRepo) Delete(ctx context.Context, guid, tenantID string) (bool, error) {
	sqlQuery, _, err := r.Builder.
		Delete("devices").
		Where("guid = ? AND tenantid = ?", guid, tenantID).
		ToSql()
	if err != nil {
		return false, fmt.Errorf("DeviceRepo - Delete - r.Builder: %w", err)
	}

	res, err := r.Pool.Exec(ctx, sqlQuery, guid, tenantID)
	if err != nil {
		return false, fmt.Errorf("DeviceRepo - Delete - r.Pool.Exec: %w", err)
	}

	return res.RowsAffected() > 0, nil
}

// Update -.
func (r *DeviceRepo) Update(ctx context.Context, d *entity.Device) (bool, error) {
	sqlQuery, args, err := r.Builder.
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

	res, err := r.Pool.Exec(ctx, sqlQuery, args...)
	if err != nil {
		return false, fmt.Errorf("DeviceRepo - Update - r.Pool.Exec: %w", err)
	}

	return res.RowsAffected() > 0, nil
}

// Insert -.
func (r *DeviceRepo) Insert(ctx context.Context, d *entity.Device) (string, error) {
	d.GUID = uuid.New().String()

	sqlQuery, args, err := r.Builder.
		Insert("devices").
		Columns("guid", "hostname", "tags", "mpsinstance", "connectionstatus", "mpsusername", "tenantid", "friendlyname", "dnssuffix", "deviceinfo", "username", "password", "usetls", "allowselfsigned").
		Values(d.GUID, d.Hostname, d.Tags, d.MpsInstance, d.ConnectionStatus, d.Mpsusername, d.TenantID, d.FriendlyName, d.DNSSuffix, d.DeviceInfo, d.Username, d.Password, d.UseTLS, d.AllowSelfSigned).
		Suffix("RETURNING xmin::text").
		ToSql()
	if err != nil {
		return "", fmt.Errorf("DeviceRepo - Insert - r.Builder: %w", err)
	}

	var version string

	err = r.Pool.QueryRow(ctx, sqlQuery, args...).Scan(&version)
	if err != nil {
		return "", fmt.Errorf("DeviceRepo - Insert - r.Pool.QueryRow: %w", err)
	}

	return version, nil
}
