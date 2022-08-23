import React, { useEffect } from 'react';
import { useSelector, RootStateOrAny } from 'react-redux';
import { useFormContext } from 'react-hook-form';
import { getRegistrationDefaultValue } from 'modules/dashboard/registration/utils';

// NOTE: need some clean up.
import LegalPersonReviewDataTable from './LegalPersonReviewDataTable';
import CertificateReviewHeader from './CertificateReviewHeader';
import CertificateReviewLayout from './CertificateReviewLayout';
import { t } from '@lingui/macro';
import Store from 'application/store';
const LegalPersonReview = () => {
  const [legalPerson, setLegalPerson] = React.useState<any>({});
  useEffect(() => {
    const getStepperData = Store.getState().stepper.data;
    const stepData = {
      ...getStepperData.entity
    };
    setLegalPerson(stepData);
  }, []);
  return (
    <CertificateReviewLayout>
      <CertificateReviewHeader step={2} title={t`Section 2: Legal Person`} />
      <LegalPersonReviewDataTable data={legalPerson} />
    </CertificateReviewLayout>
  );
};
LegalPersonReview.defaultProps = {
  data: {}
};
export default LegalPersonReview;
