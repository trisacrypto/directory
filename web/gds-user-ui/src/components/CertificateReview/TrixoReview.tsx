import React, { useEffect, Suspense } from 'react';
import { useColorModeValue } from '@chakra-ui/react';

import { t } from '@lingui/macro';
import TrixoReviewDataTable from './TrixoReviewDataTable';
import CertificateReviewHeader from './CertificateReviewHeader';
import CertificateReviewLayout from './CertificateReviewLayout';
import { useFormContext } from 'react-hook-form';
import Store from 'application/store';
import { getCurrentState } from 'application/store/selectors/stepper';
import { RootStateOrAny, useSelector } from 'react-redux';
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
