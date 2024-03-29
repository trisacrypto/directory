import { createSelector } from 'reselect';

const vaspDetailsState = (state) => state.VaspDetails;

export const getVaspDetails = createSelector(vaspDetailsState, (state) => state.data);
export const getVaspDetailsLoadingState = createSelector(
  vaspDetailsState,
  (state) => state.loading
);
export const getVaspDetailsErrorState = createSelector(vaspDetailsState, (state) => state.error);
export const getTrisaDetailsErrorState = createSelector(
  vaspDetailsState,
  (state) => state.trisaError
);
export const getIvmsErrorState = createSelector(vaspDetailsState, (state) => state.ivmsError);

export const getContacts = createSelector(getVaspDetails, (state) => state.vasp.contacts);
export const getContactErrorState = createSelector(vaspDetailsState, (state) => state.contactError);
