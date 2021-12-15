---------- Download --------------
install node.js
install angular cli (`npm install -g @angular/cli`)
install golang https://golang.org/dl/
Learn go from `https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/04.4.html`
goose migrate db (go install github.com/pressly/goose/v3/cmd/goose@latest) (go get github.com/pressly/goose/cmd/goose@v2.7.0)


---------- Rerun command to download dependencies -------- 
go mod vender
----------------------
export GO111MODULE=off
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin


-----------------------
goose -dir db/migrations mysql "root:qweqwe@tcp(localhost:3307)/app-db?parseTime=true" up  
router go-chi (go get -u github.com/go-chi/chi/v5)
xo generate entities (https://github.com/jlightning/xo)
go get -u github.com/go-sql-driver/mysql

-----------------------------------------
# docker build --tag docker-base-project .
# docker run -d -p 3000:8080 --name go-server docker-base-project
# mysqldump -u root -p app_db > test.sql