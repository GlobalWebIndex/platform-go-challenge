import PropTypes from 'prop-types'
import { createContext, useState } from 'react'
import { useLocation, useNavigate } from 'react-router-dom'
import { useQueryClient } from 'react-query'

export const AuthContext = createContext()

export const AuthProvider = ({ children }) => {
  const navigate = useNavigate()
  const { search } = useLocation()
  const queryClient = useQueryClient()

  const [token, setToken] = useState(localStorage.getItem('token') || '')

  const login = (token) => {
    localStorage.setItem('token', token)
    setToken(token)
    navigate('/')
  }

  const logout = () => {
    setToken('')
    localStorage.removeItem('token')
    queryClient.removeQueries('currentUser')
    navigate('/login')
  }

  return (
    <AuthContext.Provider value={{ token, logout, login }}>
      {children}
    </AuthContext.Provider>
  )
}

AuthProvider.propTypes = {
  children: PropTypes.any,
}

export default AuthProvider
