install node.js
install angular cli
`npm install -g @angular/cli`
install golang https://golang.org/dl/
Learn go from `https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/04.4.html`
go get -u github.com/go-chi/chi/v5
go mod vender

----------------------
export GO111MODULE=off
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin


-----------------------
goose -dir migration mysql "root:@tcp(localhost:3306)/testSchema?parseTime=true" up