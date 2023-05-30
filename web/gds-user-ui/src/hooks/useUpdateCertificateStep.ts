import { useMutation } from '@tanstack/react-query';
import { queryClient } from 'utils/react-query';
import { postCertificateStepService } from 'modules/dashboard/certificate/service/certificateService';
import type { PostCertificateMutation } from 'modules/dashboard/certificate/types';

export function useUpdateCertificateStep(): PostCertificateMutation {
  const mutation = useMutation(['update-registration-form'], postCertificateStepService, {
    onSuccess: () => {
      // queryClient.setQueryData(['fetch-certificate-step'], mutation.data);
      queryClient.invalidateQueries({ queryKey: ['fetch-certificate-step'] });
    }
  });

  return {
    updateCertificateStep: mutation.mutate,
    updatedCertificateStep: mutation.data,
    hasCertificateStepFailed: mutation.isError,
    wasCertificateStepUpdated: mutation.isSuccess,
    isUpdatingCertificateStep: mutation.isLoading,
    error: mutation.error,
    reset: mutation.reset
  };
}
