import { useState, useCallback, useRef, useEffect } from 'react'
import { useDispatch, useSelector } from 'react-redux'
import { Input, Button, Typography, Tooltip, message } from 'antd'
import { SendOutlined, ClearOutlined, RobotOutlined, UserOutlined, CopyOutlined, CheckOutlined } from '@ant-design/icons'
import { useTranslation } from 'react-i18next'
import { AppDispatch, RootState } from '../store'
import { sendMessage, clearMessages } from '../store/chatSlice'
import logo from '../assets/logo.png'

const { TextArea } = Input
const { Text, Title } = Typography

interface MessageItemProps {
  content: string
  role: 'user' | 'assistant'
  timestamp: number
}

const MessageItem = ({ content, role, timestamp }: MessageItemProps) => {
  const [copied, setCopied] = useState(false)
  const { t } = useTranslation()

  const handleCopy = useCallback(() => {
    navigator.clipboard.writeText(content)
    setCopied(true)
    message.success(t('chat.copySuccess'))
    setTimeout(() => setCopied(false), 2000)
  }, [content, t])

  return (
    <div className={`chat-message ${role}`}>
      <div className={`chat-avatar ${role}`}>
        {role === 'user' ? <UserOutlined /> : <RobotOutlined />}
      </div>
      <div className="chat-content">
        <div className="chat-header">
          <Text strong className="chat-sender" style={{ color: '#374151' }}>
            {t(`chat.sender.${role}`)}
          </Text>
          {role === 'assistant' && (
            <Tooltip title={copied ? t('chat.copied') : t('chat.copy')}>
              <Button
                type="text"
                size="small"
                icon={copied ? <CheckOutlined /> : <CopyOutlined />}
                onClick={handleCopy}
                className="copy-button"
              />
            </Tooltip>
          )}
        </div>
        <div className="chat-text">
          {content}
        </div>
        <Text type="secondary" className="chat-time" style={{ fontSize: 12, color: '#9ca3af' }}>
          {new Date(timestamp).toLocaleTimeString()}
        </Text>
      </div>
    </div>
  )
}

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
      <div style={{ padding: 40, textAlign: 'center', animation: 'fadeIn 0.3s ease-out' }}>
        <Title level={4} type="danger" style={{ marginBottom: 16 }}>
          {t('chat.errorTitle')}
        </Title>
        <Text type="danger" style={{ fontSize: 16, display: 'block', marginBottom: 16 }}>
          {error}
        </Text>
        <Button type="primary" onClick={handleClear}>
          {t('common.clear')}
        </Button>
      </div>
    )
  }

  return (
    <div className="chat-container">
      <div className="chat-messages">
        {messages.length === 0 ? (
          <div className="chat-empty">
            <div className="logo-container">
              <img
                src={logo}
                alt="Carrot Agent Logo"
                className="logo-image"
              />
            </div>
            <Title level={3} className="welcome-title">
              {t('chat.welcome')}
            </Title>
            <Text type="secondary" className="welcome-subtitle">
              {t('chat.empty')}
            </Text>
            <div className="suggestion-chips">
              <Button className="suggestion-chip" onClick={() => setInputValue(t('chat.suggestionTexts.greeting'))}>{t('chat.suggestions.greeting')}</Button>
              <Button className="suggestion-chip" onClick={() => setInputValue(t('chat.suggestionTexts.features'))}>{t('chat.suggestions.features')}</Button>
              <Button className="suggestion-chip" onClick={() => setInputValue(t('chat.suggestionTexts.code'))}>{t('chat.suggestions.code')}</Button>
            </div>
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
                <div className="loading-dots">
                  <span></span>
                  <span></span>
                  <span></span>
                </div>
                <Text type="secondary" style={{ color: '#6b7280' }}>
                  {t('chat.thinking')}
                </Text>
              </div>
            )}
            <div ref={messagesEndRef} />
          </>
        )}
      </div>

      <div className="chat-input-wrapper">
        <div className="chat-input-area">
          <TextArea
            value={inputValue}
            onChange={(e) => setInputValue(e.target.value)}
            onKeyPress={handleKeyPress}
            placeholder={t('chat.placeholder')}
            autoSize={{ minRows: 1, maxRows: 8 }}
            disabled={loading}
            className="chat-input"
          />
          <div className="chat-actions">
            <Tooltip title={t('chat.send')}>
              <Button
                type="primary"
                icon={<SendOutlined />}
                onClick={handleSend}
                disabled={inputValue.trim() === '' || loading}
                loading={loading}
                className="send-button"
              />
            </Tooltip>
            <Tooltip title={t('chat.clear')}>
              <Button
                icon={<ClearOutlined />}
                onClick={handleClear}
                disabled={messages.length === 0}
                className="clear-button"
              />
            </Tooltip>
          </div>
        </div>
        <div className="input-hint">
          <Text type="secondary" style={{ fontSize: 12 }}>
            {t('chat.inputHint')}
          </Text>
        </div>
      </div>
    </div>
  )
}

export default ChatPage
