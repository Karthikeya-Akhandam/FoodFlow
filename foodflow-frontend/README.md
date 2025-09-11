# FoodFlow Frontend

A modern, responsive web application for the FoodFlow platform - connecting food donors with organizations to reduce food waste and help communities.

## Features

- **Modern UI/UX**: Clean, professional design with green and white color scheme
- **Responsive Design**: Works seamlessly on desktop, tablet, and mobile devices
- **Authentication**: Secure login/signup with role-based access control
- **Dashboard**: Role-specific dashboards for Organizations, Collaborators, and Admins
- **Real-time Updates**: Live data updates using React Query
- **API Integration**: Complete integration with FoodFlow backend API

## Tech Stack

- **Framework**: Next.js 15 with App Router
- **Styling**: Tailwind CSS 4
- **State Management**: React Query (TanStack Query)
- **Authentication**: JWT-based with context API
- **Icons**: Lucide React
- **Fonts**: Geist Sans & Geist Mono

## Getting Started

### Prerequisites

- Node.js 18+ 
- npm, yarn, or pnpm
- FoodFlow backend running on `http://localhost:8080`

### Installation

1. Install dependencies:
```bash
npm install
# or
yarn install
# or
pnpm install
```

2. Create environment file:
```bash
cp .env.example .env.local
```

3. Update `.env.local` with your configuration:
```env
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_APP_NAME=FoodFlow
NEXT_PUBLIC_APP_VERSION=1.0.0
```

4. Run the development server:
```bash
npm run dev
# or
yarn dev
# or
pnpm dev
```

5. Open [http://localhost:3000](http://localhost:3000) in your browser.

## Project Structure

```
foodflow-frontend/
├── app/                    # Next.js App Router
│   ├── auth/              # Authentication pages
│   │   ├── login/         # Login page
│   │   └── signup/        # Signup page
│   ├── dashboard/         # Dashboard pages
│   │   ├── offers/        # Offer management
│   │   ├── claims/        # Claim management
│   │   └── profile/       # Profile management
│   ├── globals.css        # Global styles
│   ├── layout.js          # Root layout
│   └── page.js            # Landing page
├── lib/                   # Utility libraries
│   ├── api.js             # API client
│   ├── auth.js            # Authentication context
│   ├── config.js          # App configuration
│   ├── hooks.js           # Custom React Query hooks
│   └── queryClient.js     # React Query configuration
├── public/                # Static assets
└── tailwind.config.js     # Tailwind configuration
```

## User Roles

### Organizations (ORG)
- Browse and claim food donation offers
- Manage their organization profile
- Track credits and redemption history
- Create claims for remote organizations

### Collaborators (COLLAB)
- Create food donation offers
- Manage their collaborator profile
- Track tokens and earning history
- Monitor offer status and claims

### Administrators (ADMIN)
- Monitor system statistics
- Manage user accounts
- View audit logs
- Oversee platform operations

## API Integration

The frontend integrates with the FoodFlow backend API through:

- **Authentication**: JWT token-based authentication
- **Data Fetching**: React Query for caching and synchronization
- **Error Handling**: Comprehensive error handling with user-friendly messages
- **Type Safety**: Consistent data structures and validation

## Key Features

### Landing Page
- Hero section with clear value proposition
- Feature highlights
- How it works section
- Call-to-action buttons

### Authentication
- Secure login/signup forms
- Role selection during signup
- Password strength validation
- Remember me functionality

### Dashboard
- Role-specific navigation
- Real-time statistics
- Quick actions
- Recent activity feed

### Offer Management
- Create, view, and edit offers
- Advanced filtering and search
- Status tracking
- Location-based matching

### Claim Management
- Submit and track claims
- Priority scoring display
- Status updates
- Redemption confirmation

### Profile Management
- Complete profile editing
- Organization/collaborator details
- Location information
- Purpose focus selection

## Development

### Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run start` - Start production server
- `npm run lint` - Run ESLint

### Code Style

- Use functional components with hooks
- Follow React best practices
- Use TypeScript for type safety (future enhancement)
- Consistent naming conventions
- Proper error handling

## Deployment

### Vercel (Recommended)

1. Connect your GitHub repository to Vercel
2. Set environment variables in Vercel dashboard
3. Deploy automatically on push to main branch

### Other Platforms

The app can be deployed to any platform that supports Next.js:
- Netlify
- AWS Amplify
- Railway
- DigitalOcean App Platform

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is part of the FoodFlow platform. See the main repository for license information.
