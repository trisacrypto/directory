import React, { useEffect } from 'react';
import { useSelector, RootStateOrAny } from 'react-redux';

import { getRegistrationDefaultValue } from 'modules/dashboard/registration/utils';

// NOTE: need some clean up.
import LegalPersonReviewDataTable from './LegalPersonReviewDataTable';
import CertificateReviewHeader from './CertificateReviewHeader';
import CertificateReviewLayout from './CertificateReviewLayout';
import { t } from '@lingui/macro';

const LegalPersonReview = () => {
  const steps: TStep[] = useSelector((state: RootStateOrAny) => state.stepper.steps);
  const [legalPerson, setLegalPerson] = React.useState<any>({});
  useEffect(() => {
    const fetchData = async () => {
      const getStepperData = await getRegistrationDefaultValue();
      const stepData = {
        ...getStepperData.entity
      };
      setLegalPerson(stepData);
    };
    fetchData();
  }, [steps]);
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
