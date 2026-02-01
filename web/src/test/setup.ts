import '@testing-library/jest-dom'
import { beforeAll, afterEach, afterAll } from 'vitest'
import { setupServer } from 'msw/node'
import { handlers } from '../mocks/handlers'

// Polyfill for HTMLDialogElement in JSDOM
// HTMLDialogElement is not supported in JSDOM environment
if (typeof window !== 'undefined') {
  HTMLDialogElement.prototype.showModal = function (this: HTMLDialogElement) {
    this.open = true
    this.setAttribute('open', '')
  }

  HTMLDialogElement.prototype.close = function (this: HTMLDialogElement) {
    this.open = false
    this.removeAttribute('open')
  }
}

export const server = setupServer(...handlers)

beforeAll(() => server.listen())
afterEach(() => server.resetHandlers())
afterAll(() => server.close())
