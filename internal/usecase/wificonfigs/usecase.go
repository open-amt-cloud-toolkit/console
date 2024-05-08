package wificonfigs

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto"
	"github.com/open-amt-cloud-toolkit/console/pkg/consoleerrors"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
)

// UseCase -.
type UseCase struct {
	repo Repository
	log  logger.Interface
}

// New -.
func New(r Repository, log logger.Interface) *UseCase {
	return &UseCase{
		repo: r,
		log:  log,
	}
}

var (
	ErrCountNotUnique = consoleerrors.NotUniqueError{Console: consoleerrors.CreateConsoleError("WifiConfigs")}
	ErrDomainsUseCase = consoleerrors.CreateConsoleError("WificonfigsUseCase")
	ErrDatabase       = consoleerrors.DatabaseError{Console: consoleerrors.CreateConsoleError("WificonfigsUseCase")}
	ErrNotFound       = consoleerrors.NotFoundError{Console: consoleerrors.CreateConsoleError("WificonfigsUseCase")}
)

// History - getting translate history from store.
func (uc *UseCase) CheckProfileExists(ctx context.Context, profileName, tenantID string) (bool, error) {
	data, err := uc.repo.CheckProfileExists(ctx, profileName, tenantID)
	if err != nil {
		return false, ErrDatabase.Wrap("Count", "uc.repo.GetCount", err)
	}

	return data, nil
}

func (uc *UseCase) GetCount(ctx context.Context, tenantID string) (int, error) {
	count, err := uc.repo.GetCount(ctx, tenantID)
	if err != nil {
		return 0, ErrDatabase.Wrap("Count", "uc.repo.GetCount", err)
	}

	return count, nil
}

func (uc *UseCase) Get(ctx context.Context, top, skip int, tenantID string) ([]dto.WirelessConfig, error) {
	data, err := uc.repo.Get(ctx, top, skip, tenantID)
	if err != nil {
		return nil, ErrDatabase.Wrap("Get", "uc.repo.Get", err)
	}

	// iterate over the data and convert each entity to dto
	d1 := make([]dto.WirelessConfig, len(data))

	for i := range data {
		tmpEntity := data[i] // create a new variable to avoid memory aliasing
		d1[i] = *uc.entityToDTO(&tmpEntity)
	}

	return d1, nil
}

func (uc *UseCase) GetByName(ctx context.Context, profileName, tenantID string) (*dto.WirelessConfig, error) {
	data, err := uc.repo.GetByName(ctx, profileName, tenantID)
	if err != nil {
		return nil, ErrDatabase.Wrap("GetByName", "uc.repo.GetByName", err)
	}

	if data == nil {
		return nil, ErrNotFound
	}

	d1 := uc.entityToDTO(data)

	return d1, nil
}

func (uc *UseCase) Delete(ctx context.Context, profileName, tenantID string) error {
	isSuccessful, err := uc.repo.Delete(ctx, profileName, tenantID)
	if err != nil {
		return ErrDatabase.Wrap("Delete", "uc.repo.Delete", err)
	}

	if !isSuccessful {
		return ErrNotFound
	}

	return nil
}

func (uc *UseCase) Update(ctx context.Context, d *dto.WirelessConfig) (*dto.WirelessConfig, error) {
	d1 := uc.dtoToEntity(d)

	_, err := uc.repo.Update(ctx, d1)
	if err != nil {
		return nil, ErrDatabase.Wrap("Update", "uc.repo.Update", err)
	}

	updatedConfig, err := uc.repo.GetByName(ctx, d1.ProfileName, d.TenantID)
	if err != nil {
		return nil, err
	}

	d2 := uc.entityToDTO(updatedConfig)

	return d2, nil
}

func (uc *UseCase) Insert(ctx context.Context, d *dto.WirelessConfig) (*dto.WirelessConfig, error) {
	d1 := uc.dtoToEntity(d)

	_, err := uc.repo.Insert(ctx, d1)
	if err != nil {
		return nil, ErrDatabase.Wrap("Insert", "uc.repo.Insert", err)
	}

	insertedConfig, err := uc.repo.GetByName(ctx, d.ProfileName, d.TenantID)
	if err != nil {
		return nil, err
	}

	d2 := uc.entityToDTO(insertedConfig)

	return d2, nil
}

// convert dto.WirelessConfig to entity.WirelessConfig.
func (uc *UseCase) dtoToEntity(d *dto.WirelessConfig) *entity.WirelessConfig {
	// convert []int to comma separated string
	linkPolicy := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(d.LinkPolicy)), ","), "[]")

	d1 := &entity.WirelessConfig{
		ProfileName:          d.ProfileName,
		AuthenticationMethod: d.AuthenticationMethod,
		EncryptionMethod:     d.EncryptionMethod,
		SSID:                 d.SSID,
		PSKValue:             d.PSKValue,
		PSKPassphrase:        d.PSKPassphrase,
		LinkPolicy:           &linkPolicy,
		TenantID:             d.TenantID,
		IEEE8021xProfileName: d.IEEE8021xProfileName,
		Version:              d.Version,
	}

	return d1
}

// convert entity.WirelessConfig to dto.WirelessConfig.
func (uc *UseCase) entityToDTO(d *entity.WirelessConfig) *dto.WirelessConfig {
	// convert comma separated string to []int
	linkPolicyInt := []int{}

	if d.LinkPolicy != nil {
		linkPolicy := strings.Split(*d.LinkPolicy, ",")
		// convert []string to []int
		intLinkPolicy := make([]int, len(linkPolicy))

		for i, v := range linkPolicy {
			val, err := strconv.Atoi(v)
			if err != nil {
				// handle the error, e.g. log or return an error
				uc.log.Error("error converting string to int")
			}

			intLinkPolicy[i] = val
		}
	}

	d1 := &dto.WirelessConfig{
		ProfileName:          d.ProfileName,
		AuthenticationMethod: d.AuthenticationMethod,
		EncryptionMethod:     d.EncryptionMethod,
		SSID:                 d.SSID,
		PSKValue:             d.PSKValue,
		PSKPassphrase:        d.PSKPassphrase,
		LinkPolicy:           linkPolicyInt,
		TenantID:             d.TenantID,
		IEEE8021xProfileName: d.IEEE8021xProfileName,
		Version:              d.Version,
	}

	return d1
}
