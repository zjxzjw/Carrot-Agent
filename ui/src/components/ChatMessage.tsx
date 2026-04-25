import { Card, Avatar, Typography, Badge, Tooltip } from 'antd'
import { RobotOutlined, UserOutlined, ClockCircleOutlined } from '@ant-design/icons'
import type { Message } from '../store'
import './ChatMessage.css'

const { Text } = Typography

interface ChatMessageProps {
  message: Message
}

export function ChatMessage({ message }: ChatMessageProps) {
  const isUser = message.role === 'user'
  const avatar = isUser ? (
    <Avatar 
      icon={<UserOutlined />} 
      style={{ 
        background: 'linear-gradient(135deg, #007AFF 0%, #5856D6 100%)',
        boxShadow: '0 4px 12px rgba(0, 122, 255, 0.3)'
      }} 
    />
  ) : (
    <Avatar 
      icon={<RobotOutlined />} 
      style={{ 
        background: 'linear-gradient(135deg, #34C759 0%, #30D158 100%)',
        boxShadow: '0 4px 12px rgba(52, 199, 89, 0.3)'
      }} 
    />
  )

  const formatTime = (date: Date) => {
    return date.toLocaleTimeString('zh-CN', {
      hour: '2-digit',
      minute: '2-digit'
    })
  }

  return (
    <div className={`chat-message ${isUser ? 'user' : 'ai'}`}>
      <div className="message-avatar">{avatar}</div>
      <div className="message-content">
        <div className="message-header">
          <Text strong className="message-sender">
            {isUser ? '我' : 'Carrot Agent'}
          </Text>
          <Tooltip title={message.timestamp.toLocaleString()}>
            <span className="message-time">
              <ClockCircleOutlined /> {formatTime(message.timestamp)}
            </span>
          </Tooltip>
        </div>
        <Card size="small" className={`message-card ${isUser ? 'user' : 'ai'}`}>
          <Text className="message-text">{message.content}</Text>
        </Card>
      </div>
    </div>
  )
}

interface ChatListProps {
  messages: Message[]
}

export function ChatList({ messages }: ChatListProps) {
  if (messages.length === 0) {
    return (
      <div className="chat-empty">
        <RobotOutlined style={{ fontSize: 48, color: '#ccc' }} />
        <Text type="secondary">开始一段新对话吧！</Text>
      </div>
    )
  }

  return (
    <div className="chat-list">
      {messages.map((msg) => (
        <ChatMessage key={msg.id} message={msg} />
      ))}
    </div>
  )
}

interface ConnectionStatusProps {
  connected: boolean
}

export function ConnectionStatus({ connected }: ConnectionStatusProps) {
  return (
    <Badge
      status={connected ? 'success' : 'error'}
      text={connected ? '已连接' : '未连接'}
    />
  )
}