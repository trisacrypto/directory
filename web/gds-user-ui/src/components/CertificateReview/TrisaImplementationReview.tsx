import React, { useState, useEffect } from 'react';
import { useSelector } from 'react-redux';
import TrisaImplementationReviewDataTable from './TrisaImplementationReviewDataTable';
import CertificateReviewHeader from './CertificateReviewHeader';
import CertificateReviewLayout from './CertificateReviewLayout';
import { t } from '@lingui/macro';
import { getCurrentState } from 'application/store/selectors/stepper';
// import useGetStepStatusByKey from './useGetStepStatusByKey';
import RequiredElementMissing from 'components/ErrorComponent/RequiredElementMissing';
import { trisaImplementationValidationSchema } from 'modules/dashboard/certificate/lib/trisaImplementationValidationSchema';
const TrisaImplementationReview = () => {
  const currentStateValue = useSelector(getCurrentState);
  const { data: trisaData } = currentStateValue;

  console.log('trisaData', trisaData);

  const [isValid, setIsValid] = useState(false);
  const trisa = {
    mainnet: trisaData.mainnet,
    testnet: trisaData.testnet
  };
  useEffect(() => {
    const validate = async () => {
      try {
        const r = await trisaImplementationValidationSchema.validate(trisa, {
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
  }, [trisa]);

  return (
    <CertificateReviewLayout>
      <CertificateReviewHeader step={4} title={t`Section 4: TRISA Implementation`} />
      {!isValid ? <RequiredElementMissing /> : false}
      <TrisaImplementationReviewDataTable mainnet={trisa.mainnet} testnet={trisa.testnet} />
    </CertificateReviewLayout>
  );
};

export default TrisaImplementationReview;
