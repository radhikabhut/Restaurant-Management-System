# Run Instructions for Restaurant Management System

This document provides step-by-step instructions on how to set up and run both the backend and frontend of the project.

## Prerequisites

- **Go** (for the backend)
- **PostgreSQL** (for the database)
- **Flutter** (for the frontend)
- **Git** (optional, for version control)

---

## 1. Database Setup

The project uses PostgreSQL as its database. Follow these steps to set it up:

1.  **Start PostgreSQL**: Ensure your PostgreSQL service is running.
2.  **Create Database**: Create a new database named `restaurant`.
    ```bash
    createdb restaurant
    ```
3.  **Create Schema**: Run the following SQL commands to create the necessary tables. You can use `psql` or any database management tool like pgAdmin or DBeaver.

    ```sql
    CREATE TABLE users (
        id UUID PRIMARY KEY,
        name TEXT NOT NULL,
        email TEXT UNIQUE NOT NULL,
        password TEXT NOT NULL,
        role TEXT NOT NULL,
        createdAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updateAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );

    CREATE TABLE menus (
        id UUID PRIMARY KEY,
        name TEXT NOT NULL,
        price DECIMAL(10, 2) NOT NULL,
        category TEXT NOT NULL,
        isAvailable BOOLEAN DEFAULT TRUE,
        createdAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updatedAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );

    CREATE TABLE orders (
        id UUID PRIMARY KEY,
        userId UUID REFERENCES users(id),
        totalAmount DECIMAL(10, 2) NOT NULL,
        status TEXT NOT NULL, -- 'pending', 'preparing', 'completed', 'cancelled'
        tableNumber TEXT,
        createdAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updatedAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );

    CREATE TABLE order_items (
        id UUID PRIMARY KEY,
        orderId UUID REFERENCES orders(id),
        menuId UUID REFERENCES menus(id),
        quantity INTEGER NOT NULL,
        price DECIMAL(10, 2) NOT NULL,
        createdAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updatedAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );

    CREATE TABLE role_permission (
        role TEXT NOT NULL,
        permission TEXT NOT NULL,
        PRIMARY KEY (role, permission)
    );
    ```

4.  **Seed Initial Data (Recommended)**: You might want to add some initial roles and permissions.
    ```sql
    INSERT INTO role_permission (role, permission) VALUES 
    ('admin', 'menu::create'), ('admin', 'menu::update'), ('admin', 'menu::delete'), ('admin', 'menu::list'),
    ('admin', 'user::create'), ('admin', 'user::list'), ('admin', 'user::update'), ('admin', 'user::delete'),
    ('admin', 'order::create'), ('admin', 'order::list'), ('admin', 'order::update'), ('admin', 'order::delete'), ('admin', 'order::bill'),
    ('waiter', 'menu::list'), ('waiter', 'order::create'), ('waiter', 'order::list'), ('waiter', 'order::bill'),
    ('kitchen', 'menu::list'), ('kitchen', 'order::list'), ('kitchen', 'order::update');
    ```

---

## 2. Backend Setup & Run

1.  **Navigate to the backend directory**:
    ```bash
    cd backend
    ```
2.  **Configuration**: Check the `.env` file and update it with your database credentials if they differ from the default.
    ```env
    SERVER_PORT=8080
    DB_HOST=localhost
    DB_PORT=5432
    DB_USER=postgres
    DB_PASSWORD=password
    DB_NAME=restaurant
    JWT_SECRET=supersecret
    ```
3.  **Install Dependencies**:
    ```bash
    go mod tidy
    ```
4.  **Run the Server**:
    ```bash
    go run main.go
    ```
    The server should now be running on `http://localhost:8080`. You can verify this by visiting `http://localhost:8080/v1/ping` in your browser or using `curl`.

---

## 3. Frontend Setup & Run

1.  **Navigate to the frontend directory**:
    ```bash
    cd frontend
    ```
2.  **Install Dependencies**:
    ```bash
    flutter pub get
    ```
3.  **Run the App**:
    - For mobile (Android/iOS): Ensure an emulator is running or a device is connected.
    - For Web:
    ```bash
    flutter run -d chrome
    ```
    - For Desktop (Linux/macOS/Windows):
    ```bash
    flutter run -d linux # or macos/windows
    ```

The Flutter app will interact with the Go backend. Ensure the backend is running before performing actions that require API calls.
