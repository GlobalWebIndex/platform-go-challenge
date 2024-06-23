import { useState } from 'react'
import { Button, Form, Input, Space } from 'antd'

const DescriptionForm = ({ description, updateHandler, id, toggleEditing }) => {
  const [loading, setLoading] = useState(false)

  const onFinish = (values) => {
    setLoading(true)
    updateHandler(values)
      .then(() => setLoading(false))
      .catch(() => setLoading(false))
  }

  const validator = (_, value) => {
    if (value.length === 0) {
      return Promise.reject(new Error('No empty descriptions allowed'))
    }
    return Promise.resolve()
  }

  return (
    <Form
      name={`favorite-${id}`}
      layout="inline"
      onFinish={onFinish}
      initialValues={{
        description: description,
      }}
    >
      <Form.Item
        name="description"
        rules={[
          {
            validator: validator,
          },
        ]}
      >
        <Input
          type="text"
          placeholder="Please type a description"
          value={description}
        />
      </Form.Item>
      <Form.Item>
        <Space>
          <Button type="primary" htmlType="submit" loading={loading}>
            Save
          </Button>
          <Button type="primary" danger onClick={toggleEditing}>
            Cancel
          </Button>
        </Space>
      </Form.Item>
    </Form>
  )
}
export default DescriptionForm
