import React, { useEffect } from 'react';
import { useColorModeValue } from '@chakra-ui/react';
import { useSelector, RootStateOrAny } from 'react-redux';
import { TStep } from 'utils/localStorageHelper';
import useCertificateStepper from 'hooks/useCertificateStepper';
import { Trans } from '@lingui/react';
import { t } from '@lingui/macro';
import { getRegistrationDefaultValue } from 'modules/dashboard/registration/utils';
import { useFormContext } from 'react-hook-form';
import BasicDetailsReviewDataTable from './BasicDetailsReviewDataTable';
import CertificateReviewHeader from './CertificateReviewHeader';
import CertificateReviewLayout from './CertificateReviewLayout';
import Store from 'application/store';
const BasicDetailsReview = () => {
  const [basicDetail, setBasicDetail] = React.useState<any>({});
  // get basic details from the store
  useEffect(() => {
    const getStepperData = Store.getState().stepper.data;
    const stepData = {
      website: getStepperData.website,
      established_on: getStepperData.established_on,
      vasp_categories: getStepperData.vasp_categories,
      business_category: getStepperData.business_category
    };
    setBasicDetail(stepData);
  }, []);
  return (
    <CertificateReviewLayout>
      <CertificateReviewHeader step={1} title={t`Section 1: Basic Details`} />
      <BasicDetailsReviewDataTable data={basicDetail} />
    </CertificateReviewLayout>
  );
};

export default BasicDetailsReview;
