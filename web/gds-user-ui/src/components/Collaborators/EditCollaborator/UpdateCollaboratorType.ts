import type { Collaborator } from 'components/Collaborators/CollaboratorType';
export interface UpdateCollaboratorMutation {
    updateCollaborator(data: any): Promise<string | number> | void;
    reset(): void;
    collaborator?: Collaborator;
    hasCollaboratorFailed: boolean;
    wasCollaboratorUpdated: boolean;
    isUpdating: boolean;
    errorMessage?: any;
}

export type TUpdateCollaborator = Pick<Collaborator, 'roles' | 'id'>;
