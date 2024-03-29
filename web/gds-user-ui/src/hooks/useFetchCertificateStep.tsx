import { useQuery } from '@tanstack/react-query';

import { getCertificateStepService } from 'modules/dashboard/certificate/service/certificateService';
import type { GetCertificateQuery, PayloadDTO } from 'modules/dashboard/certificate/types';

export function useFetchCertificateStep(payload: PayloadDTO): GetCertificateQuery {
  const query = useQuery(
    ['fetch-certificate-step', payload.key],
    () => getCertificateStepService(payload),
    {
      enabled: !!payload.key,
      // disable caching
      // this is why we had issue with the stepper not updating
      cacheTime: 0
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
