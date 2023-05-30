import { useMutation } from '@tanstack/react-query';
import { queryClient } from 'utils/react-query';
import { deleteCertificateStepService } from 'modules/dashboard/certificate/service/certificateService';
import type { DeleteCertificateMutation } from 'modules/dashboard/certificate/types';

export function useDeleteCertificateStep(key?: string): DeleteCertificateMutation {
  const mutation = useMutation(['delete-certificate-step'], deleteCertificateStepService, {
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['fetch-certificate-step', key] });
    }
  });

  return {
    deleteCertificateStep: mutation.mutate,
    deletedCertificateStep: mutation.data,
    hasCertificateStepFailed: mutation.isError,
    wasCertificateStepDeleted: mutation.isSuccess,
    isDeletingCertificateStep: mutation.isLoading,
    error: mutation.error,
    reset: mutation.reset
  };
}
