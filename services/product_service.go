package services

import (
	"context"
	"fmt"
	"pb-backend/consts"
	"pb-backend/entities"
	"pb-backend/graph/model"
	"pb-backend/helper"
	"strings"

	"github.com/elgris/sqrl"
	"github.com/google/wire"
)

type IProductService interface {
	GetAllProducts(ctx context.Context, pagination *model.Pagination) ([]*model.ProductDto, error)
	CreateNewProduct(ctx context.Context, input model.NewProduct) (*model.ProductDto, error)
	DeleteProduct(ctx context.Context, productId int) (bool, error)
}
type ProductService struct {
	DB entities.DB
}

// define provider
var NewProductService = wire.NewSet(wire.Struct(new(ProductService), "*"), wire.Bind(new(IProductService), new(*ProductService)))

func (p *ProductService) GetAllProducts(ctx context.Context, pagination *model.Pagination) ([]*model.ProductDto, error) {
	claims, errParse := consts.CtxClaimValue(ctx)
	if errParse != nil {
		return nil, &model.MyError{Message: consts.ERR_USER_LOGIN_REQUIRED}
	}
	if claims.Role == consts.ROLE_USER_ADMIN || claims.Role == consts.ROLE_USER_SUPER_ADMIN {
		return p.getAllProductsForAdmin(ctx, pagination)
	}
	var result []*model.ProductDto
	var products []entities.Product
	productFitler := sqrl.Select("*").From("product p")
	productFitler.Where(sqrl.Eq{"p.active": true})
	err := p.DB.QueryContext(ctx, &products, productFitler)
	if err != nil {
		return nil, err
	}

	for _, product := range products {
		category := helper.ConvertToString(&product.Category)
		result = append(result, &model.ProductDto{
			ID:           product.ID,
			Name:         &product.Name,
			ProductKey:   product.ProductKey,
			Category:     &category,
			SellingPrice: product.SellingPrice,
			Number:       product.Quantity,
		})
	}
	return result, nil
}

func (p *ProductService) DeleteProduct(ctx context.Context, productId int) (bool, error) {
	queryBd := fmt.Sprintf("update product set product.active = false where id = %d", productId)
	p.DB.ExecContext(ctx, queryBd)
	return true, nil
}

func (p *ProductService) CreateNewProduct(ctx context.Context, input model.NewProduct) (*model.ProductDto, error) {
	keyName := strings.TrimSpace(input.Key)
	var existProduct int
	err := p.DB.QueryRowContext(ctx, &existProduct, sqrl.Select("count(*)").From("product").Where(sqrl.Eq{"product_key": keyName}))
	if err != nil {
		return nil, err
	}
	if existProduct > 0 {
		return nil, &model.MyError{Message: consts.ERR_DUPLICATE_PRODUCT_KEY}
	}
	newProduct := &entities.Product{
		Name:         strings.TrimSpace(input.Name),
		ProductKey:   keyName,
		Price:        input.Price,
		Category:     helper.ConvertToNullPointSqlString(input.Category),
		SellingPrice: input.SellingPrice,
		Quantity:     input.Number,
		Description:  helper.ConvertToNullPointSqlString(input.Description),
		ImageURL:     helper.ConvertToNullPointSqlString(input.ImageURL),
	}
	err = newProduct.Insert(ctx, p.DB)
	return &model.ProductDto{
		ID:         newProduct.ID,
		Name:       &newProduct.Name,
		ProductKey: newProduct.ProductKey,
	}, err
}

func (p *ProductService) getAllProductsForAdmin(ctx context.Context, pagination *model.Pagination) ([]*model.ProductDto, error) {
	var products []entities.Product
	productFitler := sqrl.Select("*").From("product p")
	productFitler.Where(sqrl.Eq{"p.active": true})
	err := p.DB.QueryContext(ctx, &products, productFitler)
	if err != nil {
		return nil, err
	}
	var result []*model.ProductDto
	for _, product := range products {
		category := helper.ConvertToString(&product.Category)
		result = append(result, &model.ProductDto{
			ID:           product.ID,
			Name:         helper.ConvertToPoinerString(product.Name),
			ProductKey:   product.ProductKey,
			Category:     &category,
			Price:        product.Price,
			SellingPrice: product.SellingPrice,
			Number:       product.Quantity,
		})
	}
	return result, nil
}
