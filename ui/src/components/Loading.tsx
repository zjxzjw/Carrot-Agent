import { Spin } from 'antd'

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
        backgroundColor: 'rgba(255, 255, 255, 0.9)',
        zIndex: 9999,
      }}>
        <Spin size={size} tip={tip} />
      </div>
    )
  }

  return (
    <div style={{ display: 'flex', justifyContent: 'center', padding: 40 }}>
      <Spin size={size} tip={tip} />
    </div>
  )
}

export default Loading