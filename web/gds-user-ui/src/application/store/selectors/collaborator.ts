import { RootStateOrAny } from 'react-redux';
import { createSelector } from 'reselect';
// import type { Collaborator } from 'components/Collaborators/CollaboratorType';
const rootState = (state: RootStateOrAny) => state.collaborators;


export const getCollaboratorState = createSelector(rootState, (state) => state);
export const setCollaborators = createSelector(rootState, (state) => state.setCollaborator);
export const getCollaborators = createSelector(rootState, (state) => state.getCollaborators);
