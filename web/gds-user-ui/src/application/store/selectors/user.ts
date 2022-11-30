import { RootStateOrAny } from 'react-redux';
import { createSelector } from 'reselect';

const rootState = (state: RootStateOrAny) => state.user;

export const getCurrentStep = createSelector(rootState, (state) => state.currentStep);
export const userSelector = createSelector(rootState, (state) => state.user);
