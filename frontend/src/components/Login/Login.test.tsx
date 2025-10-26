import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { Login } from './Login';
import axios from 'axios';

jest.mock('axios');
const mockedAxios = axios as jest.Mocked<typeof axios>;

describe('Login Component', () => {
  const mockOnLoginSuccess = jest.fn();

  beforeEach(() => {
    jest.clearAllMocks();
  });

  test('renderiza el formulario de login correctamente', () => {
    render(<Login onLoginSuccess={mockOnLoginSuccess} />);

    // Verificar que se renderiza el heading
    expect(screen.getByRole('heading', { name: /iniciar sesión/i })).toBeInTheDocument();
    
    // Verificar inputs
    expect(screen.getByLabelText(/email/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/password/i)).toBeInTheDocument();
    
    // Verificar botón de submit
    expect(screen.getByRole('button', { name: /iniciar sesión/i })).toBeInTheDocument();
  });

  test('muestra formulario de registro al cambiar modo', () => {
    render(<Login onLoginSuccess={mockOnLoginSuccess} />);

    // Click en botón para cambiar a registro
    const toggleButton = screen.getByText(/¿No tienes cuenta\? Regístrate/i);
    fireEvent.click(toggleButton);

    // Verificar que muestra heading de Registrarse
    expect(screen.getByRole('heading', { name: /registrarse/i })).toBeInTheDocument();
    
    // Verificar que aparece el campo username
    expect(screen.getByLabelText(/username/i)).toBeInTheDocument();
  });

  test('login exitoso llama a onLoginSuccess', async () => {
    const mockUser = {
      id: 1,
      email: 'test@example.com',
      username: 'testuser',
      created_at: '2025-01-01'
    };

    mockedAxios.post.mockResolvedValueOnce({ data: mockUser });

    render(<Login onLoginSuccess={mockOnLoginSuccess} />);

    // Llenar formulario
    const emailInput = screen.getByLabelText(/email/i) as HTMLInputElement;
    const passwordInput = screen.getByLabelText(/password/i) as HTMLInputElement;

    fireEvent.change(emailInput, { target: { value: 'test@example.com' } });
    fireEvent.change(passwordInput, { target: { value: '123456' } });

    // Submit
    const submitButton = screen.getByRole('button', { name: /iniciar sesión/i });
    fireEvent.click(submitButton);

    // Verificar
    await waitFor(() => {
      expect(mockOnLoginSuccess).toHaveBeenCalledWith(mockUser);
    });
  });

  test('muestra error cuando login falla', async () => {
    mockedAxios.post.mockRejectedValueOnce({
      response: {
        data: {
          error: 'Credenciales inválidas'
        }
      }
    });

    render(<Login onLoginSuccess={mockOnLoginSuccess} />);

    const emailInput = screen.getByLabelText(/email/i) as HTMLInputElement;
    const passwordInput = screen.getByLabelText(/password/i) as HTMLInputElement;

    fireEvent.change(emailInput, { target: { value: 'wrong@example.com' } });
    fireEvent.change(passwordInput, { target: { value: 'wrongpass' } });

    const submitButton = screen.getByRole('button', { name: /iniciar sesión/i });
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(screen.getByText('Credenciales inválidas')).toBeInTheDocument();
    });

    expect(mockOnLoginSuccess).not.toHaveBeenCalled();
  });

  test('deshabilita el botón mientras está cargando', async () => {
    mockedAxios.post.mockImplementation(() => 
      new Promise(resolve => setTimeout(resolve, 100))
    );

    render(<Login onLoginSuccess={mockOnLoginSuccess} />);

    const emailInput = screen.getByLabelText(/email/i) as HTMLInputElement;
    const passwordInput = screen.getByLabelText(/password/i) as HTMLInputElement;

    fireEvent.change(emailInput, { target: { value: 'test@example.com' } });
    fireEvent.change(passwordInput, { target: { value: '123456' } });

    const submitButton = screen.getByRole('button', { name: /iniciar sesión/i });
    fireEvent.click(submitButton);

    // El botón debe estar deshabilitado
    expect(submitButton).toBeDisabled();
  });
});