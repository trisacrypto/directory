import { useMutation } from '@tanstack/react-query';
import { queryClient } from 'utils/react-query';
import { deleteCertificateStepService } from 'modules/dashboard/certificate/service/certificateService';
import type { DeleteCertificateMutation } from 'modules/dashboard/certificate/types';

export function useDeleteCertificateStep(): DeleteCertificateMutation {
  const mutation = useMutation(deleteCertificateStepService, {
    onSuccess: () => {
      // queryClient.setQueryData(['fetch-certificate-step'], mutation.data);
      queryClient.invalidateQueries({ queryKey: ['fetch-certificate-step'] });
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
