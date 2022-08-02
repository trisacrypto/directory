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
  Tag,
  useColorModeValue
} from '@chakra-ui/react';
import { colors } from 'utils/theme';
import { useDispatch, useSelector, RootStateOrAny } from 'react-redux';
import { loadDefaultValueFromLocalStorage, TStep } from 'utils/localStorageHelper';
import useCertificateStepper from 'hooks/useCertificateStepper';
import { getNameIdentiferTypeLabel } from 'constants/name-identifiers';
import { getNationalIdentificationLabel } from 'constants/national-identification';
import { COUNTRIES } from 'constants/countries';
import { renderAddress } from 'utils/address-utils';
import { addressType } from 'constants/address';
import { Trans } from '@lingui/react';
import { getRegistrationDefaultValue } from 'modules/dashboard/registration/utils';

interface LegalReviewProps {
  data?: any;
}
// NOTE: need some clean up.
import LegalPersonReviewDataTable from './LegalPersonReviewDataTable';
import CertificateReviewHeader from './CertificateReviewHeader';
import CertificateReviewLayout from './CertificateReviewLayout';

const LegalPersonReview = () => {
  const steps: TStep[] = useSelector((state: RootStateOrAny) => state.stepper.steps);
  const [legalPerson, setLegalPerson] = React.useState<any>({});
  useEffect(() => {
    const fetchData = async () => {
      const getStepperData = await getRegistrationDefaultValue();
      const stepData = {
        ...getStepperData.entity
      };
      setLegalPerson(stepData);
    };
    fetchData();
  }, [steps]);
  return (
    <CertificateReviewLayout>
      <CertificateReviewHeader step={2} title="Section 2: Legal Person" />
      <LegalPersonReviewDataTable data={legalPerson} />
    </CertificateReviewLayout>
  );
};
LegalPersonReview.defaultProps = {
  data: {}
};
export default LegalPersonReview;
