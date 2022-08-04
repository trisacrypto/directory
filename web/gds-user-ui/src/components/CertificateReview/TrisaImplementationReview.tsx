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

interface TrisaImplementationReviewProps {
  mainnetData?: any;
  testnetData?: any;
}
const TrisaImplementationReview = (props: TrisaImplementationReviewProps) => {
  const { jumpToStep } = useCertificateStepper();
  const steps: TStep[] = useSelector((state: RootStateOrAny) => state.stepper.steps);
  const [trisa, setTrisa] = React.useState<any>({});
  const textColor = useColorModeValue('gray.800', '#F7F8FC');

  useEffect(() => {
    const fetchData = async () => {
      const getStepperData = await getRegistrationDefaultValue();
      const stepData = {
        mainnet: getStepperData.mainnet,
        testnet: getStepperData.testnet
      };

      setTrisa(stepData);
    };
    fetchData();
  }, [steps]);
  return (
    <CertificateReviewLayout>
      <CertificateReviewHeader step={4} title="Section 4: TRISA Implementation" />
      <TrisaImplementationReviewDataTable mainnet={trisa.mainnet} testnet={trisa.testnet} />
    </CertificateReviewLayout>
  );
};

export default TrisaImplementationReview;
