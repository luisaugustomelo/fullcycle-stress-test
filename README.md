# fullcycle-stress-test

## 📋 Descrição

Este projeto é um testador de carga (load tester) simples feito em Go.
Ele realiza requisições HTTP simultâneas para medir a performance de um serviço web.

## Estrutura do Projeto
```
fullcycle-stress-test/
├── Dockerfile
├── main.go
├── main_test.go
├── go.mod
├── go.sum
├── Makefile
└── README.md
```


## ⚙️ Como Funciona

O programa executa múltiplas requisições concorrentes para uma URL especificada, e no final gera um relatório contendo:

- Tempo total de execução.
- Total de requisições feitas.
- Número de respostas HTTP 200 OK.
- Distribuição de outros códigos HTTP (404, 500, etc.).
- Total de erros.


## 🚀 Como Rodar

### 1. Pré-requisitos
- Go instalado (>= 1.20).
- Docker instalado (opcional, caso queira rodar via container).

### 2. Rodar Localmente (sem Docker)
Clone o repositório:

```
git clone git@github.com:luisaugustomelo/fullcycle-stress-test.git
cd fullcycle-stress-test
```

Execute o programa diretamente:
```
go run main.go --url=http://google.com --requests=100 --concurrency=10
```

### 3. Rodar usando Docker

1. Build da imagem Docker:
```
docker build -t stress-tester .
```

2. Executar o stress test via Docker:
```
docker run --rm stress-tester --url=http://google.com --requests=100 --concurrency=10
```

- --url: URL alvo para o teste.
- --requests: Número total de requisições a serem feitas.
- --concurrency: Número de requisições simultâneas.

## 🧪 Como Rodar os Testes Automatizados

Usando Makefile
Este projeto já possui um Makefile configurado.

1. Para rodar todos os testes:
```
make test
```

O que o make test faz?

Executa:
```
go test -v
```

O teste automatizado cobre:

- Respostas HTTP 200 OK.
- Respostas HTTP 404 Not Found.
- Falha de conexão (host inválido ou timeout).

## 📊 Exemplo de Saída do Relatório

📊 Load Test Report
URL: http://google.com
Total requests: 100
Concurrent requests: 10
Total time: 5.834s
Successful requests (200): 100
Other status codes:
Errors: 0


