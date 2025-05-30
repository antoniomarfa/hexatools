package plugin

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const TenantIDKey = "tenant_id"

func Register(db *gorm.DB) {
	db.Callback().Query().Before("gorm:query").
		Register("multitenant:query", queryCallback)

	db.Callback().Delete().Before("gorm:delete").
		Register("multitenant:delete", deleteCallback)

	db.Callback().Update().Before("gorm:update").
		Register("multitenant:update", updateCallback)

	db.Callback().Create().Before("gorm:create").
		Register("multitenant:create", createCallback)
}

func queryCallback(tx *gorm.DB) {
	if tenantID := getTenantID(tx); tenantID != "" {
		addTenantFilter(tx, tenantID)
	}
}

func deleteCallback(tx *gorm.DB) {
	if tenantID := getTenantID(tx); tenantID != "" {
		addTenantFilter(tx, tenantID)
	}
}

func updateCallback(tx *gorm.DB) {
	if tenantID := getTenantID(tx); tenantID != "" {
		addTenantFilter(tx, tenantID)
	}
}

func createCallback(tx *gorm.DB) {
	if tenantID := getTenantID(tx); tenantID != "" {
		field := tx.Statement.Schema.LookUpField("CompanyId")
		if field != nil {
			_ = field.Set(tx.Statement.Context, tx.Statement.ReflectValue, tenantID)
		}
	}
}

func getTenantID(tx *gorm.DB) string {
	if tx.Statement == nil || tx.Statement.Context == nil {
		return ""
	}
	val := tx.Statement.Context.Value(TenantIDKey)
	if s, ok := val.(string); ok {
		return s
	}
	return ""
}

func addTenantFilter(tx *gorm.DB, tenantID string) {
	_, tenantFieldExists := tx.Statement.Schema.FieldsByDBName["company_id"]
	if tenantFieldExists {
		tx.Statement.AddClause(clause.Where{
			Exprs: []clause.Expression{
				clause.Eq{
					Column: "company_id",
					Value:  tenantID,
				},
			},
		})
	}
}
