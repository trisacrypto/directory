import { createSlice } from '@reduxjs/toolkit';
import type { Collaborator } from 'components/Collaborators/CollaboratorType';
type CollaboratorsType = Omit<Collaborator, 'user_id' | 'modified_at' | 'verified_at'>[];
type CollaboratorType = Omit<Collaborator, 'user_id' | 'modified_at' | 'verified_at'>;
export const initialValue: CollaboratorsType = [{
    name: '',
    email: '',
    roles: [],
    id: '',
    created_at: ''
}];

const collaboratorSlice: any = createSlice({
    name: 'collaborators',
    initialState: initialValue,
    reducers: {
        // get all current collaborators
        getCollaborators: (state: any, { }: any) => {
            return state;
        },

        // add collaborator to the list
        setCollaborator: (state: any, { payload }: any) => {
            state.push(payload);
        },
        // init a new collaborators array
        setCollaborators: (state: any, { payload }: any) => {
            state.collaborators = payload;
        },
        // update collaborator in the list
        updateCollaborator: (state: any, { payload }: any) => {
            state.map((collaborator: CollaboratorType) => {
                if (collaborator.id === payload.id) {
                    // eslint-disable-next-line no-param-reassign
                    collaborator = payload;
                }
            });
        },
        // delete collaborator from the list
        deleteCollaborator: (state: any, { payload }: any) => {
            state.map((collaborator: any) => {
                if (collaborator.id === payload.id) {
                    // eslint-disable-next-line no-param-reassign
                    collaborator = payload;
                }
            });
        }


    }
});

export const {
    getCollaborators,
    setCollaborator,
    updateCollaborator,
    deleteCollaborator
} = collaboratorSlice.actions;
export const collaboratorReducer = collaboratorSlice.reducer;
export const collaboratorsSelector = (state: any) => state.collaborators;
export const collaboratorInitialState = collaboratorSlice.initialState;


