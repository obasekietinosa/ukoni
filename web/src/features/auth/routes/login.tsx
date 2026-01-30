import { Link } from 'react-router-dom'
import { LoginForm } from '../components/login-form'

export function LoginRoute() {
  return (
    <div className="flex h-screen w-full items-center justify-center bg-gray-50">
      <div className="w-full max-w-md space-y-6 rounded-lg border bg-white p-6 shadow-sm">
        <div className="space-y-2 text-center">
          <h1 className="text-2xl font-bold">Login</h1>
          <p className="text-gray-500">Enter your credentials to access your account</p>
        </div>
        <LoginForm />
        <div className="text-center text-sm">
          Don&apos;t have an account?{' '}
          <Link to="/signup" className="font-medium underline hover:text-gray-900">
            Sign up
          </Link>
        </div>
      </div>
    </div>
  )
}
