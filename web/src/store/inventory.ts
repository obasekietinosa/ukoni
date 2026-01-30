import { create } from 'zustand'
import { persist } from 'zustand/middleware'

interface InventoryStore {
  activeInventoryId: string | null
  setActiveInventoryId: (id: string) => void
  clearActiveInventory: () => void
}

export const useInventoryStore = create<InventoryStore>()(
  persist(
    (set) => ({
      activeInventoryId: null,
      setActiveInventoryId: (id) => set({ activeInventoryId: id }),
      clearActiveInventory: () => set({ activeInventoryId: null }),
    }),
    {
      name: 'inventory-storage',
    }
  )
)
