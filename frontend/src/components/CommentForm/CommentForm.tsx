import React, { useState } from 'react';
import { postService } from '../../services/postService';
import './CommentForm.css';

interface CommentFormProps {
  postId: number;
  userId: number;
  onCommentCreated: () => void;
}

export const CommentForm: React.FC<CommentFormProps> = ({ postId, userId, onCommentCreated }) => {
  const [content, setContent] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      await postService.createComment(postId, { content }, userId);
      setContent('');
      onCommentCreated();
    } catch (err: any) {
      setError(err.response?.data?.error || 'Error al crear comentario');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="comment-form-container">
      <h3>Agregar Comentario</h3>
      <form onSubmit={handleSubmit}>
        <textarea
          value={content}
          onChange={(e) => setContent(e.target.value)}
          placeholder="Escribe tu comentario..."
          rows={3}
          required
          disabled={loading}
        />

        {error && <div className="error-message">{error}</div>}

        <button type="submit" disabled={loading || !content.trim()}>
          {loading ? 'Publicando...' : 'Comentar'}
        </button>
      </form>
    </div>
  );
};