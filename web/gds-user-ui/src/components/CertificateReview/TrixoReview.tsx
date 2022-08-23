import React, { useEffect, Suspense } from 'react';
import { useColorModeValue } from '@chakra-ui/react';

import { t } from '@lingui/macro';
import TrixoReviewDataTable from './TrixoReviewDataTable';
import CertificateReviewHeader from './CertificateReviewHeader';
import CertificateReviewLayout from './CertificateReviewLayout';
import { useFormContext } from 'react-hook-form';
import Store from 'application/store';
const TrixoReview: React.FC = () => {
  const [trixo, setTrixo] = React.useState<any>({});

  useEffect(() => {
    // wait 1 second before getting the state
    setTimeout(() => {
      const getStepperData = Store.getState().stepper.data;
      const stepData = {
        ...getStepperData.trixo
      };
      setTrixo(stepData);
    }, 1000);
  }, []);

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
