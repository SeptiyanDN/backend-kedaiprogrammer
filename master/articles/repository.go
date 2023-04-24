package articles

import (
	"kedaiprogrammer/kedaihelpers"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type Repository interface {
	GetAllWithCounts(tag, search string, limit, offset int, OrderColumn string, orderDirection string) ([]map[string]interface{}, int, int, error)
	Save(article Article) (Article, error)
	GetOne(articleID string) (map[string]interface{}, error)
}

type repository struct {
	db  *gorm.DB
	dbs kedaihelpers.DBStruct
}

func NewRepository(db *gorm.DB, dbs kedaihelpers.DBStruct) *repository {
	return &repository{db, dbs}
}

func (r *repository) GetAllWithCounts(tag, search string, limit, offset int, OrderColumn string, orderDirection string) ([]map[string]interface{}, int, int, error) {

	offsets := (offset - 1) * limit

	queryOrder := `ORDER BY ` + cast.ToString(OrderColumn) + ` ` + cast.ToString(orderDirection)
	queryLimit := ``

	if limit != 5 {
		queryLimit += `LIMIT ` + cast.ToString(limit) + ` OFFSET ` + cast.ToString(offsets)

	}

	queryWhere := `WHERE a.status = 1`
	if tag != "" {
		queryWhere += ` AND b.tag LIKE '%` + tag + `%' `
	}
	if search != "" {
		queryWhere += ` AND
		(
			a.title LIKE '%` + search + `%' OR
			a.description LIKE '%` + search + `%' OR
			a.body LIKE '%` + search + `%'
		)`
	}

	sql := `SELECT 
				a.article_id,
				a.body,
				a.created_at,
				a.description,
				a.main_image,
				a.publised_at,
				a.slug,
				a.title,
				a.updated_at,
				b.category_name,
				b.tag,
				c.username as author_name,
				count(*) OVER() AS total
			FROM articles as a
			LEFT JOIN categories as b on b.category_id = a.category_id
			LEFT JOIN users as c on c.uuid = a.author_id

			` + queryWhere + ` 
			` + queryOrder + ` ` + queryLimit
	rows := r.dbs.DatabaseQueryRows(sql)
	if len(rows) < 1 {
		return []map[string]interface{}{}, 0, 0, nil
	}
	baseURLS3 := viper.GetString("S3_CREDENTIALS.RESP_URL") + "/" + viper.GetString("S3_CREDENTIALS.BUCKET") + "/articles/"

	datas := []map[string]interface{}{}

	for _, v := range rows {
		v["main_image"] = baseURLS3 + cast.ToString(v["main_image"])
		datas = append(datas, v)
	}
	return datas, cast.ToInt(rows[0]["total"]), len(rows), nil
}

func (r *repository) Save(article Article) (Article, error) {
	err := r.db.Create(&article).Error
	if err != nil {
		return article, err
	}
	return article, nil
}
func (r *repository) GetOne(articleID string) (map[string]interface{}, error) {
	sql := `SELECT 
				a.article_id,
				a.body,
				a.created_at,
				a.description,
				a.main_image,
				a.publised_at,
				a.slug,
				a.title,
				a.updated_at,
				b.category_name,
				b.tag,
				c.username as author_name,
				count(*) OVER() AS total

			FROM articles as a
			LEFT JOIN categories as b on b.category_id = a.category_id
			LEFT JOIN users as c on c.uuid = a.author_id
			WHERE a.article_id ='` + articleID + `'`

	rows := r.dbs.DatabaseQuerySingleRow(sql)
	if len(rows) < 1 {
		return map[string]interface{}{}, nil
	}
	baseURLS3 := viper.GetString("S3_CREDENTIALS.RESP_URL") + "/" + viper.GetString("S3_CREDENTIALS.BUCKET") + "/articles/"

	rows["main_image"] = baseURLS3 + cast.ToString(rows["main_image"])
	return rows, nil

}
