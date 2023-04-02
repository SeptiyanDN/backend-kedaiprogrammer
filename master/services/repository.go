package services

import (
	"errors"
	"fmt"
	"kedaiprogrammer/kedaihelpers"

	"github.com/spf13/cast"
	"gorm.io/gorm"
)

type Repository interface {
	Save(service Service) (Service, error)
	GetAllWithCounts(filterValue string, size int, page int, field, dir, filterField, filterType string) ([]map[string]interface{}, int, int, error)
	GetService(id string) (map[string]interface{}, error)
}

type repository struct {
	db  *gorm.DB
	dbs kedaihelpers.DBStruct
}

func NewRepository(db *gorm.DB, dbs kedaihelpers.DBStruct) *repository {
	return &repository{db, dbs}
}

func (r *repository) Save(service Service) (Service, error) {
	err := r.db.Create(&service).Error
	fmt.Println(err)

	if err != nil {
		return service, err
	}
	return service, nil
}

func (r *repository) GetService(id string) (map[string]interface{}, error) {
	sql := `Select 
				a.service_id,
				a.service_name,
				a.is_active,
				b.business_name
			FROM services as a 
			LEFT JOIN businesses as b on b.business_id = a.business_id
			Where a.service_id = $1
	`
	row := r.dbs.DatabaseQuerySingleRow(sql, id)
	if len(row) < 1 {
		return nil, nil
	}
	return row, nil
}
func (r *repository) GetAllWithCounts(filterValue string, size int, page int, field, dir, filterField, filterType string) ([]map[string]interface{}, int, int, error) {
	offsets := (page - 1) * size
	fmt.Println(offsets)
	fmt.Println(page)
	orderField := field
	orderDirection := dir

	queryOrder := `ORDER BY ` + cast.ToString(orderField) + ` ` + cast.ToString(orderDirection)
	queryLimit := `LIMIT ` + cast.ToString(size) + ` OFFSET ` + cast.ToString(offsets)

	queryWhere := `WHERE a.is_active = true`
	if filterValue != "" {
		var operator string
		switch filterType {
		case "eq":
			operator = "="
		// case "lt":
		// 	operator = "<"
		// case "lte":
		// 	operator = "<="
		// case "gt":
		// 	operator = ">"
		// case "gte":
		// 	operator = ">="
		case "neq":
			operator = "!="
		case "like":
			operator = "LIKE"
			filterValue = "%" + filterValue + "%"
		default:
			return nil, 0, 0, errors.New("invalid filter type")
		}
		queryWhere += ` AND ` + filterField + ` ` + operator + ` '` + filterValue + `'`
	}

	sql := `SELECT 
				a.service_id,
				a.service_name,
				a.is_active,
				b.business_name,
				(
					SELECT 
						COUNT(a2.*)  
					FROM services a2 
					WHERE a2.is_active = true
				) AS total
			FROM services as a
			left join businesses as b on b.business_id = a.business_id
			` + queryWhere + ` 
			` + queryOrder + ` ` + queryLimit

	rows := r.dbs.DatabaseQueryRows(sql)

	if len(rows) < 1 {
		return nil, 0, 0, nil
	}
	return rows, cast.ToInt(rows[0]["total"]), len(rows), nil
}
