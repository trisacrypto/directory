import React, { useEffect } from 'react';
import { colors } from 'utils/theme';
import { useSelector, RootStateOrAny } from 'react-redux';
import { TStep } from 'utils/localStorageHelper';
import useCertificateStepper from 'hooks/useCertificateStepper';
import { getRegistrationDefaultValue } from 'modules/dashboard/registration/utils';
import { useFormContext } from 'react-hook-form';
import ContactsReviewDataTable from './ContactsReviewDataTable';
import CertificateReviewHeader from './CertificateReviewHeader';
import CertificateReviewLayout from './CertificateReviewLayout';
import { t } from '@lingui/macro';
import Store from 'application/store';
const ContactsReview = () => {
  const [contacts, setContacts] = React.useState<any>({});

  useEffect(() => {
    const getStepperData = Store.getState().stepper.data;
    const stepData = {
      ...getStepperData.contacts
    };
    setContacts(stepData);
  }, []);
  return (
    <CertificateReviewLayout>
      <CertificateReviewHeader step={3} title={t`Section 3: Contacts`} />
      <ContactsReviewDataTable data={contacts} />
    </CertificateReviewLayout>
  );
};

export default ContactsReview;
