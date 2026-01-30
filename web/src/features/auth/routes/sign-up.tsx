import { Link } from 'react-router-dom'
import { SignUpForm } from '../components/sign-up-form'

export function SignUpRoute() {
  return (
    <div className="flex h-screen w-full items-center justify-center bg-gray-50">
      <div className="w-full max-w-md space-y-6 rounded-lg border bg-white p-6 shadow-sm">
        <div className="space-y-2 text-center">
          <h1 className="text-2xl font-bold">Sign Up</h1>
          <p className="text-gray-500">Create an account to get started</p>
        </div>
        <SignUpForm />
        <div className="text-center text-sm">
          Already have an account?{' '}
          <Link to="/login" className="font-medium underline hover:text-gray-900">
            Login
          </Link>
        </div>
      </div>
    </div>
  )
}
