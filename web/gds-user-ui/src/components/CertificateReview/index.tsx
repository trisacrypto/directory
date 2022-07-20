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

const CertificateReview = () => {
  const toast = useToast();

  const hasReachSubmitStep: boolean = useSelector(
    (state: RootStateOrAny) => state.stepper.hasReachSubmitStep
  );
  const [isTestNetSent, setIsTestNetSent] = useState(false);
  const [isMainNetSent, setIsMainNetSent] = useState(false);
  const [result, setResult] = useState('');
  const handleSubmitRegister = async (event: React.FormEvent, network: string) => {
    event.preventDefault();
    try {
      const formValue = loadDefaultValueFromLocalStorage();
      const getMainnetObj = formValue.trisa_endpoint_mainnet;
      const getTestnetObj = formValue.trisa_endpoint_testnet;

      delete formValue.trisa_endpoint_mainnet;
      delete formValue.trisa_endpoint_testnet;
      if (network === 'testnet') {
        formValue.trisa_endpoint = getTestnetObj.endpoint;
        formValue.common_name = getTestnetObj.common_name;
      }
      if (network === 'mainnet') {
        formValue.trisa_endpoint = getMainnetObj.endpoint;
        formValue.common_name = getMainnetObj.common_name;
      }

      const response = await registrationRequest(network, mapTrixoFormForBff(formValue));

      if (response.id || response.status === ' "SUBMITTED"') {
        if (network === 'testnet') {
          setIsTestNetSent(true);
          localStorage.setItem('isTestNetSent', 'true');
        }
        if (network === 'mainnet') {
          setIsMainNetSent(true);
          localStorage.setItem('isMainNetSent', 'true');
        }
        setResult(response);
      }
    } catch (err: any) {
      // should send error to sentry

      if (!err.response.data.success) {
        toast({
          position: 'top-right',
          title: t`Error Submitting Certificate`,
          description: err.response.data.error,
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
