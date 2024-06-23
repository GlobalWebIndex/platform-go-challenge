import { useState } from 'react'

import { useQuery, useMutation, useQueryClient } from 'react-query'
import { App } from 'antd'

import Dashboard from '../ui/Dashboard/Dashboard'

import { useCurrentUser } from '../hooks/useCurrentUser'
import {
  getUserFavorites,
  removeFromFavorites,
  updateAssetDescription,
} from '../utils/network'

const FavoritesPagePage = () => {
  const { notification } = App.useApp()

  const { currentUser } = useCurrentUser()
  const queryClient = useQueryClient()

  const [paginationSpecs, setPaginationSpecs] = useState({
    page: 1,
    pageSize: 10,
  })

  const showErrorNotification = (eMsg) => {
    notification.error({
      message: 'Oops, something went wrong',
      description: eMsg,
      placement: 'topRight',
    })
  }
  const showSuccessNotification = (message) => {
    notification.success({
      message,
      placement: 'topRight',
    })
  }

  const onPaginationChange = (page, pageSize) => {
    setPaginationSpecs({ page, pageSize })
  }

  const { isLoading, isError, data, error } = useQuery({
    queryKey: [
      'favorites',
      currentUser?.id,
      paginationSpecs.page,
      paginationSpecs.pageSize,
    ],
    queryFn: () =>
      getUserFavorites(
        currentUser?.id,
        paginationSpecs.page,
        paginationSpecs.pageSize
      ),
    enabled: !!currentUser,
  })

  const updateAssetDescriptionMutation = useMutation({
    mutationFn: updateAssetDescription,
    onError: ({ message }) => {
      showErrorNotification(message)
    },
    onSuccess: async () => {
      showSuccessNotification('Description updated')
      queryClient.invalidateQueries({
        queryKey: [
          'favorites',
          currentUser?.id,
          paginationSpecs.page,
          paginationSpecs.pageSize,
        ],
      })
    },
  })

  const removeFromFavoritesMutation = useMutation({
    mutationFn: removeFromFavorites,
    onError: ({ message }) => {
      showErrorNotification(message)
    },
    onSuccess: async () => {
      showSuccessNotification('Removed from favorites')
    },
    onSettled: async () => {
      queryClient.invalidateQueries({
        queryKey: [
          'favorites',
          currentUser?.id,
          paginationSpecs.page,
          paginationSpecs.pageSize,
        ],
      })
    },
  })

  const onMainAction = (favoriteId) => {
    const payload = { favoriteId, userId: currentUser.id }
    return removeFromFavoritesMutation.mutateAsync(payload)
  }

  const onSecondaryAction = (assetId, description) => {
    const payload = { assetId, description }
    return updateAssetDescriptionMutation.mutateAsync(payload)
  }

  if (isError) {
    showErrorNotification(error)
  }

  return (
    <Dashboard
      items={data?.favorites}
      isLoading={isLoading}
      paginationHandler={onPaginationChange}
      pageSize={paginationSpecs.pageSize}
      shouldDisplayRemove={true}
      mainActionHandler={onMainAction}
      secondaryActionHandler={onSecondaryAction}
    />
  )
}

export default FavoritesPagePage
