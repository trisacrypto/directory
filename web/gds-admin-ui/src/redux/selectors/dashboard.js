import { createSelector } from 'reselect';

const summaryState = (state) => state.Summary;
const pendingVaspState = (state) => state.Vasps;

export const getSummaryData = createSelector(summaryState, (state) => state.data);
export const getSummaryLoadingState = createSelector(summaryState, (state) => state.loading);

export const getPendingVaspsData = createSelector(pendingVaspState, (state) => state.data);
export const getPendingVaspsError = createSelector(pendingVaspState, (state) => state.error);
export const getPendingVaspsLoadingState = createSelector(
  pendingVaspState,
  (state) => state.loading
);
