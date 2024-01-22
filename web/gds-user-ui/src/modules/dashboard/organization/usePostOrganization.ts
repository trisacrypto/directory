import { useMutation } from '@tanstack/react-query';
import { createOrganization } from './organizationService';

export function usePostOrganizations() {
  return useMutation<any, any>(createOrganization);
}
