import React, { useState, useEffect } from 'react';
import { useSelector } from 'react-redux';
import ContactsReviewDataTable from './ContactsReviewDataTable';
import CertificateReviewHeader from './CertificateReviewHeader';
import CertificateReviewLayout from './CertificateReviewLayout';
import { t } from '@lingui/macro';
import { getCurrentState } from 'application/store/selectors/stepper';
// import useGetStepStatusByKey from './useGetStepStatusByKey';
import RequiredElementMissing from 'components/ErrorComponent/RequiredElementMissing';

import { contactsValidationSchema } from 'modules/dashboard/certificate/lib/contactsValidationSchema';
const ContactsReview = () => {
  const currentStateValue = useSelector(getCurrentState);

  const [isValid, setIsValid] = useState(false);

  useEffect(() => {
    const validate = async () => {
      try {
        const r = await contactsValidationSchema.validate(currentStateValue.data.contacts, {
          abortEarly: false
        });
        setIsValid(true);
        console.log('r', r);
      } catch (error) {
        console.log('error', error);
        setIsValid(false);
      }
    };
    validate();
  }, [currentStateValue.data.contacts]);
  return (
    <CertificateReviewLayout>
      <CertificateReviewHeader step={3} title={t`Section 3: Contacts`} />
      {!isValid ? <RequiredElementMissing elementKey={1} /> : false}
      <ContactsReviewDataTable data={currentStateValue.data.contacts} />
    </CertificateReviewLayout>
  );
};

export default ContactsReview;
