import { useEffect } from 'react'
import { message } from 'antd'

export function useAsyncError(error: string | null | undefined, t?: (key: string) => string) {
  useEffect(() => {
    if (error) {
      message.error(t ? t('common.error') + ': ' + error : error)
    }
  }, [error, t])
}