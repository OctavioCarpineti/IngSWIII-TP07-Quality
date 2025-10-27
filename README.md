# TP07 - Quality Assurance: Manual de Uso

**Alumno:** Octavio Carpineti - Kevin Massholder 
**Materia:** Ingeniería de Software III  
**Fecha:** Octubre 2025

---

## 📋 Tabla de Contenidos

1. [Descripción del Proyecto](#descripción-del-proyecto)
2. [Requisitos Previos](#requisitos-previos)
3. [Instalación](#instalación)
4. [Ejecución del Proyecto](#ejecución-del-proyecto)
5. [Ejecución de Tests](#ejecución-de-tests)
6. [Herramientas de Calidad](#herramientas-de-calidad)
7. [Pipeline CI/CD](#pipeline-cicd)
8. [Estructura del Proyecto](#estructura-del-proyecto)

---

## 📖 Descripción del Proyecto

Mini red social desarrollada con React (frontend) y Go (backend) que implementa:

- Registro y autenticación de usuarios
- Creación, visualización y eliminación de posts
- Sistema de comentarios en posts
- Validaciones de permisos (solo el autor puede eliminar su contenido)

**Stack Tecnológico:**
- **Backend:** Go 1.24 + SQLite
- **Frontend:** React 18 + TypeScript
- **Testing:** Go testing + Jest + Cypress
- **CI/CD:** GitHub Actions + SonarCloud

---

## 🔧 Requisitos Previos

### Software Necesario

```bash
# Verificar versiones instaladas:
go version    # Debe ser 1.24 o superior
node --version # Debe ser 20 o superior
npm --version  # Debe ser 10 o superior
```

### Instalación de Dependencias (si no las tenés)

**Go:**
```bash
# macOS
brew install go

# Ubuntu/Debian
sudo apt install golang-go

# Windows
# Descargar desde: https://go.dev/dl/
```

**Node.js y npm:**
```bash
# macOS
brew install node

# Ubuntu/Debian
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt-get install -y nodejs

# Windows
# Descargar desde: https://nodejs.org/
```

---

## 📥 Instalación

### 1. Clonar el Repositorio

```bash
git clone https://github.com/OctavioCarpineti/IngSWIII-TP07-Quality.git
cd IngSWIII-TP07-Quality
```

### 2. Instalar Dependencias del Backend

```bash
cd backend
go mod download
cd ..
```

### 3. Instalar Dependencias del Frontend

```bash
cd frontend
npm install
cd ..
```

---

## 🚀 Ejecución del Proyecto

### Opción 1: Ejecución Manual (Recomendado para desarrollo)

**Terminal 1 - Backend:**
```bash
cd backend
go run cmd/api/main.go
```

El backend estará corriendo en `http://localhost:8080`

Deberías ver:
```
🚀 Servidor corriendo en http://localhost:8080
📊 Base de datos inicializada
```

**Terminal 2 - Frontend:**
```bash
cd frontend
npm start
```

El frontend estará corriendo en `http://localhost:3000`

Se abrirá automáticamente en tu navegador.

### Opción 2: Ejecución con Scripts

**Backend:**
```bash
cd backend
# Compilar
go build -o app cmd/api/main.go

# Ejecutar
./app
```

**Frontend:**
```bash
cd frontend
# Build de producción
npm run build

# Servir build (requiere serve instalado: npm install -g serve)
serve -s build -l 3000
```

---

## 🧪 Ejecución de Tests

### Tests Unitarios - Backend

```bash
cd backend

# Ejecutar todos los tests
go test ./tests/services/... -v

# Ejecutar tests con coverage
go test ./tests/services/... -v -cover -coverpkg=./internal/services/...

# Generar reporte HTML de coverage
go test ./tests/services/... -coverprofile=coverage.out -coverpkg=./internal/services/...
go tool cover -html=coverage.out

# Ver coverage en terminal
go tool cover -func=coverage.out
```

**Resultado esperado:**
```
=== RUN   TestRegister_Success
--- PASS: TestRegister_Success (0.00s)
...
PASS
coverage: 86.5% of statements in ./internal/services
ok      tp06-testing/tests/services     0.537s
```

### Tests Unitarios - Frontend

```bash
cd frontend

# Ejecutar tests en modo watch
npm test

# Ejecutar tests una vez
npm test -- --watchAll=false

# Ejecutar tests con coverage
npm test -- --coverage --watchAll=false

# Ver reporte de coverage en navegador
open coverage/lcov-report/index.html
```

**Resultado esperado:**
```
Test Suites: 8 passed, 8 total
Tests:       39 passed, 39 total
Coverage:    92.44% statements
```

### Tests E2E - Cypress

**Prerequisito: Backend y Frontend deben estar corriendo**

```bash
# Terminal 1: Backend
cd backend
go run cmd/api/main.go

# Terminal 2: Frontend  
cd frontend
npm start

# Terminal 3: Cypress
cd frontend

# Modo interactivo (recomendado)
npx cypress open
# Luego click en "E2E Testing" y seleccionar los tests

# Modo headless (para CI/CD)
npx cypress run
```

**Resultado esperado:**
```
Running:  auth.cy.js                    (1 of 4)
  ✓ 5 tests passing

Running:  posts.cy.js                   (2 of 4)
  ✓ 5 tests passing

Running:  comments.cy.js                (3 of 4)
  ✓ 4 tests passing

Running:  full-flow.cy.js               (4 of 4)
  ✓ 1 test passing

Total: 15 tests passing
```

---

## 🔍 Herramientas de Calidad

### 1. SonarCloud (Análisis Estático)

**Acceso al proyecto:**
```
URL: https://sonarcloud.io/project/overview?id=OctavioCarpineti_IngSWIII-TP07-Quality
Organization: octaviocarpineti
```

**Análisis local (opcional):**
```bash
# Requiere configuración de SONAR_TOKEN
docker run --rm \
  -e SONAR_HOST_URL="https://sonarcloud.io" \
  -e SONAR_TOKEN="tu-token" \
  -v "$(pwd):/usr/src" \
  sonarsource/sonar-scanner-cli
```

### 2. Code Coverage

**Backend:**
```bash
cd backend
go test ./tests/services/... -coverprofile=coverage.out -coverpkg=./internal/services/...

# Ver en terminal
go tool cover -func=coverage.out | grep total

# Ver en navegador
go tool cover -html=coverage.out
```

**Frontend:**
```bash
cd frontend
npm test -- --coverage --watchAll=false

# Abrir reporte HTML
open coverage/lcov-report/index.html
```

---

## 🔄 Pipeline CI/CD

### GitHub Actions

El proyecto está configurado con GitHub Actions que se ejecuta automáticamente en cada push o pull request.

**Pipeline incluye:**
1. ✅ Backend Tests (Go) + Coverage
2. ✅ Frontend Tests (React) + Coverage
3. ✅ SonarCloud Analysis
4. ✅ Cypress E2E Tests
5. ✅ Backend Build
6. ✅ Frontend Build
7. ✅ Quality Gate Summary

**Ver estado del pipeline:**
```
GitHub > Actions > CI/CD Pipeline
```

**Ejecutar pipeline manualmente:**
```bash
git commit --allow-empty -m "trigger pipeline"
git push
```

### Quality Gates Configurados

El pipeline **bloqueará** el merge si:
- ❌ Backend coverage < 70%
- ❌ Frontend coverage < 70%
- ❌ SonarCloud Quality Gate falla
- ❌ Cualquier test unitario falla
- ❌ Cualquier test E2E falla

---

## 📁 Estructura del Proyecto

```
tp07-quality/
├── backend/
│   ├── cmd/
│   │   └── api/
│   │       └── main.go              # Entry point del servidor
│   ├── internal/
│   │   ├── handlers/                # HTTP handlers (POST, GET, DELETE)
│   │   │   ├── auth_handler.go
│   │   │   ├── post_handler.go
│   │   │   └── utils.go
│   │   ├── services/                # Lógica de negocio (86.5% coverage)
│   │   │   ├── auth_service.go
│   │   │   └── post_service.go
│   │   ├── repository/              # Acceso a datos (interfaz)
│   │   │   ├── user_repository.go
│   │   │   └── post_repository.go
│   │   ├── models/                  # Estructuras de datos
│   │   │   ├── users.go
│   │   │   └── post.go
│   │   ├── database/                # Configuración BD
│   │   │   └── database.go
│   │   └── router/                  # Configuración de rutas
│   │       └── router.go
│   ├── tests/
│   │   ├── services/                # 35 tests unitarios
│   │   │   ├── auth_service_test.go
│   │   │   └── post_service_test.go
│   │   └── mocks/                   # Mocks para testing
│   │       ├── mock_user_repository.go
│   │       └── mock_post_repository.go
│   ├── go.mod
│   └── go.sum
│
├── frontend/
│   ├── src/
│   │   ├── components/              # Componentes React (92.44% coverage)
│   │   │   ├── Login/
│   │   │   │   ├── Login.tsx
│   │   │   │   ├── Login.test.tsx
│   │   │   │   └── Login.css
│   │   │   ├── PostList/
│   │   │   │   ├── PostList.tsx
│   │   │   │   ├── PostList.test.tsx
│   │   │   │   └── PostList.css
│   │   │   ├── CreatePost/
│   │   │   ├── PostDetail/
│   │   │   ├── CommentList/
│   │   │   └── CommentForm/
│   │   ├── services/                # Servicios HTTP
│   │   │   ├── authService.ts
│   │   │   ├── authService.test.ts
│   │   │   ├── postService.ts
│   │   │   └── postService.test.ts
│   │   ├── types/                   # Definiciones TypeScript
│   │   │   └── index.ts
│   │   ├── App.tsx                  # Componente principal
│   │   └── index.tsx                # Entry point
│   ├── cypress/
│   │   ├── e2e/
│   │   │   └── blog/                # 15 tests E2E
│   │   │       ├── auth.cy.js       # 5 tests
│   │   │       ├── posts.cy.js      # 5 tests
│   │   │       ├── comments.cy.js   # 4 tests
│   │   │       └── full-flow.cy.js  # 1 test
│   │   └── support/
│   │       ├── e2e.js
│   │       └── commands.js
│   ├── cypress.config.js
│   ├── package.json
│   └── package-lock.json
│
├── .github/
│   └── workflows/
│       └── ci.yml                   # Pipeline CI/CD
│
├── sonar-project.properties         # Configuración SonarCloud
├── README.md                        # Este archivo
└── decisiones.md                    # Decisiones técnicas y justificaciones
```

---

## 🐛 Troubleshooting

### Backend no inicia

```bash
# Verificar puerto 8080 disponible
lsof -i :8080
# Si está ocupado, matar el proceso:
kill -9 <PID>

# Verificar Go instalado correctamente
go version

# Limpiar y reinstalar dependencias
cd backend
rm go.sum
go mod tidy
go mod download
```

### Frontend no inicia

```bash
# Verificar puerto 3000 disponible
lsof -i :3000

# Limpiar cache y reinstalar
cd frontend
rm -rf node_modules package-lock.json
npm install

# Si falla con errores de Cypress
npm install --save-dev cypress@13.15.2
```

### Tests de Cypress fallan

```bash
# Verificar que backend y frontend estén corriendo
curl http://localhost:8080/api/health
curl http://localhost:3000

# Limpiar cache de Cypress
npx cypress cache clear
npx cypress install

# Ejecutar con logs detallados
DEBUG=cypress:* npx cypress run
```

### Pipeline falla en GitHub Actions

```bash
# Verificar logs en:
# GitHub > Actions > Click en el run fallido

# Causas comunes:
# 1. package-lock.json desincronizado
cd frontend
rm package-lock.json
npm install
git add package-lock.json
git commit -m "fix: regenerar package-lock.json"
git push

# 2. Tests fallan localmente primero
# Ejecutar todos los tests localmente antes de push
```

---

## 📊 Métricas Alcanzadas

| Métrica | Objetivo | Resultado | Estado |
|---------|----------|-----------|--------|
| Backend Coverage | ≥70% | 86.5% | ✅ |
| Frontend Coverage | ≥70% | 92.44% | ✅ |
| Total Tests | - | 89 tests | ✅ |
| SonarCloud Quality Gate | Pass | PASSED | ✅ |
| Issues Resueltos | ≥3 | 47 issues | ✅ |
| Duplications | <3% | 0.0% | ✅ |

---



**Alumno:** Octavio Carpineti - Kevin Massholder 
**GitHub:** https://github.com/OctavioCarpineti  
**Repositorio:** https://github.com/OctavioCarpineti/IngSWIII-TP07-Quality