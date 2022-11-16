import type { Collaborator } from 'components/Collaborators/CollaboratorType';
export interface DeleteCollaboratorMutation {
    deleteCollaborator(collaborator: TDeleteCollaborator | string): Promise<string | number> | void;
    reset(): void;
    collaborator?: TDeleteCollaborator;
    hasCollaboratorFailed: boolean;
    wasCollaboratorDeleted: boolean;
    isDeleting: boolean;
    errorMessage?: any;
}

export type TDeleteCollaborator = Pick<Collaborator, 'id'>;
