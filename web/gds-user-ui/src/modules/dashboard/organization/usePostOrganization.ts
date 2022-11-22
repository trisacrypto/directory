import { AxiosError } from 'axios';
import { CreateOrganisation } from './organizationService';
import { useMutation } from '@tanstack/react-query';

export function usePostOrganizations() {
  return useMutation<any, AxiosError>(CreateOrganisation);
}
