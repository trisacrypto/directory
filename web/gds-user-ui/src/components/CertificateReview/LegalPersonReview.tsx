import React, { useState, useEffect, useMemo } from 'react';
import { useSelector } from 'react-redux';

// NOTE: need some clean up.
import LegalPersonReviewDataTable from './LegalPersonReviewDataTable';
import CertificateReviewHeader from './CertificateReviewHeader';
import CertificateReviewLayout from './CertificateReviewLayout';
// import useGetStepStatusByKey from './useGetStepStatusByKey';
import RequiredElementMissing from 'components/ErrorComponent/RequiredElementMissing';
import { t } from '@lingui/macro';
import { getCurrentState } from 'application/store/selectors/stepper';
import { legalPersonValidationSchemam } from 'modules/dashboard/certificate/lib/legalPersonValidationSchema';
const LegalPersonReview = () => {
  const currentStateValue = useSelector(getCurrentState);
  const [isValid, setIsValid] = useState(false);

  const legalPerson = useMemo(() => {
    return {
      ...currentStateValue.data.entity
    };
  }, [currentStateValue.data.entity]);

  useEffect(() => {
    const validate = async () => {
      try {
        await legalPersonValidationSchemam.validate(legalPerson, { abortEarly: false });
        setIsValid(true);
      } catch (error) {
        setIsValid(false);
      }
    };
    validate();
  }, [legalPerson]);
  return (
    <CertificateReviewLayout>
      <CertificateReviewHeader step={2} title={t`Section 2: Legal Person`} />
      {!isValid ? <RequiredElementMissing elementKey={2} /> : false}
      <LegalPersonReviewDataTable data={legalPerson} />
    </CertificateReviewLayout>
  );
};
LegalPersonReview.defaultProps = {
  data: {}
};
export default LegalPersonReview;
