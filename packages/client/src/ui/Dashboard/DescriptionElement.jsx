import { Button, Space } from 'antd'

import DescriptionForm from './DescriptionForm'
import { EditOutlined } from '@ant-design/icons'
import { useState } from 'react'

const DescriptionElement = ({
  editing,
  description,
  itemId,
  setEditingHandler,
  updateHandler,
}) => {
  const toggleEditing = () => setEditingHandler(!editing)

  return editing ? (
    <DescriptionForm
      description={description}
      id={itemId}
      updateHandler={updateHandler}
      toggleEditing={toggleEditing}
    />
  ) : (
    <Space>
      <Button
        type="text"
        shape="circle"
        onClick={toggleEditing}
        icon={<EditOutlined />}
      />
      <span>{description}</span>
    </Space>
  )
}

export default DescriptionElement
