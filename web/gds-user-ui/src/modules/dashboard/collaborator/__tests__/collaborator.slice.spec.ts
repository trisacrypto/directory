// add a test for the collaborator slice reducer
import { collaboratorReducer, getCollaborators, setCollaborator, updateCollaborator, deleteCollaborator, collaboratorInitialState } from 'modules/dashboard/collaborator/collaborator.slice';
import { getCollaboratorState } from 'application/store/selectors/collaborator';

describe('collaborator reducer', () => {
    it('should handle initial state', () => {
        const expected = [{
            name: '',
            email: '',
            roles: [],
            id: '',
            created_at: ''
        }];
        expect(collaboratorReducer(undefined, { type: 'unknown' })).toEqual(expected);
    });

    it('should handle setCollaborator', () => {
        const collaborator = {
            name: 'test',
            email: 'test@localhost.dev',
            roles: ['test'],
            id: '123',
            created_at: '2021-01-01'
        };
        const actual = collaboratorReducer(collaboratorInitialState, setCollaborator(collaborator));
        expect(actual).toContainEqual(collaborator);
    });

    it('should handle updateCollaborator', () => {
        const collaborator = {
            name: 'test-update',
            email: 'test@local.dev',
            roles: ['test'],
            id: '123',
            created_at: '2021-01-01'
        };
        const actual = collaboratorReducer([collaborator], updateCollaborator(collaborator));
        expect(actual).toEqual([collaborator]);
    });

    // get all collaborators
    it('should handle getCollaborators', () => {
        const actual = collaboratorReducer(getCollaboratorState, getCollaborators());
        expect(actual).toEqual(getCollaboratorState);
    });

    it('should handle deleteCollaborator', () => {
        const collaborator = {
            id: '123'
        };
        const actual = collaboratorReducer([collaborator], deleteCollaborator(collaborator));
        expect(actual).toEqual([collaborator]);
    });
});

