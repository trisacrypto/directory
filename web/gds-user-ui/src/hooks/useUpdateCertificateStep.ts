import { useMutation } from '@tanstack/react-query';

import { postCertificateStepService } from 'modules/dashboard/certificate/service/certificateService';
import type { PostCertificateMutation } from 'modules/dashboard/certificate/types';

export function useUpdateCertificateStep(): PostCertificateMutation {
  const mutation = useMutation(['update-certificate-step'], postCertificateStepService);

  return {
    updateCertificateStep: mutation.mutate,
    certificateStep: mutation.data,
    hasCertificateStepFailed: mutation.isError,
    wasCertificateStepUpdated: mutation.isSuccess,
    isUpdatingCertificateStep: mutation.isLoading,
    error: mutation.error,
    reset: mutation.reset
  };
}
