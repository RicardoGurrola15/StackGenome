import { describe, it, expect, vi } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import FileImporter from '../components/report/FileImporter'
import { MAX_FILE_SIZE_MB } from '../lib/schema'

describe('FileImporter', () => {
  it('renders instructions', () => {
    render(<FileImporter onImport={vi.fn()} />)
    expect(screen.getByText(/Arrastra tu archivo aquí/)).toBeInTheDocument()
  })

  it('handles large file error', async () => {
    const handleImport = vi.fn()
    render(<FileImporter onImport={handleImport} />)
    
    const file = new File(['a'.repeat((MAX_FILE_SIZE_MB + 1) * 1024 * 1024)], 'huge.json', { type: 'application/json' })

    // We mock the onChange event directly since the input is hidden
    const hiddenInput = document.querySelector('input[type="file"]') as HTMLInputElement
    fireEvent.change(hiddenInput, { target: { files: [file] } })
    
    await waitFor(() => {
      expect(screen.getByText(/demasiado grande/)).toBeInTheDocument()
    })
    expect(handleImport).not.toHaveBeenCalled()
  })
})
