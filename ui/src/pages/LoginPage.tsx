import React, { useState } from 'react'
import { Form, Input, Button, Card, Alert, Typography, Space } from 'antd'
import { LockOutlined, UserOutlined } from '@ant-design/icons'
import { useTranslation } from 'react-i18next'
import { useNavigate } from 'react-router-dom'
import { useDispatch } from 'react-redux'
import { login } from '../store/authSlice'

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
      background: '#f0f2f5'
    }}>
      <Card 
        style={{ 
          width: 400, 
          padding: '24px',
          boxShadow: '0 4px 12px rgba(0,0,0,0.15)'
        }}
      >
        <div style={{ textAlign: 'center', marginBottom: '24px' }}>
          <Title level={3}>{t('login.title')}</Title>
          <Text type="secondary">{t('login.subtitle')}</Text>
        </div>
        
        {error && (
          <Alert 
            message={t('login.error.title')} 
            description={error} 
            type="error" 
            showIcon 
            style={{ marginBottom: '16px' }}
          />
        )}
        
        <Form
          name="login"
          initialValues={{ remember: true }}
          onFinish={onFinish}
        >
          <Form.Item
            name="username"
            rules={[
              { required: true, message: t('login.error.username_required') },
            ]}
          >
            <Input 
              prefix={<UserOutlined className="site-form-item-icon" />} 
              placeholder={t('login.placeholder.username')}
              size="large"
            />
          </Form.Item>
          
          <Form.Item
            name="password"
            rules={[
              { required: true, message: t('login.error.password_required') },
            ]}
          >
            <Input
              prefix={<LockOutlined className="site-form-item-icon" />}
              type="password"
              placeholder={t('login.placeholder.password')}
              size="large"
            />
          </Form.Item>
          
          <Form.Item>
            <Button 
              type="primary" 
              htmlType="submit" 
              style={{ width: '100%' }} 
              size="large"
              loading={loading}
            >
              {t('login.button.login')}
            </Button>
          </Form.Item>
          
          <Space style={{ display: 'flex', justifyContent: 'center', marginTop: '16px' }}>
            <Text type="secondary">{t('login.default_credentials')}</Text>
          </Space>
        </Form>
      </Card>
    </div>
  )
}

export default LoginPage