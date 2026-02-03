import { render, screen } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import { Logo, LogoMark } from './logo'

describe('Logo', () => {
  it('LogoMark renders correctly', () => {
    const { container } = render(<LogoMark />)
    const svg = container.querySelector('svg')
    expect(svg).toBeInTheDocument()
    expect(svg).toHaveAttribute('viewBox', '0 0 165 170')
  })

  it('Logo renders correctly with text', () => {
    const { container } = render(<Logo />)
    expect(screen.getByText('Ukoni')).toBeInTheDocument()
    const svg = container.querySelector('svg')
    expect(svg).toBeInTheDocument()
  })

  it('Logo accepts custom className', () => {
    const { container } = render(<Logo className="custom-class" />)
    const div = container.firstChild
    expect(div).toHaveClass('custom-class')
  })
})
