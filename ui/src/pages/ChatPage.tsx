import { useState, useCallback, useRef, useEffect } from 'react'
import { useDispatch, useSelector } from 'react-redux'
import { Input, Button, Typography, Spin } from 'antd'
import { SendOutlined, ClearOutlined } from '@ant-design/icons'
import { useTranslation } from 'react-i18next'
import { AppDispatch, RootState } from '../store'
import { sendMessage, clearMessages } from '../store/chatSlice'
import logo from '../assets/logo.png'

const { TextArea } = Input
const { Text } = Typography

interface MessageItemProps {
  content: string
  role: 'user' | 'assistant'
  timestamp: number
}

const MessageItem = ({ content, role, timestamp }: MessageItemProps) => (
  <div className="chat-message">
    <div className={`chat-avatar ${role}`}>
      {role === 'user' ? 'U' : 'AI'}
    </div>
    <div className="chat-content">
      <Text strong className="chat-sender">
        {role === 'user' ? '你' : 'Carrot Agent'}
      </Text>
      <div className="chat-text">
        {content}
      </div>
      <Text type="secondary" className="chat-time" style={{ fontSize: 12 }}>
        {new Date(timestamp).toLocaleTimeString()}
      </Text>
    </div>
  </div>
)

const ChatPage: React.FC = () => {
  const [inputValue, setInputValue] = useState('')
  const messagesEndRef = useRef<HTMLDivElement>(null)
  const dispatch = useDispatch<AppDispatch>()
  const { messages, loading, error } = useSelector((state: RootState) => state.chat)
  const { t } = useTranslation()

  const scrollToBottom = useCallback(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [])

  useEffect(() => {
    scrollToBottom()
  }, [messages, scrollToBottom])

  const handleSend = useCallback(() => {
    if (inputValue.trim() === '' || loading) {
      return
    }
    dispatch(sendMessage(inputValue.trim()))
    setInputValue('')
  }, [inputValue, loading, dispatch])

  const handleKeyPress = useCallback((e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleSend()
    }
  }, [handleSend])

  const handleClear = useCallback(() => {
    dispatch(clearMessages())
  }, [dispatch])

  if (error) {
    return (
      <div style={{ padding: 20, textAlign: 'center' }}>
        <Text type="danger">{error}</Text>
        <br />
        <Button type="link" onClick={handleClear}>{t('common.clear')}</Button>
      </div>
    )
  }

  return (
    <div className="chat-container">
      <div className="chat-messages">
        {messages.length === 0 ? (
          <div className="chat-empty">
            <img
              src={logo}
              alt="Carrot Agent Logo"
              style={{ width: 64, height: 64, marginBottom: 16 }}
            />
            <Text type="secondary">{t('chat.empty')}</Text>
          </div>
        ) : (
          <>
            {messages.map((item, index) => (
              <MessageItem
                key={`${item.id}-${index}`}
                content={item.content}
                role={item.role}
                timestamp={item.timestamp}
              />
            ))}
            {loading && (
              <div className="chat-loading">
                <Spin size="small" />
                <Text type="secondary">{t('chat.thinking')}</Text>
              </div>
            )}
            <div ref={messagesEndRef} />
          </>
        )}
      </div>

      <div className="chat-input-area">
        <TextArea
          value={inputValue}
          onChange={(e) => setInputValue(e.target.value)}
          onKeyPress={handleKeyPress}
          placeholder={t('chat.placeholder')}
          autoSize={{ minRows: 2, maxRows: 6 }}
          disabled={loading}
        />
        <div className="chat-actions">
          <Button
            type="primary"
            icon={<SendOutlined />}
            onClick={handleSend}
            disabled={inputValue.trim() === '' || loading}
            loading={loading}
          >
            {t('chat.send')}
          </Button>
          <Button
            icon={<ClearOutlined />}
            onClick={handleClear}
            disabled={messages.length === 0}
          >
            {t('chat.clear')}
          </Button>
        </div>
      </div>
    </div>
  )
}

export default ChatPage