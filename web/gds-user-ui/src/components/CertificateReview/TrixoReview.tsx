import React, { Suspense } from 'react';

import { t } from '@lingui/macro';
import TrixoReviewDataTable from './TrixoReviewDataTable';
import CertificateReviewHeader from './CertificateReviewHeader';
import CertificateReviewLayout from './CertificateReviewLayout';
import { getCurrentState } from 'application/store/selectors/stepper';
import { useSelector } from 'react-redux';
const TrixoReview: React.FC = () => {
  const currentStateValue = useSelector(getCurrentState);
  const trixo = {
    ...currentStateValue.data.trixo
  };

  return (
    <CertificateReviewLayout>
      <CertificateReviewHeader title={t`Section 5: TRIXO Questionnaire`} step={5} />
      <Suspense fallback={'Loading trixo data'}>
        <TrixoReviewDataTable data={trixo} />
      </Suspense>
    </CertificateReviewLayout>
  );
};

export default TrixoReview;
