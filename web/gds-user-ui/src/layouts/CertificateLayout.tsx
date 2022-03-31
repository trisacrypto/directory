import { Box, Heading, HStack, VStack } from '@chakra-ui/react';
import Card from 'components/ui/Card';
import TestNetCertificateProgressBar from 'components/TestnetProgress/TestNetCertificateProgressBar.component';
import FormButton from 'components/ui/FormButton';
import useCertificateStepper from 'hooks/useCertificateStepper';
import { FormProvider, useForm } from 'react-hook-form';
import { getCertificateRegistrationDefaultValue } from 'utils/form-references';
import { yupResolver } from '@hookform/resolvers/yup';
import { DevTool } from '@hookform/devtools';
import { RootStateOrAny, useSelector } from 'react-redux';
import _ from 'lodash';
import { certificateRegistrationValidationSchema } from 'validation-schemas';

type CertificateLayoutProps = {
  children?: React.ReactNode;
};

const fieldNamesPerSteps = {
  basicDetails: ['website', 'established_on', 'business_category', 'vasp_categories'],
  legalPerson: [
    'entity.name.name_identifiers',
    'entity.name.local_name_identifiers',
    'entity.name.phonetic_name_identifiers',
    'entity.geographic_addresses',
    'entity.country_of_registration',
    'entity.national_identification.national_identifier',
    'entity.national_identification.national_identifier_type',
    'entity.national_identification.country_of_issue',
    'entity.national_identification.registration_authority'
  ],
  contacts: [
    ...['administrative', 'technical', 'billing', 'legal'].flatMap((value) => [
      `contacts.${value}.name`,
      `contacts.${value}.email`,
      `contacts.${value}.phone`
    ])
  ],
  trisaImplementation: [
    ...['trisa_endpoint_testnet', 'trisa_endpoint_mainnet'].flatMap((value) => [
      `${value}.common_name`,
      `${value}.endpoint`
    ])
  ],
  trixoImplementation: [
    'trixo.primary_national_jurisdiction',
    'trixo.primary_regulator',
    'trixo.financial_transfers_permitted',
    'trixo.has_required_regulatory_program',
    'trixo.conducts_customer_kyc',
    'trixo.kyc_threshold',
    'trixo.kyc_threshold_currency',
    'trixo.must_comply_travel_rule',
    'trixo.compliance_threshold',
    'trixo.compliance_threshold_currency',
    'trixo.must_safeguard_pii',
    'trixo.safeguards_pii'
  ]
};

const fieldNamesPerStepsEntries = () => Object.entries(fieldNamesPerSteps);

const CertificateLayout: React.FC<CertificateLayoutProps> = ({ children }) => {
  const { nextStep, previousStep } = useCertificateStepper();
  const currentStep: number = useSelector((state: RootStateOrAny) => state.stepper.currentStep);

  const methods = useForm({
    defaultValues: getCertificateRegistrationDefaultValue(),
    resolver: yupResolver(certificateRegistrationValidationSchema),
    mode: 'onChange'
  });

  function getFieldValue(name: string) {
    return _.get(methods.getValues(), name);
  }

  function isFormCompleted() {
    const fieldsNames = fieldNamesPerStepsEntries()[currentStep - 1][1];
    return fieldsNames.every((n) => getFieldValue(n).toString());
  }

  function hasErroredField() {
    const fieldsNames = fieldNamesPerStepsEntries()[currentStep - 1][1];
    return fieldsNames.some((n: any) => methods.getFieldState(n).error);
  }

  function handleNextStepClick() {
    if (hasErroredField()) {
      // eslint-disable-next-line no-alert
      if (window.confirm('Would you like to continue ?')) {
        nextStep({ isFormCompleted: isFormCompleted(), errors: methods.formState.errors });
      }
    } else {
      nextStep({ isFormCompleted: isFormCompleted(), errors: methods.formState.errors });
    }
  }

  return (
    <>
      <Heading size="lg" mb="24px">
        Certificate Registration
      </Heading>
      <VStack spacing={3}>
        <Card maxW="100%">
          <Card.Body>
            This multi-section form is an important step in the registration and certificate
            issuance process. The information you provide will be used to verify the legal entity
            that you represent and, where appropriate, will be available to verified TRISA members
            to facilitate compliance decisions. To assist in completing the registration form, the
            form is broken into multiple sections. No information is sent until you complete Section
            6 - Review & Submit.
          </Card.Body>
        </Card>
        <Box width={'100%'}>
          <FormProvider {...methods}>
            <TestNetCertificateProgressBar />
            <DevTool control={methods.control} /> {/* set up the dev tool */}
          </FormProvider>
        </Box>
        <Box pt="27px" w="100%" mb="1rem">
          {children}
        </Box>
        <HStack width="100%" justifyContent="space-between" pt={5}>
          <FormButton onClick={() => previousStep()} borderRadius={5} w="100%" maxW="13rem">
            Previous
          </FormButton>
          <FormButton
            borderRadius={5}
            w="100%"
            maxW="13rem"
            backgroundColor="#FF7A59"
            _hover={{ backgroundColor: '#f07253' }}>
            Save & Next
          </FormButton>
          <FormButton onClick={handleNextStepClick} borderRadius={5} w="100%" maxW="13rem">
            Save & Continue Later
          </FormButton>
        </HStack>
      </VStack>
    </>
  );
};

export default CertificateLayout;
