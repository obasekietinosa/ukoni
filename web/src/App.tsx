import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { BrowserRouter, Route, Routes } from 'react-router-dom'
import { MainLayout } from '@/components/layout/main-layout'
import { ErrorBoundary } from '@/components/error-boundary'
import { LoginRoute } from '@/features/auth/routes/login'
import { SignUpRoute } from '@/features/auth/routes/sign-up'
import { RequireAuth } from '@/components/require-auth'
import { InventorySelectionRoute } from '@/features/inventory/routes/inventory-selection'
import { InventoryGuard } from '@/features/inventory/components/inventory-guard'
import { InventoryPage } from '@/features/inventory/routes/inventory-page'
import { ProductsPage } from '@/features/products/routes/products-page'
import { ProductDetailsPage } from '@/features/products/routes/product-details-page'
import { SellersPage } from '@/features/sellers/routes/sellers-page'
import { SellerDetailsPage } from '@/features/sellers/routes/seller-details-page'
import { ShoppingListsPage } from '@/features/shopping-lists/routes/shopping-lists-page'
import { ShoppingListDetailsPage } from '@/features/shopping-lists/routes/shopping-list-details-page'

const queryClient = new QueryClient()

function Home() {
  return (
    <div className="flex h-full flex-col items-center justify-center gap-4">
      <h1 className="text-3xl font-bold">Welcome to Ukoni</h1>
      <p>Dashboard coming soon.</p>
    </div>
  )
}

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <Routes>
          {/* Public Routes */}
          <Route
            path="/login"
            element={<LoginRoute />}
            errorElement={<ErrorBoundary />}
          />
          <Route
            path="/signup"
            element={<SignUpRoute />}
            errorElement={<ErrorBoundary />}
          />

          {/* Protected Routes */}
          <Route element={<RequireAuth />} errorElement={<ErrorBoundary />}>
            <Route
              path="/select-inventory"
              element={<InventorySelectionRoute />}
            />

            <Route element={<InventoryGuard />}>
              <Route element={<MainLayout />}>
                <Route path="/" element={<Home />} />
                <Route path="/inventory" element={<InventoryPage />} />
                <Route path="/products" element={<ProductsPage />} />
                <Route path="/products/:id" element={<ProductDetailsPage />} />
                <Route path="/sellers" element={<SellersPage />} />
                <Route path="/sellers/:id" element={<SellerDetailsPage />} />
                <Route path="/shopping-lists" element={<ShoppingListsPage />} />
                <Route
                  path="/shopping-lists/:id"
                  element={<ShoppingListDetailsPage />}
                />
              </Route>
            </Route>
          </Route>
        </Routes>
      </BrowserRouter>
    </QueryClientProvider>
  )
}

export default App
