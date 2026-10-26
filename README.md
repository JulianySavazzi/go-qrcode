# GO QR CODE

Api em go para gerar QR Code a partir de um link.

Endpoint para checar se o servidor está rodando localmente: [http://127.0.0.1:8080/health](http://127.0.0.1:8080/health)


```
# criar container temporário para inicializar o go.mod
# isso deve ser feito antes de fazer o build (ao fazer o setup)
docker run --rm \
  -v "$PWD":/app \
  -w /app \
  golang:1.26-alpine \
  go mod init go-qrcode
  
# subir container
docker compose up  
```