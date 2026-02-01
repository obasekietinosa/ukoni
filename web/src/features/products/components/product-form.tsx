import { useState, useEffect } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import type { Product } from '../types'

interface Props {
  initialData?: Partial<Product>
  onSubmit: (data: { name: string; brand?: string }) => void
  isLoading?: boolean
  submitLabel?: string
}

export function ProductForm({
  initialData,
  onSubmit,
  isLoading,
  submitLabel = 'Save',
}: Props) {
  const [name, setName] = useState(initialData?.name || '')
  const [brand, setBrand] = useState(initialData?.brand || '')

  useEffect(() => {
    if (initialData) {
      setName(initialData.name || '')
      setBrand(initialData.brand || '')
    }
  }, [initialData])

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!name) return
    onSubmit({ name, brand })
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="space-y-2">
        <label htmlFor="product-name" className="text-sm font-medium">
          Product Name
        </label>
        <Input
          id="product-name"
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="e.g. Tesco Whole Milk"
          required
        />
      </div>
      <div className="space-y-2">
        <label htmlFor="brand" className="text-sm font-medium">
          Brand (Optional)
        </label>
        <Input
          id="brand"
          type="text"
          value={brand}
          onChange={(e) => setBrand(e.target.value)}
          placeholder="e.g. Tesco"
        />
      </div>
      <Button type="submit" disabled={isLoading}>
        {isLoading ? 'Saving...' : submitLabel}
      </Button>
    </form>
  )
}
