import React, { useState } from 'react'
import { Form, Input, Button, Card, Alert, Typography } from 'antd'
import { LockOutlined, UserOutlined } from '@ant-design/icons'
import { useTranslation } from 'react-i18next'
import { useNavigate } from 'react-router-dom'
import { useDispatch } from 'react-redux'
import { login } from '../store/authSlice'
import logo from '../assets/logo.png'

const { Title, Text } = Typography

const LoginPage: React.FC = () => {
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const { t } = useTranslation()
  const navigate = useNavigate()
  const dispatch = useDispatch()

  const onFinish = async (values: { username: string; password: string }) => {
    setLoading(true)
    setError('')
    
    try {
      const response = await fetch('/api/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(values),
      })

      if (response.ok) {
        const data = await response.json()
        dispatch(login(data.session_id))
        navigate('/')
      } else {
        const errorData = await response.json()
        setError(errorData.error || t('login.error.invalid_credentials'))
      }
    } catch (err) {
      setError(t('login.error.network_error'))
    } finally {
      setLoading(false)
    }
  }

  return (
    <div style={{ 
      display: 'flex', 
      justifyContent: 'center', 
      alignItems: 'center', 
      minHeight: '100vh',
      background: '#f5f5f7',
      fontFamily: '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif'
    }}>
      <Card 
        style={{ 
          width: 400,
          borderRadius: '20px',
          boxShadow: '0 20px 60px rgba(0, 0, 0, 0.08)',
          border: 'none',
          background: '#ffffff'
        }}
        styles={{ body: { padding: '48px 40px' } }}
      >
        {/* Logo and Title */}
        <div style={{ textAlign: 'center', marginBottom: '40px' }}>
          <img 
            src={logo} 
            alt="Carrot Agent Logo" 
            style={{ 
              width: 72, 
              height: 72,
              marginBottom: '20px'
            }}
          />
          <Title level={2} style={{ 
            margin: 0, 
            fontSize: '28px', 
            fontWeight: 600,
            color: '#1d1d1f',
            letterSpacing: '-0.5px'
          }}>
            {t('login.title')}
          </Title>
          <Text style={{ 
            fontSize: '15px', 
            marginTop: '8px', 
            display: 'block',
            color: '#86868b'
          }}>
            {t('login.subtitle')}
          </Text>
        </div>
        
        {error && (
          <Alert 
            message={t('login.error.title')} 
            description={error} 
            type="error" 
            showIcon 
            closable
            style={{ 
              marginBottom: '24px', 
              borderRadius: '12px',
              border: 'none',
              background: '#ffebee'
            }}
          />
        )}
        
        <Form
          name="login"
          initialValues={{ 
            remember: true,
            username: 'admin',
            password: 'admin123'
          }}
          onFinish={onFinish}
          size="large"
          layout="vertical"
        >
          <Form.Item
            label={<span style={{ fontSize: '13px', fontWeight: 500, color: '#1d1d1f' }}>{t('login.placeholder.username')}</span>}
            name="username"
            rules={[
              { required: true, message: t('login.error.username_required') },
            ]}
            style={{ marginBottom: '20px' }}
          >
            <Input 
              prefix={<UserOutlined style={{ color: '#86868b' }} />} 
              placeholder={t('login.placeholder.username')}
              style={{ 
                borderRadius: '12px',
                height: '48px',
                background: '#f5f5f7',
                border: '1px solid transparent',
                transition: 'all 0.2s ease'
              }}
              onFocus={(e) => {
                e.target.style.background = '#ffffff'
                e.target.style.borderColor = '#0071e3'
              }}
              onBlur={(e) => {
                e.target.style.background = '#f5f5f7'
                e.target.style.borderColor = 'transparent'
              }}
            />
          </Form.Item>
          
          <Form.Item
            label={<span style={{ fontSize: '13px', fontWeight: 500, color: '#1d1d1f' }}>{t('login.placeholder.password')}</span>}
            name="password"
            rules={[
              { required: true, message: t('login.error.password_required') },
            ]}
            style={{ marginBottom: '28px' }}
          >
            <Input.Password
              prefix={<LockOutlined style={{ color: '#86868b' }} />}
              placeholder={t('login.placeholder.password')}
              style={{ 
                borderRadius: '12px',
                height: '48px',
                background: '#f5f5f7',
                border: '1px solid transparent',
                transition: 'all 0.2s ease'
              }}
              onFocus={(e) => {
                const input = e.target.querySelector('input')
                if (input) {
                  input.style.background = '#ffffff'
                  input.style.borderColor = '#0071e3'
                }
              }}
              onBlur={(e) => {
                const input = e.target.querySelector('input')
                if (input) {
                  input.style.background = '#f5f5f7'
                  input.style.borderColor = 'transparent'
                }
              }}
            />
          </Form.Item>
          
          <Form.Item style={{ marginBottom: 0 }}>
            <Button 
              type="primary" 
              htmlType="submit" 
              style={{ 
                width: '100%',
                height: '48px',
                borderRadius: '12px',
                fontSize: '16px',
                fontWeight: 500,
                background: '#0071e3',
                border: 'none',
                boxShadow: 'none',
                transition: 'all 0.2s ease'
              }}
              loading={loading}
              onMouseEnter={(e) => {
                e.currentTarget.style.background = '#0077ed'
              }}
              onMouseLeave={(e) => {
                e.currentTarget.style.background = '#0071e3'
              }}
            >
              {t('login.button.login')}
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  )
}

export default LoginPage