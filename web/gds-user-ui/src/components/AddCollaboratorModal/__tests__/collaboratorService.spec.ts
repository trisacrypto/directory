import { getAllCollaborators, createCollaborator } from '../CollaboratorService';

describe('ApiService', () => {
    it('should return all collaborators', async () => {
        const response = await getAllCollaborators();
        expect(response).toBeDefined();
    });

    it('should create a collaborator', async () => {
        const response = await createCollaborator({
            name: 'test',
            email: 'test@test.io',
        });

        expect(response).toBeDefined();
    });
});
