export interface Collaborator {
    id: string;
    email: string;
    user_id?: string;
    name: string;
    roles: string[];
    created_at?: string;
    modified_at?: string;
    verified_at?: string;
}
