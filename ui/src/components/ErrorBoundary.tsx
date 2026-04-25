import { Component, ReactNode, ErrorInfo } from 'react'
import { Button, Result } from 'antd'
import { ReloadOutlined } from '@ant-design/icons'

interface Props {
  children: ReactNode
  fallback?: ReactNode
}

interface State {
  hasError: boolean
  error: Error | null
}

export class ErrorBoundary extends Component<Props, State> {
  constructor(props: Props) {
    super(props)
    this.state = { hasError: false, error: null }
  }

  static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error }
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo): void {
    console.error('ErrorBoundary caught an error:', error, errorInfo)
  }

  handleReload = (): void => {
    this.setState({ hasError: false, error: null })
    window.location.reload()
  }

  render(): ReactNode {
    if (this.state.hasError) {
      if (this.props.fallback) {
        return this.props.fallback
      }

      return (
        <Result
          status="error"
          title="页面加载失败"
          subTitle={this.state.error?.message || '发生了未知错误'}
          extra={
            <Button 
              type="primary" 
              icon={<ReloadOutlined />} 
              onClick={this.handleReload}
            >
              重新加载
            </Button>
          }
        />
      )
    }

    return this.props.children
  }
}

export default ErrorBoundary