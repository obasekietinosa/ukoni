import { useState, useEffect } from 'react'
import { Link, useSearchParams } from 'react-router-dom'
import { api } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Logo } from '@/components/ui/logo'
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
  CardFooter,
} from '@/components/ui/card'

export function ResetPasswordRoute() {
  const [searchParams] = useSearchParams()
  const token = searchParams.get('token')

  const [password, setPassword] = useState('')
  const [isValidating, setIsValidating] = useState(true)
  const [isTokenValid, setIsTokenValid] = useState(false)
  const [validationError, setValidationError] = useState<string | null>(null)

  const [isSubmitting, setIsSubmitting] = useState(false)
  const [submitError, setSubmitError] = useState<string | null>(null)
  const [submitSuccess, setSubmitSuccess] = useState(false)

  useEffect(() => {
    const validateToken = async () => {
      if (!token) {
        setValidationError('No reset token provided in the URL.')
        setIsValidating(false)
        return
      }

      try {
        await api('/password-reset/validate', {
          method: 'POST',
          json: { token },
        })
        setIsTokenValid(true)
      } catch (err) {
        const message =
          err instanceof Error ? err.message : 'Invalid or expired token.'
        setValidationError(message)
      } finally {
        setIsValidating(false)
      }
    }

    validateToken()
  }, [token])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!token) return

    setIsSubmitting(true)
    setSubmitError(null)

    try {
      await api('/password-reset/reset', {
        method: 'POST',
        json: { token, password },
      })
      setSubmitSuccess(true)
    } catch (err) {
      const message =
        err instanceof Error ? err.message : 'Failed to reset password.'
      setSubmitError(message)
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <div className="flex h-screen w-full items-center justify-center p-4">
      <Card className="w-full max-w-md">
        <CardHeader className="text-center">
          <div className="flex justify-center mb-6">
            <Logo className="scale-125" iconClassName="h-12" />
          </div>
          <CardTitle>Reset Password</CardTitle>
          <CardDescription>
            {submitSuccess
              ? 'Password reset successfully'
              : 'Enter your new password below'}
          </CardDescription>
        </CardHeader>
        <CardContent>
          {isValidating ? (
            <div className="text-center text-slate-500">Validating link...</div>
          ) : !isTokenValid ? (
            <div className="space-y-4">
              <div className="text-red-500 text-center">{validationError}</div>
              <div className="text-center">
                <Link
                  to="/forgot-password"
                  className="text-sm font-medium text-electric-mint underline-offset-4 hover:underline"
                >
                  Request a new reset link
                </Link>
              </div>
            </div>
          ) : submitSuccess ? (
            <div className="text-center space-y-4">
              <div className="text-electric-mint mb-6">
                Your password has been changed successfully.
              </div>
              <Link to="/login">
                <Button className="w-full">Log In Now</Button>
              </Link>
            </div>
          ) : (
            <form onSubmit={handleSubmit} className="space-y-4">
              {submitError && (
                <div className="text-red-500 text-sm">{submitError}</div>
              )}
              <div className="space-y-2">
                <label htmlFor="password" className="text-sm font-medium">
                  New Password
                </label>
                <Input
                  id="password"
                  type="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                  minLength={8}
                />
              </div>
              <Button type="submit" className="w-full" disabled={isSubmitting}>
                {isSubmitting ? 'Resetting...' : 'Reset Password'}
              </Button>
            </form>
          )}
        </CardContent>
        {!submitSuccess && isTokenValid && (
          <CardFooter className="justify-center">
            <div className="text-sm text-slate-500">
              Remembered your password?{' '}
              <Link
                to="/login"
                className="font-medium text-electric-mint underline-offset-4 hover:underline"
              >
                Log in
              </Link>
            </div>
          </CardFooter>
        )}
      </Card>
    </div>
  )
}
