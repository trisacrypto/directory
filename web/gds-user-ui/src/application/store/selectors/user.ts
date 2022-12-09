import { RootStateOrAny } from 'react-redux';
import { createSelector } from 'reselect';

const rootState = (state: RootStateOrAny) => state.user;

export const getCurrentUserName = createSelector(rootState, (state) => state.user.name);
export const userSelector = createSelector(rootState, (state) => state.user);
