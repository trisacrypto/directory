import React, { useEffect } from 'react';
import { useSelector } from 'react-redux';
import { t } from '@lingui/macro';
import BasicDetailsReviewDataTable from './BasicDetailsReviewDataTable';
import CertificateReviewHeader from './CertificateReviewHeader';
import CertificateReviewLayout from './CertificateReviewLayout';
import { getCurrentState } from 'application/store/selectors/stepper';
// import useGetStepStatusByKey from './useGetStepStatusByKey';
import RequiredElementMissing from 'components/ErrorComponent/RequiredElementMissing';
import { basicDetailsValidationSchema } from 'modules/dashboard/certificate/lib/basicDetailsValidationSchema';
const BasicDetailsReview = () => {
  const currentStateValue = useSelector(getCurrentState);
  const [isValid, setIsValid] = React.useState(false);

  const basicDetail = {
    organization_name: currentStateValue.data.organization_name,
    website: currentStateValue.data.website,
    established_on: currentStateValue.data.established_on,
    vasp_categories: currentStateValue.data.vasp_categories,
    business_category: currentStateValue.data.business_category
  };

  useEffect(() => {
    const validate = async () => {
      try {
        const r = await basicDetailsValidationSchema.validate(basicDetail, { abortEarly: false });
        setIsValid(true);
        console.log('r', r);
      } catch (error) {
        console.log('error', error);
        setIsValid(false);
      }
    };
    validate();
  }, [currentStateValue?.data]);

  return (
    <CertificateReviewLayout>
      <CertificateReviewHeader step={1} title={t`Section 1: Basic Details`} />
      {!isValid ? <RequiredElementMissing elementKey={1} /> : false}
      <BasicDetailsReviewDataTable data={basicDetail} />
    </CertificateReviewLayout>
  );
};

export default BasicDetailsReview;
