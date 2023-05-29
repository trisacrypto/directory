import React, { useMemo } from 'react';
import { t } from '@lingui/macro';
import BasicDetailsReviewDataTable from './BasicDetailsReviewDataTable';
import CertificateReviewHeader from './CertificateReviewHeader';
import CertificateReviewLayout from './CertificateReviewLayout';
// import useGetStepStatusByKey from './useGetStepStatusByKey';
import { StepEnum } from 'types/enums';
import { useFetchCertificateStep } from 'hooks/useFetchCertificateStep';
import RequiredElementMissing from 'components/ErrorComponent/RequiredElementMissing';
const BasicDetailsReview = () => {
  const { certificateStep } = useFetchCertificateStep({
    key: StepEnum.BASIC
  });

  const hasErrors = certificateStep?.errors;

  const basicDetail = useMemo(() => {
    return {
      organization_name: certificateStep?.form?.organization_name,
      website: certificateStep?.form?.website,
      established_on: certificateStep?.form?.established_on,
      vasp_categories: certificateStep?.form?.vasp_categories,
      business_category: certificateStep?.form?.business_category
    };
  }, [certificateStep?.form]);

  return (
    <CertificateReviewLayout>
      <CertificateReviewHeader step={1} title={t`Section 1: Basic Details`} />
      {hasErrors ? <RequiredElementMissing elementKey={1} errorFields={hasErrors} /> : false}
      <BasicDetailsReviewDataTable data={basicDetail} />
    </CertificateReviewLayout>
  );
};

export default BasicDetailsReview;
