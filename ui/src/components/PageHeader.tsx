import { ReactNode } from 'react'
import { Button } from 'antd'
import { PlusOutlined } from '@ant-design/icons'

interface PageHeaderProps {
  title: string
  actionText?: string
  onAction?: () => void
  extra?: ReactNode
  showAction?: boolean
}

export const PageHeader: React.FC<PageHeaderProps> = ({
  title,
  actionText,
  onAction,
  extra,
  showAction = true,
}) => {
  return (
    <div 
      style={{ 
        display: 'flex', 
        justifyContent: 'space-between', 
        alignItems: 'center', 
        marginBottom: 24,
      }}
    >
      <h2 style={{ margin: 0 }}>{title}</h2>
      <div style={{ display: 'flex', gap: 8, alignItems: 'center' }}>
        {extra}
        {showAction && actionText && onAction && (
          <Button type="primary" icon={<PlusOutlined />} onClick={onAction}>
            {actionText}
          </Button>
        )}
      </div>
    </div>
  )
}

export default PageHeader