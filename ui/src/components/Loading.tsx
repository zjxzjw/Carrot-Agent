import { Spin } from 'antd'
import { useSelector } from 'react-redux'
import { RootState } from '../store'

interface LoadingProps {
  size?: 'small' | 'default' | 'large'
  tip?: string
  fullscreen?: boolean
}

export const Loading: React.FC<LoadingProps> = ({ 
  size = 'default', 
  tip,
  fullscreen = false 
}) => {
  const theme = useSelector((state: RootState) => state.theme.theme)
  const isDark = theme === 'dark'

  if (fullscreen) {
    return (
      <div style={{ 
        display: 'flex', 
        justifyContent: 'center', 
        alignItems: 'center', 
        height: '100vh',
        width: '100vw',
        position: 'fixed',
        top: 0,
        left: 0,
        backgroundColor: isDark ? 'rgba(15, 23, 42, 0.95)' : 'rgba(255, 255, 255, 0.95)',
        zIndex: 9999,
        backdropFilter: 'blur(8px)',
        animation: 'fadeIn 0.3s ease-out',
      }}>
        <Spin size={size} tip={tip} />
      </div>
    )
  }

  return (
    <div style={{ 
      display: 'flex', 
      justifyContent: 'center', 
      padding: 40,
      animation: 'fadeIn 0.3s ease-out'
    }}>
      <Spin size={size} tip={tip} />
    </div>
  )
}

export default Loading