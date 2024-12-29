all:build
build:front view app

front:
	echo "-> 删除旧构建"
	rm -rf admin/dist || true
	echo "-> 开始构建页面"
	cd admin && npm i && npm run build
	echo "-> 页面构建完成"

view:
	echo "-> 填充至APP页面中"
	rm -rf app/static/dist || true
	cp -r admin/dist app/static

.PHONY: app

app:
	echo "-> 构建APP"
	cd app && go mod tidy && go build -o ../out/local-cloud main.go
	echo "-> 完成"