package dataloader

import (
	"context"
	"fmt"
	"net/http"
	"pb-backend/entities"
	"time"
)

const loadersKey = "dataloaders"

type Loaders struct {
	UserById UserLoader
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), loadersKey, &Loaders{
			UserById: UserLoader{
				maxBatch: 100,
				wait:     1 * time.Millisecond,
				fetch: func(ids []int) ([]*entities.User, []error) {
					placeholders := make([]string, len(ids))
					args := make([]interface{}, len(ids))
					for i := 0; i < len(ids); i++ {
						placeholders[i] = "?"
						args[i] = i
					}
					var userResult []entities.User
					userResult = []entities.User{
						{
							ID:       1,
							Username: "haha"},
						{
							ID:       2,
							Username: "test user 2"},
					}
					users := make([]*entities.User, len(ids))
					fmt.Println("Find username ")
					var err error
					// err := db.Select(r.Context(), &userResult, sqrl.Expr("Select username from user where id in (?)", strings.Join(placeholders, ",")))
					if err != nil {
						panic(err)
					} else {
						userById := map[int]entities.User{}
						for _, user := range userResult {
							userById[user.ID] = user
						}
						for _, v := range userById {
							fmt.Println(v)
						}

						for i, id := range ids {
							yser := userById[id]
							users[i] = &yser
						}
					}
					return users, nil
				},
			},
		})
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}
