const favoriteToAssetMapper = (item) => ({
  RelatedType: item.Asset.RelatedType,
  ID: item.ID,
  AssetID: item.Asset.ID,
  Description: item.Asset.Description,
  ...(item.Asset.RelatedType === 'charts' && { Chart: item.Asset.Chart }),
  ...(item.Asset.RelatedType === 'audiences' && {
    Audience: item.Asset.Audience,
  }),
  ...(item.Asset.RelatedType === 'insights' && { Insight: item.Asset.Insight }),
})
export { favoriteToAssetMapper }
