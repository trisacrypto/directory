
export interface Collaborator {
    id: string;
    email: string;
    user_id?: string;
    name: string;
    roles: string[];
    created_at?: string;
    modified_at?: string;
    organization?: string;
    verified_at?: string;
    status?: TCollaboratorStatus;
}


export type TCollaboratorStatus = 'pending' | 'completed' | 'rejected' | 'revoked';
