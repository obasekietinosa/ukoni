import { Link } from 'react-router-dom'
import { LoginForm } from '../components/login-form'
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
  CardFooter,
} from '@/components/ui/card'

export function LoginRoute() {
  return (
    <div className="flex h-screen w-full items-center justify-center p-4">
      <Card className="w-full max-w-md">
        <CardHeader className="text-center">
          <CardTitle>Login</CardTitle>
          <CardDescription>
            Enter your credentials to access your account
          </CardDescription>
        </CardHeader>
        <CardContent>
          <LoginForm />
        </CardContent>
        <CardFooter className="justify-center">
          <div className="text-sm text-slate-500">
            Don&apos;t have an account?{' '}
            <Link
              to="/signup"
              className="font-medium text-electric-mint underline-offset-4 hover:underline"
            >
              Sign up
            </Link>
          </div>
        </CardFooter>
      </Card>
    </div>
  )
}
