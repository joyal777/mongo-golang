# Mongo-Golang
A lightweight RESTful API built with Golang and MongoDB 8.0. This project demonstrates basic CRUD operations using the official MongoDB Go Driver and the httprouter package.🚀 
**Getting Started**

1. **Prerequisites**
    Go: 1.20 or higher
    MongoDB: Version 8.0 (Community Server)
    MongoDB Compass: (Optional)  For visual data management

3. **Installation**
    Clone the repository:
    ```bash
    git clone https://github.com/joyal777/mongo-golang.git
    cd mongo-golang

4. **Install Dependencies:**
    ```bash
    go get go.mongodb.org/mongo-driver/mongo
    go get github.com/julienschmidt/httprouter

3. **Run the Application:**
    ```bash
    go run main.go

Expected Output: Successfully connected to MongoDB 8.0!

**🛠 API Endpoints**

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| **POST** | `/user` | Create a new user |
| **GET** | `/user/:id` | Retrieve user details by ID |
| **DELETE** | `/user/:id` | Remove a user from the database |

