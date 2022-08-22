import React, { useEffect } from 'react';
import { useColorModeValue } from '@chakra-ui/react';
import { useSelector, RootStateOrAny } from 'react-redux';
import { TStep } from 'utils/localStorageHelper';
import useCertificateStepper from 'hooks/useCertificateStepper';
import { t } from '@lingui/macro';
import { getRegistrationDefaultValue } from 'modules/dashboard/registration/utils';
import TrixoReviewDataTable from './TrixoReviewDataTable';
import CertificateReviewHeader from './CertificateReviewHeader';
import CertificateReviewLayout from './CertificateReviewLayout';
interface TrixoReviewProps {
  data: any;
}

const TrixoReview: React.FC<TrixoReviewProps> = ({ data }) => {
  // const textColor = useColorModeValue('gray.800', '#F7F8FC');
  // const getColorScheme = (status: string | boolean) => {
  //   if (status === 'yes' || status === true) {
  //     return 'green';
  //   } else {
  //     return 'orange';
  //   }
  // };

  console.log('[Called] TrixoReview.tsx');

  return (
    <CertificateReviewLayout>
      <CertificateReviewHeader title={t`Section 5: TRIXO Questionnaire`} step={5} />
      <TrixoReviewDataTable data={data} />
    </CertificateReviewLayout>
  );
};

export default TrixoReview;
