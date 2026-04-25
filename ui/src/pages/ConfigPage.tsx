import { useState, useEffect, useCallback } from 'react'
import { Form, Input, InputNumber, Select, Button, message, Card } from 'antd'
import { useTranslation } from 'react-i18next'
import { configService } from '../services'
import { PageHeader, Loading } from '../components'

const { Password } = Input

interface ModelList {
  provider: string
  models: string[]
}

const ConfigPage: React.FC = () => {
  const [form] = Form.useForm()
  const [loading, setLoading] = useState(false)
  const [saving, setSaving] = useState(false)
  const [models, setModels] = useState<ModelList[]>([])
  const [currentProvider, setCurrentProvider] = useState('openai')
  const { t } = useTranslation()

  useEffect(() => {
    loadConfig()
  }, [])

  const loadConfig = useCallback(async () => {
    setLoading(true)
    try {
      const configResponse = await configService.getConfig()
      const modelsResponse = await configService.getModels()

      const config: any = configResponse.data
      setCurrentProvider(config.Model.Provider)

      form.setFieldsValue({
        provider: config.Model.Provider || 'openai',
        apiKey: config.Model.APIKey || '',
        modelName: config.Model.ModelName || '',
        baseUrl: config.Model.BaseURL || '',
        temperature: config.Model.Temperature || 0.7,
        maxTokens: config.Model.MaxTokens || 4096,
      })

      setModels(modelsResponse.data || [])
    } catch (error) {
      console.error('Failed to load config:', error)
      message.error(t('common.error') + ': ' + (error instanceof Error ? error.message : '加载配置失败'))
      setModels([])
    } finally {
      setLoading(false)
    }
  }, [form, t])

  const handleSave = useCallback(async () => {
    form.validateFields().then(async (values) => {
      setSaving(true)
      try {
        await configService.updateConfig({
          model: {
            provider: values.provider,
            api_key: values.apiKey,
            model_name: values.modelName,
            base_url: values.baseUrl,
            temperature: values.temperature,
            max_tokens: values.maxTokens,
          },
        })
        message.success(t('common.success') + ': ' + t('config.updated'))
      } catch (error) {
        message.error(t('common.error') + ': ' + (error instanceof Error ? error.message : '保存配置失败'))
      } finally {
        setSaving(false)
      }
    })
  }, [form, t])

  const handleProviderChange = useCallback((value: string) => {
    setCurrentProvider(value)
    form.setFieldsValue({ modelName: '' })
  }, [form])

  const getModelsForCurrentProvider = useCallback(() => {
    const provider = models.find(m => m.provider === currentProvider)
    return provider?.models || []
  }, [models, currentProvider])

  if (loading) {
    return <Loading size="large" />
  }

  const providerOptions = models.map(m => ({ value: m.provider, label: m.provider }))

  return (
    <div>
      <PageHeader
        title={t('config.title')}
        showAction={false}
      />

      <Card
        title={t('config.modelSettings')}
        variant="outlined"
        style={{ marginBottom: 24 }}
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSave}
        >
          <Form.Item
            name="provider"
            label={t('config.provider')}
            rules={[{ required: true, message: t('config.provider') + t('common.required') }]}
          >
            <Select
              placeholder={t('config.selectProvider')}
              onChange={handleProviderChange}
              options={providerOptions}
            />
          </Form.Item>

          <Form.Item
            name="apiKey"
            label={t('config.apiKey')}
            rules={[{ required: true, message: t('config.apiKey') + t('common.required') }]}
          >
            <Password
              placeholder={t('config.enterApiKey')}
              allowClear
            />
          </Form.Item>

          <Form.Item
            name="modelName"
            label={t('config.modelName')}
            rules={[{ required: true, message: t('config.modelName') + t('common.required') }]}
          >
            <Select
              placeholder={t('config.selectModel')}
              options={getModelsForCurrentProvider().map(model => ({ value: model, label: model }))}
            />
          </Form.Item>

          <Form.Item
            name="baseUrl"
            label={t('config.baseUrl')}
          >
            <Input
              placeholder={t('config.enterBaseUrl')}
              allowClear
            />
          </Form.Item>

          <Form.Item
            name="temperature"
            label={t('config.temperature')}
            rules={[
              { required: true, message: t('config.temperature') + t('common.required') },
              { min: 0, max: 2, message: t('config.temperatureRange') }
            ]}
          >
            <InputNumber
              min={0}
              max={2}
              step={0.1}
              precision={1}
              style={{ width: '100%' }}
            />
          </Form.Item>

          <Form.Item
            name="maxTokens"
            label={t('config.maxTokens')}
            rules={[
              { required: true, message: t('config.maxTokens') + t('common.required') },
              { min: 1, max: 128000, message: t('config.maxTokensRange') }
            ]}
          >
            <InputNumber
              min={1}
              max={128000}
              style={{ width: '100%' }}
            />
          </Form.Item>

          <Form.Item>
            <Button
              type="primary"
              htmlType="submit"
              loading={saving}
            >
              {t('common.save')}
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  )
}

export default ConfigPage