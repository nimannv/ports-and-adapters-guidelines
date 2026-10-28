
# General Principles

1.  **`_test.go` Files:** Test files live alongside the code they test (e.g., `user.go` and `user_test.go` in the same directory).
2.  **`TestXxx` Functions:** Test functions must start with `Test` followed by an uppercase letter (e.g., `TestCreateUser`). They take `*testing.T` as an argument.
3.  **`t.Run` for Subtests:** Use `t.Run` to group related tests or test different scenarios for the same function.
4.  **`t.Parallel()`:** For subtests that don't share mutable state, use `t.Parallel()` to run them concurrently.
5.  **Assertions with `stretchr/testify`:**
    *   Go's standard library provides `t.Error`, `t.Fatal`, `t.Log`, etc. For **richer and more readable assertions**, use `github.com/stretchr/testify`.
    *   **`require`**: Functions like `require.NoError(t, err)` or `require.Equal(t, expected, actual)` will **stop the test immediately** (like `t.Fatal`) if the assertion fails. Use `require` for preconditions that subsequent test steps depend on.
    *   **`assert`**: Functions like `assert.Equal(t, expected, actual)` or `assert.ErrorIs(t, err, targetErr)` will **report the failure but continue executing** the test function (like `t.Error`). Use `assert` when you want to check multiple conditions even if one fails.

    ```go
    // Example testify assertion
    import (
        "testing"
        "github.com/stretchr/testify/assert"
        "github.com/stretchr/testify/require"
    )

    func TestSomething(t *testing.T) {
        // Use require for critical preconditions
        result, err := someFuncThatMightFail()
        require.NoError(t, err, "someFuncThatMightFail should not return an error") // Test stops here if err is not nil

        // Use assert for general checks
        assert.Equal(t, "expectedValue", result.Field, "Field value mismatch") // Test continues even if this fails
        assert.True(t, result.IsValid(), "Result should be valid")
    }
    ```

6.  **Mocks with `vektra/mockery`:**
    *   For isolating components, we'll often need to create mock implementations of interfaces.
    *   **`github.com/vektra/mockery`** automates mock generation from our Go interfaces.
    *   **Workflow:**
        1.  Define your interface (e.g., `port.UserRepository`).
        2.  Run `mockery` from your terminal, pointing it to your interface:
            ```bash
            mockery --name UserRepository --dir internal/port --output internal/core/application/port/outbound/mocks --filename user_repository.go --keeptree
            ```
            *(Adjust `--output` and `--filename` to place the mock where it's needed for your tests, e.g., in a `mocks` sub-directory within the `application_test` package context.)*
        3.  `mockery` will generate a file (e.g., `internal/core/application/port/outbound/mocks/user_repository.go`) containing a `MockUserRepository` struct.
        4.  In your tests, instantiate this generated mock and use its `On()` method to define expected calls and return values, and `AssertExpectations()` to verify all expectations were met.

    ```go
    // Example mockery usage (in a test file)
    import (
        "testing"
        "github.com/stretchr/testify/mock"
        "your_module/internal/application/mocks" // Import the generated mocks
        "your_module/internal/domain"
    )

    func TestMyUseCase(t *testing.T) {
        mockRepo := new(mocks.UserRepository) // Instantiate the generated mock

        // Define expected behavior: when Save is called with any User, return nil error
        mockRepo.On("Save", mock.AnythingOfType("*domain.User")).Return(nil).Once()

        // Define expected behavior: when FindByID is called with "123", return a user and no error
        userToReturn := &domain.User{ID: "123", Email: "test@example.com"}
        mockRepo.On("FindByID", "123").Return(userToReturn, nil).Once()

        // Inject the mock into your use case
        useCase := application.NewMyUseCase(mockRepo)

        // Run your test logic
        _, err := useCase.Execute(application.MyCommand{UserID: "123"})
        require.NoError(t, err)

        // Assert that all expectations set on the mock were met
        mockRepo.AssertExpectations(t)
        // You can also assert specific calls:
        mockRepo.AssertCalled(t, "Save", mock.AnythingOfType("*domain.User"))
        mockRepo.AssertNumberOfCalls(t, "FindByID", 1)
    }
    ```

---

## Testing by Layer in a Port & Adapter Architecture


### 1. Domain Layer (`internal/domain`)

*   **Approach:**
    *   Directly instantiate domain objects.
    *   Call their methods.
    *   Use `testify/require` for critical initial checks and `testify/assert` for general assertions.

**Example (e.g., `internal/domain/user.go` and `internal/domain/user_test.go`):**

```go
// internal/domain/user.go (remains the same)
package domain

import (
	"errors"
	"regexp"
)

var ErrInvalidEmail = errors.New("invalid email address")

type User struct {
	ID    string
	Email string
	Name  string
}

func NewUser(id, email, name string) (*User, error) {
	if !isValidEmail(email) {
		return nil, ErrInvalidEmail
	}
	return &User{ID: id, Email: email, Name: name}, nil
}

func (u *User) ChangeEmail(newEmail string) error {
	if !isValidEmail(newEmail) {
		return ErrInvalidEmail
	}
	u.Email = newEmail
	return nil
}

func isValidEmail(email string) bool {
	return regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`).MatchString(email)
}

// internal/domain/user_test.go
package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert" // For general assertions
	"github.com/stretchr/testify/require" // For fail-fast assertions

	"your_module/internal/domain" // Adjust module path
)

func TestNewUser(t *testing.T) {
	t.Run("Valid user creation", func(t *testing.T) {
		user, err := domain.NewUser("123", "test@example.com", "Test User")
		require.NoError(t, err) // Assert no error and stop if there is one
		require.NotNil(t, user)  // Assert user is not nil

		assert.Equal(t, "123", user.ID)
		assert.Equal(t, "test@example.com", user.Email)
		assert.Equal(t, "Test User", user.Name)
	})

	t.Run("Invalid email", func(t *testing.T) {
		user, err := domain.NewUser("123", "invalid-email", "Test User")
		assert.Error(t, err)                               // Assert an error occurred
		assert.Nil(t, user)                                 // Assert user is nil
		assert.ErrorIs(t, err, domain.ErrInvalidEmail) // Assert the specific error type
	})
}

func TestUserChangeEmail(t *testing.T) {
	user, err := domain.NewUser("1", "old@example.com", "Name")
	require.NoError(t, err) // Precondition for these tests

	t.Run("Successfully change email", func(t *testing.T) {
		err := user.ChangeEmail("new@example.com")
		require.NoError(t, err)
		assert.Equal(t, "new@example.com", user.Email)
	})

	t.Run("Fail to change email with invalid format", func(t *testing.T) {
		originalEmail := user.Email // Capture current state for rollback check
		err := user.ChangeEmail("bad-email")
		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrInvalidEmail)
		assert.Equal(t, originalEmail, user.Email, "email should not have changed on error") // Ensure no side effects
	})
}
```

### 2. Application Layer (`internal/application`)

*   **Approach:**
    *   **Generate mocks for outbound ports** using `mockery`.
    *   Inject these generated mocks into the application service constructor.
    *   Use `mock.On().Return()` to define mock behavior.
    *   Call the application service's method.
    *   Use `testify/require` and `testify/assert` for assertions.
    *   **Crucially, call `mock.AssertExpectations(t)`** to verify that all expected calls to the mocks actually happened.

**Example (e.g., `internal/port/user_repository.go`, `internal/application/create_user.go`, and `internal/application/create_user_test.go`):**

```go
// internal/port/user_repository.go (remains the same)
package port

import "your_module/internal/domain"

type UserRepository interface {
	Save(user *domain.User) error
	FindByID(id string) (*domain.User, error)
}

// internal/application/create_user.go (remains the same)
package application

import (
	"errors"
	"your_module/internal/domain"
	"your_module/internal/port"
)

var ErrUserAlreadyExists = errors.New("user with this ID already exists")

type CreateUserCommand struct {
	ID    string
	Email string
	Name  string
}

type CreateUserUseCase struct {
	userRepo port.UserRepository
}

func NewCreateUserUseCase(userRepo port.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{userRepo: userRepo}
}

func (uc *CreateUserUseCase) Execute(cmd CreateUserCommand) (*domain.User, error) {
	// Check if user already exists (example business rule)
	existingUser, err := uc.userRepo.FindByID(cmd.ID)
	if err != nil && err.Error() != "user not found" { // Mock will return this specific error
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	user, err := domain.NewUser(cmd.ID, cmd.Email, cmd.Name)
	if err != nil {
		return nil, err
	}

	if err := uc.userRepo.Save(user); err != nil {
		return nil, err
	}
	return user, nil
}

// --- GENERATE MOCKS ---
// From your project root, run:
// mockery --name UserRepository --dir internal/port --output internal/application/mocks --filename user_repository.go --keeptree

// This will create a file like: internal/application/mocks/user_repository.go
// It will contain a struct like:
// type UserRepository struct {
//    mock.Mock
// }
// ... with generated methods.

// internal/application/create_user_test.go (UPDATED)
package application_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock" // For mock.Anything, mock.AnythingOfType, etc.
	"github.com/stretchr/testify/require"

	"your_module/internal/application"
	"your_module/internal/application/mocks" // Import the generated mock
	"your_module/internal/domain"
)

func TestCreateUserUseCase(t *testing.T) {
	t.Run("Successfully creates a user", func(t *testing.T) {
		mockRepo := new(mocks.UserRepository) // Use the generated mock

		// Define mock behavior:
		// 1. FindByID is called once with cmd.ID, returns nil and "user not found" error
		mockRepo.On("FindByID", mock.AnythingOfType("string")).Return(
			(*domain.User)(nil), errors.New("user not found"),
		).Once()
		// 2. Save is called once with any *domain.User, returns nil error
		mockRepo.On("Save", mock.AnythingOfType("*domain.User")).Return(nil).Once()

		useCase := application.NewCreateUserUseCase(mockRepo)

		cmd := application.CreateUserCommand{
			ID:    "user-123",
			Email: "test@example.com",
			Name:  "Test User",
		}

		user, err := useCase.Execute(cmd)
		require.NoError(t, err)
		require.NotNil(t, user)
		assert.Equal(t, cmd.ID, user.ID)
		assert.Equal(t, cmd.Email, user.Email)

		mockRepo.AssertExpectations(t) // Verify all defined expectations were met
	})

	t.Run("Fails if user already exists", func(t *tester.T) {
		mockRepo := new(mocks.UserRepository)
		existingUser, _ := domain.NewUser("user-123", "existing@example.com", "Existing User")

		// Define mock behavior:
		// 1. FindByID is called once with cmd.ID, returns an existing user
		mockRepo.On("FindByID", mock.AnythingOfType("string")).Return(existingUser, nil).Once()
		// 2. Save should NOT be called
		mockRepo.On("Save", mock.AnythingOfType("*domain.User")).Return(errors.New("should not be called")).Maybe() // Use Maybe to prevent panic if called unexpectedly

		useCase := application.NewCreateUserUseCase(mockRepo)

		cmd := application.CreateUserCommand{
			ID:    "user-123",
			Email: "test@example.com",
			Name:  "Test User",
		}

		user, err := useCase.Execute(cmd)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.ErrorIs(t, err, application.ErrUserAlreadyExists)

		mockRepo.AssertExpectations(t) // Verify expectations
		mockRepo.AssertNotCalled(t, "Save", mock.Anything) // Explicitly ensure Save was not called
	})

	t.Run("Fails for invalid email in command", func(t *testing.T) {
		mockRepo := new(mocks.UserRepository)

		// Define mock behavior:
		// 1. FindByID is called once, returns not found (pre-check is still valid)
		mockRepo.On("FindByID", mock.AnythingOfType("string")).Return(
			(*domain.User)(nil), errors.New("user not found"),
		).Once()
		// 2. Save should NOT be called as domain.NewUser will fail first
		mockRepo.On("Save", mock.AnythingOfType("*domain.User")).Return(errors.New("should not be called")).Maybe()

		useCase := application.NewCreateUserUseCase(mockRepo)

		cmd := application.CreateUserCommand{
			ID:    "user-123",
			Email: "bad-email", // Invalid email
			Name:  "Test User",
		}

		user, err := useCase.Execute(cmd)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.ErrorIs(t, err, domain.ErrInvalidEmail)

		mockRepo.AssertExpectations(t)
		mockRepo.AssertCalled(t, "FindByID", cmd.ID) // FindByID might still be called first
		mockRepo.AssertNotCalled(t, "Save", mock.Anything)
	})

	t.Run("Fails if repository save fails", func(t *testing.T) {
		mockRepo := new(mocks.UserRepository)
		repoError := errors.New("database connection failed")

		mockRepo.On("FindByID", mock.AnythingOfType("string")).Return(
			(*domain.User)(nil), errors.New("user not found"),
		).Once()
		mockRepo.On("Save", mock.AnythingOfType("*domain.User")).Return(repoError).Once()

		useCase := application.NewCreateUserUseCase(mockRepo)

		cmd := application.CreateUserCommand{
			ID:    "user-123",
			Email: "test@example.com",
			Name:  "Test User",
		}

		user, err := useCase.Execute(cmd)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.ErrorIs(t, err, repoError)

		mockRepo.AssertExpectations(t)
	})
}
```

### 3. Port Layer (`internal/port`)

*   **What to test:** Nothing. This layer only contains interfaces, which define contracts but no behavior. There's no logic to test here.

### 4. Adapter Layer

#### A. Primary (Inbound) Adapters (e.g., HTTP, gRPC) (`adapters/primary/http`)

*   **Approach (HTTP example):**
    *   **Generate mocks for inbound ports** (Application Services) using `mockery`.
    *   Use `net/http/httptest` to simulate HTTP requests.
    *   Use `mock.On().Return()` to define mock behavior for the application service.
    *   Use `testify/assert` to check HTTP status, headers, and body.
    *   Call `mock.AssertExpectations(t)` to verify the application service was called correctly.

**Example (e.g., `adapters/primary/http/user_handler.go` and `adapters/primary/http/user_handler_test.go`):**

```go
// internal/application/ports.go (just an example for the interface the handler needs)
package application

import (
	"your_module/internal/domain"
)

type CreateUserUseCasePort interface { // Interface for the inbound port
	Execute(cmd CreateUserCommand) (*domain.User, error)
}

// adapters/primary/http/user_handler.go (remains mostly same)
package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"your_module/internal/application"
	"your_module/internal/domain"
)

type UserHandler struct {
	createUserUseCase application.CreateUserUseCasePort // Depends on the port interface
}

func NewUserHandler(createUserUC application.CreateUserUseCasePort) *UserHandler {
	return &UserHandler{createUserUseCase: createUserUC}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	cmd := application.CreateUserCommand{
		ID:    reqBody.ID,
		Email: reqBody.Email,
		Name:  reqBody.Name,
	}

	user, err := h.createUserUseCase.Execute(cmd) // Calls the inbound port
	if err != nil {
		if errors.Is(err, domain.ErrInvalidEmail) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, application.ErrUserAlreadyExists) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// --- GENERATE MOCKS ---
// From your project root, run:
// mockery --name CreateUserUseCasePort --dir internal/application --output adapters/primary/http/mocks --filename create_user_usecase_port.go --keeptree

// This will create a file like: adapters/primary/http/mocks/create_user_usecase_port.go
// It will contain a struct like:
// type CreateUserUseCasePort struct {
//    mock.Mock
// }
// ...

// adapters/primary/http/user_handler_test.go (UPDATED)
package http_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert