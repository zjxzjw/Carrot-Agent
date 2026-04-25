import { ReactNode } from 'react'
import { Modal, Form, FormProps } from 'antd'

interface ModalFormProps<T = Record<string, unknown>> {
  open: boolean
  title: string
  form: ReturnType<typeof Form.useForm>[0]
  onOk: () => void
  onCancel: () => void
  okText?: string
  cancelText?: string
  children: ReactNode
  formProps?: FormProps<T>
  width?: number
}

export function ModalForm<T = Record<string, unknown>>({
  open,
  title,
  form,
  onOk,
  onCancel,
  okText,
  cancelText,
  children,
  width = 520,
}: ModalFormProps<T>) {
  return (
    <Modal
      title={title}
      open={open}
      onOk={onOk}
      onCancel={onCancel}
      okText={okText}
      cancelText={cancelText}
      destroyOnClose
      width={width}
    >
      <Form
        form={form}
        layout="vertical"
        preserve={false}
      >
        {children}
      </Form>
    </Modal>
  )
}