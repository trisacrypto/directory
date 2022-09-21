import { RootStateOrAny } from 'react-redux';
import { createSelector } from 'reselect';

const rootState = (state: RootStateOrAny) => state.stepper;

export const getCurrentStep = createSelector(rootState, (state) => state.currentStep);
export const getSteps = createSelector(rootState, (state) => state.steps);
export const getLastStep = createSelector(rootState, (state) => state.lastStep);
export const resetStepper = createSelector(rootState, (state) => state.clearStepper);
export const getHasReachSubmitStep = createSelector(rootState, (state) => state.hasReachSubmitStep);
export const getCurrentState = createSelector(rootState, (state) => state);
export const getTestNetSubmittedStatus = createSelector(
  rootState,
  (state) => state.testnetSubmitted
);
export const getMainNetSubmittedStatus = createSelector(
  rootState,
  (state) => state.mainnetSubmitted
);
export const getCertificateData = createSelector(rootState, (state) => state.data);
export const getHasReachedSubmitStep = createSelector(
  rootState,
  (state) => state.hasReachSubmitStep
);
export const getHasReachedReviewStep = createSelector(
  rootState,
  (state) => state.hasReachReviewStep
);
