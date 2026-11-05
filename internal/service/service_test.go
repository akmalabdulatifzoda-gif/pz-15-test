package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// ---------- Stub: заранее известные данные, ничего не проверяет ----------

type stubRepo struct {
	users map[string]User
}

func (r stubRepo) ByEmail(email string) (User, error) {
	u, ok := r.users[email]
	if !ok {
		return User{}, ErrNotFound
	}
	return u, nil
}

func newStubService() *Service {
	return New(stubRepo{users: map[string]User{
		"akmal@example.com": {ID: 1, Email: "akmal@example.com"},
		"test@mirea.ru":     {ID: 42, Email: "test@mirea.ru"},
	}})
}

func TestService_FindIDByEmail(t *testing.T) {
	svc := newStubService()

	t.Run("найден", func(t *testing.T) {
		id, err := svc.FindIDByEmail("akmal@example.com")
		require.NoError(t, err)
		assert.Equal(t, int64(1), id)
	})

	t.Run("не найден", func(t *testing.T) {
		id, err := svc.FindIDByEmail("nobody@example.com")
		assert.ErrorIs(t, err, ErrNotFound)
		assert.Zero(t, id)
	})
}

func TestService_FindIDByEmail_Table(t *testing.T) {
	svc := newStubService()
	tests := []struct {
		name    string
		email   string
		wantID  int64
		wantErr error
	}{
		{"точное совпадение", "test@mirea.ru", 42, nil},
		{"регистр не важен", "Test@MIREA.ru", 42, nil},
		{"пробелы по краям обрезаются", "  akmal@example.com\n", 1, nil},
		{"неизвестный пользователь", "ghost@mirea.ru", 0, ErrNotFound},
		{"пустая строка", "", 0, ErrInvalidEmail},
		{"только пробелы", "   ", 0, ErrInvalidEmail},
		{"нет @", "akmal.example.com", 0, ErrInvalidEmail},
		{"нет локальной части", "@example.com", 0, ErrInvalidEmail},
		{"нет домена", "akmal@", 0, ErrInvalidEmail},
		{"домен без точки", "akmal@localhost", 0, ErrInvalidEmail},
		{"две @", "a@b@example.com", 0, ErrInvalidEmail},
		{"пробел внутри", "ak mal@example.com", 0, ErrInvalidEmail},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := svc.FindIDByEmail(tt.email)
			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Zero(t, id)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantID, id)
		})
	}
}

// ---------- Ручной mock: запоминает вызовы, чтобы проверить взаимодействие ----------

type spyRepo struct {
	calls []string
	user  User
	err   error
}

func (r *spyRepo) ByEmail(email string) (User, error) {
	r.calls = append(r.calls, email)
	return r.user, r.err
}

func TestService_PassesNormalizedEmailToRepo(t *testing.T) {
	repo := &spyRepo{user: User{ID: 7}}
	id, err := New(repo).FindIDByEmail("  Akmal@Example.COM ")
	require.NoError(t, err)
	assert.Equal(t, int64(7), id)
	assert.Equal(t, []string{"akmal@example.com"}, repo.calls, "в репозиторий должен уйти нормализованный email, ровно один раз")
}

func TestService_InvalidEmailDoesNotHitRepo(t *testing.T) {
	repo := &spyRepo{}
	_, err := New(repo).FindIDByEmail("not-an-email")
	assert.ErrorIs(t, err, ErrInvalidEmail)
	assert.Empty(t, repo.calls, "при некорректном email обращения к БД быть не должно")
}

// ---------- testify/mock: ожидания вызовов и аргументов ----------

type mockRepo struct{ mock.Mock }

func (m *mockRepo) ByEmail(email string) (User, error) {
	args := m.Called(email)
	return args.Get(0).(User), args.Error(1)
}

func TestService_RepoErrorIsWrapped(t *testing.T) {
	dbErr := errors.New("connection refused")
	repo := new(mockRepo)
	repo.On("ByEmail", "akmal@example.com").Return(User{}, dbErr).Once()

	id, err := New(repo).FindIDByEmail("akmal@example.com")

	require.Error(t, err)
	assert.ErrorIs(t, err, dbErr, "исходная ошибка доступна через errors.Is")
	assert.Contains(t, err.Error(), "find user by email", "ошибка дополнена контекстом")
	assert.Zero(t, id)
	repo.AssertExpectations(t)
}

func TestService_WithTestifyMock_Found(t *testing.T) {
	repo := new(mockRepo)
	repo.On("ByEmail", "test@mirea.ru").Return(User{ID: 42, Email: "test@mirea.ru"}, nil).Once()

	id, err := New(repo).FindIDByEmail("TEST@mirea.ru")

	require.NoError(t, err)
	assert.Equal(t, int64(42), id)
	repo.AssertExpectations(t)
	repo.AssertNumberOfCalls(t, "ByEmail", 1)
}

// ---------- NormalizeEmail отдельно ----------

func TestNormalizeEmail(t *testing.T) {
	got, err := NormalizeEmail(" User@Mail.RU ")
	require.NoError(t, err)
	assert.Equal(t, "user@mail.ru", got)

	_, err = NormalizeEmail("bad")
	require.ErrorIs(t, err, ErrInvalidEmail)
	assert.Contains(t, err.Error(), `"bad"`, "в сообщении указан исходный ввод")
}

func BenchmarkFindIDByEmail(b *testing.B) {
	svc := newStubService()
	for i := 0; i < b.N; i++ {
		_, _ = svc.FindIDByEmail("Test@MIREA.ru")
	}
}
