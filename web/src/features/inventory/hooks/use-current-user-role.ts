import { useQuery } from '@tanstack/react-query'
import { useAuthStore } from '@/store/auth'
import { getMembers } from '../api'

export function useCurrentUserRole(inventoryId: string | null) {
  const user = useAuthStore((state) => state.user)

  const { data: members, isLoading, error } = useQuery({
    queryKey: ['members', inventoryId],
    queryFn: () => getMembers(inventoryId!),
    enabled: !!inventoryId && !!user,
  })

  const membership = members?.find((m) => m.user_id === user?.id)

  return {
    role: membership?.role,
    isLoading,
    error,
    membership,
  }
}
