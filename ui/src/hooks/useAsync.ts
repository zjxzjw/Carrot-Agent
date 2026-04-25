import { useState, useCallback } from 'react'

export function useAsync<T, A extends unknown[]>(
  asyncAction: (...args: A) => Promise<T>,
  deps: unknown[] = []
): {
  data: T | null
  loading: boolean
  error: string | null
  execute: (...args: A) => Promise<T>
} {
  const [data, setData] = useState<T | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const execute = useCallback(async (...args: A) => {
    setLoading(true)
    setError(null)
    try {
      const result = await asyncAction(...args)
      setData(result)
      return result
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '操作失败'
      setError(errorMessage)
      throw err
    } finally {
      setLoading(false)
    }
  }, [asyncAction, ...deps])

  return { data, loading, error, execute }
}

export default useAsync