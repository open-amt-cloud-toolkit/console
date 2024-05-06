package postgresdb

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	"github.com/open-amt-cloud-toolkit/console/pkg/consoleerrors"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
	"github.com/open-amt-cloud-toolkit/console/pkg/postgres"
)

// DeviceRepo -.
type DeviceRepo struct {
	*postgres.DB
	log logger.Interface
}

var (
	ErrDeviceDatabase  = consoleerrors.DatabaseError{Console: consoleerrors.CreateConsoleError("DeviceRepo")}
	ErrDeviceNotUnique = consoleerrors.DatabaseError{Console: consoleerrors.CreateConsoleError("DeviceRepo")}
)

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
		return 0, ErrDeviceDatabase.Wrap("GetCount", "r.Builder: ", err)
	}

	var count int

	err = r.Pool.QueryRow(ctx, sqlQuery, tenantID).Scan(&count)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}

		return 0, ErrDeviceDatabase.Wrap("GetCount", "r.Pool.QueryRow", err)
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
		return nil, ErrDeviceDatabase.Wrap("Get", "r.Builder: ", err)
	}

	rows, err := r.Pool.Query(ctx, sqlQuery, tenantID)
	if err != nil {
		return nil, ErrDeviceDatabase.Wrap("Get", "r.Pool.Query", err)
	}

	defer rows.Close()

	domains := make([]entity.Device, 0)

	for rows.Next() {
		d := entity.Device{}

		err = rows.Scan(&d.GUID, &d.Hostname, &d.Tags, &d.MpsInstance, &d.ConnectionStatus, &d.Mpsusername, &d.TenantID, &d.FriendlyName, &d.DNSSuffix, &d.DeviceInfo, &d.Username, &d.Password, &d.UseTLS, &d.AllowSelfSigned)
		if err != nil {
			return nil, ErrDeviceDatabase.Wrap("Get", "rows.Scan: ", err)
		}

		domains = append(domains, d)
	}

	return domains, nil
}

// GetByID -.
func (r *DeviceRepo) GetByID(ctx context.Context, guid, tenantID string) (*entity.Device, error) {
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
		return nil, ErrDeviceDatabase.Wrap("Get", "r.Builder: ", err)
	}

	rows, err := r.Pool.Query(ctx, sqlQuery, guid, tenantID)
	if err != nil {
		return nil, ErrDeviceDatabase.Wrap("Get", "r.Pool.Query", err)
	}

	defer rows.Close()

	devices := make([]*entity.Device, 0)

	for rows.Next() {
		d := &entity.Device{}

		err = rows.Scan(&d.GUID, &d.Hostname, &d.Tags, &d.MpsInstance, &d.ConnectionStatus, &d.Mpsusername, &d.TenantID, &d.FriendlyName, &d.DNSSuffix, &d.DeviceInfo, &d.Username, &d.Password, &d.UseTLS, &d.AllowSelfSigned)
		if err != nil {
			return d, ErrDeviceDatabase.Wrap("Get", "rows.Scan: ", err)
		}

		devices = append(devices, d)
	}

	if len(devices) == 0 {
		return nil, nil
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
		return []string{}, ErrDeviceDatabase.Wrap("GetDistinctTags", "r.Builder: ", err)
	}

	rows, err := r.Pool.Query(ctx, sqlQuery, tenantID)
	if err != nil {
		return []string{}, ErrDeviceDatabase.Wrap("GetDistinctTags", "r.Pool.Query", err)
	}

	defer rows.Close()

	tags := make([]string, 0)

	for rows.Next() {
		var tag string

		err = rows.Scan(&tag)
		if err != nil {
			return []string{tag}, ErrDeviceDatabase.Wrap("GetDistinctTags", "rows.Scan: ", err)
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
		return []entity.Device{}, ErrDeviceDatabase.Wrap("GetByTags", "r.Builder: ", err)
	}

	rows, err := r.Pool.Query(ctx, sqlQuery, tags, tenantID)
	if err != nil {
		return []entity.Device{}, ErrDeviceDatabase.Wrap("GetByTags", "r.Pool.Query", err)
	}

	defer rows.Close()

	devices := make([]entity.Device, 0)

	for rows.Next() {
		d := entity.Device{}

		err = rows.Scan(&d.GUID, &d.Hostname, &d.Tags, &d.MpsInstance, &d.ConnectionStatus, &d.Mpsusername, &d.TenantID, &d.FriendlyName, &d.DNSSuffix, &d.DeviceInfo)
		if err != nil {
			return []entity.Device{d}, ErrDeviceDatabase.Wrap("GetByTags", "rows.Scan", err)
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
		return false, ErrDeviceDatabase.Wrap("Delete", "r.Builder", err)
	}

	res, err := r.Pool.Exec(ctx, sqlQuery, guid, tenantID)
	if err != nil {
		return false, ErrDeviceDatabase.Wrap("Delete", "r.Pool.Exec", err)
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
		return false, ErrDeviceDatabase.Wrap("Update", "r.Builder", err)
	}

	res, err := r.Pool.Exec(ctx, sqlQuery, args...)
	if err != nil {
		return false, ErrDeviceDatabase.Wrap("Update", "r.Pool.Exec", err)
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
		return "", ErrDeviceDatabase.Wrap("Insert", "r.Builder", err)
	}

	var version string

	err = r.Pool.QueryRow(ctx, sqlQuery, args...).Scan(&version)
	if err != nil {
		if postgres.CheckNotUnique(err) {
			return "", ErrDeviceNotUnique
		}

		return "", ErrDeviceDatabase.Wrap("Insert", "r.Pool.QueryRow", err)
	}

	return version, nil
}
