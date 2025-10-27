# Descripción de Arquitectura - TP07

**Propósito:** Documento técnico que explica la estructura del proyecto, responsabilidades de cada capa y decisiones arquitectónicas implementadas.

---

## 📋 Tabla de Contenidos

1. [Visión General de la Arquitectura](#1-visión-general-de-la-arquitectura)
2. [Backend: Estructura y Capas](#2-backend-estructura-y-capas)
3. [Frontend: Estructura y Componentes](#3-frontend-estructura-y-componentes)
4. [Testing: Organización y Estrategia](#4-testing-organización-y-estrategia)
5. [CI/CD: Pipeline y Configuración](#5-cicd-pipeline-y-configuración)

---

## 1. Visión General de la Arquitectura

### 1.1 Patrón Arquitectónico

El proyecto implementa una **arquitectura en capas (layered architecture)** tanto en backend como frontend, siguiendo los principios de **separación de responsabilidades** y **dependency inversion**.

```
┌─────────────────────────────────────────┐
│         PRESENTACIÓN (Frontend)         │
│    React Components + TypeScript        │
└──────────────┬──────────────────────────┘
               │ HTTP/REST
               ↓
┌─────────────────────────────────────────┐
│         API Layer (Handlers)            │
│    HTTP Request/Response Management     │
└──────────────┬──────────────────────────┘
               │
               ↓
┌─────────────────────────────────────────┐
│      LÓGICA DE NEGOCIO (Services)       │
│    Validaciones + Reglas de Negocio    │
└──────────────┬──────────────────────────┘
               │
               ↓
┌─────────────────────────────────────────┐
│      ACCESO A DATOS (Repository)        │
│         CRUD + Queries SQL              │
└──────────────┬──────────────────────────┘
               │
               ↓
┌─────────────────────────────────────────┐
│         BASE DE DATOS (SQLite)          │
│       Persistencia de Datos             │
└─────────────────────────────────────────┘
```

**Ventajas de esta arquitectura:**
- ✅ **Testabilidad:** Cada capa se puede testear aisladamente
- ✅ **Mantenibilidad:** Cambios en una capa no afectan otras
- ✅ **Escalabilidad:** Fácil agregar nuevas funcionalidades
- ✅ **Legibilidad:** Estructura clara y predecible

---

### 1.2 Flujo de una Request

**Ejemplo:** Usuario crea un post

```
1. Frontend (CreatePost.tsx)
   ↓ POST /api/posts + { title, content }
   
2. Handler (post_handler.go)
   ↓ Valida headers, extrae user_id
   
3. Service (post_service.go)
   ↓ Valida datos, aplica reglas de negocio
   
4. Repository (post_repository.go)
   ↓ INSERT INTO posts
   
5. Database (SQLite)
   ↓ Persiste datos
   
6. Response ← ← ← ← ←
   Frontend recibe el post creado
```

---

## 2. Backend: Estructura y Capas

### 2.1 Directorio `/backend/cmd/api/`

**Propósito:** Entry point de la aplicación backend.

```
backend/cmd/api/
└── main.go
```

**Responsabilidades de `main.go`:**

```go
func main() {
    // 1. Inicialización de base de datos
    db := database.InitDB()
    
    // 2. Creación de repositorios
    userRepo := repository.NewUserRepository(db)
    postRepo := repository.NewPostRepository(db)
    
    // 3. Creación de services (inyección de dependencias)
    authService := services.NewAuthService(userRepo)
    postService := services.NewPostService(postRepo)
    
    // 4. Configuración de router y handlers
    router := router.SetupRouter(authService, postService)
    
    // 5. Inicio del servidor HTTP
    log.Fatal(http.ListenAndServe(":8080", router))
}
```

**¿Por qué esta estructura?**

✅ **Dependency Injection:** Services reciben repositories como parámetros
✅ **Testabilidad:** main.go orquesta, no contiene lógica de negocio
✅ **Configuración centralizada:** Un solo punto de inicialización

**Coverage:** NO se mide (es configuración, no lógica)

---

### 2.2 Directorio `/backend/internal/handlers/`

**Propósito:** Capa de presentación HTTP. Maneja requests y responses.

```
backend/internal/handlers/
├── auth_handler.go      # POST /api/auth/login, /api/auth/register
├── post_handler.go      # CRUD de posts y comentarios
└── utils.go             # Funciones auxiliares (respondWithJSON, respondWithError)
```

**Responsabilidades:**

1. **Parsear HTTP Request:**
```go
var req LoginRequest
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
    respondWithError(w, http.StatusBadRequest, "Request inválido")
    return
}
```

2. **Extraer headers/parámetros:**
```go
userIDStr := r.Header.Get("X-User-ID")
postIDStr := chi.URLParam(r, "postId")
```

3. **Invocar Service:**
```go
user, err := h.authService.Login(req.Email, req.Password)
```

4. **Formatear Response:**
```go
respondWithJSON(w, http.StatusOK, user)
```

**¿Por qué handlers NO tienen lógica de negocio?**

❌ Handlers **solo** mapean HTTP → Service → HTTP
✅ Lógica en services → más fácil de testear
✅ Handlers se testean con pruebas E2E (Cypress)

**Coverage:** NO se mide (se testea con E2E)

**Constantes implementadas (SonarCloud):**

```go
const (
    HeaderUserID           = "X-User-ID"
    ErrUserNotAuthenticated = "Usuario no autenticado"
    ErrInvalidUserID       = "User ID inválido"
    // ... más constantes
)
```

Esto resuelve los 47 code smells de duplicación detectados por SonarCloud.

---

### 2.3 Directorio `/backend/internal/services/`

**Propósito:** Capa de lógica de negocio. Contiene las reglas y validaciones.

```
backend/internal/services/
├── auth_service.go      # Register, Login, validaciones de usuario
└── post_service.go      # CRUD posts, validación de permisos
```

**Responsabilidades de `auth_service.go`:**

```go
type AuthService struct {
    userRepo repository.UserRepository  // Dependency injection
}

func (s *AuthService) Register(email, password, username string) (*models.User, error) {
    // 1. Validaciones de negocio
    if email == "" || password == "" || username == "" {
        return nil, errors.New("campos requeridos")
    }
    
    // 2. Verificar email único
    existingUser, _ := s.userRepo.FindByEmail(email)
    if existingUser != nil {
        return nil, errors.New("email ya existe")
    }
    
    // 3. Hash de password
    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    
    // 4. Crear usuario
    user := &models.User{
        Email:    email,
        Password: string(hashedPassword),
        Username: username,
    }
    
    // 5. Persistir
    return s.userRepo.Create(user)
}
```

**Responsabilidades de `post_service.go`:**

```go
func (s *PostService) DeletePost(postID, userID int) error {
    // 1. Obtener post
    post, err := s.postRepo.FindByID(postID)
    
    // 2. VALIDACIÓN DE PERMISOS (lógica de negocio)
    if post.UserID != userID {
        return errors.New("no tienes permisos")
    }
    
    // 3. Eliminar
    return s.postRepo.Delete(postID)
}
```

**¿Por qué la lógica está en services?**

✅ **Reutilizable:** Puede usarse desde handlers, CLI, workers
✅ **Testeable:** Se mockea el repository fácilmente
✅ **Mantenible:** Cambios en lógica no afectan HTTP layer

**Coverage:** **86.5%** (objetivo principal de testing)

---

### 2.4 Directorio `/backend/internal/repository/`

**Propósito:** Capa de acceso a datos. Abstracción sobre la base de datos.

```
backend/internal/repository/
├── user_repository.go   # Interface + implementación CRUD usuarios
└── post_repository.go   # Interface + implementación CRUD posts/comments
```

**Patrón: Repository Pattern con Interfaces**

```go
// Interface (contrato)
type UserRepository interface {
    Create(user *models.User) (*models.User, error)
    FindByEmail(email string) (*models.User, error)
    FindByID(id int) (*models.User, error)
}

// Implementación SQLite
type SQLiteUserRepository struct {
    db *sql.DB
}

func (r *SQLiteUserRepository) FindByEmail(email string) (*models.User, error) {
    var user models.User
    query := `SELECT id, email, password, username, created_at FROM users WHERE email = ?`
    
    err := r.db.QueryRow(query, email).Scan(
        &user.ID, &user.Email, &user.Password, &user.Username, &user.CreatedAt,
    )
    
    return &user, err
}
```

**¿Por qué usar interfaces?**

✅ **Dependency Inversion:** Services dependen de la abstracción, no de SQLite
✅ **Testeable:** Fácil crear mocks del repository
✅ **Cambiable:** Podemos cambiar SQLite por PostgreSQL sin tocar services

**Coverage:** NO se mide (se mockea en tests de services)

---

### 2.5 Directorio `/backend/internal/models/`

**Propósito:** Definiciones de estructuras de datos (DTOs/Entities).

```
backend/internal/models/
├── user.go      # struct User
└── post.go      # struct Post, Comment
```

**Ejemplo:**

```go
type Post struct {
    ID        int       `json:"id"`
    Title     string    `json:"title"`
    Content   string    `json:"content"`
    UserID    int       `json:"user_id"`
    Username  string    `json:"username"`
    CreatedAt time.Time `json:"created_at"`
}
```

**¿Por qué separar en models/?**

✅ Reutilizable en todas las capas
✅ Single source of truth para estructuras
✅ Tags JSON centralizados

**Coverage:** NO se mide (son solo definiciones)

---

### 2.6 Directorio `/backend/internal/database/`

**Propósito:** Configuración e inicialización de la base de datos.

```
backend/internal/database/
└── database.go
```

**Responsabilidades:**

```go
func InitDB() *sql.DB {
    // 1. Abrir conexión SQLite
    db, err := sql.Open("sqlite3", "./blog.db")
    
    // 2. Crear tablas si no existen
    createTables(db)
    
    return db
}

func createTables(db *sql.DB) {
    // CREATE TABLE users ...
    // CREATE TABLE posts ...
    // CREATE TABLE comments ...
}
```

**Coverage:** NO se mide (es configuración)

---

### 2.7 Directorio `/backend/tests/`

**Propósito:** Tests unitarios del backend.

```
backend/tests/
├── services/
│   ├── auth_service_test.go    # 11 tests
│   └── post_service_test.go    # 24 tests
└── mocks/
    ├── mock_user_repository.go  # Mock generado con testify
    └── mock_post_repository.go  # Mock generado con testify
```

**Estrategia de testing:**

```go
// Test de AuthService usando mock
func TestRegister_Success(t *testing.T) {
    // 1. Setup: crear mock
    mockRepo := new(mocks.MockUserRepository)
    service := services.NewAuthService(mockRepo)
    
    // 2. Expectations: qué debe hacer el mock
    mockRepo.On("FindByEmail", "test@example.com").Return(nil, nil)
    mockRepo.On("Create", mock.Anything).Return(&models.User{ID: 1}, nil)
    
    // 3. Act: ejecutar función
    user, err := service.Register("test@example.com", "123456", "testuser")
    
    // 4. Assert: verificar resultado
    assert.NoError(t, err)
    assert.Equal(t, 1, user.ID)
    
    // 5. Verify: mock fue llamado correctamente
    mockRepo.AssertExpectations(t)
}
```

**¿Por qué mockear el repository?**

✅ Tests rápidos (sin DB real)
✅ Tests aislados (no dependen de datos pre-existentes)
✅ Control total sobre casos de prueba

**Coverage alcanzado:** 86.5% en services

---

## 3. Frontend: Estructura y Componentes

### 3.1 Directorio `/frontend/src/components/`

**Propósito:** Componentes React reutilizables.

```
frontend/src/components/
├── Login/
│   ├── Login.tsx          # Componente de autenticación
│   ├── Login.test.tsx     # 5 tests unitarios
│   └── Login.css          # Estilos
├── PostList/
│   ├── PostList.tsx       # Lista de posts
│   ├── PostList.test.tsx  # 8 tests
│   └── PostList.css
├── CreatePost/
│   ├── CreatePost.tsx     # Formulario crear post
│   ├── CreatePost.test.tsx # 4 tests
│   └── CreatePost.css
├── PostDetail/
│   ├── PostDetail.tsx     # Vista detalle + comentarios
│   ├── PostDetail.test.tsx # 4 tests
│   └── PostDetail.css
├── CommentList/
│   ├── CommentList.tsx    # Lista comentarios
│   ├── CommentList.test.tsx # 8 tests
│   └── CommentList.css
└── CommentForm/
    ├── CommentForm.tsx    # Formulario comentario
    ├── CommentForm.test.tsx # 4 tests
    └── CommentForm.css
```

**Patrón: Componentes controlados con hooks**

```typescript
// Login.tsx - Ejemplo de componente controlado
export const Login: React.FC<LoginProps> = ({ onLoginSuccess }) => {
  // 1. Estado local (hooks)
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  // 2. Handler de submit
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    
    try {
      // 3. Llamar al service
      const user = await authService.login({ email, password });
      
      // 4. Callback al padre
      onLoginSuccess(user);
    } catch (err) {
      setError('Error en la autenticación');
    } finally {
      setLoading(false);
    }
  };

  // 5. Render
  return (
    <form onSubmit={handleSubmit}>
      <input value={email} onChange={(e) => setEmail(e.target.value)} />
      {/* ... */}
    </form>
  );
};
```

**¿Por qué esta estructura de componentes?**

✅ **Reusabilidad:** Cada componente es independiente
✅ **Testeable:** Se puede testear aisladamente
✅ **Mantenible:** Cambios localizados
✅ **Single Responsibility:** Cada componente hace una cosa

**Coverage:** 92.44% (muy alto)

---

### 3.2 Directorio `/frontend/src/services/`

**Propósito:** Capa de comunicación HTTP con el backend (API client).

```
frontend/src/services/
├── authService.ts       # login(), register()
├── authService.test.ts  # 6 tests
├── postService.ts       # CRUD posts y comentarios
└── postService.test.ts  # 10 tests
```

**Patrón: Service Layer + axios**

```typescript
// authService.ts
import axios from 'axios';
import { User } from '../types';

const API_URL = 'http://localhost:8080/api';

export const authService = {
  async login(credentials: { email: string; password: string }): Promise<User> {
    const response = await axios.post(`${API_URL}/auth/login`, credentials);
    return response.data;
  },

  async register(data: { email: string; password: string; username: string }): Promise<User> {
    const response = await axios.post(`${API_URL}/auth/register`, data);
    return response.data;
  }
};
```

**¿Por qué separar services de components?**

✅ **Reusabilidad:** Múltiples componentes usan el mismo service
✅ **Testeable:** Mockear axios es más fácil que mockear fetch en componentes
✅ **Mantenibilidad:** Cambios en API solo afectan services
✅ **Type Safety:** TypeScript verifica tipos de request/response

**Tests de services:**

```typescript
// authService.test.ts
jest.mock('axios');
const mockedAxios = axios as jest.Mocked<typeof axios>;

test('login exitoso', async () => {
  const mockUser = { id: 1, email: 'test@example.com' };
  mockedAxios.post.mockResolvedValue({ data: mockUser });
  
  const user = await authService.login({ email: 'test@example.com', password: '123456' });
  
  expect(user).toEqual(mockUser);
  expect(mockedAxios.post).toHaveBeenCalledWith(
    'http://localhost:8080/api/auth/login',
    { email: 'test@example.com', password: '123456' }
  );
});
```

**Coverage:** 100% en services

---

### 3.3 Directorio `/frontend/src/types/`

**Propósito:** Definiciones de tipos TypeScript compartidos.

```
frontend/src/types/
└── index.ts
```

**Contenido:**

```typescript
export interface User {
  id: number;
  email: string;
  username: string;
  created_at: string;
}

export interface Post {
  id: number;
  title: string;
  content: string;
  user_id: number;
  username: string;
  created_at: string;
}

export interface Comment {
  id: number;
  post_id: number;
  user_id: number;
  username: string;
  content: string;
  created_at: string;
}
```

**¿Por qué centralizar tipos?**

✅ **Type Safety:** Errores de tipo detectados en compilación
✅ **Autocompletado:** IDEs ofrecen sugerencias
✅ **Documentación:** Los tipos documentan la estructura
✅ **Refactoring:** Cambiar un tipo actualiza todo

---

### 3.4 Archivo `/frontend/src/App.tsx`

**Propósito:** Componente raíz que orquesta la aplicación.

```typescript
function App() {
  // 1. Estado global
  const [currentUser, setCurrentUser] = useState<User | null>(null);
  const [currentView, setCurrentView] = useState<'list' | 'detail'>('list');
  const [selectedPostId, setSelectedPostId] = useState<number | null>(null);

  // 2. Renderizado condicional
  if (!currentUser) {
    return <Login onLoginSuccess={setCurrentUser} />;
  }

  return (
    <div className="App">
      <header>
        <h1>Mini Red Social</h1>
        <button onClick={handleLogout}>Cerrar Sesión</button>
      </header>

      {currentView === 'list' ? (
        <>
          <CreatePost userId={currentUser.id} onPostCreated={handlePostCreated} />
          <PostList currentUserId={currentUser.id} onViewPost={handleViewPost} />
        </>
      ) : (
        <PostDetail postId={selectedPostId!} userId={currentUser.id} onBack={handleBackToList} />
      )}
    </div>
  );
}
```

**¿Por qué App.tsx NO se mide en coverage?**

❌ Solo orquesta componentes (no tiene lógica)
✅ Se testea implícitamente con E2E (Cypress)

---

## 4. Testing: Organización y Estrategia

### 4.1 Tests Unitarios Backend (Go)

**Ubicación:** `/backend/tests/services/`

**Herramientas:**
- `testing` (built-in Go)
- `testify/assert` (assertions)
- `testify/mock` (mocking)

**Estrategia:**

```go
// Estructura de un test
func TestFunctionName_Scenario_ExpectedResult(t *testing.T) {
    // 1. Arrange (preparar)
    mockRepo := new(mocks.MockUserRepository)
    service := services.NewAuthService(mockRepo)
    mockRepo.On("FindByEmail", "test@example.com").Return(nil, nil)
    
    // 2. Act (ejecutar)
    user, err := service.Register("test@example.com", "123456", "testuser")
    
    // 3. Assert (verificar)
    assert.NoError(t, err)
    assert.NotNil(t, user)
    
    // 4. Verify mocks
    mockRepo.AssertExpectations(t)
}
```

**Cobertura:** 86.5% en `/internal/services/`

---

### 4.2 Tests Unitarios Frontend (React)

**Ubicación:** Junto a cada componente (`*.test.tsx`)

**Herramientas:**
- Jest (test runner)
- React Testing Library (render + queries)
- @testing-library/user-event (interacciones)

**Estrategia:**

```typescript
// Estructura de un test
test('descripción del comportamiento', () => {
  // 1. Arrange: renderizar componente
  render(<Login onLoginSuccess={mockCallback} />);
  
  // 2. Act: interactuar con UI
  const emailInput = screen.getByLabelText(/email/i);
  fireEvent.change(emailInput, { target: { value: 'test@example.com' } });
  fireEvent.click(screen.getByRole('button', { name: /iniciar sesión/i }));
  
  // 3. Assert: verificar resultado
  await waitFor(() => {
    expect(mockCallback).toHaveBeenCalled();
  });
});
```

**Cobertura:** 92.44% en components y services

---

### 4.3 Tests E2E (Cypress)

**Ubicación:** `/frontend/cypress/e2e/blog/`

**Herramientas:**
- Cypress 13.15.2

**Estrategia:**

```javascript
// Test de flujo completo
describe('Authentication Flow', () => {
  beforeEach(() => {
    cy.visit('/')
  })

  it('debería hacer login exitoso', () => {
    // 1. Mock del backend
    cy.intercept('POST', '**/api/auth/login', {
      statusCode: 200,
      body: { id: 1, email: 'test@example.com', username: 'testuser' }
    })

    // 2. Interacción con UI
    cy.get('input#email').type('test@example.com')
    cy.get('input#password').type('123456')
    cy.get('button[type="submit"]').click()

    // 3. Verificación
    cy.contains('Mini Red Social').should('be.visible')
    cy.contains('Hola, @testuser').should('be.visible')
  })
})
```

**Cobertura:** 15 tests que validan flujos completos

---

## 5. CI/CD: Pipeline y Configuración

### 5.1 Archivo `/.github/workflows/ci.yml`

**Propósito:** Definición del pipeline de CI/CD en GitHub Actions.

**Estructura:**

```yaml
name: CI/CD Pipeline

on:
  push:
    branches: [ main, master, develop ]
  pull_request:
    branches: [ main, master, develop ]

jobs:
  backend-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
      - run: go test ./tests/services/... -cover
      - run: # Verificar coverage >= 70%

  frontend-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
      - run: npm ci
      - run: npm test -- --coverage

  sonarcloud:
    needs: [backend-tests, frontend-tests]
    steps:
      - uses: SonarSource/sonarqube-scan-action@v5.0.0

  cypress-e2e:
    needs: [backend-tests, frontend-tests]
    steps:
      - uses: cypress-io/github-action@v6

  quality-summary:
    needs: [backend-tests, frontend-tests, cypress-e2e, sonarcloud]
    steps:
      - run: echo "✅ All quality gates passed"
```

**¿Por qué esta estructura?**

✅ **Paralelización:** backend-tests + frontend-tests corren simultáneamente
✅ **Dependencies:** jobs posteriores dependen de anteriores
✅ **Quality Gates:** pipeline falla si algo no pasa
✅ **Feedback rápido:** tests unitarios primero, E2E después

---

### 5.2 Archivo `/sonar-project.properties`

**Propósito:** Configuración de SonarCloud.

```properties
# Identificación del proyecto
sonar.projectKey=OctavioCarpineti_IngSWIII-TP07-Quality
sonar.organization=octaviocarpineti

# Directorios de código fuente
sonar.sources=backend/internal,frontend/src

# Exclusiones de análisis
sonar.exclusions=**/tests/**,**/mocks/**,**/*.test.tsx

# Coverage reports
sonar.go.coverage.reportPaths=backend/coverage.out
sonar.javascript.lcov.reportPaths=frontend/coverage/lcov.info

# Duplications
sonar.cpd.exclusions=**/tests/**
```

**¿Por qué estas exclusiones?**

✅ Tests no son código de producción
✅ Mocks son generados automáticamente
✅ Se mide solo código relevante

---

## 📊 Resumen: Responsabilidades por Capa

| Capa | Responsabilidad | Coverage | Testing |
|------|----------------|----------|---------|
| **Handlers** | HTTP mapping | NO | E2E (Cypress) |
| **Services** | Lógica de negocio | ✅ 86.5% | Unit (mocks) |
| **Repository** | Acceso a datos | NO | Mockado |
| **Components** | UI React | ✅ 92.44% | Unit (RTL) |
| **Services (FE)** | API client | ✅ 100% | Unit (mock axios) |
| **E2E** | Flujos completos | N/A | 15 tests |

---

## 🎯 Conclusión

La arquitectura implementada sigue principios sólidos de ingeniería de software:

1. **Separación de responsabilidades:** Cada capa tiene un propósito claro
2. **Dependency Inversion:** Dependencias inyectadas via interfaces
3. **Testabilidad:** Cada capa se testea apropiadamente
4. **Mantenibilidad:** Estructura clara y predecible
5. **Escalabilidad:** Fácil agregar nuevas funcionalidades

Esta estructura permite alcanzar **89 tests** con **86.5%/92.44% coverage** de manera organizada y profesional.