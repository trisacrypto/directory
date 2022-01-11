import { createSelector } from 'reselect'

const vaspDetailsState = state => state.VaspDetails

export const getVaspDetails = createSelector(vaspDetailsState, state => state.data)
export const getVaspDetailsLoadingState = createSelector(vaspDetailsState, state => state.loading)
export const getVaspDetailsErrorState = createSelector(vaspDetailsState, state => state.error)