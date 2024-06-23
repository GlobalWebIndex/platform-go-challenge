import PropTypes from 'prop-types'
import { createContext, useEffect, useState } from 'react'
import { useQuery } from 'react-query'

import axios from 'axios'

import { useAuth } from '../hooks/useAuth'

export const CurrentUserContext = createContext()

const fetchCurrentUser = async () => {
  try {
    const token = localStorage.getItem('token')

    const response = await axios.get('http://localhost:3000/api/v1/users', {
      ...(token && {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }),
    })

    return response.data
  } catch (e) {
    console.error('Error:', e.message)
    throw e.message
  }
}

export const CurrentUserProvider = ({ children }) => {
  const [currentUser, setCurrentUser] = useState(null)

  const auth = useAuth()
  const token = auth.token

  const { data, isLoading: loadingUser } = useQuery(
    'currentUser',
    fetchCurrentUser,
    {
      enabled: !!token,
      retries: false,
    }
  )

  useEffect(() => {
    if (!loadingUser && data) {
      const { currentUser } = data
      setCurrentUser(currentUser)
    }
  }, [data])

  useEffect(() => {
    if (!token && currentUser) {
      setCurrentUser(null)
    }
  }, [token])

  return (
    <CurrentUserContext.Provider value={{ currentUser, loadingUser }}>
      {children}
    </CurrentUserContext.Provider>
  )
}

CurrentUserProvider.propTypes = {
  children: PropTypes.any,
}

export default CurrentUserProvider
