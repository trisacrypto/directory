import React, { useState, useEffect, useMemo } from 'react';
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

  const [isValid, setIsValid] = useState(false);
  const trisa = useMemo(() => {
    return {
      mainnet: trisaData.mainnet,
      testnet: trisaData.testnet
    };
  }, [trisaData.mainnet, trisaData.testnet]);

  useEffect(() => {
    const validate = async () => {
      try {
        const test = await trisaImplementationValidationSchema.validate(trisa, {
          abortEarly: false
        });
        console.log('[TrisaImplementationReview] test', test);
        setIsValid(true);
      } catch (error) {
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
