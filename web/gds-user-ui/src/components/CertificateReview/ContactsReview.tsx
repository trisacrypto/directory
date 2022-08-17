import React, { useEffect } from 'react';
import { colors } from 'utils/theme';
import { useSelector, RootStateOrAny } from 'react-redux';
import { TStep } from 'utils/localStorageHelper';
import useCertificateStepper from 'hooks/useCertificateStepper';
import { getRegistrationDefaultValue } from 'modules/dashboard/registration/utils';

import ContactsReviewDataTable from './ContactsReviewDataTable';
import CertificateReviewHeader from './CertificateReviewHeader';
import CertificateReviewLayout from './CertificateReviewLayout';
import { t } from '@lingui/macro';

const ContactsReview = () => {
  const { jumpToStep } = useCertificateStepper();
  const steps: TStep[] = useSelector((state: RootStateOrAny) => state.stepper.steps);
  const [contacts, setContacts] = React.useState<any>({});
  useEffect(() => {
    const fetchData = async () => {
      const getStepperData = await getRegistrationDefaultValue();
      const stepData = {
        ...getStepperData.contacts
      };
      setContacts(stepData);
    };
    fetchData();
  }, [steps]);

  return (
    <CertificateReviewLayout>
      <CertificateReviewHeader step={3} title={t`Section 3: Contacts`} />
      <ContactsReviewDataTable data={contacts} />
    </CertificateReviewLayout>
  );
};

export default ContactsReview;
