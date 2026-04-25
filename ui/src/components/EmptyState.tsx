import { Button } from 'antd'
import { PlusOutlined, InboxOutlined } from '@ant-design/icons'

interface EmptyStateProps {
  title?: string
  description?: string
  actionText?: string
  onAction?: () => void
}

export const EmptyState: React.FC<EmptyStateProps> = ({
  title = '暂无数据',
  description,
  actionText,
  onAction
}) => {
  return (
    <div style={{
      display: 'flex',
      flexDirection: 'column',
      alignItems: 'center',
      justifyContent: 'center',
      padding: '60px 20px',
      textAlign: 'center',
    }}>
      <InboxOutlined style={{ fontSize: 64, marginBottom: 16, color: '#d9d9d9' }} />
      <h3 style={{ marginBottom: 8, color: '#333' }}>{title}</h3>
      {description && (
        <p style={{ color: '#999', marginBottom: 24, maxWidth: 300 }}>
          {description}
        </p>
      )}
      {actionText && onAction && (
        <Button type="primary" icon={<PlusOutlined />} onClick={onAction}>
          {actionText}
        </Button>
      )}
    </div>
  )
}

export default EmptyState