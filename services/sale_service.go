package services

import (
	"context"
	"fmt"
	"pb-backend/entities"
	"pb-backend/graph/model"
	"pb-backend/helper"

	"github.com/elgris/sqrl"
	"github.com/google/wire"
)

type ISaleService interface {
	GetOverviewUsersSales(ctx context.Context, fitler *model.OverviewUserSaleFilter, pagination *model.Pagination) ([]*model.OverviewUserSaleDto, error)
}
type SaleService struct {
	DB entities.DB
}

// define provider
var NewSaleService = wire.NewSet(wire.Struct(new(SaleService), "*"), wire.Bind(new(ISaleService), new(*SaleService)))

func (s *SaleService) GetOverviewUsersSales(ctx context.Context, fitler *model.OverviewUserSaleFilter, pagination *model.Pagination) ([]*model.OverviewUserSaleDto, error) {
	// Todo
	var users []*model.OverviewUserSaleDto
	stss := sqrl.Select(" u.username UserName, u.email UserEmail, u.role_label as UserRole, count(s.id) TotalSaledProduct, sum(s.price) EarningMoney").
		From("sale s").Join("user u on u.id = s.fk_user")
	stss.Where(sqrl.And{
		sqrl.Eq{"s.active": "true"},
		sqrl.Eq{"s.sale_status": entities.SaleStatusSaled},
	})
	s.DB.AddPagination(stss, pagination)
	if fitler != nil {
		if fitler.UserName != nil {
			stss.Where(sqrl.Eq{"username": fitler.UserName})
		}
		if fitler.DateTime != nil {
			stss.Where(sqrl.GtOrEq{"s.saled_date": helper.BeginningOfMonth(*fitler.DateTime)})
		}
	}
	stss.GroupBy("s.id, s.price")
	err := s.DB.QueryContext(ctx, &users, stss)
	// if there is an error opening the connection, handle it
	if err != nil {
		fmt.Print(err.Error())
		return nil, err
	}
	return users, nil
}
