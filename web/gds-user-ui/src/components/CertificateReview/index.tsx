/* eslint-disable prefer-reflect */

import React, { useState } from 'react';
import { useToast } from '@chakra-ui/react';

import { RootStateOrAny, useSelector } from 'react-redux';
import ReviewSubmit from 'components/ReviewSubmit';
import { registrationRequest } from 'modules/dashboard/certificate/service';
import { loadDefaultValueFromLocalStorage } from 'utils/localStorageHelper';
import { t } from '@lingui/macro';
import ReviewsSummary from './ReviewsSummary';
import { mapTrixoFormForBff } from 'utils/utils';
import {
  submitMainnetRegistration,
  submitTestnetRegistration
} from 'modules/dashboard/registration/service';
import useCertificateStepper from 'hooks/useCertificateStepper';

const CertificateReview = () => {
  const toast = useToast();
  const { testnetSubmissionState, mainnetSubmissionState } = useCertificateStepper();

  const hasReachSubmitStep: boolean = useSelector(
    (state: RootStateOrAny) => state.stepper.hasReachSubmitStep
  );
  const [isTestNetSent, setIsTestNetSent] = useState(false);
  const [isMainNetSent, setIsMainNetSent] = useState(false);
  const [result, setResult] = useState('');
  const handleSubmitRegister = async (event: React.FormEvent, network: string) => {
    event.preventDefault();
    try {
      if (network === 'testnet') {
        const response = await submitTestnetRegistration();
        if (response.status === 200) {
          setIsTestNetSent(true);
          testnetSubmissionState();
          setResult(response?.data);
        }
      }
      if (network === 'mainnet') {
        const response = await submitMainnetRegistration();
        if (response?.status === 200) {
          setIsMainNetSent(true);
          mainnetSubmissionState();
          setResult(response?.data);
        }
      }
    } catch (err: any) {
      console.log('[err catched]', err);
      if (!err?.response?.data?.success) {
        console.log('[err catched]', err?.response.data.error);
        toast({
          position: 'top-right',
          title: t`Error Submitting Certificate`,
          description: t`${err?.response?.data?.error}`,
          status: 'error',
          duration: 5000,
          isClosable: true
        });
      } else {
        console.error('something went wrong');
      }
    }
  };

  if (!hasReachSubmitStep) {
    return <ReviewsSummary />;
  }

  return (
    <ReviewSubmit
      onSubmitHandler={handleSubmitRegister}
      isTestNetSent={isTestNetSent}
      isMainNetSent={isMainNetSent}
      result={result}
    />
  );
};

export default CertificateReview;
