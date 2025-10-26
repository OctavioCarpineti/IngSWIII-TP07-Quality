package services

import (
	"errors"
	"testing"

	"tp06-testing/internal/models"
	"tp06-testing/internal/services"
	"tp06-testing/tests/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestCreatePost_Success prueba la creación exitosa de un post
func TestCreatePost_Success(t *testing.T) {
	// ARRANGE
	mockRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	postService := services.NewPostService(mockRepo, mockUserRepo)

	// ← AGREGAR ESTO
	existingUser := &models.User{
		ID:       1,
		Email:    "test@example.com",
		Username: "testuser",
	}
	mockUserRepo.On("FindByID", 1).Return(existingUser, nil)
	// ← FIN

	// Configurar mock: Create debe ejecutarse correctamente
	mockRepo.On("Create", mock.AnythingOfType("*models.Post")).Return(nil)

	req := &models.CreatePostRequest{
		Title:   "Test Post",
		Content: "This is a test post",
	}

	// ACT
	post, err := postService.CreatePost(req, 1)

	// ASSERT
	assert.NoError(t, err)
	assert.NotNil(t, post)
	assert.Equal(t, "Test Post", post.Title)
	assert.Equal(t, "This is a test post", post.Content)

	// Verificar que se llamaron los métodos del mock
	mockRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t) // ← AGREGAR ESTO TAMBIÉN
}

// TestCreatePost_UserNotFound: el userId no existe -> error
func TestCreatePost_UserNotFound(t *testing.T) {
	// ARRANGE
	mockRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	postService := services.NewPostService(mockRepo, mockUserRepo)

	// Configurar mock: FindByID del user devuelve nil (no existe)
	mockUserRepo.On("FindByID", 999).Return(nil, nil)

	req := &models.CreatePostRequest{
		Title:   "Test Post",
		Content: "This is a test post",
	}

	// ACT
	post, err := postService.CreatePost(req, 999)

	// ASSERT
	assert.Error(t, err)
	assert.Nil(t, post)
	assert.Equal(t, "usuario no encontrado", err.Error())

	mockUserRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Create")
}

// TestCreatePost_RepoError: el repositorio falla al crear -> se propaga error
func TestCreatePost_RepoError(t *testing.T) {
	// ARRANGE
	mockRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	postService := services.NewPostService(mockRepo, mockUserRepo)

	// Usuario existe
	existingUser := &models.User{ID: 1, Email: "u@u.com", Username: "u"}
	mockUserRepo.On("FindByID", 1).Return(existingUser, nil)

	// El repo Create falla
	mockRepo.On("Create", mock.AnythingOfType("*models.Post")).Return(errors.New("db error"))

	req := &models.CreatePostRequest{
		Title:   "Test Post",
		Content: "This is a test post",
	}

	// ACT
	post, err := postService.CreatePost(req, 1)

	// ASSERT
	assert.Error(t, err)
	assert.Nil(t, post)
	assert.Equal(t, "db error", err.Error())

	mockRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

// TestCreatePost_TitleVacio: validación previa falla si title vacío
func TestCreatePost_TitleVacio(t *testing.T) {
	// ARRANGE
	mockRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	postService := services.NewPostService(mockRepo, mockUserRepo)

	req := &models.CreatePostRequest{
		Title:   "", // título vacío
		Content: "Contenido",
	}

	// ACT
	post, err := postService.CreatePost(req, 1)

	// ASSERT
	assert.Error(t, err)
	assert.Nil(t, post)
	assert.Equal(t, "el título es requerido", err.Error())
	// No debe llamar al repo ni al userRepo
	mockRepo.AssertNotCalled(t, "Create")
	mockUserRepo.AssertNotCalled(t, "FindByID")
}

// TestCreatePost_ContentVacio: validación previa falla si content vacío
func TestCreatePost_ContentVacio(t *testing.T) {
	// ARRANGE
	mockRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	postService := services.NewPostService(mockRepo, mockUserRepo)

	req := &models.CreatePostRequest{
		Title:   "Test Post",
		Content: "", // content vacío
	}

	// ACT
	post, err := postService.CreatePost(req, 1)

	// ASSERT
	assert.Error(t, err)
	assert.Nil(t, post)
	assert.Equal(t, "el contenido es requerido", err.Error())

	mockRepo.AssertNotCalled(t, "Create")
	mockUserRepo.AssertNotCalled(t, "FindByID")
}

// TestDeletePost_Success prueba eliminación exitosa por el autor
func TestDeletePost_Success(t *testing.T) {
	// ARRANGE
	mockRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	postService := services.NewPostService(mockRepo, mockUserRepo)

	existingPost := &models.Post{
		ID:       1,
		Title:    "Test Post",
		Content:  "Content",
		UserID:   1, // El autor es el usuario 1
		Username: "testuser",
	}

	// Configurar mocks
	mockRepo.On("FindByID", 1).Return(existingPost, nil)
	mockRepo.On("Delete", 1).Return(nil)

	// ACT: El usuario 1 elimina su propio post
	err := postService.DeletePost(1, 1)

	// ASSERT
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

// TestDeletePost_PostNoExiste prueba eliminar post inexistente
func TestDeletePost_PostNoExiste(t *testing.T) {
	// ARRANGE
	mockRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	postService := services.NewPostService(mockRepo, mockUserRepo)

	// Post no existe
	mockRepo.On("FindByID", 999).Return(nil, nil)

	// ACT
	err := postService.DeletePost(999, 1)

	// ASSERT
	assert.Error(t, err)
	assert.Equal(t, "post no encontrado", err.Error())

	// NO debe intentar eliminar
	mockRepo.AssertNotCalled(t, "Delete")
}

// TestDeletePost_NoEsAutor prueba que solo el autor puede eliminar
func TestDeletePost_NoEsAutor(t *testing.T) {
	// ARRANGE
	mockRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	postService := services.NewPostService(mockRepo, mockUserRepo)

	existingPost := &models.Post{
		ID:       1,
		Title:    "Test Post",
		Content:  "Content",
		UserID:   1, // El autor es el usuario 1
		Username: "testuser",
	}

	mockRepo.On("FindByID", 1).Return(existingPost, nil)

	// ACT: El usuario 2 intenta eliminar el post del usuario 1
	err := postService.DeletePost(1, 2)

	// ASSERT
	assert.Error(t, err)
	assert.Equal(t, "no tienes permiso para eliminar este post", err.Error())

	// NO debe llamar a Delete porque no tiene permiso
	mockRepo.AssertNotCalled(t, "Delete")
	mockRepo.AssertExpectations(t)
}

// TestDeleteComment_Success prueba eliminación exitosa por el autor
func TestDeleteComment_Success(t *testing.T) {
	// ARRANGE
	mockRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	postService := services.NewPostService(mockRepo, mockUserRepo)

	existingPost := &models.Post{
		ID:       1,
		Title:    "Test Post",
		UserID:   1,
		Username: "testuser",
	}

	existingUser := &models.User{
		ID:       1,
		Email:    "test@example.com",
		Username: "testuser",
	}

	// Configurar mocks
	mockRepo.On("FindByID", 1).Return(existingPost, nil)
	mockUserRepo.On("FindByID", 1).Return(existingUser, nil)
	mockRepo.On("DeleteComment", 1, 10, 1).Return(nil)

	// ACT: El usuario 1 elimina su propio comentario
	err := postService.DeleteComment(1, 10, 1)

	// ASSERT
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

// TestDeleteComment_PostNoExiste prueba eliminar comentario en post inexistente
func TestDeleteComment_PostNoExiste(t *testing.T) {
	// ARRANGE
	mockRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	postService := services.NewPostService(mockRepo, mockUserRepo)

	// Post no existe
	mockRepo.On("FindByID", 999).Return(nil, nil)

	// ACT
	err := postService.DeleteComment(999, 10, 1)

	// ASSERT
	assert.Error(t, err)
	assert.Equal(t, "post no encontrado", err.Error())

	// NO debe intentar eliminar
	mockRepo.AssertNotCalled(t, "DeleteComment")
}

// TestDeleteComment_UsuarioNoExiste prueba eliminar con usuario inexistente
func TestDeleteComment_UsuarioNoExiste(t *testing.T) {
	// ARRANGE
	mockRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	postService := services.NewPostService(mockRepo, mockUserRepo)

	existingPost := &models.Post{
		ID:       1,
		Title:    "Test Post",
		UserID:   1,
		Username: "testuser",
	}

	mockRepo.On("FindByID", 1).Return(existingPost, nil)
	mockUserRepo.On("FindByID", 999).Return(nil, nil)

	// ACT
	err := postService.DeleteComment(1, 10, 999)

	// ASSERT
	assert.Error(t, err)
	assert.Equal(t, "usuario no encontrado", err.Error())
	mockRepo.AssertNotCalled(t, "DeleteComment")
}

// TestDeleteComment_NoEsAutor prueba que solo el autor puede eliminar su comentario
func TestDeleteComment_NoEsAutor(t *testing.T) {
	// ARRANGE
	mockRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	postService := services.NewPostService(mockRepo, mockUserRepo)

	existingPost := &models.Post{
		ID:       1,
		Title:    "Test Post",
		UserID:   1,
		Username: "testuser",
	}

	existingUser := &models.User{
		ID:       2,
		Email:    "other@example.com",
		Username: "otheruser",
	}

	mockRepo.On("FindByID", 1).Return(existingPost, nil)
	mockUserRepo.On("FindByID", 2).Return(existingUser, nil)

	// Usuario 2 intenta eliminar comentario del usuario 1
	mockRepo.On("DeleteComment", 1, 10, 2).Return(errors.New("no tienes permiso para eliminar este comentario o no existe"))

	// ACT
	err := postService.DeleteComment(1, 10, 2)

	// ASSERT
	assert.Error(t, err)
	assert.Equal(t, "no tienes permiso para eliminar este comentario o no existe", err.Error())
	mockRepo.AssertExpectations(t)
}

// TestGetAllPosts_Success prueba obtener todos los posts
func TestGetAllPosts_Success(t *testing.T) {
	// ARRANGE
	mockPostRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	postService := services.NewPostService(mockPostRepo, mockUserRepo)

	mockPosts := []*models.Post{
		{ID: 1, Title: "Post 1", Content: "Content 1", UserID: 1},
		{ID: 2, Title: "Post 2", Content: "Content 2", UserID: 2},
	}
	mockPostRepo.On("FindAll").Return(mockPosts, nil)

	// ACT
	posts, err := postService.GetAllPosts()

	// ASSERT
	assert.NoError(t, err)
	assert.Len(t, posts, 2)
	mockPostRepo.AssertExpectations(t)
}

// TestGetPostByID_Success prueba obtener un post por ID
func TestGetPostByID_Success(t *testing.T) {
	// ARRANGE
	mockPostRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	postService := services.NewPostService(mockPostRepo, mockUserRepo)

	mockPost := &models.Post{
		ID:      1,
		Title:   "Test Post",
		Content: "Test Content",
		UserID:  1,
	}
	mockPostRepo.On("FindByID", 1).Return(mockPost, nil)

	// ACT
	post, err := postService.GetPostByID(1)

	// ASSERT
	assert.NoError(t, err)
	assert.NotNil(t, post)
	assert.Equal(t, "Test Post", post.Title)
	mockPostRepo.AssertExpectations(t)
}

// TestCreateComment_Success prueba crear un comentario
func TestCreateComment_Success(t *testing.T) {
	// ARRANGE
	mockPostRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	postService := services.NewPostService(mockPostRepo, mockUserRepo)

	mockPost := &models.Post{ID: 1, Title: "Post", UserID: 1}
	mockUser := &models.User{ID: 2, Username: "commenter"}

	mockPostRepo.On("FindByID", 1).Return(mockPost, nil)
	mockUserRepo.On("FindByID", 2).Return(mockUser, nil)
	mockPostRepo.On("CreateComment", mock.AnythingOfType("*models.Comment")).Return(nil)

	req := &models.CreateCommentRequest{
		Content: "Great post!",
	}

	// ACT
	comment, err := postService.CreateComment(1, req, 2)

	// ASSERT
	assert.NoError(t, err)
	assert.NotNil(t, comment)
	assert.Equal(t, "Great post!", comment.Content)
	assert.Equal(t, 1, comment.PostID)
	assert.Equal(t, 2, comment.UserID)
	mockPostRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

// TestGetCommentsByPostID_Success prueba obtener comentarios de un post
func TestGetCommentsByPostID_Success(t *testing.T) {
	// ARRANGE
	mockPostRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	postService := services.NewPostService(mockPostRepo, mockUserRepo)

	mockPost := &models.Post{ID: 1, Title: "Post", UserID: 1}
	mockComments := []*models.Comment{
		{ID: 1, PostID: 1, UserID: 1, Content: "Comment 1"},
		{ID: 2, PostID: 1, UserID: 2, Content: "Comment 2"},
	}

	mockPostRepo.On("FindByID", 1).Return(mockPost, nil)
	mockPostRepo.On("FindCommentsByPostID", 1).Return(mockComments, nil)

	// ACT
	comments, err := postService.GetCommentsByPostID(1)

	// ASSERT
	assert.NoError(t, err)
	assert.Len(t, comments, 2)
	mockPostRepo.AssertExpectations(t)
}

// ========== Tests adicionales para GetAllPosts ==========

// TestGetAllPosts_Empty prueba cuando no hay posts
func TestGetAllPosts_Empty(t *testing.T) {
	// ARRANGE
	mockPostRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	postService := services.NewPostService(mockPostRepo, mockUserRepo)

	mockPostRepo.On("FindAll").Return(nil, nil)

	// ACT
	posts, err := postService.GetAllPosts()

	// ASSERT
	assert.NoError(t, err)
	assert.NotNil(t, posts)
	assert.Len(t, posts, 0)
	mockPostRepo.AssertExpectations(t)
}

// ========== Tests adicionales para GetPostByID ==========

// TestGetPostByID_InvalidID prueba con ID inválido (negativo o cero)
func TestGetPostByID_InvalidID(t *testing.T) {
	// ARRANGE
	mockPostRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	postService := services.NewPostService(mockPostRepo, mockUserRepo)

	// ACT
	post, err := postService.GetPostByID(0)

	// ASSERT
	assert.Error(t, err)
	assert.Nil(t, post)
	assert.Equal(t, "id inválido", err.Error())
	mockPostRepo.AssertNotCalled(t, "FindByID")
}

// TestGetPostByID_NotFound prueba cuando el post no existe
func TestGetPostByID_NotFound(t *testing.T) {
	// ARRANGE
	mockPostRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	postService := services.NewPostService(mockPostRepo, mockUserRepo)

	mockPostRepo.On("FindByID", 999).Return(nil, nil)

	// ACT
	post, err := postService.GetPostByID(999)

	// ASSERT
	assert.Error(t, err)
	assert.Nil(t, post)
	assert.Equal(t, "post no encontrado", err.Error())
	mockPostRepo.AssertExpectations(t)
}

// ========== Tests adicionales para CreateComment ==========

// TestCreateComment_EmptyContent prueba comentario con contenido vacío
func TestCreateComment_EmptyContent(t *testing.T) {
	// ARRANGE
	mockPostRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	postService := services.NewPostService(mockPostRepo, mockUserRepo)

	req := &models.CreateCommentRequest{
		Content: "",
	}

	// ACT
	comment, err := postService.CreateComment(1, req, 1)

	// ASSERT
	assert.Error(t, err)
	assert.Nil(t, comment)
	assert.Equal(t, "el contenido del comentario es requerido", err.Error())
	mockPostRepo.AssertNotCalled(t, "FindByID")
	mockPostRepo.AssertNotCalled(t, "CreateComment")
}

// TestCreateComment_PostNotFound prueba comentar en post inexistente
func TestCreateComment_PostNotFound(t *testing.T) {
	// ARRANGE
	mockPostRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	postService := services.NewPostService(mockPostRepo, mockUserRepo)

	mockPostRepo.On("FindByID", 999).Return(nil, nil)

	req := &models.CreateCommentRequest{
		Content: "Great post!",
	}

	// ACT
	comment, err := postService.CreateComment(999, req, 1)

	// ASSERT
	assert.Error(t, err)
	assert.Nil(t, comment)
	assert.Equal(t, "post no encontrado", err.Error())
	mockPostRepo.AssertExpectations(t)
	mockPostRepo.AssertNotCalled(t, "CreateComment")
}

// TestCreateComment_UserNotFound prueba comentar con usuario inexistente
func TestCreateComment_UserNotFound(t *testing.T) {
	// ARRANGE
	mockPostRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	postService := services.NewPostService(mockPostRepo, mockUserRepo)

	mockPost := &models.Post{ID: 1, Title: "Post", UserID: 1}
	mockPostRepo.On("FindByID", 1).Return(mockPost, nil)
	mockUserRepo.On("FindByID", 999).Return(nil, nil)

	req := &models.CreateCommentRequest{
		Content: "Great post!",
	}

	// ACT
	comment, err := postService.CreateComment(1, req, 999)

	// ASSERT
	assert.Error(t, err)
	assert.Nil(t, comment)
	assert.Equal(t, "usuario no encontrado", err.Error())
	mockPostRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockPostRepo.AssertNotCalled(t, "CreateComment")
}

// ========== Tests adicionales para GetCommentsByPostID ==========

// TestGetCommentsByPostID_PostNotFound prueba obtener comentarios de post inexistente
func TestGetCommentsByPostID_PostNotFound(t *testing.T) {
	// ARRANGE
	mockPostRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	postService := services.NewPostService(mockPostRepo, mockUserRepo)

	mockPostRepo.On("FindByID", 999).Return(nil, nil)

	// ACT
	comments, err := postService.GetCommentsByPostID(999)

	// ASSERT
	assert.Error(t, err)
	assert.Nil(t, comments)
	assert.Equal(t, "post no encontrado", err.Error())
	mockPostRepo.AssertExpectations(t)
	mockPostRepo.AssertNotCalled(t, "FindCommentsByPostID")
}

// TestGetCommentsByPostID_Empty prueba cuando no hay comentarios
func TestGetCommentsByPostID_Empty(t *testing.T) {
	// ARRANGE
	mockPostRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	postService := services.NewPostService(mockPostRepo, mockUserRepo)

	mockPost := &models.Post{ID: 1, Title: "Post", UserID: 1}
	mockPostRepo.On("FindByID", 1).Return(mockPost, nil)
	mockPostRepo.On("FindCommentsByPostID", 1).Return(nil, nil)

	// ACT
	comments, err := postService.GetCommentsByPostID(1)

	// ASSERT
	assert.NoError(t, err)
	assert.NotNil(t, comments)
	assert.Len(t, comments, 0)
	mockPostRepo.AssertExpectations(t)
}
