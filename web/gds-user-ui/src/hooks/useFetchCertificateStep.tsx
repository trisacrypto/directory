import { useQuery } from '@tanstack/react-query';

import { getCertificateStepService } from 'modules/dashboard/certificate/service/certificateService';
import type { GetCertificateQuery, PayloadDTO } from 'modules/dashboard/certificate/types';

export function useFetchCertificateStep(payload: PayloadDTO): GetCertificateQuery {
  const query = useQuery(
    ['fetch-certficate-step', payload.key],
    () => getCertificateStepService(payload),
    {
      enabled: !!payload.key,
      retry: false
    }
  );
  return {
    getCertificateStep: query.refetch,
    certificateStep: query.data,
    error: query.error,
    hasCertificateStepFailed: query.isError,
    wasCertificateStepFetched: query.isSuccess,
    isFetchingCertificateStep: query.isLoading
  };
}
