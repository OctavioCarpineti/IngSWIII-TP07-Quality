import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { CommentForm } from './CommentForm';
import { postService } from '../../services/postService';

jest.mock('../../services/postService');
const mockedPostService = postService as jest.Mocked<typeof postService>;

describe('CommentForm Component', () => {
    const mockOnCommentCreated = jest.fn();

    beforeEach(() => {
        jest.clearAllMocks();
    });

    test('renderiza el formulario correctamente', () => {
        render(<CommentForm postId={1} userId={1} onCommentCreated={mockOnCommentCreated} />);

        expect(screen.getByText('Agregar Comentario')).toBeInTheDocument();
        expect(screen.getByPlaceholderText(/escribe tu comentario/i)).toBeInTheDocument();
        expect(screen.getByRole('button', { name: /comentar/i })).toBeInTheDocument();
    });

    test('crea comentario exitosamente', async () => {
        mockedPostService.createComment.mockResolvedValueOnce({
            id: 1,
            post_id: 1,
            user_id: 1,
            username: 'testuser',
            content: 'Great post!',
            created_at: '2024-01-01'
        });

        render(<CommentForm postId={1} userId={1} onCommentCreated={mockOnCommentCreated} />);

        const textarea = screen.getByPlaceholderText(/escribe tu comentario/i);
        const submitButton = screen.getByRole('button', { name: /comentar/i });

        fireEvent.change(textarea, { target: { value: 'Great post!' } });
        fireEvent.click(submitButton);

        await waitFor(() => {
            expect(mockedPostService.createComment).toHaveBeenCalledWith(
                1,
                { content: 'Great post!' },
                1
            );
            expect(mockOnCommentCreated).toHaveBeenCalled();
        });
    });

    test('muestra error cuando falla la creación', async () => {
        mockedPostService.createComment.mockRejectedValueOnce({
            response: { data: { error: 'Error al crear comentario' } }
        });

        render(<CommentForm postId={1} userId={1} onCommentCreated={mockOnCommentCreated} />);

        fireEvent.change(screen.getByPlaceholderText(/escribe tu comentario/i), {
            target: { value: 'Test comment' }
        });
        fireEvent.click(screen.getByRole('button', { name: /comentar/i }));

        await waitFor(() => {
            expect(screen.getByText('Error al crear comentario')).toBeInTheDocument();
        });

        expect(mockOnCommentCreated).not.toHaveBeenCalled();
    });

    test('botón deshabilitado cuando el textarea está vacío', () => {
        render(<CommentForm postId={1} userId={1} onCommentCreated={mockOnCommentCreated} />);

        const submitButton = screen.getByRole('button', { name: /comentar/i });
        expect(submitButton).toBeDisabled();
    });
});