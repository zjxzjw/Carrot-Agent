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
  <div style={{ 
    display: 'flex', 
    alignItems: 'flex-start', 
    marginBottom: 16,
    animation: 'fadeIn 0.3s ease-in-out',
  }}>
    <div style={{ 
      width: 40, 
      height: 40, 
      borderRadius: '50%', 
      display: 'flex', 
      alignItems: 'center', 
      justifyContent: 'center',
      marginRight: 12,
      flexShrink: 0,
      backgroundColor: role === 'user' ? '#1890ff' : '#f0f0f0',
      color: role === 'user' ? '#fff' : '#000',
      fontWeight: 'bold',
    }}>
      {role === 'user' ? 'U' : 'AI'}
    </div>
    <div style={{ flex: 1, minWidth: 0 }}>
      <Text strong style={{ display: 'block', marginBottom: 4 }}>
        {role === 'user' ? '你' : 'Carrot Agent'}
      </Text>
      <div style={{ 
        lineHeight: 1.6, 
        whiteSpace: 'pre-wrap',
        wordBreak: 'break-word',
      }}>
        {content}
      </div>
      <Text type="secondary" style={{ fontSize: 12, display: 'block', marginTop: 4 }}>
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
    <div style={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
      <style>{`
        @keyframes fadeIn {
          from { opacity: 0; transform: translateY(10px); }
          to { opacity: 1; transform: translateY(0); }
        }
      `}</style>

      <div style={{ 
        flex: 1, 
        overflowY: 'auto', 
        padding: 16,
        backgroundColor: '#fafafa',
        borderRadius: 8,
        marginBottom: 16,
      }}>
        {messages.length === 0 ? (
          <div style={{ 
            display: 'flex', 
            flexDirection: 'column', 
            alignItems: 'center', 
            justifyContent: 'center',
            height: '100%',
            color: '#999',
          }}>
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
              <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                <Spin size="small" />
                <Text type="secondary">{t('chat.thinking')}</Text>
              </div>
            )}
            <div ref={messagesEndRef} />
          </>
        )}
      </div>

      <div style={{ display: 'flex', gap: 12, alignItems: 'flex-end' }}>
        <TextArea
          value={inputValue}
          onChange={(e) => setInputValue(e.target.value)}
          onKeyPress={handleKeyPress}
          placeholder={t('chat.placeholder')}
          autoSize={{ minRows: 2, maxRows: 6 }}
          disabled={loading}
          style={{ flex: 1 }}
        />
        <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
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