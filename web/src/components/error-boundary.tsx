import { useRouteError } from 'react-router-dom'

export function ErrorBoundary() {
  const error = useRouteError() as any
  console.error(error)

  return (
    <div className="flex h-screen w-full flex-col items-center justify-center gap-4">
      <h1 className="text-3xl font-bold">Oops!</h1>
      <p className="text-gray-500">Sorry, an unexpected error has occurred.</p>
      <p className="font-mono text-sm bg-gray-100 p-2 rounded">
        {error?.statusText || error?.message || 'Unknown error'}
      </p>
    </div>
  )
}
