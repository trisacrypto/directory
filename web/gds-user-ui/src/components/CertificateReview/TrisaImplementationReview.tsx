import React, { FC, useEffect } from 'react';
import {
  Stack,
  Box,
  Text,
  Heading,
  Table,
  Tbody,
  Tr,
  Td,
  Button,
  Divider,
  useColorModeValue
} from '@chakra-ui/react';
import { colors } from 'utils/theme';
import { useDispatch, useSelector, RootStateOrAny } from 'react-redux';
import { loadDefaultValueFromLocalStorage, TStep } from 'utils/localStorageHelper';
import useCertificateStepper from 'hooks/useCertificateStepper';
import { Trans } from '@lingui/react';
import { getRegistrationDefaultValue } from 'modules/dashboard/registration/utils';
import TrisaImplementationReviewDataTable from './TrisaImplementationReviewDataTable';
import CertificateReviewHeader from './CertificateReviewHeader';
import CertificateReviewLayout from './CertificateReviewLayout';
import { MdSystemUpdateAlt } from 'react-icons/md';
import { t } from '@lingui/macro';
import { useFormContext } from 'react-hook-form';
import Store from 'application/store';
const TrisaImplementationReview = () => {
  const [trisa, setTrisa] = React.useState<any>({});
  const textColor = useColorModeValue('gray.800', '#F7F8FC');

  useEffect(() => {
    const getStepperData = Store.getState().stepper.data;
    const stepData = {
      mainnet: getStepperData.mainnet,
      testnet: getStepperData.testnet
    };

    setTrisa(stepData);
  }, []);
  return (
    <CertificateReviewLayout>
      <CertificateReviewHeader step={4} title={t`Section 4: TRISA Implementation`} />
      <TrisaImplementationReviewDataTable mainnet={trisa.mainnet} testnet={trisa.testnet} />
    </CertificateReviewLayout>
  );
};

export default TrisaImplementationReview;
