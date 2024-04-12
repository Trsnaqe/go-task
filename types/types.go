package types

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	CreateUser(User) error
	UpdateUser(userID int, updates UpdateUserPayload) error
	ChangePassword(userID int, oldPassword string, newPassword string) error
}

type TaskStore interface {
	GetTasks() ([]Task, error)
	CreateTask(Task) error
	UpdateTask(taskID int, updates UpdateTaskPayload) error
	GetTaskByID(taskID int) (*Task, error)
	DeleteTask(taskID int) error
	ProgressTask(taskID int) error
	GetTasksByStatus(status TaskStatus) ([]Task, error)
}
type TaskStatus string

const (
	StatusPending    TaskStatus = "pending"
	StatusInProgress TaskStatus = "in_progress"
	StatusCompleted  TaskStatus = "completed"
)

type Task struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	CreatedAt   string     `json:"created_at"`
	UpdatedAt   string     `json:"updated_at"`
}

type UpdateTaskPayload struct {
	Title       *string     `json:"title" validate:"omitempty,min=3,max=32"`
	Description *string     `json:"description" validate:"omitempty,min=3,max=255"`
	Status      *TaskStatus `json:"status" validate:"omitempty,oneof=pending in_progress completed"`
}

type UpdateUserPayload struct {
	Email        *string `json:"email" validate:"omitempty,email"`
	Password     *string `json:"password" validate:"omitempty,min=6,max=32"`
	RefreshToken *string `json:"refresh_token" validate:"omitempty"`
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RegisterUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=32"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=32"`
}

type CreateTaskPayload struct {
	Title       string     `json:"title" validate:"required,min=3,max=32"`
	Description string     `json:"description" validate:"required,min=3,max=255"`
	Status      TaskStatus `json:"status" validate:"required,oneof=pending in_progress completed"`
}

type ChangePasswordPayload struct {
	OldPassword string `json:"old_password" validate:"required,min=6,max=32"`
	NewPassword string `json:"new_password" validate:"required,min=6,max=32"`
}

type User struct {
	ID           int     `json:"id"`
	Email        string  `json:"email"`
	Password     string  `json:"password"`
	RefreshToken *string `json:"refresh_token"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

type ErrorResponse struct {
	Error      string `json:"error"`
	StatusCode int    `json:"status_code,omitempty"`
}

type contextKey string

const UserKey contextKey = "userID"
