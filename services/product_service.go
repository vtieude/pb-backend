package services

import (
	"context"
	"encoding/base64"
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
	GetProductDetail(ctx context.Context, id int) (*model.ProductDto, error)
	GetAllProducts(ctx context.Context, pagination *model.Pagination) ([]*model.ProductDto, error)
	CreateNewProduct(ctx context.Context, input model.ProductInputModel) (*model.ProductDto, error)
	EditProduct(ctx context.Context, input model.ProductInputModel) (bool, error)
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

func (p *ProductService) GetProductDetail(ctx context.Context, id int) (*model.ProductDto, error) {
	product, err := entities.ProductByID(ctx, p.DB, id)
	if err != nil {
		return nil, err
	}
	category := helper.ConvertToString(&product.Category)
	claims, errParse := consts.CtxClaimValue(ctx)
	if errParse != nil {
		return nil, &model.MyError{Message: consts.ERR_USER_LOGIN_REQUIRED}
	}
	price := 0.0
	if claims.Role == consts.ROLE_USER_ADMIN || claims.Role == consts.ROLE_USER_SUPER_ADMIN {
		price = product.Price
	}
	var senc *string
	var prefix *string
	if product.ImageBase64 != nil {
		convertString := base64.StdEncoding.EncodeToString(product.ImageBase64)
		senc = helper.ConvertToPoinerString(convertString)
		prefix = helper.ConvertToPoinerString(helper.ConvertToString(&product.ImagePrefix))
	}
	return &model.ProductDto{
		ID:           product.ID,
		Name:         &product.Name,
		ProductKey:   product.ProductKey,
		Category:     &category,
		SellingPrice: product.SellingPrice,
		Price:        price,
		Number:       product.Quantity,
		Description:  &product.Description.String,
		ImageBase64:  senc,
		ImagePrefix:  prefix,
	}, err
}

func (p *ProductService) DeleteProduct(ctx context.Context, productId int) (bool, error) {
	queryBd := fmt.Sprintf("update product set product.active = false where id = %d", productId)
	p.DB.ExecContext(ctx, queryBd)
	return true, nil
}

func (p *ProductService) EditProduct(ctx context.Context, input model.ProductInputModel) (bool, error) {
	if input.ID == nil {
		return false, nil
	}
	product, err := entities.ProductByID(ctx, p.DB, *input.ID)
	if err != nil {
		return false, err
	}
	product.Category = helper.ConvertToNullPointSqlString(input.Category)
	product.Name = input.Name
	product.ProductKey = input.Key
	product.SellingPrice = input.SellingPrice
	product.Quantity = input.Number
	product.Price = input.Price
	return true, product.Update(ctx, p.DB)
}

func (p *ProductService) CreateNewProduct(ctx context.Context, input model.ProductInputModel) (*model.ProductDto, error) {
	keyName := strings.TrimSpace(input.Key)
	var existProduct int
	err := p.DB.QueryRowContext(ctx, &existProduct, sqrl.Select("count(*)").From("product").Where(sqrl.Eq{"product_key": keyName}))
	if err != nil {
		return nil, err
	}
	if existProduct > 0 {
		return nil, &model.MyError{Message: consts.ERR_DUPLICATE_PRODUCT_KEY}
	}
	var blob []byte
	var imagePrefix string
	if input.ImageBase64 != nil {
		imageString := *input.ImageBase64
		imagePrefix = imageString[:strings.IndexByte(imageString, ',')+1]
		blob, err = base64.StdEncoding.DecodeString(imageString[strings.IndexByte(imageString, ',')+1:])
		if err != nil {
			fmt.Printf("Error decoding Base64 encoded data %v", err)
		}
	} else {
		imagePrefix = ""
	}
	newProduct := &entities.Product{
		Name:         strings.TrimSpace(input.Name),
		ProductKey:   keyName,
		Price:        input.Price,
		Category:     helper.ConvertToNullPointSqlString(input.Category),
		SellingPrice: input.SellingPrice,
		Quantity:     input.Number,
		Description:  helper.ConvertToNullPointSqlString(input.Description),
		ImagePrefix:  helper.ConvertToNullPointSqlString(&imagePrefix),
		ImageBase64:  blob,
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
