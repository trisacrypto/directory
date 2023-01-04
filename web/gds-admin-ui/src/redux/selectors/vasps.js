const { createSelector } = require('reselect');

const Vasps = (state) => state.Vasps;

const getAllVasps = createSelector(Vasps, (state) => state.data);
const getVaspsLoadingState = createSelector(Vasps, (state) => state.loading);

export { getAllVasps, getVaspsLoadingState };
