import { Heading, HStack, Stack, Text, Link } from '@chakra-ui/react';

import FormLayout from 'layouts/FormLayout';
import LegalForm from './LegalForm';
import { useSelector } from 'react-redux';
import { getCurrentStep, getSteps } from 'application/store/selectors/stepper';
import { getStepStatus } from 'utils/utils';
import { SectionStatus } from 'components/SectionStatus';
import { Trans } from '@lingui/react';

const LegalPerson: React.FC = () => {
  const steps = useSelector(getSteps);
  const currentStep = useSelector(getCurrentStep);
  const stepStatus = getStepStatus(steps, currentStep);

  return (
    <Stack spacing={7} mt="2rem">
      <HStack>
        <Heading size="md" pr={3} ml={2}>
          <Trans id={'Section 2: Legal Person'}>Section 2: Legal Person</Trans>
        </Heading>
        {stepStatus ? <SectionStatus status={stepStatus} /> : null}
      </HStack>
      <FormLayout>
        <Text>
          <Trans id="Please enter the information that identifies your organization as a Legal Person. This form represents the">
            Please enter the information that identifies your organization as a Legal Person. This
            form represents the
          </Trans>{' '}
          <Link isExternal href="https://intervasp.org/" color={'blue'} fontWeight={'bold'}>
            {' '}
            <Trans id="IVMS 101 data structure">IVMS 101 data structure</Trans>
          </Link>{' '}
          <Trans
            id={
              'for legal persons and is strongly suggested for use as KYC or CDD (Know your Counterparty) information exchanged in TRISA transfers.'
            }>
            for legal persons and is strongly suggested for use as KYC or CDD (Know your
            Counterparty) information exchanged in TRISA transfers.
          </Trans>
        </Text>
      </FormLayout>
      <LegalForm />
    </Stack>
  );
};

export default LegalPerson;
