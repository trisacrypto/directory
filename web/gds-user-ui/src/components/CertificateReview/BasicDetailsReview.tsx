import React, { useEffect } from 'react';
import { useColorModeValue } from '@chakra-ui/react';
import { useSelector, RootStateOrAny } from 'react-redux';
import { TStep } from 'utils/localStorageHelper';
import useCertificateStepper from 'hooks/useCertificateStepper';
import { Trans } from '@lingui/react';
import { t } from '@lingui/macro';
import { getRegistrationDefaultValue } from 'modules/dashboard/registration/utils';

import BasicDetailsReviewDataTable from './BasicDetailsReviewDataTable';
import CertificateReviewHeader from './CertificateReviewHeader';
import CertificateReviewLayout from './CertificateReviewLayout';

const BasicDetailsReview = () => {
  const { jumpToStep } = useCertificateStepper();
  const steps: TStep[] = useSelector((state: RootStateOrAny) => state.stepper.steps);
  const [basicDetail, setBasicDetail] = React.useState<any>({});

  useEffect(() => {
    const fetchData = async () => {
      const getStepperData = await getRegistrationDefaultValue();
      const stepData = {
        website: getStepperData.website,
        established_on: getStepperData.established_on,
        vasp_categories: getStepperData.vasp_categories,
        business_category: getStepperData.business_category
      };
      setBasicDetail(stepData);
    };
    fetchData();
  }, [steps]);
  return (
    <CertificateReviewLayout>
      <CertificateReviewHeader step={1} title={t`Section 1: Basic Details`} />
      <BasicDetailsReviewDataTable data={basicDetail} />
    </CertificateReviewLayout>
  );
};

export default BasicDetailsReview;
