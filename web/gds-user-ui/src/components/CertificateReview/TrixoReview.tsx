import React, { useEffect } from 'react';
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
  Tag,
  TagLabel,
  useColorModeValue
} from '@chakra-ui/react';
import { colors } from 'utils/theme';
import { useSelector, RootStateOrAny } from 'react-redux';
import { loadDefaultValueFromLocalStorage, TStep } from 'utils/localStorageHelper';
import useCertificateStepper from 'hooks/useCertificateStepper';
import { COUNTRIES } from 'constants/countries';
import { currencyFormatter } from 'utils/utils';
import { Trans } from '@lingui/react';
import { t } from '@lingui/macro';
import TrixoReviewDataTable from './TrixoReviewDataTable';
import CertificateReviewHeader from './CertificateReviewHeader';
import CertificateReviewLayout from './CertificateReviewLayout';
interface TrixoReviewProps {}

const TrixoReview: React.FC<TrixoReviewProps> = (props) => {
  const { jumpToStep } = useCertificateStepper();
  const steps: TStep[] = useSelector((state: RootStateOrAny) => state.stepper.steps);
  const [trixo, setTrixo] = React.useState<any>({});
  const textColor = useColorModeValue('gray.800', '#F7F8FC');
  const getColorScheme = (status: string | boolean) => {
    if (status === 'yes' || status === true) {
      return 'green';
    } else {
      return 'orange';
    }
  };
  useEffect(() => {
    const getStepperData = loadDefaultValueFromLocalStorage();
    const stepData = {
      ...getStepperData.trixo
    };
    setTrixo(stepData);
  }, [steps]);

  return (
    <CertificateReviewLayout>
      <CertificateReviewHeader title="Section 5: TRIXO Questionnaire" step={5} />
      <TrixoReviewDataTable />
    </CertificateReviewLayout>
  );
};

export default TrixoReview;
