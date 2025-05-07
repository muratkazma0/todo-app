# Todo App API

A RESTful API for managing todos and todo items, built with Go and Gin framework.


## API Endpoints

### Authentication
- `POST /login` - Login and get JWT token

### Users
- `POST /api/users` - Create a new user
- `GET /api/users/:id` - Get user by ID
- `PUT /api/users/:id` - Update user
- `DELETE /api/users/:id` - Delete user

### Todos
- `POST /api/todos` - Create a new todo
- `GET /api/todos` - Get all todos
- `GET /api/todos/:id` - Get todo by ID
- `PUT /api/todos/:id` - Update todo
- `DELETE /api/todos/:id` - Delete todo

### Todo Items
- `POST /api/todos/items/:todo_id` - Create a new todo item
- `GET /api/todos/items/:todo_id` - Get all items for a todo
- `PUT /api/todos/items/:todo_id/:item_id` - Update todo item
- `DELETE /api/todos/items/:todo_id/:item_id` - Delete todo item

## Features

- JWT-based authentication
- Role-based authorization (admin/user)
- Soft delete functionality
- Todo completion tracking
- Todo item completion tracking
- Completion percentage calculation
- Admin-specific features

## Default Users

For testing purposes, the following users are created by default:

1. Admin User:
   - Email: admin@example.com
   - Password: admin123
   - Role: admin

2. Regular User:
   - Email: user@example.com
   - Password: user123
   - Role: user

## Getting Started

1. Clone the repository
2. Install dependencies: `go mod download`
3. Run the application: `go run main.go`
4. The server will start on `http://localhost:8080`

## Deployment
The project is hosted in a private GitHub repository at: [https://github.com/muratkazma0/todo-app]

## Development Notes

The project follows a chronological development approach, starting from basic setup and gradually adding more complex features. Each major step is documented with detailed commit messages that explain the changes and additions made to the project.

The development process ensures that:
1. Basic functionality is implemented first
2. Security features are added early
3. Core business logic is implemented
4. API layer is properly structured
5. Admin features are added
6. Testing and documentation are completed
7. Final optimizations are made

This approach allows for a working system at each stage of development while maintaining code quality and proper documentation.

## Varsayılan Kullanıcılar

Uygulama ilk çalıştırıldığında aşağıdaki varsayılan kullanıcılar otomatik olarak oluşturulur:

### Admin Kullanıcısı
- **Username**: `admin`
- **Password**: `admin123`
- **Role**: `admin`
- **Varsayılan Todo**: "Admin Todo"
- **Varsayılan Todo Item**: "Admin Todo Item"

### Normal Kullanıcı
- **Username**: `user`
- **Password**: `user123`
- **Role**: `user`
- **Varsayılan Todo**: "User Todo"
- **Varsayılan Todo Item**: "User Todo Item"

## API Endpoints

### Authentication

#### Login
- **URL**: `/login`
- **Method**: `POST`
- **Auth Required**: No
- **Body**:
  ```json
  {
    "username": "string",
    "password": "string"
  }
  ```
- **Success Response**:
  ```json
  {
    "token": "string",
    "user": {
      "id": "integer",
      "username": "string",
      "role": "string"
    }
  }
  ```
- **Error Response**:
  ```json
  {
    "error": "invalid credentials"
  }
  ```

### Users

#### Create User
- **URL**: `/api/users`
- **Method**: `POST`
- **Auth Required**: No
- **Body**:
  ```json
  {
    "username": "string",
    "password": "string",
    "role": "string"
  }
  ```
- **Success Response**: `201 Created`
  ```json
  {
    "id": "integer",
    "username": "string",
    "role": "string"
  }
  ```

#### Get All Users
- **URL**: `/api/users`
- **Method**: `GET`
- **Auth Required**: Yes
- **Admin Required**: Yes
- **Success Response**: `200 OK`
  ```json
  [
    {
      "id": "integer",
      "username": "string",
      "role": "string"
    }
  ]
  ```

#### Get User by ID
- **URL**: `/api/users/:id`
- **Method**: `GET`
- **Auth Required**: Yes
- **URL Parameters**: `id=[integer]`
- **Success Response**: `200 OK`
  ```json
  {
    "id": "integer",
    "username": "string",
    "role": "string"
  }
  ```

#### Get User by Username
- **URL**: `/api/users/username/:username`
- **Method**: `GET`
- **Auth Required**: Yes
- **URL Parameters**: `username=[string]`
- **Success Response**: `200 OK`
  ```json
  {
    "id": "integer",
    "username": "string",
    "role": "string"
  }
  ```

#### Update User
- **URL**: `/api/users/:id`
- **Method**: `PUT`
- **Auth Required**: Yes
- **URL Parameters**: `id=[integer]`
- **Body**:
  ```json
  {
    "username": "string",
    "password": "string",
    "role": "string"
  }
  ```
- **Success Response**: `200 OK`
  ```json
  {
    "id": "integer",
    "username": "string",
    "role": "string"
  }
  ```

#### Delete User
- **URL**: `/api/users/:id`
- **Method**: `DELETE`
- **Auth Required**: Yes
- **Admin Required**: Yes
- **URL Parameters**: `id=[integer]`
- **Success Response**: `200 OK`
  ```json
  {
    "message": "user deleted"
  }
  ```

### Todos

#### Create Todo
- **URL**: `/api/todos`
- **Method**: `POST`
- **Auth Required**: Yes
- **Body**:
  ```json
  {
    "title": "string",
    "description": "string"
  }
  ```
- **Success Response**: `201 Created`
  ```json
  {
    "id": "integer",
    "title": "string",
    "description": "string",
    "completed": "boolean",
    "user_id": "integer",
    "created_at": "datetime",
    "updated_at": "datetime"
  }
  ```

#### Get All Todos
- **URL**: `/api/todos`
- **Method**: `GET`
- **Auth Required**: Yes
- **Success Response**: `200 OK`
  ```json
  [
    {
      "id": "integer",
      "title": "string",
      "description": "string",
      "completed": "boolean",
      "user_id": "integer",
      "created_at": "datetime",
      "updated_at": "datetime"
    }
  ]
  ```
- **Notes**: 
  - Normal kullanıcılar sadece kendi todolarını görür
  - Admin tüm todoları görür (silinmiş olanlar dahil)

#### Get Todo by ID
- **URL**: `/api/todos/:id`
- **Method**: `GET`
- **Auth Required**: Yes
- **URL Parameters**: `id=[integer]`
- **Success Response**: `200 OK`
  ```json
  {
    "id": "integer",
    "title": "string",
    "description": "string",
    "completed": "boolean",
    "user_id": "integer",
    "created_at": "datetime",
    "updated_at": "datetime"
  }
  ```
- **Notes**: 
  - Normal kullanıcılar sadece kendi todolarını görebilir
  - Admin tüm todoları görebilir (silinmiş olanlar dahil)

#### Update Todo
- **URL**: `/api/todos/:id`
- **Method**: `PUT`
- **Auth Required**: Yes
- **URL Parameters**: `id=[integer]`
- **Body**:
  ```json
  {
    "title": "string",
    "description": "string",
    "completed": "boolean"
  }
  ```
- **Success Response**: `200 OK`
  ```json
  {
    "id": "integer",
    "title": "string",
    "description": "string",
    "completed": "boolean",
    "user_id": "integer",
    "created_at": "datetime",
    "updated_at": "datetime"
  }
  ```
- **Notes**: 
  - Normal kullanıcılar sadece kendi todolarını güncelleyebilir
  - Admin tüm todoları güncelleyebilir

#### Delete Todo
- **URL**: `/api/todos/:id`
- **Method**: `DELETE`
- **Auth Required**: Yes
- **URL Parameters**: `id=[integer]`
- **Success Response**: `200 OK`
  ```json
  {
    "message": "todo deleted"
  }
  ```
- **Notes**: 
  - Normal kullanıcılar sadece kendi todolarını silebilir
  - Admin tüm todoları silebilir
  - Silme işlemi soft delete olarak gerçekleşir

### Todo Items

#### Create Todo Item
- **URL**: `/api/todos/items/:todo_id`
- **Method**: `POST`
- **Auth Required**: Yes
- **URL Parameters**: `todo_id=[integer]`
- **Body**:
  ```json
  {
    "title": "string",
    "description": "string"
  }
  ```
- **Success Response**: `201 Created`
  ```json
  {
    "id": "integer",
    "title": "string",
    "description": "string",
    "completed": "boolean",
    "todo_id": "integer",
    "created_at": "datetime",
    "updated_at": "datetime"
  }
  ```

#### Get Todo Items
- **URL**: `/api/todos/items/:todo_id`
- **Method**: `GET`
- **Auth Required**: Yes
- **URL Parameters**: `todo_id=[integer]`
- **Success Response**: `200 OK`
  ```json
  [
    {
      "id": "integer",
      "title": "string",
      "description": "string",
      "completed": "boolean",
      "todo_id": "integer",
      "created_at": "datetime",
      "updated_at": "datetime"
    }
  ]
  ```
- **Notes**: 
  - Normal kullanıcılar sadece kendi todo itemlarını görür
  - Admin tüm todo itemları görür (silinmiş olanlar dahil)

#### Update Todo Item
- **URL**: `/api/todos/items/:todo_id/:item_id`
- **Method**: `PUT`
- **Auth Required**: Yes
- **URL Parameters**: 
  - `todo_id=[integer]`
  - `item_id=[integer]`
- **Body**:
  ```json
  {
    "title": "string",
    "description": "string",
    "completed": "boolean"
  }
  ```
- **Success Response**: `200 OK`
  ```json
  {
    "id": "integer",
    "title": "string",
    "description": "string",
    "completed": "boolean",
    "todo_id": "integer",
    "created_at": "datetime",
    "updated_at": "datetime"
  }
  ```
- **Notes**: 
  - Normal kullanıcılar sadece kendi todo itemlarını güncelleyebilir
  - Admin tüm todo itemları güncelleyebilir

#### Delete Todo Item
- **URL**: `/api/todos/items/:todo_id/:item_id`
- **Method**: `DELETE`
- **Auth Required**: Yes
- **URL Parameters**: 
  - `todo_id=[integer]`
  - `item_id=[integer]`
- **Success Response**: `200 OK`
  ```json
  {
    "message": "todo item deleted"
  }
  ```
- **Notes**: 
  - Normal kullanıcılar sadece kendi todo itemlarını silebilir
  - Admin tüm todo itemları silebilir
  - Silme işlemi soft delete olarak gerçekleşir

## Authentication

Tüm korumalı rotalar için JWT token gereklidir. Token'ı HTTP header'da şu şekilde göndermelisiniz:

```
Authorization: Bearer <token>
```

## Error Responses

Tüm hata durumlarında aşağıdaki formatta yanıt alırsınız:

```json
{
  "error": "error message"
}
```

Yaygın hata kodları:
- `400 Bad Request`: Geçersiz istek parametreleri
- `401 Unauthorized`: Geçersiz veya eksik token
- `403 Forbidden`: Yetkisiz erişim
- `404 Not Found`: Kaynak bulunamadı
- `500 Internal Server Error`: Sunucu hatası 
