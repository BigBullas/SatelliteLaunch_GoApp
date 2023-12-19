package repository

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"RIP_lab1/internal/models"
)

type Repository struct {
	db *gorm.DB
}

// DeleteStarById implements api.Repo.
func (*Repository) DeleteStarById(starId int) error {
	panic("unimplemented")
}

// GetStarByID implements api.Repo.
func (*Repository) GetStarByID(threatId int) (models.Product, error) {
	panic("unimplemented")
}

// GetStarsByNameFilter implements api.Repo.
func (*Repository) GetStarsByNameFilter(substring string) ([]models.Product, error) {
	panic("unimplemented")
}

func NewRepo(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Repository{
		db: db,
	}, nil
}

func (r *Repository) GetProductByID(id uint) (*models.Product, error) {
	product := &models.Product{}

	err := r.db.First(product, "id = ?", "1").Error // find product with id = 1
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *Repository) CreateProduct(product models.Product) error {
	return r.db.Create(product).Error
}
