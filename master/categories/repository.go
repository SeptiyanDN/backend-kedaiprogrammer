package categories

import (
	"fmt"
	"kedaiprogrammer/kedaihelpers"

	"github.com/spf13/cast"
	"gorm.io/gorm"
)

type Repository interface {
	Save(category Category) (Category, error)
	GetAllWithCounts(search string, limit, offset int, OrderColumn string, orderDirection string) ([]map[string]interface{}, int, int, error)
	GetCategory(id string) (map[string]interface{}, error)
}

type repository struct {
	db  *gorm.DB
	dbs kedaihelpers.DBStruct
}

func NewRepository(db *gorm.DB, dbs kedaihelpers.DBStruct) *repository {
	return &repository{db, dbs}
}

func (r *repository) Save(category Category) (Category, error) {
	err := r.db.Create(&category).Error
	fmt.Println(err)

	if err != nil {
		return category, err
	}
	return category, nil
}

func (r *repository) GetCategory(id string) (map[string]interface{}, error) {
	sql := `Select 
				a.category_id,
				a.category_name,
				a.slug,
				a.is_active,
				b.service_name,
				c.business_name
			FROM categories as a 
			LEFT JOIN services as b on b.service_id = a.service_id
			LEFT JOIN businesses as c on c.business_id = b.business_id
			Where a.category_id = $1
	`
	row := r.dbs.DatabaseQuerySingleRow(sql, id)
	if len(row) < 1 {
		return nil, nil
	}
	return row, nil
}

func (r *repository) GetAllWithCounts(search string, limit, offset int, OrderColumn string, orderDirection string) ([]map[string]interface{}, int, int, error) {
	offsets := (offset - 1) * limit

	queryOrder := `ORDER BY ` + cast.ToString(OrderColumn) + ` ` + cast.ToString(orderDirection)
	queryLimit := `LIMIT ` + cast.ToString(limit) + ` OFFSET ` + cast.ToString(offsets)

	queryWhere := `WHERE a.is_active = true`
	if search != "" {
		queryWhere += ` AND
		(
			a.category_name LIKE '%` + search + `%'
		)`
	}

	sql := `SELECT 
				a.category_id,
				a.category_name,
				a.slug,
				a.is_active,
				b.service_name,
				c.business_name,
				(
					SELECT 
						COUNT(a2.*)  
					FROM categories a2 
					WHERE a2.is_active = true
				) AS total
			FROM categories as a
			LEFT JOIN services as b on b.service_id = a.service_id
			LEFT JOIN businesses as c on c.business_id = b.business_id

			` + queryWhere + ` 
			` + queryOrder + ` ` + queryLimit

	rows := r.dbs.DatabaseQueryRows(sql)
	if len(rows) < 1 {
		return nil, 0, 0, nil
	}
	return rows, cast.ToInt(rows[0]["total"]), len(rows), nil
}
