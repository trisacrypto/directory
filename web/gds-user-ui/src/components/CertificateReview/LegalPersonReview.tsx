import React from 'react';
import LegalPersonReviewDataTable from './LegalPersonReviewDataTable';
import CertificateReviewHeader from './CertificateReviewHeader';
import CertificateReviewLayout from './CertificateReviewLayout';
import RequiredElementMissing from 'components/ErrorComponent/RequiredElementMissing';
import { t } from '@lingui/macro';
import { StepEnum } from 'types/enums';
import { useFetchCertificateStep } from 'hooks/useFetchCertificateStep';

const LegalPersonReview = () => {
  const { certificateStep } = useFetchCertificateStep({
    key: StepEnum.LEGAL
  });

  const hasErrors = certificateStep?.errors;

  return (
    <CertificateReviewLayout>
      <CertificateReviewHeader step={2} title={t`Section 2: Legal Person`} />
      {hasErrors ? <RequiredElementMissing elementKey={2} /> : false}
      <LegalPersonReviewDataTable data={certificateStep?.form?.entity} />
    </CertificateReviewLayout>
  );
};

export default LegalPersonReview;
