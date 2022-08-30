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
import { getCurrentState } from 'application/store/selectors/stepper';

const TrisaImplementationReview = () => {
  const currentStateValue = useSelector(getCurrentState);
  const { data: trisaData } = currentStateValue;

  const trisa = {
    mainnet: trisaData.mainnet,
    testnet: trisaData.testnet
  };

  return (
    <CertificateReviewLayout>
      <CertificateReviewHeader step={4} title={t`Section 4: TRISA Implementation`} />
      <TrisaImplementationReviewDataTable mainnet={trisa.mainnet} testnet={trisa.testnet} />
    </CertificateReviewLayout>
  );
};

export default TrisaImplementationReview;
