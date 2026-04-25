import { useState, useEffect, useCallback } from 'react'
import { useDispatch, useSelector } from 'react-redux'
import { Table, Modal, Form, Input, Select, message } from 'antd'
import { useTranslation } from 'react-i18next'
import { AppDispatch, RootState } from '../store'
import { fetchSkills, createSkill } from '../store/skillsSlice'
import { PageHeader, Loading, EmptyState } from '../components'
import { MEMORY_TYPE_OPTIONS } from '../types'

const { TextArea } = Input

const SkillsPage: React.FC = () => {
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [form] = Form.useForm()
  const [selectedType, setSelectedType] = useState('')
  const dispatch = useDispatch<AppDispatch>()
  const { skills, loading, error } = useSelector((state: RootState) => state.skills)
  const skillsData = skills || []
  const { t } = useTranslation()

  useEffect(() => {
    dispatch(fetchSkills())
  }, [dispatch])

  useEffect(() => {
    if (error) {
      message.error(error)
    }
  }, [error])

  const handleCreateSkill = useCallback(() => {
    form.validateFields().then((values) => {
      dispatch(createSkill(values))
        .unwrap()
        .then(() => {
          message.success(t('skills.success'))
          setIsModalOpen(false)
          form.resetFields()
          dispatch(fetchSkills())
        })
        .catch((err: string) => {
          message.error(err || t('common.error'))
        })
    })
  }, [form, dispatch, t])

  const handleOpenModal = useCallback(() => {
    setIsModalOpen(true)
  }, [])

  const handleCloseModal = useCallback(() => {
    setIsModalOpen(false)
    form.resetFields()
  }, [form])

  const filteredSkills = selectedType
    ? skillsData.filter(skill => skill.platforms?.includes(selectedType))
    : skillsData

  const columns = [
    {
      title: t('skills.name'),
      dataIndex: 'name',
      key: 'name',
      width: 150,
    },
    {
      title: t('skills.description'),
      dataIndex: 'description',
      key: 'description',
      ellipsis: true,
    },
    {
      title: 'Version',
      dataIndex: 'version',
      key: 'version',
      width: 100,
    },
    {
      title: 'Platforms',
      dataIndex: 'platforms',
      key: 'platforms',
      width: 120,
      render: (platforms: string) => platforms || '-',
    },
    {
      title: 'Created At',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 180,
    },
  ]

  if (loading && skillsData.length === 0) {
    return <Loading size="large" />
  }

  return (
    <div>
      <PageHeader
        title={t('skills.title')}
        actionText={t('skills.create')}
        onAction={handleOpenModal}
        extra={
          <Select
            value={selectedType}
            onChange={setSelectedType}
            style={{ width: 140 }}
            options={MEMORY_TYPE_OPTIONS}
            placeholder={t('skills.filter')}
          />
        }
      />

      {filteredSkills.length === 0 ? (
        <EmptyState
          title={t('skills.empty')}
          description={t('skills.emptyDesc')}
          actionText={t('skills.create')}
          onAction={handleOpenModal}
        />
      ) : (
        <Table
          columns={columns}
          dataSource={filteredSkills}
          rowKey="id"
          loading={loading}
          pagination={{ pageSize: 10, showSizeChanger: true }}
        />
      )}

      <Modal
        title={t('skills.create')}
        open={isModalOpen}
        onOk={handleCreateSkill}
        onCancel={handleCloseModal}
        okText={t('common.create')}
        cancelText={t('common.cancel')}
        destroyOnClose
      >
        <Form form={form} layout="vertical" preserve={false}>
          <Form.Item
            name="name"
            label={t('skills.name')}
            rules={[{ required: true, message: t('skills.name') + t('common.required') }]}
          >
            <Input placeholder={t('skills.name')} maxLength={50} showCount />
          </Form.Item>
          <Form.Item
            name="description"
            label={t('skills.description')}
            rules={[{ required: true, message: t('skills.description') + t('common.required') }]}
          >
            <Input placeholder={t('skills.description')} maxLength={200} showCount />
          </Form.Item>
          <Form.Item
            name="content"
            label={t('skills.content')}
            rules={[{ required: true, message: t('skills.content') + t('common.required') }]}
          >
            <TextArea rows={6} placeholder={t('skills.content') + '...'} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default SkillsPage