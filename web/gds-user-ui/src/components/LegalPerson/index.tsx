import { Heading, HStack, Stack, Text, Link, chakra } from '@chakra-ui/react';
import CountryOfRegistration from 'components/CountryOfRegistration';
import FormLayout from 'layouts/FormLayout';
import NameIdentifiers from '../NameIdentifiers';
import NationalIdentification from '../NameIdentification';
import Address from 'components/Addresses';
import { useSelector } from 'react-redux';
import { getCurrentStep, getSteps } from 'application/store/selectors/stepper';
import { getStepStatus } from 'utils/utils';
import { SectionStatus } from 'components/SectionStatus';
import { Trans } from '@lingui/react';
import { FormProvider, useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import { legalPersonValidationSchemam } from 'modules/dashboard/certificate/lib/legalPersonValidationSchema';
import useCertificateStepper from 'hooks/useCertificateStepper';
import StepButtons from 'components/StepsButtons';
const LegalPerson: React.FC = () => {
  const { previousStep } = useCertificateStepper();
  const steps = useSelector(getSteps);
  const currentStep = useSelector(getCurrentStep);
  const stepStatus = getStepStatus(steps, currentStep);

  const resolver = yupResolver(legalPersonValidationSchemam);
  const methods = useForm({
    defaultValues: {},
    resolver,

    mode: 'onChange'
  });

  const handleNextStepClick = () => {
    console.log('[] methods values', methods.getValues());
  };

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
      <FormProvider {...methods}>
        <chakra.form onSubmit={methods.handleSubmit(handleNextStepClick)}>
          <NameIdentifiers />
          <Address />
          <CountryOfRegistration />
          <NationalIdentification />
          <StepButtons handlePreviousStep={previousStep} />
        </chakra.form>
      </FormProvider>
    </Stack>
  );
};

export default LegalPerson;
