import { SimpleDashboardLayout } from 'layouts';
import { Box, Heading, HStack, VStack } from '@chakra-ui/react';
import Card from 'components/ui/Card';
import TestNetCertificateProgressBar from 'components/TestnetProgress/TestNetCertificateProgressBar.component';
import FormButton from 'components/ui/FormButton';
import useCertificateStepper from 'hooks/useCertificateStepper';
import { FormProvider, useForm, useFormState } from 'react-hook-form';
import { getCertificateRegistrationDefaultValue } from 'modules/dashboard/certificate/lib/form-references';
import { yupResolver } from '@hookform/resolvers/yup';
import { DevTool } from '@hookform/devtools';
import { RootStateOrAny, useSelector } from 'react-redux';
import _ from 'lodash';
import { certificateRegistrationValidationSchema } from './lib/certificate-registration-validation-schema';

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

const Certificate: React.FC = () => {
  const { nextStep, previousStep } = useCertificateStepper();
  const currentStep: number = useSelector((state: RootStateOrAny) => state.stepper.currentStep);
  const lastStep: number = useSelector((state: RootStateOrAny) => state.stepper.lastStep);
  const hasReachSubmitStep: boolean = useSelector(
    (state: RootStateOrAny) => state.stepper.hasReachSubmitStep
  );
  const current = currentStep === lastStep ? lastStep - 1 : currentStep;
  console.log('current', current);
  const methods = useForm({
    defaultValues: getCertificateRegistrationDefaultValue(),
    resolver: yupResolver(certificateRegistrationValidationSchema),
    mode: 'onChange'
  });

  function getFieldValue(name: string) {
    return _.get(methods.getValues(), name);
  }

  function isFormCompleted() {
    const fieldsNames = fieldNamesPerStepsEntries()[current - 1][1];
    return fieldsNames.every((n) => getFieldValue(n).toString());
  }

  function getCurrentFormValue() {
    const fieldsNames = fieldNamesPerStepsEntries()[current - 1][1];
    return fieldsNames.reduce((acc, n) => ({ ...acc, [n]: getFieldValue(n) }), {});
  }

  function hasErroredField() {
    const fieldsNames = fieldNamesPerStepsEntries()[current - 1][1];
    return fieldsNames.some((n: any) => methods.getFieldState(n).error);
  }

  function handleNextStepClick() {
    if (hasErroredField()) {
      console.log('last step errr');
      // i think we should not use alert here , but we need to find a way to display the error message
      // eslint-disable-next-line no-alert
      if (window.confirm('Some requirement are missing , Would you like to continue ?')) {
        nextStep({
          isFormCompleted: isFormCompleted(),
          errors: methods.formState.errors,
          formValues: getCurrentFormValue()
        });
      }
    } else {
      nextStep({
        isFormCompleted: isFormCompleted(),
        formValues: getCurrentFormValue()
      });
    }
  }
  const handlePreviousStep = () => {
    previousStep();
  };

  return (
    // <DashboardLayout>
    //   <CertificateLayout>
    //     <BasicDetails />
    //   </CertificateLayout>
    // </DashboardLayout>
    <SimpleDashboardLayout>
      <>
        <Heading size="lg" mb="24px">
          Certificate Registration
        </Heading>
        <VStack spacing={3}>
          <Card maxW="100%" bg={'white'}>
            <Card.Body>
              This multi-section form is an important step in the registration and certificate
              issuance process. The information you provide will be used to verify the legal entity
              that you represent and, where appropriate, will be available to verified TRISA members
              to facilitate compliance decisions. To assist in completing the registration form, the
              form is broken into multiple sections. No information is sent until you complete
              Section 6 - Review & Submit.
            </Card.Body>
          </Card>
          <Box width={'100%'}>
            <FormProvider {...methods}>
              <TestNetCertificateProgressBar />
              <DevTool control={methods.control} /> {/* setting up the hook form dev tool */}
            </FormProvider>
          </Box>
          {!hasReachSubmitStep && (
            <HStack width="100%" spacing={4} justifyContent={'center'} pt={4}>
              <FormButton
                onClick={handlePreviousStep}
                isDisabled={currentStep === 1}
                borderRadius={5}
                w="100%"
                maxW="13rem">
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
                {currentStep === lastStep ? 'Finish & submit' : 'Save & Continue Later'}
              </FormButton>
            </HStack>
          )}
        </VStack>
      </>
    </SimpleDashboardLayout>
  );
};

export default Certificate;
