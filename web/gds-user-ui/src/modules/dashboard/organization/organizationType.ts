export interface Organization {
    id: string;
    name: string;
    domain?: string;
    createdAt?: string;
    refreshToken?: string;
}

export type NewOrganization = Pick<Organization, 'id' | 'name' | 'domain'>;

export type OrganizationMutation = {
    createOrganization(organization: NewOrganization): void;
    reset(): void;
    organization?: Organization;
    hasOrganizationFailed: boolean;
    wasOrganizationCreated: boolean;
    isCreating: boolean;
    errorMessage?: any;
};

export type OrganizationQuery = {
    getAllOrganizations(): void;
    reset?(): void;
    organizations?: Organization[];
    hasOrganizationFailed: boolean;
    wasOrganizationFetched: boolean;
    isFetching: boolean;
    errorMessage?: any;
};

