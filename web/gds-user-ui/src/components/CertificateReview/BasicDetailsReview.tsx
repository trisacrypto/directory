import React, { useEffect } from 'react';
import { useColorModeValue } from '@chakra-ui/react';
import { useSelector } from 'react-redux';
import { TStep } from 'utils/localStorageHelper';
import useCertificateStepper from 'hooks/useCertificateStepper';
import { Trans } from '@lingui/react';
import { t } from '@lingui/macro';
import { getRegistrationDefaultValue } from 'modules/dashboard/registration/utils';
import { useFormContext } from 'react-hook-form';
import BasicDetailsReviewDataTable from './BasicDetailsReviewDataTable';
import CertificateReviewHeader from './CertificateReviewHeader';
import CertificateReviewLayout from './CertificateReviewLayout';
import { getCurrentState } from 'application/store/selectors/stepper';
const BasicDetailsReview = () => {
  const currentStateValue = useSelector(getCurrentState);

  const basicDetail = {
    website: currentStateValue.data.website,
    established_on: currentStateValue.data.established_on,
    vasp_categories: currentStateValue.data.vasp_categories,
    business_category: currentStateValue.data.business_category
  };

  return (
    <CertificateReviewLayout>
      <CertificateReviewHeader step={1} title={t`Section 1: Basic Details`} />
      <BasicDetailsReviewDataTable data={basicDetail} />
    </CertificateReviewLayout>
  );
};

export default BasicDetailsReview;
