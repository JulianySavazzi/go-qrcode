```
# criar container temporário para inicializar o go.mod
# isso deve ser feito antes de fazer o build (ao fazer o setup)
docker run --rm \
  -v "$PWD":/app \
  -w /app \
  golang:1.26-alpine \
  go mod init go-qrcode
```