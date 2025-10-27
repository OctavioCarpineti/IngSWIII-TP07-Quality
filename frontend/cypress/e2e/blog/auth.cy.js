describe('Authentication Flow', () => {
    beforeEach(() => {
        cy.visit('/')
    })

    it('debería mostrar el formulario de login por defecto', () => {
        cy.get('h2').should('contain', 'Iniciar Sesión')
        cy.get('input#email').should('be.visible')
        cy.get('input#password').should('be.visible')
        cy.get('button[type="submit"]').should('contain', 'Iniciar Sesión')
    })

    it('debería cambiar entre login y registro', () => {
        // Cambiar a registro
        cy.contains('¿No tienes cuenta? Regístrate').click()
        cy.get('h2').should('contain', 'Registrarse')
        cy.get('input#username').should('be.visible')
        cy.get('button[type="submit"]').should('contain', 'Registrarse')

        // Volver a login
        cy.contains('¿Ya tienes cuenta? Inicia sesión').click()
        cy.get('h2').should('contain', 'Iniciar Sesión')
        cy.get('input#username').should('not.exist')
    })

    it('debería mostrar error con credenciales inválidas', () => {
        cy.intercept('POST', '**/api/auth/login', {
            statusCode: 401,
            body: { error: 'Credenciales inválidas' }
        })

        cy.get('input#email').type('invalid@example.com')
        cy.get('input#password').type('wrongpass')
        cy.get('button[type="submit"]').click()

        cy.get('.error-message').should('be.visible')
            .and('contain', 'Credenciales inválidas')
    })

    it('debería hacer login exitoso', () => {
        cy.intercept('POST', '**/api/auth/login', {
            statusCode: 200,
            body: {
                id: 1,
                email: 'test@example.com',
                username: 'testuser'
            }
        }).as('loginRequest')

        cy.intercept('GET', '**/api/posts', {
            statusCode: 200,
            body: []
        })

        cy.get('input#email').type('test@example.com')
        cy.get('input#password').type('123456')
        cy.get('button[type="submit"]').click()

        cy.wait('@loginRequest')

        // Verificar que muestra la app
        cy.contains('Mini Red Social').should('be.visible')
        cy.contains('Hola, @testuser').should('be.visible')
    })

    it('debería registrarse exitosamente', () => {
        cy.intercept('POST', '**/api/auth/register', {
            statusCode: 201,
            body: {
                id: 2,
                email: 'newuser@example.com',
                username: 'newuser'
            }
        }).as('registerRequest')

        cy.intercept('GET', '**/api/posts', {
            statusCode: 200,
            body: []
        })

        // Cambiar a modo registro
        cy.contains('¿No tienes cuenta? Regístrate').click()

        cy.get('input#email').type('newuser@example.com')
        cy.get('input#username').type('newuser')
        cy.get('input#password').type('123456')
        cy.get('button[type="submit"]').click()

        cy.wait('@registerRequest')

        // Verificar que muestra la app
        cy.contains('Hola, @newuser').should('be.visible')
    })
})