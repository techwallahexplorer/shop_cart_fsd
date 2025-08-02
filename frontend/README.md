# Shopping Cart Frontend

This is the React frontend for the Shopping Cart application.

## Prerequisites
- Node.js 18+
- npm or yarn

## Setup
1. Install dependencies:
   ```
   npm install
   # or
   yarn install
   ```
2. Start the development server:
   ```
   npm run dev
   # or
   yarn dev
   ```
   The app will be available at `http://localhost:3000`.

## Features
- User login (with username and password)
- List all items
- Add items to cart
- Checkout (convert cart to order)
- View cart items and order history (alerts)

## Notes
- The frontend expects the backend to be running on `http://localhost:8080`.
- All cart and order requests require the user's token (managed automatically after login).
