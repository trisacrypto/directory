export interface getCollaborators {
    getAllCollaborators(): void;
    collaborators: any;
    hasCollaboratorsFailed: boolean;
    wasCollaboratorsFetched: boolean;
    isFetchingCollaborators: boolean;
    errorMessage?: any;
}
