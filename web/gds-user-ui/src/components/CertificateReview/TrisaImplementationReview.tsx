import React from 'react';
import TrisaImplementationReviewDataTable from './TrisaImplementationReviewDataTable';
import CertificateReviewHeader from './CertificateReviewHeader';
import CertificateReviewLayout from './CertificateReviewLayout';
import { t } from '@lingui/macro';
import RequiredElementMissing from 'components/ErrorComponent/RequiredElementMissing';
import { StepEnum } from 'types/enums';
import { useFetchCertificateStep } from 'hooks/useFetchCertificateStep';

const TrisaImplementationReview = () => {
  const { certificateStep } = useFetchCertificateStep({
    key: StepEnum.TRISA
  });

  const hasErrors = certificateStep?.errors;

  return (
    <CertificateReviewLayout>
      <CertificateReviewHeader step={4} title={t`Section 4: TRISA Implementation`} />
      {hasErrors ? <RequiredElementMissing elementKey={4} /> : false}
      <TrisaImplementationReviewDataTable
        mainnet={certificateStep?.form?.mainnet}
        testnet={certificateStep?.form?.testnet}
      />
    </CertificateReviewLayout>
  );
};

export default TrisaImplementationReview;
