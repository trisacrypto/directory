import Store from 'application/store';


export const mapCollaborators = (arr: any) => {
    return arr.map((collaborator: any) => {
        return {
            ...collaborator,
            created_at: new Date(collaborator.created_at),
        };
    });
};

// lis tCollaborators by recent date
export const sortCollaboratorsByRecentDate = (arr: any) => {
    const refactoredArr = mapCollaborators(arr);
    return refactoredArr.sort((a: any, b: any) => b.created_at.getTime() - a.created_at.getTime());
};

// is collaborator current user
export const isCurrentUser = (collaboratorEmail: string): boolean => {
    return collaboratorEmail === Store.getState()?.user?.user?.email;
};


