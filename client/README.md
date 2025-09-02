This is a [Next.js](https://nextjs.org) project bootstrapped with [`create-next-app`](https://nextjs.org/docs/app/api-reference/cli/create-next-app).

## Getting Started

First, run the development server:

```bash
npm run dev
# or
yarn dev
# or
pnpm dev
# or
bun dev
```

Open [http://localhost:3000](http://localhost:3000) with your browser to see the result.

You can start editing the page by modifying `app/page.tsx`. The page auto-updates as you edit the file.

This project uses [`next/font`](https://nextjs.org/docs/app/building-your-application/optimizing/fonts) to automatically optimize and load [Geist](https://vercel.com/font), a new font family for Vercel.

## Authentication Setup

This application uses JWT-based authentication with the following features:

### Environment Variables
Create a `.env` file based on `.env.example`:
```bash
cp .env.example .env
```

### Authentication Flow
1. Users can register at `/register`
2. Users can login at `/login`
3. Protected routes (like `/dashboard` and `/profile`) require authentication
4. The auth context manages user state and tokens

### API Integration
The frontend communicates with the Go backend API at `http://localhost:8080` by default.

## Project Structure
- `app/contexts/auth.context.tsx` - Authentication context provider
- `app/utils/auth.service.ts` - Authentication service for API calls
- `app/types/auth.types.ts` - TypeScript types for authentication
- `app/components/ProtectedRoute.tsx` - Component to protect authenticated routes
- `app/login/page.tsx` - Login page
- `app/register/page.tsx` - Registration page
- `app/dashboard/page.tsx` - User dashboard
- `app/profile/page.tsx` - User profile page

## Learn More

To learn more about Next.js, take a look at the following resources:

- [Next.js Documentation](https://nextjs.org/docs) - learn about Next.js features and API.
- [Learn Next.js](https://nextjs.org/learn) - an interactive Next.js tutorial.

You can check out [the Next.js GitHub repository](https://github.com/vercel/next.js) - your feedback and contributions are welcome!

## Deploy on Vercel

The easiest way to deploy your Next.js app is to use the [Vercel Platform](https://vercel.com/new?utm_medium=default-template&filter=next.js&utm_source=create-next-app&utm_campaign=create-next-app-readme) from the creators of Next.js.

Check out our [Next.js deployment documentation](https://nextjs.org/docs/app/building-your-application/deploying) for more details.