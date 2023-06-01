---------- Download --------------
install node.js
install angular cli (`npm install -g @angular/cli`)
install golang https://golang.org/dl/
Learn go from `https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/04.4.html`
goose migrate db (go install github.com/pressly/goose/v3/cmd/goose@latest) (go get github.com/pressly/goose/cmd/goose@v2.7.0)


---------- Rerun command to download dependencies -------- 
go mod vender
go mod tidy
----------------------
export GO111MODULE=off
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin 
export PATH="$GOPATH/bin:$PATH"


-----------------------
goose -dir db/migrations mysql "root:qweqwe@tcp(localhost:3307)/app_db?parseTime=true" up  
router go-chi (go get -u github.com/go-chi/chi/v5)
xo generate entities (https://github.com/jlightning/xo)
go get -u github.com/go-sql-driver/mysql
go get github.com/vektah/dataloaden
go get -u github.com/xo/xo
-----------------------------------------
# docker build --tag docker-base-project .
# docker run -d -p 3000:8080 --name go-server docker-base-project
# mysqldump -u root -p app_db > test.sql
xo schema 'mysql://root:qweqwe@localhost:3307/app_db?parseTime=true&columnsWithAlias=true' -o entities -e goose_db_version  --src templates/
xo schema 'mysql://root:qweqwe@localhost:3307/app_db?parseTime=true&columnsWithAlias=true' -o entities -e goose_db_version -e *.created_at -e *.updated_at --src templates/
go run github.com/99designs/gqlgen

go install github.com/google/wire/cmd/wire@latest
wire ./wiregen
// Loader avoid n+1
go run github.com/vektah/dataloaden UserLoader int *../entities.User

#For Mac with xo schema issue
unsetopt nomatch
------------------------------ Using ----------------------------------------
-- Download go lib ----
go mod tidy
go get github.com/99designs/gqlgen

-- Generate new model schema
go run github.com/99designs/gqlgen
-- Generate dependencies
wire ./wiregen ( go install github.com/google/wire/cmd/wire@latest)
-- Generate Entity database (go get -u github.com/xo/xo)
xo schema 'mysql://root:qweqwe@localhost:3307/app_db?parseTime=true&columnsWithAlias=true' -o entities -e goose_db_version  --src templates/
-- xo schema 'mysql://root:qweqwe@localhost:3307/app_db?parseTime=true&columnsWithAlias=true' -o entities -e goose_db_version -e *.created_at -e *.updated_at --src templates/

-- Goose database
goose -dir db/migrations mysql "root:qweqwe@tcp(localhost:3307)/app_db?parseTime=true" up  
