package services

import (
	"errors"
	"strings"

	"tp06-testing/internal/models"
	"tp06-testing/internal/repository"
)

// Constantes para mensajes de error
const (
	ErrUserNotFound = "usuario no encontrado"
	ErrPostNotFound = "post no encontrado"
)

// PostService maneja la lógica de posts y comentarios
type PostService struct {
	postRepo repository.PostRepository
	userRepo repository.UserRepository
}

// NewPostService crea una nueva instancia
func NewPostService(postRepo repository.PostRepository, userRepo repository.UserRepository) *PostService {
	return &PostService{
		postRepo: postRepo,
		userRepo: userRepo,
	}
}

// CreatePost crea un nuevo post
func (s *PostService) CreatePost(req *models.CreatePostRequest, userID int) (*models.Post, error) {
	if strings.TrimSpace(req.Title) == "" {
		return nil, errors.New("el título es requerido")
	}

	if len(strings.TrimSpace(req.Title)) < 3 {
		return nil, errors.New("el título debe tener al menos 3 caracteres")
	}

	if strings.TrimSpace(req.Content) == "" {
		return nil, errors.New("el contenido es requerido")
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New(ErrUserNotFound)
	}

	post := &models.Post{
		Title:   strings.TrimSpace(req.Title),
		Content: strings.TrimSpace(req.Content),
		UserID:  userID,
	}

	err = s.postRepo.Create(post)
	if err != nil {
		return nil, err
	}

	post.Username = user.Username

	return post, nil
}

// GetAllPosts obtiene todos los posts del sistema.
// Retorna una lista vacía si no hay posts, nunca retorna nil.
func (s *PostService) GetAllPosts() ([]*models.Post, error) {
	posts, err := s.postRepo.FindAll()
	if err != nil {
		return nil, err
	}

	// Si no hay posts, devolver lista vacía (no es error)
	if posts == nil {
		return []*models.Post{}, nil
	}

	return posts, nil
}

// GetPostByID obtiene un post específico
func (s *PostService) GetPostByID(id int) (*models.Post, error) {
	if id <= 0 {
		return nil, errors.New("id inválido")
	}

	post, err := s.postRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if post == nil {
		return nil, errors.New(ErrPostNotFound)
	}

	return post, nil
}

// DeletePost elimina un post (solo el autor puede hacerlo)
func (s *PostService) DeletePost(postID int, userID int) error {
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return err
	}
	if post == nil {
		return errors.New(ErrPostNotFound)
	}

	if post.UserID != userID {
		return errors.New("no tienes permiso para eliminar este post")
	}

	return s.postRepo.Delete(postID)
}

// CreateComment agrega un comentario a un post
func (s *PostService) CreateComment(postID int, req *models.CreateCommentRequest, userID int) (*models.Comment, error) {
	if strings.TrimSpace(req.Content) == "" {
		return nil, errors.New("el contenido del comentario es requerido")
	}

	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, errors.New(ErrPostNotFound)
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New(ErrUserNotFound)
	}

	comment := &models.Comment{
		PostID:  postID,
		UserID:  userID,
		Content: strings.TrimSpace(req.Content),
	}

	err = s.postRepo.CreateComment(comment)
	if err != nil {
		return nil, err
	}

	comment.Username = user.Username

	return comment, nil
}

// GetCommentsByPostID obtiene todos los comentarios de un post
func (s *PostService) GetCommentsByPostID(postID int) ([]*models.Comment, error) {
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, errors.New(ErrPostNotFound)
	}

	comments, err := s.postRepo.FindCommentsByPostID(postID)
	if err != nil {
		return nil, err
	}

	if comments == nil {
		return []*models.Comment{}, nil
	}

	return comments, nil
}

func (s *PostService) DeleteComment(postID int, commentID int, userID int) error {
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return err
	}
	if post == nil {
		return errors.New(ErrPostNotFound)
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New(ErrUserNotFound)
	}

	return s.postRepo.DeleteComment(postID, commentID, userID)
}
