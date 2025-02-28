package repository

import (
	"schedule_table/internal/model/dao"

	"gorm.io/gorm"
)

type tables interface {
	*dao.User | *dao.Schedule | *dao.Calendar | *dao.Employee | *dao.Task | *dao.Leave
}

type Repository[table tables] interface {
	FindOne(conds ...any) (table, error)
	FindOneWithAggregate(preload []string, conds ...any) (table, error)
	FindMany(conds ...any) ([]table, error)
	FindManyWithAggregate(preload []string, conds ...any) ([]table, error)
	FindManyPagination(offset int, limit int, conds ...any) ([]table, error)
	FindManyPaginationWithAggregate(offset int, limit int, preload []string, conds ...any) ([]table, error)
	Create(insert table) error
	CreateMany(insert []table) error
	Save(data table) error
	Update(data table) error
	UpdateMany(data []table) error
	UpdateColumn(id string, column string, value any) error
	Delete(conds ...any) error
	ForceDelete(conds ...any) error
	Count(query any, conds ...any) int64
	IsExist(query any, conds ...any) bool
}

// implement new factory create repository
type repository[table tables] struct {
	db *gorm.DB
}

func (repo *repository[table]) FindOne(conds ...any) (table, error) {
	var result table
	if err := repo.db.First(&result, conds...).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *repository[table]) FindOneWithAggregate(preload []string, conds ...any) (table, error) {
	var result table
	db := repo.db
	for _, p := range preload {
		db = db.Preload(p)
	}
	if err := db.First(&result, conds...).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *repository[table]) FindMany(conds ...any) ([]table, error) {
	var result []table

	if err := repo.db.Find(&result, conds...).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *repository[table]) FindManyWithAggregate(preload []string, conds ...any) ([]table, error) {
	var result []table
	db := repo.db
	for _, p := range preload {
		db = db.Preload(p)
	}
	if err := db.Find(&result, conds...).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *repository[table]) FindManyPagination(offset int, limit int, conds ...any) ([]table, error) {
	var result []table

	if err := repo.db.Offset(offset).Limit(limit).Find(&result, conds...).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *repository[table]) FindManyPaginationWithAggregate(offset int, limit int, preload []string, conds ...any) ([]table, error) {
	var result []table
	db := repo.db
	for _, p := range preload {
		db = db.Preload(p)
	}
	if err := db.Offset(offset).Limit(limit).Find(&result, conds...).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *repository[table]) Create(insert table) error {
	return repo.db.Create(insert).Error
}

func (repo *repository[table]) CreateMany(insert []table) error {
	return repo.db.Create(insert).Error
}

func (repo *repository[table]) Save(data table) error {
	return repo.db.Save(data).Error
}

func (repo *repository[table]) Update(data table) error {
	return repo.db.Updates(data).Error
}

func (repo *repository[table]) UpdateMany(data []table) error {
	return repo.db.Updates(data).Error
}

func (repo *repository[table]) UpdateColumn(id string, column string, value any) error {
	return repo.db.Model(new(table)).Where("id = ?", id).Update(column, value).Error
}

func (repo *repository[table]) Delete(conds ...any) error {
	return repo.db.Delete(new(table), conds...).Error
}

func (repo *repository[table]) ForceDelete(conds ...any) error {
	return repo.db.Unscoped().Delete(new(table), conds...).Error
}

func (repo *repository[table]) IsExist(query any, conds ...any) bool {
	var count int64
	if err := repo.db.Limit(1).Model(new(table)).Where(query, conds...).Count(&count).Error; err != nil {
		return false
	}
	return count > 0
}

func (repo *repository[table]) Count(query any, conds ...any) int64 {
	var count int64
	if err := repo.db.Model(new(table)).Where(query, conds...).Count(&count).Error; err != nil {
		return 0
	}
	return count
}

func NewRepository[table tables](db *gorm.DB) Repository[table] {
	return &repository[table]{db: db}
}
