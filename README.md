# 🍽️ FoodFlow

A modern food sharing platform that connects food donors (collaborators) with organizations that help feed communities. Built with Go backend and React frontend.

## 🌟 Features

### For Collaborators (Food Donors)
- Create food donations with detailed descriptions
- Set expiry dates and special instructions
- Earn tokens for contributions
- Track donation history

### For Organizations (Food Recipients)
- Browse and claim available food donations
- Manage credits based on community impact
- View organization statistics
- Proxy assignment for remote operations

### For Administrators
- Reset monthly credits for all organizations
- Assign proxy organizations
- Monitor system-wide statistics
- Manage organization relationships

## 🏗️ Architecture

- **Backend**: Go with Gin framework
- **Database**: PostgreSQL with GORM ORM
- **Frontend**: React with Vite
- **Styling**: Tailwind CSS
- **Authentication**: JWT tokens
- **CORS**: Configured for cross-origin requests

## 📋 Prerequisites

- Go 1.21 or higher
- Node.js 16 or higher
- PostgreSQL 12 or higher
- npm or yarn

## 🚀 Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/Karthikeya-Akhandam/FoodFlow.git
cd foodshare
```

### 2. Backend Setup

```bash
# Install Go dependencies
go mod tidy

# Set up environment variables
export DATABASE_DSN="host=localhost user=postgres password=your_password dbname=foodshare port=5432 sslmode=disable TimeZone=Asia/Kolkata"
export PORT=8080

# Run the backend
go run main.go
```

The backend will start on `http://localhost:8080`

### 3. Frontend Setup

```bash
# Navigate to frontend directory
cd web

# Install dependencies
npm install

# Create environment file
echo "VITE_API_URL=http://localhost:8080" > .env

# Start development server
npm run dev
```

The frontend will start on `http://localhost:5173`

## 🗄️ Database Setup

### PostgreSQL Configuration

1. Create a PostgreSQL database:
```sql
CREATE DATABASE foodshare;
```

2. Update the connection string in your environment:
```bash
export DATABASE_DSN="host=localhost user=postgres password=your_password dbname=foodshare port=5432 sslmode=disable TimeZone=Asia/Kolkata"
```

3. The application will automatically create tables on first run.

### Default Admin Account

A default admin account is created automatically:
- **Email**: `admin@local`
- **Password**: `admin123`

## 🔧 Configuration

### Backend Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_DSN` | PostgreSQL connection string | Local PostgreSQL |
| `PORT` | Server port | `8080` |
| `JWT_SECRET` | JWT signing secret | Hardcoded (change in production) |

### Frontend Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `VITE_API_URL` | Backend API URL | `http://localhost:8080` |

## 📚 API Endpoints

### Authentication
- `POST /auth/signup` - User registration
- `POST /auth/login` - User login

### Public
- `GET /donations` - List available donations

### Protected (Requires JWT)
- `POST /donations` - Create donation (collaborators only)
- `POST /donations/:id/claim` - Claim donation (organizations only)
- `GET /me` - Get current user profile

### Admin Only
- `POST /admin/reset-credits` - Reset all organization credits
- `POST /admin/assign-proxy` - Assign proxy organization

## 🎨 User Roles

### 👤 Collaborator
- Create food donations
- Earn tokens for contributions
- View available donations

### 🏢 Organization
- Claim food donations using credits
- View credit balance and statistics
- Update people fed statistics

### 👑 Administrator
- Manage all organizations
- Reset monthly credits
- Assign proxy relationships
- Monitor system statistics

## 🛠️ Development

### Backend Development

```bash
# Run with hot reload (requires air)
go install github.com/cosmtrek/air@latest
air

# Run tests
go test ./...

# Build for production
go build -o foodshare main.go
```

### Frontend Development

```bash
cd web

# Development server
npm run dev

# Build for production
npm run build

# Preview production build
npm run start
```

## 📦 Production Deployment

### Backend Deployment

1. Build the binary:
```bash
go build -o foodshare main.go
```

2. Set production environment variables:
```bash
export DATABASE_DSN="your_production_database_url"
export PORT=8080
export JWT_SECRET="your_secure_jwt_secret"
```

3. Run the application:
```bash
./foodshare
```

### Frontend Deployment

1. Build the frontend:
```bash
cd web
npm run build
```

2. Serve the `dist` folder with any static file server (nginx, Apache, etc.)

## 🔒 Security Considerations

- Change the default JWT secret in production
- Use environment variables for sensitive configuration
- Implement proper CORS policies
- Add rate limiting for API endpoints
- Use HTTPS in production
- Regularly update dependencies

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

If you encounter any issues or have questions:

1. Check the [Issues](https://github.com/your-repo/issues) page
2. Create a new issue with detailed information
3. Contact the development team

## 🗺️ Roadmap

- [ ] Real-time notifications
- [ ] Mobile app (React Native)
- [ ] Advanced analytics dashboard
- [ ] Integration with food delivery services
- [ ] Multi-language support
- [ ] Advanced reporting features

## 🙏 Acknowledgments

- [Gin Web Framework](https://gin-gonic.com/)
- [GORM](https://gorm.io/)
- [React](https://reactjs.org/)
- [Tailwind CSS](https://tailwindcss.com/)
- [Vite](https://vitejs.dev/)

---

**Made with ❤️ for the community**
