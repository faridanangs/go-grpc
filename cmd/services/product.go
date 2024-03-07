package services

import (
	"context"
	"go_grpc_yt/cmd/helpers"
	"go_grpc_yt/pb/pagination"
	prodPB "go_grpc_yt/pb/product"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type ProductService struct {
	prodPB.UnimplementedProductServiceServer
	DB *gorm.DB
}

func (p *ProductService) GetProducts(ctx context.Context, param *prodPB.Page) (*prodPB.Products, error) {
	var products []*prodPB.Product
	var pagination pagination.Pagination

	var page int64 = 1
	if param.GetPage() != 0 {
		page = param.GetPage()
	}

	sql := p.DB.Table("products AS p").
		Joins("LEFT JOIN categories AS c ON c.id = p.category_id").
		Select("p.id, p.name, p.price, p.stock, c.id AS category_id, c.name AS category_name")

	offset, limit := helpers.Pagination(sql, &pagination, page)

	rows, err := sql.Offset(int(offset)).Limit(int(limit)).Rows()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var product prodPB.Product
		var category prodPB.Category

		rows.Scan(&product.Id, &product.Name, &product.Price, &product.Stock, &category.Id, &category.Name)

		product.Category = &category
		products = append(products, &product)
	}

	response := prodPB.Products{
		Pagination: &pagination,
		Data:       products,
	}
	return &response, nil
}

func (p *ProductService) GetProduct(ctx context.Context, id *prodPB.ID) (*prodPB.Product, error) {
	var product prodPB.Product
	var category prodPB.Category

	row := p.DB.Table("products AS p").
		Joins("LEFT JOIN categories AS c ON c.id = p.category_id").
		Select("p.id, p.name, p.price, p.stock, c.id AS category_id, c.name AS category_name").
		Where("p.id = ?", id.GetId()).
		Row()

	if err := row.Scan(&product.Id, &product.Name, &product.Price, &product.Stock, &category.Id, &category.Name); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	product.Category = &category

	return &product, nil
}

func (p *ProductService) CreateProduct(ctx context.Context, prod *prodPB.Product) (*prodPB.ID, error) {
	var response prodPB.ID

	if err := p.DB.Transaction(func(tx *gorm.DB) error {
		category := &prodPB.Category{
			Id:   0,
			Name: prod.GetCategory().GetName(),
		}
		if err := tx.Debug().Table("categories").
			Where("name = ?", prod.GetName()).
			FirstOrCreate(&category).Error; err != nil {
			return err
		}

		product := struct {
			ID         uint64
			Name       string
			Price      float64
			Stock      uint32
			CategoryID uint64
		}{
			ID:         prod.GetId(),
			Name:       prod.GetName(),
			Price:      prod.GetPrice(),
			Stock:      prod.GetStock(),
			CategoryID: category.GetId(),
		}

		if err := tx.Debug().Table("products").Create(&product).Error; err != nil {
			return err
		}

		response.Id = product.ID
		return nil
	}); err != nil {
		return nil, err
	}

	return &response, nil
}

func (p *ProductService) DeleteProduct(ctx context.Context, id *prodPB.ID) (*prodPB.Status, error) {
	var response prodPB.Status

	if err := p.DB.Debug().Table("products").Where("id = ?", id.GetId()).Delete(nil).Error; err != nil {
		return nil, err
	}

	response.Status = 1

	return &response, nil
}

func (p *ProductService) UpdateProduct(ctx context.Context, prod *prodPB.Product) (*prodPB.Status, error) {
	var response prodPB.Status

	if err := p.DB.Transaction(func(tx *gorm.DB) error {

		category := &prodPB.Category{
			Id:   0,
			Name: prod.Category.GetName(),
		}

		if err := tx.Debug().Table("categories").
			Where("name = ?", category.GetName()).
			FirstOrCreate(&category).Error; err != nil {
			return err
		}

		product := struct {
			ID         uint64
			Name       string
			Price      float64
			Stock      uint32
			CategoryID uint64
		}{
			ID:         prod.GetId(),
			Name:       prod.GetName(),
			Price:      prod.GetPrice(),
			Stock:      prod.GetStock(),
			CategoryID: category.GetId(),
		}
		if err := tx.Debug().Table("products").Where("id = ?", prod.GetId()).Updates(&product).Error; err != nil {
			return err
		}

		response.Status = 1
		return nil

	}); err != nil {
		return nil, err
	}
	return &response, nil
}
