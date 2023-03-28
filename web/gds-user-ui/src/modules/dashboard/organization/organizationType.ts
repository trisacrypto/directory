export interface Organization {
  id: string;
  name: string;
  domain?: string;
  createdAt?: string;
  refreshToken?: string;
}

export interface OrganizationResponse {
  count: number;
  organizations: Organization[];
  page: number;
  page_size: number;
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
  organizations: OrganizationResponse;
  hasOrganizationFailed: boolean;
  wasOrganizationFetched: boolean;
  isFetching: boolean;
  errorMessage?: any;
};

export type OrganizationPagination = {
  page: number;
  pageSize: number;
};
