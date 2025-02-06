package tenant_repository

import (
	"context"
	"wms-service/internal/db"

	"gorm.io/gorm/clause"

	"fmt"
	"sync"
	"wms-service/internal/domain"
	error3 "wms-service/pkg/error"
	tenant_error "wms-service/pkg/error"

	"github.com/omniful/go_commons/db/sql/postgres"
	error2 "github.com/omniful/go_commons/error"
	"gorm.io/gorm"
)

type Repository struct {
	db *postgres.DbCluster
}

var repo *Repository
var repoOnce sync.Once

func NewRepository(db *postgres.DbCluster) *Repository {
	repoOnce.Do(func() {
		repo = &Repository{
			db: db,
		}
	})

	return repo
}

func (r *Repository) CreateTenant(c context.Context, tenant *domain.Tenant) error2.CustomError {
	result := r.db.GetMasterDB(c).Create(&tenant)
	if resultErr := result.Error; resultErr != nil {
		err := error2.NewCustomError(error3.SqlCreateError, fmt.Sprintf("Could not create tenant for condition  : %+v, err: %v", tenant, resultErr))
		return err
	}

	return error2.CustomError{}
}

func (r *Repository) GetTenants(
	ctx context.Context,
	condition map[string]interface{},
	countScopes []func(db *gorm.DB) *gorm.DB,
	scopes ...func(db *gorm.DB) *gorm.DB,
) (tenants []*domain.Tenant, count int64, err error2.CustomError) {
	return db.GetPaginatedDataT[domain.Tenant](ctx, r.db.GetSlaveDB(ctx), domain.Tenant{}, condition, countScopes, scopes...)
}

func (r *Repository) GetTenantByCondition(ctx context.Context, condition map[string]interface{}, scopes ...func(db *gorm.DB) *gorm.DB) (tenants *[]domain.Tenant, err error2.CustomError) {
	tenantResult := r.db.GetSlaveDB(ctx).Scopes(scopes...).Where(condition).Find(&tenants)
	if resultErr := tenantResult.Error; resultErr != nil {
		err = error2.NewCustomError(error3.SqlFetchError, fmt.Sprintf("Could not get tenants for condition: %+v, err: %v", condition, resultErr))
		return
	}

	return
}

func (r *Repository) GetTenant(
	ctx context.Context,
	condition map[string]interface{},
) (tenant domain.Tenant, cusErr error2.CustomError) {
	tenants, cusErr := r.GetTenantByCondition(ctx, condition)
	if cusErr.Exists() {
		return
	}

	if len(*tenants) == 0 {
		cusErr = error3.InvalidRequest(ctx, "TenantNotFound")
		return
	}

	return (*tenants)[0], cusErr
}

func (r *Repository) GetTenantByConditionCount(ctx context.Context, condition map[string]interface{}) (count int64, err error2.CustomError) {
	result := r.db.GetSlaveDB(ctx).Model(&domain.Tenant{}).Where(condition).Count(&count)
	if resultErr := result.Error; resultErr != nil {
		err = error2.NewCustomError(tenant_error.SqlFetchError, "Could not get tenant")
	}

	return
}

func (r *Repository) UpdateTenantByCondition(c context.Context, update map[string]interface{}, condition map[string]interface{}) (tenants []*domain.Tenant, err error2.CustomError) {
	result := r.db.GetMasterDB(c).Model(&tenants).Where(condition).Clauses(clause.Returning{}).Updates(update)
	if resultErr := result.Error; resultErr != nil {
		err = error2.NewCustomError(error3.SqlUpdateError, fmt.Sprintf("Could not update tenants for condition: %+v, err: %v", condition, resultErr))
		return
	}

	if rowsEffected := result.RowsAffected; rowsEffected == 0 {
		err = error2.NewCustomError(error3.NoRowsAffectedError, "NO ROWS UPDATED")
		return
	}

	return
}

func (r *Repository) CreateTenantPricingPlan(c context.Context, tenantSubscription *domain.TenantPricingPlan) (err error2.CustomError) {
	result := r.db.GetMasterDB(c).Create(&tenantSubscription)
	if resultErr := result.Error; resultErr != nil {
		err = error2.NewCustomError(error3.SqlCreateError, fmt.Sprintf("Could not create tenant pricing plan condition  : %+v, err: %v", tenantSubscription, resultErr))
		return
	}

	return error2.CustomError{}
}

func (r *Repository) GetTenantPricingPlans(c context.Context, tenantID uint64) (response []*domain.TenantPricingPlan, err error2.CustomError) {
	result := r.db.GetSlaveDB(c).Where("tenant_id= ?", tenantID).Find(&response)
	if resultErr := result.Error; resultErr != nil {
		err = error2.NewCustomError(error3.SqlFetchError, fmt.Sprintf("Could not get tenant pricing plan for tenantID: %+v, err: %v", tenantID, resultErr))
		return
	}
	return
}

func (r *Repository) UpdateTenantPricingPlanByCondition(c context.Context, tenantPricingPlan *domain.TenantPricingPlan, condition map[string]interface{}) (err error2.CustomError) {
	result := r.db.GetMasterDB(c).Model(&domain.TenantPricingPlan{}).Where(condition).Updates(tenantPricingPlan)
	if resultErr := result.Error; resultErr != nil {
		err = error2.NewCustomError(error3.SqlUpdateError, fmt.Sprintf("Could not update tenant pricing plan for condition  : %+v, err: %v", condition, resultErr))
		return
	}

	if rowsEffected := result.RowsAffected; rowsEffected == 0 {
		err = error2.NewCustomError(error3.NoRowsAffectedError, "NO ROWS UPDATED")
	}

	return
}

func (r *Repository) CreateTenantDomain(c context.Context, request *domain.TenantDomain) error2.CustomError {
	result := r.db.GetMasterDB(c).Create(&request)
	if resultErr := result.Error; resultErr != nil {
		err := error2.NewCustomError(error3.SqlCreateError, fmt.Sprintf("Could not create tenants domain for condition  : %+v, err: %v", request, resultErr))
		return err
	}

	return error2.CustomError{}
}

func (r *Repository) GetTenantDomainByCondition(c context.Context, conditionName map[string]interface{}) ([]*domain.TenantDomain, error2.CustomError) {
	response := make([]*domain.TenantDomain, 0)
	result := r.db.GetSlaveDB(c).Where(conditionName).Find(&response)
	if resultErr := result.Error; resultErr != nil {
		err := error2.NewCustomError(error3.SqlFetchError, fmt.Sprintf("Could not get tenants for the condition: %+v, err: %v", conditionName, resultErr))
		return nil, err
	}
	return response, error2.CustomError{}
}

func (r *Repository) CreateTenantUser(ctx context.Context, tenantUser *domain.TenantUser) (err error2.CustomError) {
	// Create tenant user
	result := r.db.GetMasterDB(ctx).Model(domain.TenantUser{}).Create(&tenantUser)
	if resultErr := result.Error; resultErr != nil {
		err = error2.NewCustomError(error3.SqlCreateError, fmt.Sprintf("Could not create tenant user for condition  : %+v, err: %v", tenantUser, resultErr))
		return
	}
	return
}

func (r *Repository) UpdateTenantUser(ctx context.Context, condition map[string]interface{}, updates map[string]interface{}) (err error2.CustomError) {
	result := r.db.GetMasterDB(ctx).Model(domain.TenantUser{}).Where(condition).Updates(updates)
	if resultErr := result.Error; resultErr != nil {
		err = error2.NewCustomError(error3.SqlUpdateError, fmt.Sprintf("Could not update tenant user for condition  : %+v, err: %v", condition, resultErr))
		return
	}

	if result.RowsAffected == 0 {
		err = error2.NewCustomError(error3.NoRowsAffectedError, "no updates made for requested tenant_user object")
		return
	}
	return
}

func (r *Repository) UpdateTenantUsers(ctx context.Context, condition map[string]interface{}, updates map[string]interface{}) (tenantUsers []*domain.TenantUser, err error2.CustomError) {
	result := r.db.GetMasterDB(ctx).Model(domain.TenantUser{}).Where(condition).Updates(updates).Scan(&tenantUsers)
	if resultErr := result.Error; resultErr != nil {
		err = error2.NewCustomError(error3.SqlUpdateError, fmt.Sprintf("Could not update tenant user for condition  : %+v, err: %v", condition, resultErr))
		return
	}

	if result.RowsAffected == 0 {
		err = error2.NewCustomError(error3.NoRowsAffectedError, "no updates made for requested tenant_user object")
		return
	}
	return
}

func (r *Repository) GetTenantUserByCondition(ctx context.Context, condition map[string]interface{}, scopes ...func(db *gorm.DB) *gorm.DB) (tenantUsers []*domain.TenantUser, err error2.CustomError) {
	result := r.db.GetMasterDB(ctx).Model(domain.TenantUser{}).Where(condition).Scopes(scopes...).Find(&tenantUsers)
	if resultErr := result.Error; resultErr != nil {
		err = error2.NewCustomError(error3.SqlFetchError, fmt.Sprintf("Could not get tenants users for condition: %+v, err: %v", condition, resultErr))
		return
	}
	return
}

func (r *Repository) GetTenantUserWithPagination(c context.Context, condition map[string]interface{}, countScopes []func(db *gorm.DB) *gorm.DB, scopes ...func(db *gorm.DB) *gorm.DB) (tenantUsers []*domain.TenantUser, count int64, cusErr error2.CustomError) {
	tenantUsers, count, cusErr = db.GetPaginatedDataT[domain.TenantUser](c, r.db.GetSlaveDB(c), domain.TenantUser{}, condition, countScopes, scopes...)
	if cusErr.Exists() {
		return nil, 0, cusErr
	}
	return
}

func (r *Repository) GetTenantUserByConditionCount(ctx context.Context, condition map[string]interface{}, scopes ...func(db *gorm.DB) *gorm.DB) (count int64, err error2.CustomError) {
	result := r.db.GetSlaveDB(ctx).Model(domain.TenantUser{}).Where(condition).Scopes(scopes...).Count(&count)
	if resultErr := result.Error; resultErr != nil {
		err = error2.NewCustomError(error3.SqlFetchError, fmt.Sprintf("Could not get tenants users for condition: %+v, err: %v", condition, resultErr))
		return
	}
	return
}

func (r *Repository) UpdateTenantUserByCondition(ctx context.Context, condition map[string]interface{}, tenantUser *domain.TenantUser) (tenantUsers []*domain.TenantUser, err error2.CustomError) {
	result := r.db.GetMasterDB(ctx).Model(&tenantUsers).Clauses(clause.Returning{}).Where(condition).Updates(&tenantUser)
	if resultErr := result.Error; resultErr != nil {
		err = error2.NewCustomError(error3.SqlUpdateError, fmt.Sprintf("Could not update tenant user for condition  : %+v, err: %v", condition, resultErr))
		return
	}
	if result.RowsAffected == 0 {
		err = error2.NewCustomError(error3.NoRowsAffectedError, "no updates made for requested tenant_user object")
		return
	}

	return
}

func (r *Repository) UpdatePickerStatus(ctx context.Context, condition map[string]interface{}, updates map[string]interface{}) (err error2.CustomError) {
	result := r.db.GetMasterDB(ctx).Model(&domain.TenantUser{}).Where(condition).Updates(&updates)
	if resultErr := result.Error; resultErr != nil {
		err = error2.NewCustomError(error3.SqlUpdateError, fmt.Sprintf("Could not update picker status for condition  : %+v, err: %v", condition, resultErr))
		return
	}

	return
}

func (r *Repository) GroupTenantUsersByRoleIDs(ctx context.Context) (map[uint64]uint64, error2.CustomError) {
	var result []struct {
		RoleID uint64
		Count  int
	}

	err := r.db.GetSlaveDB(ctx).Model(domain.TenantUser{}).
		Select("unnest(role_ids) as role_id, count(*) as count").
		Group("role_id").
		Order("count desc").
		Where("deleted_at is NULL").
		Scan(&result).Error
	if err != nil {
		cusErr := error2.NewCustomError(error3.BadRequest, err.Error())
		return nil, cusErr
	}

	resultMap := make(map[uint64]uint64)
	for _, r := range result {
		resultMap[r.RoleID] = uint64(r.Count)
	}
	return resultMap, error2.CustomError{}
}

func (r *Repository) GetAllTenantsCount(c context.Context, conditions map[string]interface{}) (count int64, err error2.CustomError) {
	result := r.db.GetSlaveDB(c).Model(domain.Tenant{}).Where(conditions).Count(&count)
	if resultErr := result.Error; resultErr != nil {
		err = error2.NewCustomError(tenant_error.SqlFetchError, "Could not get tenent")
		return
	}
	return
}
