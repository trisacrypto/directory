import React, { useState, useEffect } from 'react';

import { t } from '@lingui/macro';
import TrixoReviewDataTable from './TrixoReviewDataTable';
import CertificateReviewHeader from './CertificateReviewHeader';
import CertificateReviewLayout from './CertificateReviewLayout';
import RequiredElementMissing from 'components/ErrorComponent/RequiredElementMissing';
import { StepEnum } from 'types/enums';
import { useFetchCertificateStep } from 'hooks/useFetchCertificateStep';
// import { ErrorBoundary } from '@sentry/react';
const TrixoReview: React.FC = () => {
  const { certificateStep, wasCertificateStepFetched } = useFetchCertificateStep({
    key: StepEnum.TRIXO
  });

  const [trixoData, setTrixoData] = useState([]);

  const hasErrors = certificateStep?.errors;

  useEffect(() => {
    if (
      wasCertificateStepFetched &&
      certificateStep?.form?.trixo &&
      certificateStep?.step === StepEnum.TRIXO
    ) {
      setTrixoData(certificateStep?.form?.trixo);
    }
  }, [certificateStep?.form, wasCertificateStepFetched, certificateStep?.step]);

  return (
    <CertificateReviewLayout>
      <CertificateReviewHeader title={t`Section 5: TRIXO Questionnaire`} step={5} />
      {hasErrors ? <RequiredElementMissing elementKey={5} errorFields={hasErrors} /> : false}

      {/* <ErrorBoundary
        fallback={
          <div>
            <p>Something went wrong</p>
          </div>
        }> */}
      <TrixoReviewDataTable data={trixoData} />
      {/* </ErrorBoundary> */}
    </CertificateReviewLayout>
  );
};

export default TrixoReview;
