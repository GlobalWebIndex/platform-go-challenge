import { useContext } from 'react'
import { CurrentUserContext } from '../providers/CurrentUserProvider'

export const useCurrentUser = () => {
  return useContext(CurrentUserContext)
}
