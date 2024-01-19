import { CreateOrganisation } from './organizationService';
import { useMutation } from '@tanstack/react-query';

export function usePostOrganizations() {
  return useMutation<any, any>(CreateOrganisation);
}
