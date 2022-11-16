import type { Collaborator } from 'components/Collaborators/CollaboratorType';
export interface CollaboratorMutation {
    createCollaborator(collaborator: NewCollaborator): void;
    reset(): void;
    collaborator?: Collaborator;
    hasCollaboratorFailed: boolean;
    wasCollaboratorCreated: boolean;
    isCreating: boolean;
    errorMessage?: any;
}

export type NewCollaborator = Pick<Collaborator, 'name' | 'email'>;
