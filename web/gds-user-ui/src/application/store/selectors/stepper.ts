import { RootStateOrAny } from 'react-redux';
import { createSelector } from 'reselect';

const rootState = (state: RootStateOrAny) => state.stepper;

export const getCurrentStep = createSelector(rootState, (state) => state.currentStep);
export const getSteps = createSelector(rootState, (state) => state.steps);
export const getLastStep = createSelector(rootState, (state) => state.lastStep);
export const resetStepper = createSelector(rootState, (state) => state.clearStepper);
