# Snooker Assistant Server

## Introduction
The Snooker Assistant Server is a server-side application designed for snooker games, aimed at providing players with an efficient and reliable game management platform. This application supports user registration, login, creation and management of game rooms, as well as real-time notification features, enhancing the user experience.

## Features
- **User Management**: Supports user registration, login, information updates, and avatar uploads.
- **Game Management**: Allows users to create, delete, and query game rooms, supporting various game types.
- **Real-time Notification System**: Implements real-time notifications via WebSocket, ensuring users receive timely updates on games and system messages.
- **Database Support**: Uses Redis for caching and session storage, and MySQL as the main database.
- **Object Storage**: Integrates MinIO object storage for storing user-uploaded files and game data.

## Tech Stack
- **Programming Language**: Go
- **Web Framework**: Gin
- **Database**: MySQL
- **Cache**: Redis
- **Object Storage**: MinIO
- **ORM**: GORM

## Installation
1. **Clone the Project**:
   ```bash
   git clone https://github.com/superwhys/snooker-assistant-server.git
   cd snooker-assistant-server
   ```

2. **Install Dependencies**:
   ```bash
   go mod tidy
   ```

3. **Configuration File**:
   - Create a `config.yaml` file and add the necessary configurations, such as database connection information, Redis configuration, and MinIO configuration. Example configuration is as follows:
   ```yaml:content/config.yaml
   conf:
     wxAppId: wx018d25746881635b
     wxAppSecret: 3a3027d6cb756d73fbb448770f4d3c5e
     avatarDir: "avatar"
   redisAuth:
     server: localhost:6379
     db: 10
   mysqlAuth:
     instance: localhost:3306
     database: snooker
     username: root
     password: yang4869
   minioAuth:
     endpoint: 10.11.43.115:9000 
     accessKey: jRdUastUrl69Nn1Zv6CZ 
     secretKey: EooHGXHOp8CFDZlKoczzKtJ3WO8pNv9167PDuwx5
     bucket: snooker-assistant
   ```

## Usage
1. **Start the Server**:
   ```bash
   go run main.go
   ```

2. **Access the API Documentation**:
   - After the server starts, you can view the API documentation by visiting `http://localhost:8080/docs` to get available endpoints and usage examples.

3. **Test the API**:
   - You can use Postman or curl to test the API. Here are some example requests:
     - **User Registration**:
       ```bash
       curl -X POST http://localhost:8080/v1/user/register -d '{"username": "testuser", "password": "password123"}' -H "Content-Type: application/json"
       ```
     - **Get Game List**:
       ```bash
       curl -X GET http://localhost:8080/v1/game/list
       ```

## Contribution
Contributions of any kind are welcome! If you find issues or have suggestions for improvements, please submit issues or pull requests. Please follow these steps:
1. Fork this repository
2. Create your feature branch (`git checkout -b feature/YourFeature`)
3. Commit your changes (`git commit -m 'Add some feature'`)
4. Push to the branch (`git push origin feature/YourFeature`)
5. Create a new Pull Request

## License
This project is licensed under the MIT License. For more details, please see the [LICENSE](LICENSE) file.

## Contact Information
If you have any questions or suggestions, please contact the project maintainer: [Your Name] ([Your Email]).
