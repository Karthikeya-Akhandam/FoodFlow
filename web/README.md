# FoodShare Frontend

A React frontend for the FoodShare application.

## Setup

1. Install dependencies:
```bash
npm install
```

2. Create a `.env` file in the web directory with:
```
VITE_API_URL=http://localhost:8080
```

3. Start the development server:
```bash
npm run dev
```

## Features

- User authentication (login/signup)
- Role-based access (collaborator/organization)
- Create donations (collaborators)
- Claim donations (organizations)
- Responsive design with Tailwind CSS

## API Integration

The frontend communicates with the Go backend API. Make sure the backend is running on the configured port (default: 8080).
