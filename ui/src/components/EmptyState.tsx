import { Button } from 'antd'
import { PlusOutlined, InboxOutlined } from '@ant-design/icons'
import { useSelector } from 'react-redux'
import { RootState } from '../store'

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
  const theme = useSelector((state: RootState) => state.theme.theme)
  const isDark = theme === 'dark'

  return (
    <div style={{
      display: 'flex',
      flexDirection: 'column',
      alignItems: 'center',
      justifyContent: 'center',
      padding: '80px 20px',
      textAlign: 'center',
      animation: 'fadeInUp 0.4s ease-out',
    }}>
      <InboxOutlined style={{ 
        fontSize: 72, 
        marginBottom: 20, 
        color: isDark ? '#475569' : '#d1d5db' 
      }} />
      <h3 style={{ 
        marginBottom: 12, 
        color: isDark ? '#e2e8f0' : '#374151',
        fontSize: 20,
        fontWeight: 600
      }}>{title}</h3>
      {description && (
        <p style={{ 
          color: isDark ? '#94a3b8' : '#6b7280', 
          marginBottom: 28, 
          maxWidth: 350,
          lineHeight: 1.6,
          fontSize: 15
        }}>
          {description}
        </p>
      )}
      {actionText && onAction && (
        <Button 
          type="primary" 
          icon={<PlusOutlined />} 
          onClick={onAction}
          style={{
            height: 44,
            paddingLeft: 24,
            paddingRight: 24,
            fontSize: 15,
            borderRadius: 12,
            background: 'linear-gradient(135deg, #3b82f6 0%, #2563eb 100%)',
            boxShadow: '0 4px 12px rgba(59, 130, 246, 0.3)',
            border: 'none',
          }}
        >
          {actionText}
        </Button>
      )}
    </div>
  )
}

export default EmptyState