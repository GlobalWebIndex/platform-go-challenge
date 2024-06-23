import { useState } from 'react'
import { List, Button } from 'antd'

import DescriptionElement from './DescriptionElement'
import { StarOutlined, DeleteOutlined } from '@ant-design/icons'

import { favoriteToAssetMapper } from '../../utils/itemMappers'
import {
  generateItemAvatar,
  generateItemContent,
  generateItemTitle,
} from './generators'

const ListItem = ({
  item,
  mainActionHandler,
  shouldDisplayRemove,
  secondaryActionHandler,
}) => {
  const convertedItem = shouldDisplayRemove ? favoriteToAssetMapper(item) : item
  const { ID, Description, RelatedType } = convertedItem
  const [loading, setLoading] = useState(false)
  const [editing, setEditing] = useState(false)

  const onUpdateHandler = ({ description }) => {
    const whichId = shouldDisplayRemove
      ? convertedItem.AssetID
      : convertedItem.ID
    setEditing(false)
    return secondaryActionHandler(whichId, description)
  }

  const onClickHandler = () => {
    setLoading(true)
    mainActionHandler(ID)
      .then(() => setLoading(false))
      .catch(() => setLoading(false))
  }

  const actions = !shouldDisplayRemove
    ? [
        <Button
          type="primary"
          loading={loading}
          onClick={onClickHandler}
          icon={<StarOutlined />}
          key={`add-favorite-${ID}`}
        >
          Add to Favorites
        </Button>,
      ]
    : [
        <Button
          type="primary"
          danger
          loading={loading}
          icon={<DeleteOutlined />}
          onClick={onClickHandler}
          key={`remove-favorite-${ID}`}
        >
          Remove from Favorites
        </Button>,
      ]

  return (
    <List.Item key={ID} actions={actions}>
      <List.Item.Meta
        avatar={generateItemAvatar(RelatedType)}
        title={generateItemTitle(RelatedType)}
        description={
          <DescriptionElement
            description={Description}
            itemId={ID}
            editing={editing}
            setEditingHandler={setEditing}
            updateHandler={onUpdateHandler}
          />
        }
      />
      {generateItemContent(RelatedType, convertedItem)}
    </List.Item>
  )
}

export default ListItem
