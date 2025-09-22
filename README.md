FiscalGo API

Uma API RESTful em Go projetada para ajudar profissionais liberais e pequenos empresários a gerir as suas finanças fiscais, controlando receitas e despesas dedutíveis para a declaração do Imposto de Renda.

🎯 Sobre o ProjetoEste projeto nasceu da necessidade de profissionais, como dentistas e médicos, de terem uma ferramenta simples para organizar o seu Livro-Caixa.

A aplicação permite o registo de receitas (recibos emitidos) e despesas (notas fiscais de gastos), com o objetivo de fornecer um balanço claro para otimizar a declaração de impostos e manter a conformidade com a Receita Federal.

A arquitetura segue os princípios de software limpo, com uma clara separação de camadas (handler, service, repository) e um forte foco em testes automatizados e integração contínua.



✨ Funcionalidades PrincipaisGestão de Utilizadores: Cadastro e autenticação segura de utilizadores (profissionais).Gestão de Clientes: CRUD completo para os clientes (pacientes) de cada profissional.Registo de Receitas: Lançamento de rendimentos recebidos de clientes.Registo de Despesas: Lançamento de despesas dedutíveis, com suporte opcional para upload de comprovativos (notas fiscais).


Armazenamento de Ficheiros: Integração com o MinIO para armazenamento seguro de imagens de notas fiscais.Cache de Performance: Utilização do Redis para cache de consultas frequentes, aliviando a carga sobre o banco de dados.Busca Flexível: Endpoints de API que permitem a busca e filtragem de dados por múltiplos critérios.



🛠️ Tecnologias UtilizadasLinguagem: GoFramework Web: FiberBanco de Dados: PostgreSQLORM: GORMCache: RedisArmazenamento de Objetos: MinIOContainerização: Docker & Docker ComposeMigrations: golang-migrate/migrateTestes: testify (mock & assert)CI/CD: GitHub Actions


🚀 Como Começar (Ambiente Local)Para executar o projeto localmente, você precisará ter o Docker e o Docker Compose instalados.1. Clonar o Repositóriogit clone [https://github.com/Henrique-Rmc/fiscalgo.git](https://github.com/Henrique-Rmc/fiscalgo.git)
cd fiscalgo
2. Configurar Variáveis de AmbienteCrie um ficheiro .env na raiz do projeto, copiando o exemplo de .env.example (se existir) ou usando o modelo abaixo.# Configuração do Banco de Dados PostgreSQL
DB_HOST=db
DB_PORT=5432
DB_USER=user
DB_PASSWORD=sua_senha_segura
DB_NAME=fiscalgo_db

# (Estas variáveis são usadas pelo serviço do Postgres no Docker Compose)
POSTGRES_USER=user
POSTGRES_PASSWORD=sua_senha_segura
POSTGRES_DB=fiscalgo_db

# Configuração do MinIO
MINIO_ENDPOINT=minio:9000
MINIO_ENDPOINT_PUBLIC=http://localhost:9000 # URL para acesso externo
MINIO_ACCESS_KEY=userminio
MINIO_SECRET_KEY=sua_senha_minio_segura
BUCKET_NAME=fiscal-images

# Configuração do Redis
REDIS_ADDR=redis:6379
3. Subir a Infraestrutura com Docker Compose

Este comando irá construir a imagem da sua aplicação Go e iniciar todos os serviços (Postgres, MinIO, Redis) em segundo plano.docker compose up -d --build


A aplicação estará a ser executada em http://localhost:3000.4. 

Rodar as Migrations e o Seeder (Opcional)A aplicação está configurada para rodar as migrations automaticamente na inicialização. 

Para popular o banco de dados com dados de teste, execute o seeder:docker compose run --rm app ./main --seed
(Substitua app pelo nome do serviço da sua aplicação no docker-compose.yml, se for diferente)📖 API EndpointsA API segue os padrões RESTful. 



# Rodar os testes e gerar o relatório de cobertura
go test ./... -coverpkg=./... -coverprofile=coverage.out

# Visualizar o relatório de cobertura num browser
go tool cover -html coverage.out
🏗️ Pipeline de CI/CDEste projeto utiliza GitHub Actions para automação. O workflow, definido em .github/workflows/go-ci.yml, é acionado a cada push ou pull request para a branch main e executa as seguintes tarefas:Faz o checkout do código.Inicia os serviços (Postgres, Redis, MinIO) usando Docker Compose.Executa a suíte de testes completa.Calcula a cobertura de testes e falha o build se estiver abaixo de um limite pré-definido.
