import { useState, useRef, useEffect } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { MessageSquare, X, Send, Bot, Loader2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card } from '@/components/ui/card'
import { useInventoryStore } from '@/store/inventory'
import { sendMessage, type ActionResult } from '../api/chat'
import { cn } from '@/lib/utils'
import { ActionRenderer } from './action-renderer'

type Message = {
  id: string
  role: 'user' | 'agent'
  content: string
  actions?: ActionResult[]
}

export function ChatInterface() {
  const [isOpen, setIsOpen] = useState(false)
  const [messages, setMessages] = useState<Message[]>([])
  const [input, setInput] = useState('')
  const scrollRef = useRef<HTMLDivElement>(null)

  const activeInventoryId = useInventoryStore(
    (state) => state.activeInventoryId
  )
  const queryClient = useQueryClient()

  const { mutate, isPending } = useMutation({
    mutationFn: sendMessage,
    onSuccess: (data) => {
      setMessages((prev) => [
        ...prev,
        {
          id: Date.now().toString(),
          role: 'agent',
          content: data.response,
          actions: data.actions,
        },
      ])
      // Invalidate queries to refresh data
      queryClient.invalidateQueries({ queryKey: ['products'] })
      queryClient.invalidateQueries({ queryKey: ['inventory'] })
    },
    onError: (error) => {
      setMessages((prev) => [
        ...prev,
        {
          id: Date.now().toString(),
          role: 'agent',
          content: 'Sorry, I encountered an error. Please try again.',
        },
      ])
      console.error(error)
    },
  })

  // Scroll to bottom when messages change
  useEffect(() => {
    if (scrollRef.current) {
      scrollRef.current.scrollTop = scrollRef.current.scrollHeight
    }
  }, [messages, isOpen])

  const handleSend = () => {
    if (!input.trim() || !activeInventoryId) return

    const userMessage: Message = {
      id: Date.now().toString(),
      role: 'user',
      content: input,
    }

    setMessages((prev) => [...prev, userMessage])
    setInput('')

    mutate({ prompt: userMessage.content, inventory_id: activeInventoryId })
  }

  // Only show if inventory is selected
  if (!activeInventoryId) return null

  return (
    <div className="fixed bottom-6 right-6 z-50 flex flex-col items-end gap-4">
      {isOpen ? (
        <Card className="flex h-[500px] w-[350px] flex-col shadow-xl animate-in fade-in slide-in-from-bottom-10 overflow-hidden border-slate-200 hover:scale-100 hover:shadow-xl bg-white">
          <div className="flex items-center justify-between border-b bg-slate-50 px-4 py-3">
            <div className="flex items-center gap-2 font-medium text-slate-700">
              <Bot className="h-4 w-4 text-primary" />
              <span>AI Assistant</span>
            </div>
            <Button
              variant="ghost"
              size="icon"
              className="h-6 w-6 text-slate-500 hover:text-slate-900"
              onClick={() => setIsOpen(false)}
            >
              <X className="h-4 w-4" />
            </Button>
          </div>

          <div
            className="flex-1 overflow-y-auto p-4 space-y-4 bg-white"
            ref={scrollRef}
          >
            {messages.length === 0 && (
              <div className="flex flex-col items-center justify-center h-full text-center text-slate-500 space-y-4">
                <div className="p-3 bg-slate-100 rounded-full">
                  <Bot className="h-8 w-8 text-slate-400" />
                </div>
                <div className="space-y-1">
                  <p className="font-medium text-slate-900">How can I help?</p>
                  <p className="text-xs max-w-[200px] mx-auto">
                    Try asking to add products, find recipes, or check
                    inventory.
                  </p>
                </div>
              </div>
            )}

            {messages.map((msg) => (
              <div
                key={msg.id}
                className={cn(
                  'max-w-[85%] rounded-2xl px-3 py-2 text-sm',
                  msg.role === 'user'
                    ? 'ml-auto bg-primary text-primary-foreground rounded-tr-none'
                    : 'bg-slate-100 text-slate-800 rounded-tl-none'
                )}
              >
                {msg.content}
                {msg.actions && <ActionRenderer actions={msg.actions} />}
              </div>
            ))}

            {isPending && (
              <div className="bg-slate-100 w-max rounded-2xl rounded-tl-none px-3 py-2">
                <Loader2 className="h-4 w-4 animate-spin text-slate-500" />
              </div>
            )}
          </div>

          <div className="p-3 border-t bg-white">
            <form
              onSubmit={(e) => {
                e.preventDefault()
                handleSend()
              }}
              className="flex gap-2"
            >
              <Input
                value={input}
                onChange={(e) => setInput(e.target.value)}
                placeholder="Message..."
                disabled={isPending}
                className="flex-1"
              />
              <Button
                type="submit"
                size="icon"
                disabled={isPending || !input.trim()}
              >
                <Send className="h-4 w-4" />
              </Button>
            </form>
          </div>
        </Card>
      ) : (
        <Button
          onClick={() => setIsOpen(true)}
          size="icon"
          className="h-12 w-12 rounded-full shadow-lg hover:scale-105 transition-transform"
        >
          <MessageSquare className="h-6 w-6" />
        </Button>
      )}
    </div>
  )
}
