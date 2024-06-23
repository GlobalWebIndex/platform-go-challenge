import axios from 'axios'

const getAssets = async (page = 1, pageSize = 10, type = null) => {
  try {
    const token = localStorage.getItem('token')

    if (!token) {
      throw Error('missing token')
    }
    const res = await axios.get(
      `http://localhost:3000/api/v1/assets?page=${page}&limit=${pageSize}${
        type ? `&type=${type}` : ''
      }`,
      {
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
      }
    )

    return res.data
  } catch (e) {
    throw new Error(e.message)
  }
}

const getUserFavorites = async (userId, page = 1, pageSize = 10) => {
  try {
    const token = localStorage.getItem('token')

    if (!token) {
      throw Error('missing token')
    }
    const res = await axios.get(
      `http://localhost:3000/api/v1/users/${userId}/favorites?page=${page}&limit=${pageSize}`,
      {
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
      }
    )

    return res.data
  } catch (e) {
    throw new Error(e.message)
  }
}

const addToFavorites = async ({ assetId, userId }) => {
  try {
    const token = localStorage.getItem('token')

    if (!token) {
      throw Error('missing token')
    }
    const res = await axios.post(
      `http://localhost:3000/api/v1/users/${userId}/favorites`,
      { id: assetId },
      {
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
      }
    )

    return res.data
  } catch (e) {
    let msg = e.message

    if (
      e.response?.data?.error?.includes(
        'ERROR: duplicate key value violates unique constraint'
      )
    ) {
      msg = 'already in favorites of user'
    }
    throw new Error(msg)
  }
}

const updateAssetDescription = async ({ assetId, description }) => {
  try {
    const token = localStorage.getItem('token')

    if (!token) {
      throw Error('missing token')
    }

    const res = await axios.patch(
      `http://localhost:3000/api/v1/assets/${assetId}`,
      { description },
      {
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
      }
    )

    return res.data
  } catch (e) {
    throw new Error(e.message)
  }
}

const removeFromFavorites = async ({ favoriteId, userId }) => {
  try {
    const token = localStorage.getItem('token')

    if (!token) {
      throw Error('missing token')
    }

    const res = await axios.delete(
      `http://localhost:3000/api/v1/users/${userId}/favorites/${favoriteId}`,
      {
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
      }
    )

    return res.data
  } catch (e) {
    let msg = e.message
    throw new Error(msg)
  }
}
export {
  getAssets,
  getUserFavorites,
  addToFavorites,
  removeFromFavorites,
  updateAssetDescription,
}
