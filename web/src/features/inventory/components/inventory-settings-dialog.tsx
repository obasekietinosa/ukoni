import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { useInventoryStore } from '@/store/inventory'
import { getInventorySettings, updateInventorySettings } from '../api'
import { Loader2 } from 'lucide-react'
import type { InventorySettings } from '../types'

interface Props {
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function InventorySettingsDialog({ open, onOpenChange }: Props) {
  const activeInventoryId = useInventoryStore(
    (state) => state.activeInventoryId
  )
  const queryClient = useQueryClient()

  const { data: settings, isLoading } = useQuery({
    queryKey: ['inventory-settings', activeInventoryId],
    queryFn: () => getInventorySettings(activeInventoryId!),
    enabled: !!activeInventoryId && open,
  })

  const mutation = useMutation({
    mutationFn: (data: Partial<InventorySettings>) =>
      updateInventorySettings(activeInventoryId!, data),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ['inventory-settings', activeInventoryId],
      })
      onOpenChange(false)
    },
  })

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    if (!activeInventoryId) return

    const formData = new FormData(e.currentTarget)
    const provider = formData.get('llm_provider') as string
    const apiKey = formData.get('llm_api_key') as string

    mutation.mutate({
      llm_provider: provider || undefined,
      llm_api_key: apiKey || undefined,
    })
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Inventory Settings</DialogTitle>
        </DialogHeader>

        {isLoading && !settings ? (
          <div className="flex justify-center p-8">
            <Loader2 className="h-6 w-6 animate-spin text-gray-500" />
          </div>
        ) : (
          <form
            onSubmit={handleSubmit}
            className="space-y-4 py-4"
            key={settings?.updated_at}
          >
            {mutation.error && (
              <div className="rounded-md bg-red-50 p-3 text-sm text-red-500">
                {mutation.error instanceof Error
                  ? mutation.error.message
                  : 'Failed to update settings'}
              </div>
            )}

            <div className="space-y-2">
              <label
                htmlFor="llm_provider"
                className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
              >
                LLM Provider
              </label>
              <Input
                id="llm_provider"
                name="llm_provider"
                defaultValue={settings?.llm_provider}
                placeholder="e.g. openai"
              />
            </div>

            <div className="space-y-2">
              <label
                htmlFor="llm_api_key"
                className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
              >
                LLM API Key
              </label>
              <Input
                id="llm_api_key"
                name="llm_api_key"
                type="password"
                defaultValue={settings?.llm_api_key}
                placeholder="sk-..."
              />
            </div>

            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => onOpenChange(false)}
              >
                Cancel
              </Button>
              <Button type="submit" disabled={mutation.isPending}>
                {mutation.isPending && (
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                )}
                Save Changes
              </Button>
            </DialogFooter>
          </form>
        )}
      </DialogContent>
    </Dialog>
  )
}
